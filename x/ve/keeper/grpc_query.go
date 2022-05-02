package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/merlion-zone/merlion/x/ve/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) TotalVotingPower(c context.Context, msg *types.QueryTotalVotingPowerRequest) (*types.QueryTotalVotingPowerResponse, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	if msg.AtTime > 0 {
		if msg.AtBlock > 0 {
			return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be specified")
		}
	} else if msg.AtBlock > 0 {
		if msg.AtBlock > ctx.BlockHeight() {
			return nil, status.Error(codes.InvalidArgument, "invalid block")
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be zero")
	}

	power := k.GetTotalVotingPower(ctx, msg.AtTime, msg.AtBlock)

	return &types.QueryTotalVotingPowerResponse{
		Power: power,
	}, nil
}

func (k Keeper) GetTotalVotingPower(ctx sdk.Context, atTime uint64, atBlock int64) sdk.Int {
	epoch := k.GetEpoch(ctx)
	var power sdk.Int

	if atTime > 0 {
		pointLast := k.GetCheckpoint(ctx, epoch)
		if atTime <= pointLast.Timestamp {
			// in the past
			targetEpoch := k.findTimeEpoch(ctx, types.EmptyVeID, atTime, epoch)
			point := k.GetCheckpoint(ctx, targetEpoch)

			var dt int64
			if targetEpoch < epoch {
				pointNext := k.GetCheckpoint(ctx, targetEpoch+1)
				if point.Block != pointNext.Block {
					dt = int64(atTime - point.Timestamp)
				}
			} else if point.Block != ctx.BlockHeight() {
				dt = int64(atTime - point.Timestamp)
			}

			power = point.Bias.Sub(point.Slope.MulRaw(dt))

		} else {
			// in the future
			ti := types.RegulatedUnixTime(pointLast.Timestamp)
			for {
				ti = types.NextRegulatedUnixTime(ti)
				var slopeChange sdk.Int
				if ti > atTime {
					ti = atTime
					slopeChange = sdk.ZeroInt()
				} else {
					slopeChange = k.GetSlopeChange(ctx, ti)
				}
				pointLast.Bias = pointLast.Bias.Sub(pointLast.Slope.MulRaw(int64(ti - pointLast.Timestamp)))
				if ti == atTime {
					break
				}
				pointLast.Slope = pointLast.Slope.Add(slopeChange)
				pointLast.Timestamp = ti
			}

			power = pointLast.Bias
		}

	} else if atBlock > 0 {
		targetEpoch := k.findBlockEpoch(ctx, types.EmptyVeID, atBlock, epoch)
		point := k.GetCheckpoint(ctx, targetEpoch)

		var dt int64
		if targetEpoch < epoch {
			pointNext := k.GetCheckpoint(ctx, targetEpoch+1)
			if point.Block != pointNext.Block {
				dt = (atBlock - point.Block) * int64(pointNext.Timestamp-point.Timestamp) / (pointNext.Block - point.Block)
			}
		} else if point.Block != ctx.BlockHeight() {
			dt = (atBlock - point.Block) * (ctx.BlockTime().Unix() - int64(point.Timestamp)) / (ctx.BlockHeight() - point.Block)
		}

		power = point.Bias.Sub(point.Slope.MulRaw(dt))

	}

	if power.IsNegative() {
		power = sdk.ZeroInt()
	}

	return power
}

func (k Keeper) VotingPower(c context.Context, msg *types.QueryVotingPowerRequest) (*types.QueryVotingPowerResponse, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	if !k.nftKeeper.HasNFT(ctx, types.VeNftClass.Id, msg.VeId) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidVeID, "invalid ve id: %s", msg.VeId)
	}

	veID := types.Uint64FromVeID(msg.VeId)

	if msg.AtTime > 0 {
		if msg.AtBlock > 0 {
			return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be specified")
		}
	} else if msg.AtBlock > 0 {
		// always in the past
		if msg.AtBlock > ctx.BlockHeight() {
			return nil, status.Error(codes.InvalidArgument, "invalid block")
		}
	} else {
		return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be zero")
	}

	power := k.GetVotingPower(ctx, veID, msg.AtTime, msg.AtBlock)

	return &types.QueryVotingPowerResponse{
		Power: power,
	}, nil
}

func (k Keeper) GetVotingPower(ctx sdk.Context, veID uint64, atTime uint64, atBlock int64) sdk.Int {
	epoch := k.GetEpoch(ctx)
	userEpoch := k.GetUserEpoch(ctx, veID)
	var power sdk.Int

	if atTime > 0 {
		userPoint := k.GetUserCheckpoint(ctx, veID, userEpoch)
		if atTime <= userPoint.Timestamp {
			// in the past
			targetUserEpoch := k.findTimeEpoch(ctx, veID, atTime, userEpoch)
			userPoint = k.GetUserCheckpoint(ctx, veID, targetUserEpoch)
		} else {
			// in the future
		}
		power = userPoint.Bias.Sub(userPoint.Slope.MulRaw(int64(atTime - userPoint.Timestamp)))

	} else if atBlock > 0 {
		// find timestamp through system checkpoint history
		var blockTimestamp uint64
		{
			targetEpoch := k.findBlockEpoch(ctx, types.EmptyVeID, atBlock, epoch)
			point := k.GetCheckpoint(ctx, targetEpoch)

			var dt int64
			if targetEpoch < epoch {
				pointNext := k.GetCheckpoint(ctx, targetEpoch+1)
				if point.Block != pointNext.Block {
					dt = (atBlock - point.Block) * int64(pointNext.Timestamp-point.Timestamp) / (pointNext.Block - point.Block)
				}
			} else if point.Block != ctx.BlockHeight() {
				dt = (atBlock - point.Block) * (ctx.BlockTime().Unix() - int64(point.Timestamp)) / (ctx.BlockHeight() - point.Block)
			}

			blockTimestamp = point.Timestamp + uint64(dt)
		}

		targetUserEpoch := k.findBlockEpoch(ctx, veID, atBlock, userEpoch)
		userPoint := k.GetUserCheckpoint(ctx, veID, targetUserEpoch)

		power = userPoint.Bias.Sub(userPoint.Slope.MulRaw(int64(blockTimestamp - userPoint.Timestamp)))

	}

	if power.IsNegative() {
		power = sdk.ZeroInt()
	}

	return power
}

func (k Keeper) Params(c context.Context, msg *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

// findTimeEpoch finds approximate epoch for specified timestamp
func (k Keeper) findTimeEpoch(ctx sdk.Context, veID uint64, timestamp uint64, maxEpoch uint64) uint64 {
	min := uint64(0)
	max := maxEpoch
	// binary search
	for {
		if min >= max {
			break
		}
		mid := (min + max + 1) / 2

		var midTimestamp uint64
		if veID != types.EmptyVeID {
			midTimestamp = k.GetUserCheckpoint(ctx, veID, mid).Timestamp
		} else {
			midTimestamp = k.GetCheckpoint(ctx, mid).Timestamp
		}

		if midTimestamp <= timestamp {
			min = mid
		} else {
			max = mid - 1
		}
	}
	return min
}

// findBlockEpoch finds approximate epoch for specified block
func (k Keeper) findBlockEpoch(ctx sdk.Context, veID uint64, block int64, maxEpoch uint64) uint64 {
	min := uint64(0)
	max := maxEpoch
	// binary search
	for {
		if min >= max {
			break
		}
		mid := (min + max + 1) / 2

		var midBlock int64
		if veID != types.EmptyVeID {
			midBlock = k.GetUserCheckpoint(ctx, veID, mid).Block
		} else {
			midBlock = k.GetCheckpoint(ctx, mid).Block
		}

		if midBlock <= block {
			min = mid
		} else {
			max = mid - 1
		}
	}
	return min
}
