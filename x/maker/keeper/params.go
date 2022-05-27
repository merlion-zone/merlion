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

// BackingRatioStep is backing ratio adjust step
func (k Keeper) BackingRatioStep(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyBackingRatioStep, &res)
	return
}

// BackingRatioPriceBand is price band within which backing ratio will not be adjusted
func (k Keeper) BackingRatioPriceBand(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyBackingRatioPriceBand, &res)
	return
}

// BackingRatioCooldownPeriod is minimum cooldown period after which backing ratio can be adjusted
func (k Keeper) BackingRatioCooldownPeriod(ctx sdk.Context) (res int64) {
	k.paramstore.Get(ctx, types.KeyBackingRatioCooldownPeriod, &res)
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
