package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (k Keeper) SetBackingRiskParams(ctx sdk.Context, params types.BackingRiskParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBackingParams)
	bz := k.cdc.MustMarshal(&params)
	store.Set([]byte(params.BackingDenom), bz)
}

func (k Keeper) SetCollateralRiskParams(ctx sdk.Context, params types.CollateralRiskParams) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCollateralParams)
	bz := k.cdc.MustMarshal(&params)
	store.Set([]byte(params.CollateralDenom), bz)
}

func (k Keeper) IsBackingRegistered(ctx sdk.Context, denom string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBackingParams)
	return store.Has([]byte(denom))
}

func (k Keeper) GetBackingRiskParams(ctx sdk.Context, denom string) (types.BackingRiskParams, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBackingParams)
	bz := store.Get([]byte(denom))
	var params types.BackingRiskParams
	if len(bz) == 0 {
		return params, false
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params, true
}

func (k Keeper) GetAllBackingRiskParams(ctx sdk.Context) []types.BackingRiskParams {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixBackingParams)
	defer iterator.Close()

	var allParams []types.BackingRiskParams
	for ; iterator.Valid(); iterator.Next() {
		var params types.BackingRiskParams
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		allParams = append(allParams, params)
	}

	return allParams
}

func (k Keeper) IsCollateralRegistered(ctx sdk.Context, denom string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCollateralParams)
	return store.Has([]byte(denom))
}

func (k Keeper) GetCollateralRiskParams(ctx sdk.Context, denom string) (types.CollateralRiskParams, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCollateralParams)
	bz := store.Get([]byte(denom))
	var params types.CollateralRiskParams
	if len(bz) == 0 {
		return params, false
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params, true
}

func (k Keeper) GetAllCollateralRiskParams(ctx sdk.Context) []types.CollateralRiskParams {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixCollateralParams)
	defer iterator.Close()

	var allParams []types.CollateralRiskParams
	for ; iterator.Valid(); iterator.Next() {
		var params types.CollateralRiskParams
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		allParams = append(allParams, params)
	}

	return allParams
}
