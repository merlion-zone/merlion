package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
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

func (k Keeper) CollateralRatio(c context.Context, req *types.QueryCollateralRatioRequest) (*types.QueryCollateralRatioResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryCollateralRatioResponse{
		CollateralRatio: k.GetCollateralRatio(ctx),
		LastUpdateBlock: k.GetCollateralRatioLastBlock(ctx),
	}, nil
}

func (k Keeper) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) MintBySwapRequirement(c context.Context, req *types.QueryMintBySwapRequirementRequest) (*types.QueryMintBySwapRequirementResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	return k.getMintBySwapRequirement(ctx, req.MintTarget, req.BackingDenom, req.FullCollateral)
}

func (k Keeper) getMintBySwapRequirement(ctx sdk.Context, mintTarget sdk.Coin, backingDenom string, fullCollateral bool) (*types.QueryMintBySwapRequirementResponse, error) {
	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return nil, err
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// prevent minting if mer price is below the lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(k.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralRatio := k.GetCollateralRatio(ctx)

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
	}
	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
	}

	mintOut := mintTarget
	mintFee := computeFee(mintOut, backingParams.MintFee)
	mintTotal := mintOut.Add(mintFee)

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	poolBacking.MerMinted = poolBacking.MerMinted.Add(mintTotal)
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		return nil, sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
	}

	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(merlion.MicroUSDTarget)
	backingIn := sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionIn := sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	if collateralRatio.GTE(sdk.OneDec()) || fullCollateral {
		// full/over collateralized, or user selects full collateralization
		backingIn.Amount = mintTotalInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if collateralRatio.IsZero() {
		// algorithmic
		lionIn.Amount = mintTotalInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingIn.Amount = mintTotalInUSD.Mul(collateralRatio).QuoRoundUp(backingPrice).RoundInt()
		lionIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
	}

	return &types.QueryMintBySwapRequirementResponse{
		BackingIn: backingIn,
		LionIn:    lionIn,
		MintFee:   mintFee,
	}, nil
}

func (k Keeper) MintBySwapCapacity(c context.Context, req *types.QueryMintBySwapCapacityRequest) (*types.QueryMintBySwapCapacityResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	backingDenom := req.BackingAvail.Denom

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return nil, err
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// prevent minting if mer price is below the lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(k.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralRatio := k.GetCollateralRatio(ctx)

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
	}
	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
	}

	backingAvailInUSD := backingPrice.MulInt(req.BackingAvail.Amount)
	lionAvailInUSD := lionPrice.MulInt(req.LionAvail.Amount)

	mintTotalInUSD := sdk.ZeroDec()
	backingIn := sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionIn := sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	if collateralRatio.GTE(sdk.OneDec()) || req.LionAvail.IsZero() {
		// full/over collateralized, or user selects full collateralization
		mintTotalInUSD = backingAvailInUSD
		backingIn.Amount = req.BackingAvail.Amount
	} else if collateralRatio.IsZero() {
		// algorithmic
		mintTotalInUSD = lionAvailInUSD
		lionIn.Amount = req.LionAvail.Amount
	} else {
		// fractional
		mintTotalWithBackingInUSD := backingAvailInUSD.Quo(collateralRatio)
		mintTotalWithLionInUSD := lionAvailInUSD.Quo(sdk.OneDec().Sub(collateralRatio))
		if mintTotalWithBackingInUSD.LT(mintTotalWithLionInUSD) {
			mintTotalInUSD = mintTotalWithBackingInUSD
			backingIn.Amount = req.BackingAvail.Amount
			lionIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).Quo(lionPrice).RoundInt()
		} else {
			mintTotalInUSD = mintTotalWithLionInUSD
			lionIn.Amount = req.LionAvail.Amount
			backingIn.Amount = mintTotalInUSD.Mul(collateralRatio).Quo(backingPrice).RoundInt()
		}
	}

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return nil, err
	}
	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
	}

	mintTotalAmount := mintTotalInUSD.Quo(merPrice)
	poolBacking.MerMinted = poolBacking.MerMinted.AddAmount(mintTotalAmount.RoundInt())
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		return nil, sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
	}

	mintFeeRate := sdk.ZeroDec()
	if backingParams.MintFee != nil {
		mintFeeRate = *backingParams.MintFee
	}
	mintOutAmount := mintTotalAmount.Quo(sdk.OneDec().Add(mintFeeRate))
	mintFeeAmount := mintOutAmount.Mul(mintFeeRate)

	mintOut := sdk.NewCoin(merlion.MicroUSDDenom, mintOutAmount.RoundInt())
	mintFee := sdk.NewCoin(merlion.MicroUSDDenom, mintFeeAmount.RoundInt())

	return &types.QueryMintBySwapCapacityResponse{
		BackingIn: backingIn,
		LionIn:    lionIn,
		MintOut:   mintOut,
		MintFee:   mintFee,
	}, nil
}
