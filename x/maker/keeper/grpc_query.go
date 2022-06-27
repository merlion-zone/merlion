package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/maker/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) AllBackingRiskParams(c context.Context, req *types.QueryAllBackingRiskParamsRequest) (*types.QueryAllBackingRiskParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryAllBackingRiskParamsResponse{
		RiskParams: k.GetAllBackingRiskParams(ctx),
	}, nil
}

func (k Keeper) AllCollateralRiskParams(c context.Context, req *types.QueryAllCollateralRiskParamsRequest) (*types.QueryAllCollateralRiskParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryAllCollateralRiskParamsResponse{
		RiskParams: k.GetAllCollateralRiskParams(ctx),
	}, nil
}

func (k Keeper) AllBackingPools(c context.Context, req *types.QueryAllBackingPoolsRequest) (*types.QueryAllBackingPoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryAllBackingPoolsResponse{
		BackingPools: k.GetAllPoolBacking(ctx),
	}, nil
}

func (k Keeper) AllCollateralPools(c context.Context, req *types.QueryAllCollateralPoolsRequest) (*types.QueryAllCollateralPoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryAllCollateralPoolsResponse{
		CollateralPools: k.GetAllPoolCollateral(ctx),
	}, nil
}

func (k Keeper) BackingPool(c context.Context, req *types.QueryBackingPoolRequest) (*types.QueryBackingPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.GetPoolBacking(ctx, req.BackingDenom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "backing pool with backing denom '%s'", req.BackingDenom)
	}

	return &types.QueryBackingPoolResponse{
		BackingPool: pool,
	}, nil
}

func (k Keeper) CollateralPool(c context.Context, req *types.QueryCollateralPoolRequest) (*types.QueryCollateralPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	pool, found := k.GetPoolCollateral(ctx, req.CollateralDenom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "collateral pool with collateral denom '%s'", req.GetCollateralDenom())
	}

	return &types.QueryCollateralPoolResponse{
		CollateralPool: pool,
	}, nil
}

func (k Keeper) CollateralOfAccount(c context.Context, req *types.QueryCollateralOfAccountRequest) (*types.QueryCollateralOfAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	account, err := sdk.AccAddressFromBech32(req.Account)
	if err != nil {
		return nil, err
	}

	collateral, found := k.GetAccountCollateral(ctx, account, req.CollateralDenom)
	if !found {
		return nil, status.Errorf(codes.NotFound, "collateral with collateral denom '%s' of account '%s'", req.CollateralDenom, account)
	}

	return &types.QueryCollateralOfAccountResponse{
		AccountCollateral: collateral,
	}, nil
}

func (k Keeper) TotalBacking(c context.Context, req *types.QueryTotalBackingRequest) (*types.QueryTotalBackingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	total, _ := k.GetTotalBacking(ctx)

	return &types.QueryTotalBackingResponse{
		TotalBacking: total,
	}, nil
}

func (k Keeper) TotalCollateral(c context.Context, req *types.QueryTotalCollateralRequest) (*types.QueryTotalCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	total, _ := k.GetTotalCollateral(ctx)

	return &types.QueryTotalCollateralResponse{
		TotalCollateral: total,
	}, nil
}

func (k Keeper) BackingRatio(c context.Context, req *types.QueryBackingRatioRequest) (*types.QueryBackingRatioResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryBackingRatioResponse{
		BackingRatio:    k.GetBackingRatio(ctx),
		LastUpdateBlock: k.GetBackingRatioLastBlock(ctx),
	}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) EstimateMintBySwapIn(c context.Context, req *types.EstimateMintBySwapInRequest) (*types.EstimateMintBySwapInResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	backingIn, lionIn, mintFee, err := k.estimateMintBySwapIn(ctx, req.MintOut, req.BackingDenom, req.FullBacking)
	if err != nil {
		return nil, err
	}

	return &types.EstimateMintBySwapInResponse{
		BackingIn: backingIn,
		LionIn:    lionIn,
		MintFee:   mintFee,
	}, nil
}

func (k Keeper) EstimateMintBySwapOut(c context.Context, req *types.EstimateMintBySwapOutRequest) (*types.EstimateMintBySwapOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	backingIn, lionIn, mintOut, mintFee, err := k.estimateMintBySwapOut(ctx, req.BackingInMax, req.LionInMax, req.FullBacking)
	if err != nil {
		return nil, err
	}

	return &types.EstimateMintBySwapOutResponse{
		BackingIn: backingIn,
		LionIn:    lionIn,
		MintOut:   mintOut,
		MintFee:   mintFee,
	}, nil
}

func (k Keeper) EstimateBurnBySwapIn(c context.Context, req *types.EstimateBurnBySwapInRequest) (*types.EstimateBurnBySwapInResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	burnIn, backingOut, lionOut, burnFee, err := k.estimateBurnBySwapIn(ctx, req.BackingOutMax, req.LionOutMax)
	if err != nil {
		return nil, err
	}

	return &types.EstimateBurnBySwapInResponse{
		BurnIn:     burnIn,
		BackingOut: backingOut,
		LionOut:    lionOut,
		BurnFee:    burnFee,
	}, nil
}

func (k Keeper) EstimateBurnBySwapOut(c context.Context, req *types.EstimateBurnBySwapOutRequest) (*types.EstimateBurnBySwapOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	backingOut, lionOut, burnFee, err := k.estimateBurnBySwapOut(ctx, req.BurnIn, req.BackingDenom)
	if err != nil {
		return nil, err
	}

	return &types.EstimateBurnBySwapOutResponse{
		BackingOut: backingOut,
		LionOut:    lionOut,
		BurnFee:    burnFee,
	}, nil
}

func (k Keeper) EstimateBuyBackingIn(c context.Context, req *types.EstimateBuyBackingInRequest) (*types.EstimateBuyBackingInResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	lionIn, buybackFee, err := k.estimateBuyBackingIn(ctx, req.BackingOut)
	if err != nil {
		return nil, err
	}

	return &types.EstimateBuyBackingInResponse{
		LionIn:     lionIn,
		BuybackFee: buybackFee,
	}, nil
}

func (k Keeper) EstimateBuyBackingOut(c context.Context, req *types.EstimateBuyBackingOutRequest) (*types.EstimateBuyBackingOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	backingOut, buybackFee, err := k.estimateBuyBackingOut(ctx, req.LionIn, req.BackingDenom)
	if err != nil {
		return nil, err
	}

	return &types.EstimateBuyBackingOutResponse{
		BackingOut: backingOut,
		BuybackFee: buybackFee,
	}, nil
}

func (k Keeper) EstimateSellBackingIn(c context.Context, req *types.EstimateSellBackingInRequest) (*types.EstimateSellBackingInResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	backingIn, sellbackFee, err := k.estimateSellBackingIn(ctx, req.LionOut, req.BackingDenom)
	if err != nil {
		return nil, err
	}

	return &types.EstimateSellBackingInResponse{
		BackingIn:   backingIn,
		SellbackFee: sellbackFee,
	}, nil
}

func (k Keeper) EstimateSellBackingOut(c context.Context, req *types.EstimateSellBackingOutRequest) (*types.EstimateSellBackingOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	lionOut, sellbackFee, err := k.estimateSellBackingOut(ctx, req.BackingIn)
	if err != nil {
		return nil, err
	}

	return &types.EstimateSellBackingOutResponse{
		LionOut:     lionOut,
		SellbackFee: sellbackFee,
	}, nil
}

func (k Keeper) EstimateMintByCollateralIn(c context.Context, req *types.EstimateMintByCollateralInRequest) (*types.EstimateMintByCollateralInResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	account, err := sdk.AccAddressFromBech32(req.Account)
	if err != nil {
		return nil, err
	}

	collateralIn, lionIn, mintFee, _, _, _, err := k.estimateMintByCollateralIn(ctx, account, req.MintOut, req.CollateralDenom, req.Ltv)
	if err != nil {
		return nil, err
	}

	return &types.EstimateMintByCollateralInResponse{
		CollateralIn: collateralIn,
		LionIn:       lionIn,
		MintFee:      mintFee,
	}, nil
}

func (k Keeper) EstimateMintByCollateralOut(c context.Context, req *types.EstimateMintByCollateralOutRequest) (*types.EstimateMintByCollateralOutResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	account, err := sdk.AccAddressFromBech32(req.Account)
	if err != nil {
		return nil, err
	}

	lionIn, mintOut, mintFee, _, _, _, err := k.estimateMintByCollateralOut(ctx, account, req.CollateralIn, req.Ltv)
	if err != nil {
		return nil, err
	}

	return &types.EstimateMintByCollateralOutResponse{
		LionIn:  lionIn,
		MintOut: mintOut,
		MintFee: mintFee,
	}, nil
}
