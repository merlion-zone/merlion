package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (k Keeper) SetTotalBacking(ctx sdk.Context, pool types.TotalBacking) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(types.KeyPrefixBackingTotal, bz)
}

func (k Keeper) GetTotalBacking(ctx sdk.Context) (types.TotalBacking, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixBackingTotal)
	var total types.TotalBacking
	if len(bz) == 0 {
		return total, false
	}
	k.cdc.MustUnmarshal(bz, &total)
	return total, true
}

func (k Keeper) SetPoolBacking(ctx sdk.Context, pool types.PoolBacking) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBackingPool)
	bz := k.cdc.MustMarshal(&pool)
	store.Set([]byte(pool.Backing.Denom), bz)
}

func (k Keeper) GetPoolBacking(ctx sdk.Context, denom string) (types.PoolBacking, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixBackingPool)
	bz := store.Get([]byte(denom))
	var pool types.PoolBacking
	if len(bz) == 0 {
		return pool, false
	}
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

func (k Keeper) GetAllPoolBacking(ctx sdk.Context) []types.PoolBacking {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixBackingPool)
	defer iterator.Close()

	var allParams []types.PoolBacking
	for ; iterator.Valid(); iterator.Next() {
		var params types.PoolBacking
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		allParams = append(allParams, params)
	}

	return allParams
}

func (k Keeper) SetTotalCollateral(ctx sdk.Context, pool types.TotalCollateral) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(types.KeyPrefixCollateralTotal, bz)
}

func (k Keeper) GetTotalCollateral(ctx sdk.Context) (types.TotalCollateral, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefixCollateralTotal)
	var total types.TotalCollateral
	if len(bz) == 0 {
		return total, false
	}
	k.cdc.MustUnmarshal(bz, &total)
	return total, true
}

func (k Keeper) SetPoolCollateral(ctx sdk.Context, pool types.PoolCollateral) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCollateralPool)
	bz := k.cdc.MustMarshal(&pool)
	store.Set([]byte(pool.Collateral.Denom), bz)
}

func (k Keeper) GetPoolCollateral(ctx sdk.Context, denom string) (types.PoolCollateral, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCollateralPool)
	bz := store.Get([]byte(denom))
	var pool types.PoolCollateral
	if len(bz) == 0 {
		return pool, false
	}
	k.cdc.MustUnmarshal(bz, &pool)
	return pool, true
}

func (k Keeper) GetAllPoolCollateral(ctx sdk.Context) []types.PoolCollateral {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.KeyPrefixCollateralPool)
	defer iterator.Close()

	var allParams []types.PoolCollateral
	for ; iterator.Valid(); iterator.Next() {
		var params types.PoolCollateral
		k.cdc.MustUnmarshal(iterator.Value(), &params)

		allParams = append(allParams, params)
	}

	return allParams
}

func (k Keeper) SetAccountCollateral(ctx sdk.Context, addr sdk.AccAddress, col types.AccountCollateral) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCollateralAccount)
	bz := k.cdc.MustMarshal(&col)
	store.Set(keyByAddrDenom(types.KeyPrefixCollateralAccount, addr, col.Collateral.Denom), bz)
}

func (k Keeper) GetAccountCollateral(ctx sdk.Context, addr sdk.AccAddress, denom string) (types.AccountCollateral, bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefixCollateralPool)
	bz := store.Get(keyByAddrDenom(types.KeyPrefixCollateralAccount, addr, denom))
	var collateral types.AccountCollateral
	if len(bz) == 0 {
		return collateral, false
	}
	k.cdc.MustUnmarshal(bz, &collateral)
	return collateral, true
}

func keyByAddrDenom(prefix []byte, addr sdk.AccAddress, denom string) (key []byte) {
	key = append(prefix, address.MustLengthPrefix(addr)...)
	return append(key, []byte(denom)...)
}
