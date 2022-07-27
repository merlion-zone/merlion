package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/ethereum/go-ethereum/common"
	testkeeper "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestKeeper_Params(t *testing.T) {
	keeper, ctx := testkeeper.Erc20Keeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}

func TestKeeper_TokenPairs(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		pairs  = []types.TokenPair{
			types.NewTokenPair(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), "USDT", 0),
			types.NewTokenPair(common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), "USDC", 0)}
		req = &types.QueryTokenPairsRequest{
			Pagination: &query.PageRequest{Key: nil, Limit: 2, CountTotal: true},
		}
	)
	for _, p := range pairs {
		k.SetTokenPair(ctx, p)
	}
	_, err := k.TokenPairs(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err, status.Error(codes.InvalidArgument, "empty request"))

	res, err := k.TokenPairs(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.Equal(t, uint64(len(pairs)), res.Pagination.Total)
	require.Equal(t, pairs, res.TokenPairs)
}

func TestKeeper_TokenPair(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		pairs  = []types.TokenPair{
			types.NewTokenPair(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), "USDT", 0),
			types.NewTokenPair(common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), "USDC", 0)}
		req = &types.QueryTokenPairRequest{
			Token: pairs[0].Erc20Address,
		}
	)

	_, err := k.TokenPair(sdk.WrapSDKContext(ctx), nil)
	require.Error(t, err, status.Error(codes.InvalidArgument, "empty request"))

	_, err = k.TokenPair(sdk.WrapSDKContext(ctx), &types.QueryTokenPairRequest{Token: "bad"})
	require.Error(t, err, status.Errorf(
		codes.InvalidArgument,
		"invalid format for token %s, should be either hex ('0x...') address or cosmos denom", "bad",
	))

	_, err = k.TokenPair(sdk.WrapSDKContext(ctx), &types.QueryTokenPairRequest{Token: "DAI"})
	require.Error(t, err, status.Errorf(codes.NotFound, "token pair with token '%s'", "DAI"))

	for _, p := range pairs {
		k.SetTokenPair(ctx, p)
	}

	_, err = k.TokenPair(sdk.WrapSDKContext(ctx), req)
	require.Error(t, err, status.Errorf(codes.NotFound, "token pair with token '%s'", req.Token))

	for _, p := range pairs {
		k.SetERC20Map(ctx, common.HexToAddress(p.Erc20Address), p.GetID())
	}

	res, err := k.TokenPair(sdk.WrapSDKContext(ctx), req)
	require.NoError(t, err)
	require.Equal(t, res.TokenPair, pairs[0])
}
