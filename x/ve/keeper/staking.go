package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

func (k Keeper) SlashLockedAmountByUser(ctx sdk.Context, veID uint64, amount sdk.Int) {
	// burn amount
	coin := sdk.NewCoin(k.LockDenom(ctx), amount)
	err := k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		panic(err)
	}

	// update total locked
	totalLocked := k.GetTotalLockedAmount(ctx)
	totalLocked = totalLocked.Sub(amount)
	k.SetTotalLockedAmount(ctx, totalLocked)

	locked := k.GetLockedAmountByUser(ctx, veID)

	lockedOld := locked
	locked.Amount = lockedOld.Amount.Sub(amount)

	// update locked for veID
	k.SetLockedAmountByUser(ctx, veID, locked)

	// regulate user checkpoint for veID,
	// also regulate system checkpoint
	k.RegulateUserCheckpoint(ctx, veID, lockedOld, locked)
}

func (k Keeper) SetGetDelegatedAmountByUser(getDelegatedAmount func(ctx sdk.Context, veID uint64) sdk.Int) {
	k.getDelegatedAmount = getDelegatedAmount
}
