package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

// RegulateUserCheckpoint regulates user checkpoint.
// veID:       must be valid ve id
// lockedOld:
//             Amount: can be zero
//             End: can be expired or zero
// lockedNew:
//             Amount: can be zero
//             End: must be in the future or be zero
func (k Keeper) RegulateUserCheckpoint(ctx sdk.Context, veID uint64, lockedOld types.LockedBalance, lockedNew types.LockedBalance) {
	// check whether timestamp is regulated
	types.CheckRegulatedUnixTime(lockedOld.End)
	types.CheckRegulatedUnixTime(lockedNew.End)

	// use block time as now timestamp
	now := uint64(ctx.BlockTime().Unix())

	// user point initialized with zero values
	userPointOld := types.Checkpoint{
		Bias:  sdk.ZeroInt(),
		Slope: sdk.ZeroInt(),
	}
	userPointNew := types.Checkpoint{
		Bias:  sdk.ZeroInt(),
		Slope: sdk.ZeroInt(),
	}

	// calculate slope and bias from now on,
	// kept at zero after the unlocking time
	if lockedOld.End > now && lockedOld.Amount.IsPositive() {
		userPointOld.Slope = lockedOld.Amount.QuoRaw(types.MaxLockTime)
		userPointOld.Bias = userPointOld.Slope.MulRaw(int64(lockedOld.End - now))
	}
	if lockedNew.End > now && lockedNew.Amount.IsPositive() {
		// slope is always proportional to locked amount
		userPointNew.Slope = lockedNew.Amount.QuoRaw(types.MaxLockTime)
		// bias represents the voting power at the present:
		//     slope * (end - now)
		// so that at any future time t, voting power will decay to:
		//     bias - slope * (t - now)
		userPointNew.Bias = userPointNew.Slope.MulRaw(int64(lockedNew.End - now))
	}

	// regulate system checkpoint history
	userSlopeChange := userPointNew.Slope.Sub(userPointOld.Slope)
	userBiasChange := userPointNew.Bias.Sub(userPointOld.Bias)
	k.regulateCheckpoint(ctx, userSlopeChange, userBiasChange)

	slopeChangeOld := k.GetSlopeChange(ctx, lockedOld.End)
	slopeChangeNew := k.GetSlopeChange(ctx, lockedNew.End)

	// Schedule future slope changes at any unlocking time of any user's ve,
	// since the slope should be reset to zero after unlocking time is up.
	// Actually equivalent to, for the second segment of the piecewise linear function,
	// the slope should be zero.
	if lockedOld.End > now {
		// old slope waw subtracted, so here add it back to cancel it
		slopeChangeOld = slopeChangeOld.Add(userPointOld.Slope)
		if lockedNew.End == lockedOld.End {
			// subtract new slope
			slopeChangeOld = slopeChangeOld.Sub(userPointNew.Slope)
		}
		k.SetSlopeChange(ctx, lockedOld.End, slopeChangeOld)
	}
	if lockedNew.End > now {
		if lockedNew.End > lockedOld.End {
			// subtract new slope, i.e., the slope is reset to zero
			slopeChangeNew = slopeChangeNew.Sub(userPointNew.Slope)
			k.SetSlopeChange(ctx, lockedNew.End, slopeChangeNew)
		} else {
			// has been handled in slopeChangeOld
		}
	}

	// increase user epoch
	userEpoch := k.GetUserEpoch(ctx, veID) + 1
	k.SetUserEpoch(ctx, veID, userEpoch)

	// set new user checkpoint
	userPointNew.Timestamp = now
	userPointNew.Block = ctx.BlockHeight()
	k.SetUserCheckpoint(ctx, veID, userEpoch, userPointNew)
}

func (k Keeper) RegulateCheckpoint(ctx sdk.Context) {
	now := uint64(ctx.BlockTime().Unix())
	epoch := k.GetEpoch(ctx)
	pointLast := k.GetCheckpoint(ctx, epoch)
	if now-pointLast.Timestamp >= types.RegulatedPeriod {
		k.regulateCheckpoint(ctx, sdk.ZeroInt(), sdk.ZeroInt())
	}
}

func (k Keeper) regulateCheckpoint(ctx sdk.Context, userSlopeChange, userBiasChange sdk.Int) {
	now := uint64(ctx.BlockTime().Unix())

	epoch := k.GetEpoch(ctx)

	pointLast := types.Checkpoint{
		Bias:      sdk.ZeroInt(),
		Slope:     sdk.ZeroInt(),
		Timestamp: now,
		Block:     ctx.BlockHeight(),
	}
	if epoch > 0 {
		pointLast = k.GetCheckpoint(ctx, epoch)
	}

	timeLast := pointLast.Timestamp

	pointLastInitial := pointLast
	blockSlope := sdk.ZeroDec()
	if timeLast < now {
		// block increasing slope, proportional to lapse of time
		// here use Dec to maintain sufficient precision
		blockSlope = sdk.NewDecFromInt(sdk.NewInt(ctx.BlockHeight() - pointLast.Block)).QuoInt64(int64(now - timeLast))
	}

	ti := types.RegulatedUnixTime(timeLast)
	i := 0
	for {
		i++
		if i > 2 {
			// since checkpoint will always be regulated when past the regulated period at end block,
			// this loop will not exceed two rounds
			panic("broken checkpoint regulation")
		}

		ti = types.NextRegulatedUnixTime(ti)

		var slopeChange sdk.Int
		if ti > now {
			// since ti is regulated, now must not be at regulated time
			// so set slope change as zero
			slopeChange = sdk.ZeroInt()
			// up to the current time
			ti = now
		} else {
			// ti is at regulated time
			slopeChange = k.GetSlopeChange(ctx, ti)
		}

		// calculate new bias and slope
		pointLast.Bias = pointLast.Bias.Sub(pointLast.Slope.MulRaw(int64(ti - timeLast)))
		if pointLast.Bias.IsNegative() {
			// can happen
			pointLast.Bias = sdk.ZeroInt()
		}
		pointLast.Slope = pointLast.Slope.Add(slopeChange)
		if pointLast.Slope.IsNegative() {
			// cannot happen, just in case
			pointLast.Slope = sdk.ZeroInt()
		}

		timeLast = ti
		pointLast.Timestamp = ti
		// calculate block approximately
		pointLast.Block = pointLastInitial.Block + blockSlope.MulInt64(int64(ti-pointLastInitial.Timestamp)).TruncateInt64()

		// increase epoch
		epoch += 1

		if ti == now {
			pointLast.Block = ctx.BlockHeight()
			break // break loop
		} else {
			// set new checkpoint
			k.SetCheckpoint(ctx, epoch, pointLast)
		}
	}

	// TODO: delete slope changes in the past, since they will be no longer used

	// set new last epoch
	k.SetEpoch(ctx, epoch)

	// add new change at now to the new last point
	if userBiasChange.IsPositive() {
		pointLast.Bias = pointLast.Bias.Add(userBiasChange)
		if pointLast.Bias.IsNegative() {
			pointLast.Bias = sdk.ZeroInt()
		}
	}
	if userSlopeChange.IsPositive() {
		pointLast.Slope = pointLast.Slope.Add(userSlopeChange)
		if pointLast.Slope.IsNegative() {
			pointLast.Slope = sdk.ZeroInt()
		}
	}

	// set new checkpoint
	k.SetCheckpoint(ctx, epoch, pointLast)
}

func (k Keeper) SetEpoch(ctx sdk.Context, epoch uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(epoch)
	store.Set(types.EpochKey(), bz)
}

func (k Keeper) GetEpoch(ctx sdk.Context) uint64 {
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

func (k Keeper) SetUserEpoch(ctx sdk.Context, veID uint64, epoch uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(epoch)
	store.Set(types.UserEpochKey(veID), bz)
}

func (k Keeper) GetUserEpoch(ctx sdk.Context, veID uint64) uint64 {
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
