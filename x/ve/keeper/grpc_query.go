package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) TotalVotingPower(ctx context.Context, request *types.QueryTotalVotingPowerRequest) (*types.QueryTotalVotingPowerRequest, error) {
	// TODO implement me
	panic("implement me")
}

func (k Keeper) VotingPower(ctx context.Context, request *types.QueryVotingPowerRequest) (*types.QueryVotingPowerRequest, error) {
	// TODO implement me
	panic("implement me")
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
