package keeper

import (
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/tharsis/evmos/v3/contracts"
)

func (k Keeper) IsDenomForErc20(denom string) bool {
	denomSplit := strings.Split(denom, "/")
	if len(denomSplit) != 2 {
		return false
	}
	if denomSplit[0] != types.DenomPrefix {
		return false
	}
	return common.IsHexAddress(denomSplit[1])
}

func (k Keeper) getContractByDenom(ctx sdk.Context, denom string) (common.Address, bool) {
	denomSplit := strings.Split(denom, "/")
	if len(denomSplit) == 2 && denomSplit[0] == types.DenomPrefix && common.IsHexAddress(denomSplit[1]) {
		return common.HexToAddress(denomSplit[1]), true
	} else {
		id := k.GetTokenPairID(ctx, denom)
		if len(id) == 0 {
			return common.Address{}, false
		}
		pair, found := k.GetTokenPair(ctx, id)
		if !found {
			panic(sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token pair '%s' with denom '%s' not found", id, denom))
		}
		return pair.GetERC20Contract(), true
	}
}

func (k Keeper) GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	if k.IsDenomForErc20(denom) {
		contract, found := k.getContractByDenom(ctx, denom)
		if !found {
			return banktypes.Metadata{}, false
		}

		strContract := contract.String()

		erc20Data, err := k.QueryERC20(ctx, contract)
		if err != nil {
			return banktypes.Metadata{}, false
		}

		// Base denomination
		base := types.CreateDenom(strContract)

		// Create a bank denom metadata based on the ERC20 token ABI details
		metadata := banktypes.Metadata{
			Description: types.CreateDenomDescription(strContract),
			Base:        base,
			DenomUnits: []*banktypes.DenomUnit{
				{
					Denom:    base,
					Exponent: 0,
				},
			},
			Name:    types.CreateDenom(strContract),
			Symbol:  erc20Data.Symbol,
			Display: base,
		}

		// Only append metadata if decimals > 0, otherwise validation fails
		if erc20Data.Decimals > 0 {
			nameSanitized := types.SanitizeERC20Name(erc20Data.Name)
			metadata.DenomUnits = append(
				metadata.DenomUnits,
				&banktypes.DenomUnit{
					Denom:    nameSanitized,
					Exponent: uint32(erc20Data.Decimals),
				},
			)
			metadata.Display = nameSanitized
		}

		if err := metadata.Validate(); err != nil {
			return banktypes.Metadata{}, false
		}

		return metadata, true
	} else {
		return k.bankKeeper.GetDenomMetaData(ctx, denom)
	}
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	contract, found := k.getContractByDenom(ctx, denom)
	if !found {
		return sdk.Coin{}
	}
	amt, err := k.totalSupply(ctx, contract)
	if err != nil {
		return sdk.Coin{}
	}
	return sdk.NewCoin(denom, sdk.NewIntFromBigInt(amt))
}

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	contract, found := k.getContractByDenom(ctx, denom)
	if !found {
		return sdk.Coin{}
	}
	amt, err := k.balanceOf(ctx, contract, common.BytesToAddress(addr))
	if err != nil {
		return sdk.Coin{}
	}
	return sdk.NewCoin(denom, sdk.NewIntFromBigInt(amt))
}

func (k Keeper) SplitCoinsByErc20(amt sdk.Coins) (nativeCoins sdk.Coins, nativeErc20Tokens sdk.Coins) {
	for _, coin := range amt {
		if k.IsDenomForErc20(coin.Denom) {
			nativeErc20Tokens = nativeErc20Tokens.Add(coin)
		} else {
			nativeCoins = nativeCoins.Add(coin)
		}
	}
	return
}

func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, nativeCoins sdk.Coins, nativeErc20Tokens sdk.Coins) (err error) {
	erc20 := contracts.ERC20BurnableContract.ABI
	from := common.BytesToAddress(fromAddr)
	to := common.BytesToAddress(toAddr)

	var tokenContracts []common.Address
	var amounts []*big.Int

	for _, coin := range nativeCoins {
		id := k.GetTokenPairID(ctx, coin.Denom)
		if len(id) == 0 {
			return sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token pair with denom '%s' not found", coin.Denom)
		}
		pair, found := k.GetTokenPair(ctx, id)
		if !found {
			panic(sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token pair '%s' with denom '%s' not found", id, coin.Denom))
		}
		tokenContracts = append(tokenContracts, pair.GetERC20Contract())
		amounts = append(amounts, coin.Amount.BigInt())
	}

	for _, coin := range nativeErc20Tokens {
		denomSplit := strings.Split(coin.Denom, "/")
		contract := common.HexToAddress(denomSplit[1]) // had checked preceding
		tokenContracts = append(tokenContracts, contract)
		amounts = append(amounts, coin.Amount.BigInt())
	}

	for i, contract := range tokenContracts {
		amt := amounts[i]
		if len(fromAddr) == 0 {
			// Mint
			_, err = k.CallEVM(ctx, erc20, types.ModuleAddress, contract, "mint", to, amt)
		} else if len(toAddr) == 0 {
			// Burn
			_, err = k.CallEVM(ctx, erc20, from, contract, "burn", amt)
		} else {
			// Transfer
			_, err = k.CallEVM(ctx, erc20, from, contract, "transfer", to, amt)
		}
	}

	return err
}
