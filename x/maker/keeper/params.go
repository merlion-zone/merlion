package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// CollateralRatioStep is collateral ratio adjust step
func (k Keeper) CollateralRatioStep(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyCollateralRatioStep, &res)
	return
}

// CollateralRatioPriceBand is price band within which collateral ratio will not be adjusted
func (k Keeper) CollateralRatioPriceBand(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyCollateralRatioPriceBand, &res)
	return
}

// CollateralRatioCooldownPeriod is minimum cooldown period after which collateral ratio can be adjusted
func (k Keeper) CollateralRatioCooldownPeriod(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyCollateralRatioCooldownPeriod, &res)
	return
}

// MintPriceBias is mint price bias ratio
func (k Keeper) MintPriceBias(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMintPriceBias, &res)
	return
}

// BurnPriceBias is burn price bias ratio
func (k Keeper) BurnPriceBias(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMintPriceBias, &res)
	return
}

// RecollateralizeBonus is recollateralization bonus ratio
func (k Keeper) RecollateralizeBonus(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyRecollateralizeBonus, &res)
	return
}

// LiquidationCommissionFee is liquidation commission fee ratio
func (k Keeper) LiquidationCommissionFee(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyLiquidationCommissionFee, &res)
	return
}
