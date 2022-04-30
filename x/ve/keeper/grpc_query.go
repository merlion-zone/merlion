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

	epoch := k.GetEpoch(ctx)
	var power sdk.Int

	if msg.AtTime > 0 {
		if msg.AtBlock > 0 {
			return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be specified")
		}

		pointLast := k.GetCheckpoint(ctx, epoch)
		if msg.AtTime <= pointLast.Timestamp {
			// in the past
			targetEpoch := k.findTimeEpoch(ctx, types.EmptyVeID, msg.AtTime, epoch)
			point := k.GetCheckpoint(ctx, targetEpoch)

			var dt int64
			if targetEpoch < epoch {
				pointNext := k.GetCheckpoint(ctx, targetEpoch+1)
				if point.Block != pointNext.Block {
					dt = int64(msg.AtTime - point.Timestamp)
				}
			} else if point.Block != ctx.BlockHeight() {
				dt = int64(msg.AtTime - point.Timestamp)
			}

			power = point.Bias.Sub(point.Slope.MulRaw(dt))

		} else {
			// in the future
			ti := types.RegulatedUnixTime(pointLast.Timestamp)
			for {
				ti = types.NextRegulatedUnixTime(ti)
				var slopeChange sdk.Int
				if ti > msg.AtTime {
					ti = msg.AtTime
					slopeChange = sdk.ZeroInt()
				} else {
					slopeChange = k.GetSlopeChange(ctx, ti)
				}
				pointLast.Bias = pointLast.Bias.Sub(pointLast.Slope.MulRaw(int64(ti - pointLast.Timestamp)))
				if ti == msg.AtTime {
					break
				}
				pointLast.Slope = pointLast.Slope.Add(slopeChange)
				pointLast.Timestamp = ti
			}

			power = pointLast.Bias
		}

	} else if msg.AtBlock > 0 {
		if msg.AtBlock > ctx.BlockHeight() {
			return nil, status.Error(codes.InvalidArgument, "invalid block")
		}

		targetEpoch := k.findBlockEpoch(ctx, types.EmptyVeID, msg.AtBlock, epoch)
		point := k.GetCheckpoint(ctx, targetEpoch)

		var dt int64
		if targetEpoch < epoch {
			pointNext := k.GetCheckpoint(ctx, targetEpoch+1)
			if point.Block != pointNext.Block {
				dt = (msg.AtBlock - point.Block) * int64(pointNext.Timestamp-point.Timestamp) / (pointNext.Block - point.Block)
			}
		} else if point.Block != ctx.BlockHeight() {
			dt = (msg.AtBlock - point.Block) * (ctx.BlockTime().Unix() - int64(point.Timestamp)) / (ctx.BlockHeight() - point.Block)
		}

		power = point.Bias.Sub(point.Slope.MulRaw(dt))

	} else {
		return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be zero")
	}

	if power.IsNegative() {
		power = sdk.ZeroInt()
	}

	return &types.QueryTotalVotingPowerResponse{
		Power: power,
	}, nil
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

	epoch := k.GetEpoch(ctx)
	userEpoch := k.GetUserEpoch(ctx, veID)
	var power sdk.Int

	if msg.AtTime > 0 {
		if msg.AtBlock > 0 {
			return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be specified")
		}

		userPoint := k.GetUserCheckpoint(ctx, veID, userEpoch)
		if msg.AtTime <= userPoint.Timestamp {
			// in the past
			targetUserEpoch := k.findTimeEpoch(ctx, veID, msg.AtTime, userEpoch)
			userPoint = k.GetUserCheckpoint(ctx, veID, targetUserEpoch)
		} else {
			// in the future
		}
		power = userPoint.Bias.Sub(userPoint.Slope.MulRaw(int64(msg.AtTime - userPoint.Timestamp)))

	} else if msg.AtBlock > 0 {
		// always in the past
		if msg.AtBlock > ctx.BlockHeight() {
			return nil, status.Error(codes.InvalidArgument, "invalid block")
		}

		// find timestamp through system checkpoint history
		var blockTimestamp uint64
		{
			targetEpoch := k.findBlockEpoch(ctx, types.EmptyVeID, msg.AtBlock, epoch)
			point := k.GetCheckpoint(ctx, targetEpoch)

			var dt int64
			if targetEpoch < epoch {
				pointNext := k.GetCheckpoint(ctx, targetEpoch+1)
				if point.Block != pointNext.Block {
					dt = (msg.AtBlock - point.Block) * int64(pointNext.Timestamp-point.Timestamp) / (pointNext.Block - point.Block)
				}
			} else if point.Block != ctx.BlockHeight() {
				dt = (msg.AtBlock - point.Block) * (ctx.BlockTime().Unix() - int64(point.Timestamp)) / (ctx.BlockHeight() - point.Block)
			}

			blockTimestamp = point.Timestamp + uint64(dt)
		}

		targetUserEpoch := k.findBlockEpoch(ctx, veID, msg.AtBlock, userEpoch)
		userPoint := k.GetUserCheckpoint(ctx, veID, targetUserEpoch)

		power = userPoint.Bias.Sub(userPoint.Slope.MulRaw(int64(blockTimestamp - userPoint.Timestamp)))

	} else {
		return nil, status.Error(codes.InvalidArgument, "at time and at block cannot both be zero")
	}

	if power.IsNegative() {
		power = sdk.ZeroInt()
	}

	return &types.QueryVotingPowerResponse{
		Power: power,
	}, nil
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
