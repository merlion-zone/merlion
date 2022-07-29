package keeper

import (
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k Keeper) InputOutputCoins(ctx sdk.Context, inputs []banktypes.Input, outputs []banktypes.Output) error {
	// Safety check ensuring that when sending coins the keeper must maintain the
	// Check supply invariant and validity of coins.
	if err := banktypes.ValidateInputsOutputs(inputs, outputs); err != nil {
		return err
	}

	moduleAddress := authtypes.NewModuleAddress(banktypes.ModuleName)

	for _, in := range inputs {
		inAddress, err := sdk.AccAddressFromBech32(in.Address)
		if err != nil {
			return err
		}

		nativeCoins, nativeErc20Tokens := k.erc20Keeper().SplitCoinsByErc20(in.Coins)

		err = k.SubUnlockedCoins(ctx, inAddress, nativeCoins)
		if err != nil {
			return err
		}

		// emit coin spent event
		ctx.EventManager().EmitEvent(
			banktypes.NewCoinSpentEvent(inAddress, nativeErc20Tokens),
		)

		err = k.erc20Keeper().SendCoins(ctx, inAddress, moduleAddress, nativeCoins, nativeErc20Tokens)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(banktypes.AttributeKeySender, in.Address),
			),
		)
	}

	for _, out := range outputs {
		outAddress, err := sdk.AccAddressFromBech32(out.Address)
		if err != nil {
			return err
		}

		nativeCoins, nativeErc20Tokens := k.erc20Keeper().SplitCoinsByErc20(out.Coins)

		err = k.AddCoins(ctx, outAddress, nativeCoins)
		if err != nil {
			return err
		}

		// emit coin received event
		ctx.EventManager().EmitEvent(
			banktypes.NewCoinReceivedEvent(outAddress, nativeErc20Tokens),
		)

		err = k.erc20Keeper().SendCoins(ctx, moduleAddress, outAddress, nativeCoins, nativeErc20Tokens)
		if err != nil {
			return err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				banktypes.EventTypeTransfer,
				sdk.NewAttribute(banktypes.AttributeKeyRecipient, out.Address),
				sdk.NewAttribute(sdk.AttributeKeyAmount, out.Coins.String()),
			),
		)

		// Create account if recipient does not exist.
		accExists := k.ak.HasAccount(ctx, outAddress)
		if !accExists {
			defer telemetry.IncrCounter(1, "new", "account")
			k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, outAddress))
		}
	}

	return nil
}

func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	nativeCoins, nativeErc20Tokens := k.erc20Keeper().SplitCoinsByErc20(amt)

	err := k.BaseKeeper.SubUnlockedCoins(ctx, fromAddr, nativeCoins)
	if err != nil {
		return err
	}

	err = k.BaseKeeper.AddCoins(ctx, toAddr, nativeCoins)
	if err != nil {
		return err
	}

	// emit coin spent event
	ctx.EventManager().EmitEvent(
		banktypes.NewCoinSpentEvent(fromAddr, nativeErc20Tokens),
	)

	// emit coin received event
	ctx.EventManager().EmitEvent(
		banktypes.NewCoinReceivedEvent(toAddr, nativeErc20Tokens),
	)

	err = k.erc20Keeper().SendCoins(ctx, fromAddr, toAddr, nativeCoins, nativeErc20Tokens)
	if err != nil {
		return err
	}

	// Create account if recipient does not exist.
	accExists := k.ak.HasAccount(ctx, toAddr)
	if !accExists {
		defer telemetry.IncrCounter(1, "new", "account")
		k.ak.SetAccount(ctx, k.ak.NewAccountWithAddress(ctx, toAddr))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			banktypes.EventTypeTransfer,
			sdk.NewAttribute(banktypes.AttributeKeyRecipient, toAddr.String()),
			sdk.NewAttribute(banktypes.AttributeKeySender, fromAddr.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, amt.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(banktypes.AttributeKeySender, fromAddr.String()),
		),
	})

	return nil
}
