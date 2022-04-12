package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/merlion-zone/merlion/x/maker/types"
)

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

func keyByAddrDenom(prefix []byte, addr sdk.AccAddress, denom string) (key []byte) {
	key = append(prefix, address.MustLengthPrefix(addr)...)
	return append(key, []byte(denom)...)
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
