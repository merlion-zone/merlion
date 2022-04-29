package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

func (k Keeper) SetTotalLockedAmount(ctx sdk.Context, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{amount})
	store.Set(types.TotalLockedAmountKey(), bz)
}

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

func (k Keeper) SetLockedAmountByUser(ctx sdk.Context, veID uint64, amount types.LockedBalance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&amount)
	store.Set(types.LockedAmountByUserKey(veID), bz)
}

func (k Keeper) GetLockedAmountByUser(ctx sdk.Context, veID uint64) types.LockedBalance {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.LockedAmountByUserKey(veID))
	if bz == nil {
		return types.LockedBalance{
			Amount: merlion.ZeroInt,
			End:    0,
		}
	}
	var amount types.LockedBalance
	k.cdc.MustUnmarshal(bz, &amount)
	return amount
}
