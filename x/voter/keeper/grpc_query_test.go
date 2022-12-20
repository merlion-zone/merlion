package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testkeeper "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/x/voter/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.VoterKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	response, err := keeper.Params(wctx, nil)
	require.Nil(t, response)
	require.Equal(t, status.Error(codes.InvalidArgument, "invalid request"), err)

	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err = keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
