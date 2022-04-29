package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

func (k Keeper) SaveNftClass(ctx sdk.Context) error {
	return k.nftKeeper.SaveClass(ctx, types.VeNftClass)
}

func (k Keeper) SetNextVeID(ctx sdk.Context, nextVeID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(nextVeID)
	store.Set(types.NextVeNftIDKey(), bz)
}

func (k Keeper) GetNextVeID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextVeNftIDKey())
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}
