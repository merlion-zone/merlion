package keeper

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	evmtypes "github.com/tharsis/ethermint/x/evm/types"

	"github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/tharsis/evmos/v3/contracts"
)

// EvmHooks wrapper struct for erc20 keeper
type EvmHooks struct {
	k Keeper
}

var _ evmtypes.EvmHooks = EvmHooks{}

// EvmHooks returns the wrapper struct
func (k Keeper) EvmHooks() EvmHooks {
	return EvmHooks{k}
}

// PostTxProcessing implements evmtypes.EvmHooks.PostTxProcessing
func (h EvmHooks) PostTxProcessing(
	ctx sdk.Context,
	from common.Address,
	to *common.Address,
	receipt *ethtypes.Receipt,
) error {
	erc20 := contracts.ERC20MinterBurnerDecimalsContract.ABI

	// We only care about event Transfer(sender, recipient, amount)
	for i, log := range receipt.Logs {
		if len(log.Topics) < 3 {
			continue
		}

		eventID := log.Topics[0] // event ID

		event, err := erc20.EventByID(eventID)
		if err != nil {
			// invalid event for ERC20
			h.k.Logger(ctx).Error("failed to get event by id", "error", err.Error())
			continue
		}

		if event.Name != types.ERC20EventTransfer {
			h.k.Logger(ctx).Debug("emitted event", "name", event.Name, "signature", event.Sig)
			continue
		}

		transferEvent, err := erc20.Unpack(event.Name, log.Data)
		if err != nil {
			h.k.Logger(ctx).Error("failed to unpack transfer event", "error", err.Error())
			continue
		}

		if len(transferEvent) == 0 {
			continue
		}

		tokens, ok := transferEvent[0].(*big.Int)
		// safety check and ignore if amount not positive
		if !ok || tokens == nil || tokens.Sign() != 1 {
			continue
		}

		contractAddr := log.Address
		from := common.BytesToAddress(log.Topics[1].Bytes())
		to := common.BytesToAddress(log.Topics[2].Bytes())

		var pair types.TokenPair
		id := h.k.GetERC20Map(ctx, contractAddr)
		if len(id) == 0 {
			// No token is registered for the caller contract,
			// so it must be native erc20 token.

			pair, err = h.k.RegisterERC20(ctx, contractAddr)
			if err != nil {
				h.k.Logger(ctx).Error("failed to register token pair", "error", err.Error())
				continue
			}
		} else {
			var found bool
			pair, found = h.k.GetTokenPair(ctx, id)
			if !found {
				panic(sdkerrors.Wrapf(types.ErrTokenPairNotFound, "token pair '%s' with contract '%s' not found", id, contractAddr))
			}
		}

		if pair.ContractOwner != types.OWNER_MODULE {
			// Only native coin needs bank operations
			continue
		}

		// Create the corresponding sdk.Coin that is paired with ERC20
		coins := sdk.Coins{{Denom: pair.Denom, Amount: sdk.NewIntFromBigInt(tokens)}}

		zeroAddr := common.Address{}
		zeroFrom := bytes.Equal(from.Bytes(), zeroAddr.Bytes())
		zeroTo := bytes.Equal(to.Bytes(), zeroAddr.Bytes())

		if zeroFrom && !zeroTo {
			err = h.k.bankKeeper.MintCoins(ctx, types.ModuleName, coins)
			if err == nil {
				err = h.k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(to.Bytes()), coins)
			}
		} else if !zeroFrom && zeroTo {
			err = h.k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(from.Bytes()), types.ModuleName, coins)
			if err == nil {
				err = h.k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
			}
		} else if !zeroFrom && !zeroTo {
			err = h.k.bankKeeper.SendCoins(ctx, sdk.AccAddress(from.Bytes()), sdk.AccAddress(to.Bytes()), coins)
		}

		if err != nil {
			h.k.Logger(ctx).Error(
				"failed to process EVM hook for erc20 -> mint/burn/send coins",
				"txHash", receipt.TxHash.Hex(), "logIndex", i,
				"coin", pair.Denom, "contract", contractAddr, "error", err.Error(),
			)
			return err
		}
	}

	return nil
}
