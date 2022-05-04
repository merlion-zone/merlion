package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/voter/types"
)

func (k Keeper) SetTotalVotes(ctx sdk.Context, votes sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{votes})
	store.Set(types.TotalVotesKey(), bz)
}

func (k Keeper) GetTotalVotes(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TotalVotesKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var votes sdk.IntProto
	k.cdc.MustUnmarshal(bz, &votes)
	return votes.Int
}

func (k Keeper) SetTotalVotesByUser(ctx sdk.Context, veID uint64, votes sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{votes})
	store.Set(types.TotalVotesByUserKey(veID), bz)
}

func (k Keeper) GetTotalVotesByUser(ctx sdk.Context, veID uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.TotalVotesByUserKey(veID))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var votes sdk.IntProto
	k.cdc.MustUnmarshal(bz, &votes)
	return votes.Int
}

func (k Keeper) DeleteTotalVotesByUser(ctx sdk.Context, veID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.TotalVotesByUserKey(veID))
}

func (k Keeper) SetPoolWeightedVotes(ctx sdk.Context, poolDenom string, votes sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{votes})
	store.Set(types.PoolWeightedVotesKey(poolDenom), bz)
}

func (k Keeper) GetPoolWeightedVotes(ctx sdk.Context, poolDenom string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PoolWeightedVotesKey(poolDenom))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var votes sdk.IntProto
	k.cdc.MustUnmarshal(bz, &votes)
	return votes.Int
}

func (k Keeper) SetPoolWeightedVotesByUser(ctx sdk.Context, veID uint64, poolDenom string, votes sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{votes})
	store.Set(types.PoolWeightedVotesByUserKey(veID, poolDenom), bz)
}

func (k Keeper) GetPoolWeightedVotesByUser(ctx sdk.Context, veID uint64, poolDenom string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PoolWeightedVotesByUserKey(veID, poolDenom))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var votes sdk.IntProto
	k.cdc.MustUnmarshal(bz, &votes)
	return votes.Int
}

func (k Keeper) DeletePoolWeightedVotesByUser(ctx sdk.Context, veID uint64, poolDenom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PoolWeightedVotesByUserKey(veID, poolDenom))
}

func (k Keeper) SetIndex(ctx sdk.Context, index sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{index})
	store.Set(types.IndexKey(), bz)
}

func (k Keeper) GetIndex(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.IndexKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var index sdk.IntProto
	k.cdc.MustUnmarshal(bz, &index)
	return index.Int
}

func (k Keeper) SetIndexAtLastUpdatedByGauge(ctx sdk.Context, poolDenom string, index sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{index})
	store.Set(types.IndexAtLastUpdatedByGaugeKey(poolDenom), bz)
}

func (k Keeper) GetIndexAtLastUpdatedByGauge(ctx sdk.Context, poolDenom string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.IndexAtLastUpdatedByGaugeKey(poolDenom))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var index sdk.IntProto
	k.cdc.MustUnmarshal(bz, &index)
	return index.Int
}

func (k Keeper) SetClaimableRewardByGauge(ctx sdk.Context, poolDenom string, claimable sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{claimable})
	store.Set(types.ClaimableRewardByGaugeKey(poolDenom), bz)
}

func (k Keeper) GetClaimableRewardByGauge(ctx sdk.Context, poolDenom string) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ClaimableRewardByGaugeKey(poolDenom))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var claimable sdk.IntProto
	k.cdc.MustUnmarshal(bz, &claimable)
	return claimable.Int
}
