package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/vesting/types"
)

func (k Keeper) AllocateAtGenesis(ctx sdk.Context) {
}

// SetAllocationAddresses sets allocation target addresses
func (k Keeper) SetAllocationAddresses(ctx sdk.Context, addresses types.AllocationAddresses) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&addresses)
	store.Set(types.AllocationAddrKey(), bz)
}

// GetAllocationAddresses gets allocation target addresses
func (k Keeper) GetAllocationAddresses(ctx sdk.Context) types.AllocationAddresses {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AllocationAddrKey())
	if bz == nil {
		return types.AllocationAddresses{}
	}
	var addresses types.AllocationAddresses
	k.cdc.MustUnmarshal(bz, &addresses)
	return addresses
}

// SetAirdropTotalAmount sets airdrop total amount
func (k Keeper) SetAirdropTotalAmount(ctx sdk.Context, total sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{total})
	store.Set(types.AirdropsTotalAmountKey(), bz)
}

// GetAirdropTotalAmount gets airdrop total amount
func (k Keeper) GetAirdropTotalAmount(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropsTotalAmountKey())
	if bz == nil {
		return sdk.ZeroInt()
	}
	var total sdk.IntProto
	k.cdc.MustUnmarshal(bz, &total)
	return total.Int
}

// SetAirdrop sets airdrop target
func (k Keeper) SetAirdrop(ctx sdk.Context, acc sdk.AccAddress, airdrop types.Airdrop) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&airdrop)
	store.Set(types.AirdropsKey(acc), bz)
}

// GetAirdrop gets airdrop target
func (k Keeper) GetAirdrop(ctx sdk.Context, acc sdk.AccAddress) types.Airdrop {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropsKey(acc))
	if bz == nil {
		return types.Airdrop{}
	}
	var airdrop types.Airdrop
	k.cdc.MustUnmarshal(bz, &airdrop)
	return airdrop
}

// DeleteAirdrop deletes airdrop target
func (k Keeper) DeleteAirdrop(ctx sdk.Context, acc sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AirdropsKey(acc))
}

// IterateAirdrops iterates airdrop targets
func (k Keeper) IterateAirdrops(ctx sdk.Context, handler func(airdrop types.Airdrop) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixAirdrops)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var airdrop types.Airdrop
		k.cdc.MustUnmarshal(iter.Value(), &airdrop)
		if handler(airdrop) {
			break
		}
	}
}

// SetAirdropCompleted sets completed airdrop target
func (k Keeper) SetAirdropCompleted(ctx sdk.Context, acc sdk.AccAddress, airdrop types.Airdrop) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&airdrop)
	store.Set(types.AirdropsCompletedKey(acc), bz)
}

// GetAirdropCompleted gets completed airdrop target
func (k Keeper) GetAirdropCompleted(ctx sdk.Context, acc sdk.AccAddress) types.Airdrop {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AirdropsCompletedKey(acc))
	if bz == nil {
		return types.Airdrop{}
	}
	var airdrop types.Airdrop
	k.cdc.MustUnmarshal(bz, &airdrop)
	return airdrop
}

// IterateAirdropsCompleted iterates completed airdrop targets
func (k Keeper) IterateAirdropsCompleted(ctx sdk.Context, handler func(airdrop types.Airdrop) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.KeyPrefixAirdropsCompleted)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var airdrop types.Airdrop
		k.cdc.MustUnmarshal(iter.Value(), &airdrop)
		if handler(airdrop) {
			break
		}
	}
}
