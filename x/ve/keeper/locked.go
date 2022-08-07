package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

// SetTotalLockedAmount sets total locked amount
func (k Keeper) SetTotalLockedAmount(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: amount})
	store.Set(types.TotalLockedAmountKey(), bz)
}

// GetTotalLockedAmount gets total locked amount
func (k Keeper) GetTotalLockedAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TotalLockedAmountKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.IntProto
	k.cdc.MustUnmarshal(bz, &amount)
	return amount.Int
}

// SetLockedAmountByUser sets locked amount of the specified ve
func (k Keeper) SetLockedAmountByUser(ctx sdk.Context, veID uint64, amount types.LockedBalance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&amount)
	store.Set(types.LockedAmountByUserKey(veID), bz)
}

// GetLockedAmountByUser Gets locked amount of the specified ve
func (k Keeper) GetLockedAmountByUser(ctx sdk.Context, veID uint64) types.LockedBalance {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LockedAmountByUserKey(veID))
	if bz == nil {
		return types.NewLockedBalance()
	}
	var amount types.LockedBalance
	k.cdc.MustUnmarshal(bz, &amount)
	return amount
}

// DeleteLockedAmountByUser deletes locked amount of the specified ve
func (k Keeper) DeleteLockedAmountByUser(ctx sdk.Context, veID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.LockedAmountByUserKey(veID))
}
