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
	backingIn, lionIn, mintFee, err := k.mintBySwapRequirement(ctx, req.MintTarget, req.BackingDenom, req.FullCollateral)
	if err != nil {
		return nil, err
	}

	return &types.QueryMintBySwapRequirementResponse{
		BackingIn: backingIn,
		LionIn:    lionIn,
		MintFee:   mintFee,
	}, nil
}

func (k Keeper) mintBySwapRequirement(ctx sdk.Context, mintTarget sdk.Coin, backingDenom string, fullCollateral bool) (backingIn sdk.Coin, lionIn sdk.Coin, mintFee sdk.Coin, err error) {
	backingIn = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionIn = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	mintFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return
	}

	// prevent minting if mer price is below the lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(k.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		err = sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	mintFee = computeFee(mintTarget, backingParams.MintFee)
	mintTotal := mintTarget.Add(mintFee)

	_, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	poolBacking.MerMinted = poolBacking.MerMinted.Add(mintTotal)
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		err = sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
		return
	}

	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	collateralRatio := k.GetCollateralRatio(ctx)
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
		err = sdkerrors.Wrapf(types.ErrBackingCeiling, "backing over ceiling")
		return
	}

	return
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

	collateralRatio := k.GetCollateralRatio(ctx)
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
		mintTotalWithBackingInUSD := backingAvailInUSD.QuoRoundUp(collateralRatio)
		mintTotalWithLionInUSD := lionAvailInUSD.QuoRoundUp(sdk.OneDec().Sub(collateralRatio))
		if mintTotalWithBackingInUSD.LT(mintTotalWithLionInUSD) {
			mintTotalInUSD = mintTotalWithBackingInUSD
			backingIn.Amount = req.BackingAvail.Amount
			lionIn.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).QuoRoundUp(lionPrice).RoundInt()
		} else {
			mintTotalInUSD = mintTotalWithLionInUSD
			lionIn.Amount = req.LionAvail.Amount
			backingIn.Amount = mintTotalInUSD.Mul(collateralRatio).QuoRoundUp(backingPrice).RoundInt()
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

	mintTotalValue := mintTotalInUSD.QuoRoundUp(merPrice)
	poolBacking.MerMinted = poolBacking.MerMinted.AddAmount(mintTotalValue.RoundInt())
	if backingParams.MaxMerMint != nil && poolBacking.MerMinted.Amount.GT(*backingParams.MaxMerMint) {
		return nil, sdkerrors.Wrapf(types.ErrMerCeiling, "mer over ceiling")
	}

	mintFeeRate := sdk.ZeroDec()
	if backingParams.MintFee != nil {
		mintFeeRate = *backingParams.MintFee
	}
	mintOutValue := mintTotalValue.QuoRoundUp(sdk.OneDec().Add(mintFeeRate))
	mintFeeValue := mintOutValue.Mul(mintFeeRate)

	mintOut := sdk.NewCoin(merlion.MicroUSDDenom, mintOutValue.RoundInt())
	mintFee := sdk.NewCoin(merlion.MicroUSDDenom, mintFeeValue.RoundInt())

	return &types.QueryMintBySwapCapacityResponse{
		BackingIn: backingIn,
		LionIn:    lionIn,
		MintOut:   mintOut,
		MintFee:   mintFee,
	}, nil
}

func (k Keeper) QueryBurnBySwap(c context.Context, req *types.QueryBurnBySwapRequest) (*types.QueryBurnBySwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	backingOut, lionOut, burnFee, err := k.queryBurnBySwap(ctx, req.BurnTarget, req.BackingDenom)
	if err != nil {
		return nil, err
	}

	return &types.QueryBurnBySwapResponse{
		BackingOut: backingOut,
		LionOut:    lionOut,
		BurnFee:    burnFee,
	}, nil
}

func (k Keeper) queryBurnBySwap(ctx sdk.Context, burnTarget sdk.Coin, backingDenom string) (backingOut sdk.Coin, lionOut sdk.Coin, burnFee sdk.Coin, err error) {
	backingOut = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	burnFee = sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}
	merPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return
	}

	// check price upper bound
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(k.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		err = sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	burnFee = computeFee(burnTarget, backingParams.BurnFee)
	burnActual := burnTarget.Sub(burnFee)
	burnActualInUSD := burnActual.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	collateralRatio := k.GetCollateralRatio(ctx)
	if collateralRatio.GTE(sdk.OneDec()) {
		// full/over collateralized
		backingOut.Amount = burnActualInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if collateralRatio.IsZero() {
		// algorithmic
		lionOut.Amount = burnActualInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingOut.Amount = burnActualInUSD.Mul(collateralRatio).QuoRoundUp(backingPrice).RoundInt()
		lionOut.Amount = burnActualInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	moduleOwnedBacking := k.bankKeeper.GetBalance(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), backingDenom)
	if moduleOwnedBacking.IsLT(backingOut) {
		err = sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) < balance(%s)", backingOut, moduleOwnedBacking)
		return
	}

	return
}

func (k Keeper) QueryBuyBacking(c context.Context, req *types.QueryBuyBackingRequest) (*types.QueryBuyBackingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	backingOut, buybackFee, err := k.queryBuyBacking(ctx, req.LionIn, req.BackingDenom)
	if err != nil {
		return nil, err
	}

	return &types.QueryBuyBackingResponse{
		BackingOut: backingOut,
		BuybackFee: buybackFee,
	}, nil
}

func (k Keeper) queryBuyBacking(ctx sdk.Context, lionIn sdk.Coin, backingDenom string) (backingOut sdk.Coin, buybackFee sdk.Coin, err error) {
	backingOut = sdk.NewCoin(backingDenom, sdk.ZeroInt())
	buybackFee = sdk.NewCoin(backingDenom, sdk.ZeroInt())

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	totalBacking, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}
	if !totalBacking.MerMinted.IsPositive() {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	collateralRatio := k.GetCollateralRatio(ctx)
	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(collateralRatio).TruncateInt()

	totalBackingValue, err := k.totalBackingInUSD(ctx)
	if err != nil {
		return
	}

	availableExcessBackingValue := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if requiredBackingValue.LT(totalBackingValue.Amount) {
		availableExcessBackingValue.Amount = totalBackingValue.Amount.Sub(requiredBackingValue)
	}

	lionInValue := lionIn.Amount.ToDec().Mul(lionPrice)
	if lionInValue.TruncateInt().GT(availableExcessBackingValue.Amount) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	backingOut = sdk.NewCoin(backingDenom, lionInValue.Quo(backingPrice).TruncateInt())
	if poolBacking.Backing.IsLT(backingOut) {
		err = sdkerrors.Wrap(types.ErrBackingCoinInsufficient, "insufficient available backing coin")
		return
	}

	buybackFee = computeFee(backingOut, backingParams.BuybackFee)
	backingOut = backingOut.Sub(buybackFee)

	return
}

func (k Keeper) QuerySellBacking(c context.Context, req *types.QuerySellBackingRequest) (*types.QuerySellBackingResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	lionOut, sellbackFee, err := k.querySellBacking(ctx, req.BackingIn)
	if err != nil {
		return nil, err
	}

	return &types.QuerySellBackingResponse{
		LionOut:     lionOut,
		SellbackFee: sellbackFee,
	}, nil
}

func (k Keeper) querySellBacking(ctx sdk.Context, backingIn sdk.Coin) (lionOut sdk.Coin, sellbackFee sdk.Coin, err error) {
	lionOut = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())
	sellbackFee = sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt())

	backingDenom := backingIn.Denom

	// get prices in usd
	backingPrice, err := k.oracleKeeper.GetExchangeRate(ctx, backingDenom)
	if err != nil {
		return
	}
	lionPrice, err := k.oracleKeeper.GetExchangeRate(ctx, merlion.AttoLionDenom)
	if err != nil {
		return
	}

	backingParams, found := k.GetBackingRiskParams(ctx, backingDenom)
	if !found {
		err = sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", backingDenom)
		return
	}
	if !backingParams.Enabled {
		err = sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", backingDenom)
		return
	}

	totalBacking, poolBacking, err := k.getBacking(ctx, backingDenom)
	if err != nil {
		return
	}

	poolBacking.Backing = poolBacking.Backing.Add(backingIn)
	if backingParams.MaxBacking != nil && poolBacking.Backing.Amount.GT(*backingParams.MaxBacking) {
		err = sdkerrors.Wrap(types.ErrBackingCeiling, "over ceiling")
		return
	}

	collateralRatio := k.GetCollateralRatio(ctx)
	requiredBackingValue := totalBacking.MerMinted.Amount.ToDec().Mul(collateralRatio).TruncateInt()

	totalBackingValue, err := k.totalBackingInUSD(ctx)
	if err != nil {
		return
	}

	availableMissingBackingValue := sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt())
	if requiredBackingValue.GT(totalBackingValue.Amount) {
		availableMissingBackingValue.Amount = requiredBackingValue.Sub(totalBackingValue.Amount)
	}
	availableLionOut := availableMissingBackingValue.Amount.ToDec().Quo(lionPrice)

	bonusRatio := k.RecollateralizeBonus(ctx)
	lionMint := sdk.NewCoin(merlion.AttoLionDenom, backingIn.Amount.ToDec().Mul(backingPrice).Quo(lionPrice).TruncateInt())
	bonus := computeFee(lionMint, &bonusRatio)
	sellbackFee = computeFee(lionMint, backingParams.RecollateralizeFee)

	lionMint = lionMint.Add(bonus)
	if lionMint.Amount.ToDec().GT(availableLionOut) {
		err = sdkerrors.Wrap(types.ErrLionCoinInsufficient, "insufficient available lion coin")
		return
	}

	lionOut = lionMint.Sub(sellbackFee)
	return
}
