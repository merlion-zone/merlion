package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

// RegulateUserCheckpoint regulates user checkpoint.
// veID:       must be valid ve id
// lockedOld:
//             Amount: can be zero
//             End: can be zero
// lockedNew:
//             Amount: can be zero
//             End: can be zero
func (k Keeper) RegulateUserCheckpoint(ctx sdk.Context, veID uint64, lockedOld types.LockedBalance, lockedNew types.LockedBalance) {
	types.CheckRegulatedUnixTime(lockedOld.End)
	types.CheckRegulatedUnixTime(lockedNew.End)

	now := uint64(ctx.BlockTime().Unix())

	userPointOld := types.Checkpoint{
		Bias:  sdk.ZeroInt(),
		Slope: sdk.ZeroInt(),
	}
	userPointNew := types.Checkpoint{
		Bias:  sdk.ZeroInt(),
		Slope: sdk.ZeroInt(),
	}

	if lockedOld.End > now && lockedOld.Amount.IsPositive() {
		userPointOld.Slope = lockedOld.Amount.QuoRaw(types.MaxLockTime)
		userPointOld.Bias = userPointOld.Slope.MulRaw(int64(lockedOld.End - now))
	}
	if lockedNew.End > now && lockedNew.Amount.IsPositive() {
		userPointNew.Slope = lockedNew.Amount.QuoRaw(types.MaxLockTime)
		userPointNew.Bias = userPointNew.Slope.MulRaw(int64(lockedNew.End - now))
	}

	slopeChangeOld := k.GetSlopeChange(ctx, lockedOld.End)
	slopeChangeNew := k.GetSlopeChange(ctx, lockedNew.End)

	k.regulateCheckpoint(ctx, userPointNew.Slope.Sub(userPointOld.Slope), userPointNew.Bias.Sub(userPointOld.Bias))

	if lockedOld.End > now {
		slopeChangeOld = slopeChangeOld.Add(userPointOld.Slope)
		if lockedNew.End == lockedOld.End {
			slopeChangeOld = slopeChangeOld.Sub(userPointNew.Slope)
		}
		k.SetSlopeChange(ctx, lockedOld.End, slopeChangeOld)
	}
	if lockedNew.End > now {
		if lockedNew.End > lockedOld.End {
			slopeChangeNew = slopeChangeNew.Sub(userPointNew.Slope)
			k.SetSlopeChange(ctx, lockedNew.End, slopeChangeNew)
		}
	}

	userEpoch := k.GetCurrentUserEpoch(ctx, veID) + 1
	k.SetCurrentUserEpoch(ctx, veID, userEpoch)
	userPointNew.Timestamp = now
	userPointNew.Block = ctx.BlockHeight()
	k.SetUserCheckpoint(ctx, veID, userEpoch, userPointNew)
}

func (k Keeper) RegulateCheckpoint(ctx sdk.Context) {
	now := uint64(ctx.BlockTime().Unix())
	epoch := k.GetCurrentEpoch(ctx)
	pointLast := k.GetCheckpoint(ctx, epoch)
	if now-pointLast.Timestamp >= merlion.SecondsPerWeek {
		k.regulateCheckpoint(ctx, sdk.ZeroInt(), sdk.ZeroInt())
	}
}

func (k Keeper) regulateCheckpoint(ctx sdk.Context, userSlopeChange, userBiasChange sdk.Int) {
	epoch := k.GetCurrentEpoch(ctx)

	now := uint64(ctx.BlockTime().Unix())
	pointLast := types.Checkpoint{
		Bias:      sdk.ZeroInt(),
		Slope:     sdk.ZeroInt(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
		Block:     ctx.BlockHeight(),
	}
	if epoch > 0 {
		pointLast = k.GetCheckpoint(ctx, epoch)
	}

	timeLast := pointLast.Timestamp

	ti := types.RegulatedUnixTime(timeLast)
	ti = types.NextRegulatedUnixTime(ti)
	slopeChange := sdk.ZeroInt()
	if ti > now {
		ti = now
	} else {
		slopeChange = k.GetSlopeChange(ctx, ti)
	}
	pointLast.Bias = pointLast.Bias.Sub(pointLast.Slope.MulRaw(int64(ti - timeLast)))
	if userBiasChange.IsPositive() {
		pointLast.Bias = pointLast.Bias.Add(userBiasChange)
	}
	if pointLast.Bias.IsNegative() {
		// can happen
		pointLast.Bias = sdk.ZeroInt()
	}
	pointLast.Slope = pointLast.Slope.Add(slopeChange)
	if userSlopeChange.IsPositive() {
		pointLast.Slope = pointLast.Slope.Add(userSlopeChange)
	}
	if pointLast.Slope.IsNegative() {
		// cannot happen, just in case
		pointLast.Slope = sdk.ZeroInt()
	}

	if ti == now {
		pointLast.Block = ctx.BlockHeight()
	}

	epoch += 1
	k.SetCurrentEpoch(ctx, epoch)
	k.SetCheckpoint(ctx, epoch, pointLast)
}

func (k Keeper) SetCurrentEpoch(ctx sdk.Context, epoch uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(epoch)
	store.Set(types.EpochKey(), bz)
}

func (k Keeper) GetCurrentEpoch(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.EpochKey())
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetCheckpoint(ctx sdk.Context, epoch uint64, point types.Checkpoint) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&point)
	store.Set(types.PointKey(epoch), bz)
}

func (k Keeper) GetCheckpoint(ctx sdk.Context, epoch uint64) types.Checkpoint {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PointKey(epoch))
	if bz == nil {
		return types.Checkpoint{
			Bias:  sdk.ZeroInt(),
			Slope: sdk.ZeroInt(),
		}
	}
	var point types.Checkpoint
	k.cdc.MustUnmarshal(bz, &point)
	return point
}

func (k Keeper) SetCurrentUserEpoch(ctx sdk.Context, veID uint64, epoch uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(epoch)
	store.Set(types.UserEpochKey(veID), bz)
}

func (k Keeper) GetCurrentUserEpoch(ctx sdk.Context, veID uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UserEpochKey(veID))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetUserCheckpoint(ctx sdk.Context, veID uint64, epoch uint64, point types.Checkpoint) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&point)
	store.Set(types.UserPointKey(veID, epoch), bz)
}

func (k Keeper) GetUserCheckpoint(ctx sdk.Context, veID uint64, epoch uint64) types.Checkpoint {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.UserPointKey(veID, epoch))
	if bz == nil {
		return types.Checkpoint{
			Bias:  sdk.ZeroInt(),
			Slope: sdk.ZeroInt(),
		}
	}
	var point types.Checkpoint
	k.cdc.MustUnmarshal(bz, &point)
	return point
}

func (k Keeper) SetSlopeChange(ctx sdk.Context, timestamp uint64, slopeChange sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{slopeChange})
	store.Set(types.SlopeChangeKey(timestamp), bz)
}

func (k Keeper) GetSlopeChange(ctx sdk.Context, timestamp uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.SlopeChangeKey(timestamp))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var slopeChange sdk.IntProto
	k.cdc.MustUnmarshal(bz, &slopeChange)
	return slopeChange.Int
}
