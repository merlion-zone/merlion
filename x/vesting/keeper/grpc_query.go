package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/merlion-zone/merlion/x/vesting/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Airdrops(c context.Context, msg *types.QueryAirdropsRequest) (*types.QueryAirdropsResponse, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	keyPrefix := types.KeyPrefixAirdrops
	if msg.Completed {
		keyPrefix = types.KeyPrefixAirdropsCompleted
	}

	var airdrops []types.Airdrop
	store := ctx.KVStore(k.storeKey)
	valStore := prefix.NewStore(store, keyPrefix)
	pageRes, err := query.FilteredPaginate(valStore, msg.Pagination, func(key []byte, value []byte, accumulate bool) (bool, error) {
		var airdrop types.Airdrop
		k.cdc.MustUnmarshal(value, &airdrop)

		if accumulate {
			airdrops = append(airdrops, airdrop)
		}

		return true, nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAirdropsResponse{
		Airdrops:   airdrops,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) Airdrop(c context.Context, msg *types.QueryAirdropRequest) (*types.QueryAirdropResponse, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	targetAddr, err := sdk.AccAddressFromBech32(msg.TargetAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var airdrop types.Airdrop
	if !msg.Completed {
		airdrop = k.GetAirdrop(ctx, targetAddr)
	} else {
		airdrop = k.GetAirdropCompleted(ctx, targetAddr)
	}
	if airdrop.Empty() {
		return nil, status.Error(codes.NotFound, "airdrop target not found")
	}

	return &types.QueryAirdropResponse{
		Airdrop: airdrop,
	}, nil
}

func (k Keeper) Params(c context.Context, msg *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
