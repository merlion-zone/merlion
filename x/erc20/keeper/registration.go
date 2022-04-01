package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/tharsis/evmos/v3/contracts"
)

// RegisterCoin deploys an erc20 contract and creates the token pair for the existing cosmos coin
func (k Keeper) RegisterCoin(ctx sdk.Context, denom string) (*types.TokenPair, error) {
	coinMetadata, found := k.bankKeeper.GetDenomMetaData(ctx, denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrEVMDenom, "cannot get metadata of denom %s", denom)
	}

	// Prohibit denominations that contain the "lion" denom
	if strings.Contains(coinMetadata.Base, "lion") {
		return nil, sdkerrors.Wrapf(types.ErrEVMDenom, "cannot register the EVM denomination %s", coinMetadata.Base)
	}

	// Check if the denomination already registered
	if k.IsDenomRegistered(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", coinMetadata.Base)
	}

	// Check if the coin exists by ensuring the supply is set
	if !k.bankKeeper.HasSupply(ctx, coinMetadata.Base) {
		return nil, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins,
			"base denomination '%s' cannot have a supply of 0", coinMetadata.Base,
		)
	}

	addr, err := k.DeployERC20Contract(ctx, coinMetadata)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to create wrapped coin denom metadata for ERC20")
	}

	pair := types.NewTokenPair(addr, coinMetadata.Base, types.OWNER_MODULE)
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())

	return &pair, nil
}

// DeployERC20Contract creates and deploys an ERC20 contract on the EVM with the
// erc20 module account as owner.
func (k Keeper) DeployERC20Contract(
	ctx sdk.Context,
	coinMetadata banktypes.Metadata,
) (common.Address, error) {
	decimals := uint8(coinMetadata.DenomUnits[0].Exponent)
	ctorArgs, err := contracts.ERC20MinterBurnerDecimalsContract.ABI.Pack(
		"",
		coinMetadata.Name,
		coinMetadata.Symbol,
		decimals,
	)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(types.ErrABIPack, "coin metadata is invalid %s: %s", coinMetadata.Name, err.Error())
	}

	data := make([]byte, len(contracts.ERC20MinterBurnerDecimalsContract.Bin)+len(ctorArgs))
	copy(data[:len(contracts.ERC20MinterBurnerDecimalsContract.Bin)], contracts.ERC20MinterBurnerDecimalsContract.Bin)
	copy(data[len(contracts.ERC20MinterBurnerDecimalsContract.Bin):], ctorArgs)

	nonce, err := k.accountKeeper.GetSequence(ctx, types.ModuleAddress.Bytes())
	if err != nil {
		return common.Address{}, err
	}

	contractAddr := crypto.CreateAddress(types.ModuleAddress, nonce)
	_, err = k.CallEVMWithData(ctx, types.ModuleAddress, nil, data)
	if err != nil {
		return common.Address{}, sdkerrors.Wrapf(err, "failed to deploy contract for %s", coinMetadata.Name)
	}

	return contractAddr, nil
}

// RegisterERC20 registers the token pair between the coin and the ERC20
func (k Keeper) RegisterERC20(ctx sdk.Context, contract common.Address) (types.TokenPair, error) {
	if k.IsERC20Registered(ctx, contract) {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrTokenPairAlreadyExists, "token ERC20 contract already registered: %s", contract.String())
	}

	metadata, ok := k.GetDenomMetaData(ctx, types.CreateDenom(contract.String()))
	if !ok {
		return types.TokenPair{}, sdkerrors.Wrap(types.ErrInternalTokenPair, "failed to get denom metadata for ERC20")
	}

	pair := types.NewTokenPair(contract, metadata.Name, types.OWNER_EXTERNAL)
	k.SetTokenPair(ctx, pair)
	k.SetDenomMap(ctx, pair.Denom, pair.GetID())
	k.SetERC20Map(ctx, common.HexToAddress(pair.Erc20Address), pair.GetID())
	return pair, nil
}
