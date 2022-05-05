package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
	nftkeeper "github.com/cosmos/cosmos-sdk/x/nft/keeper"
	nft "github.com/cosmos/cosmos-sdk/x/nft/module"
	"github.com/merlion-zone/merlion/x/ve/types"
)

type NftAppModule struct {
	nft.AppModule
	keeper NftKeeper
}

func NewNftAppModule(module nft.AppModule, keeper NftKeeper) NftAppModule {
	return NftAppModule{AppModule: module, keeper: keeper}
}

func (am NftAppModule) RegisterServices(cfg module.Configurator) {
	nfttypes.RegisterMsgServer(cfg.MsgServer(), am.keeper)
	nfttypes.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

type NftKeeper struct {
	nftkeeper.Keeper
	veKeeper func() Keeper
}

func NewNftKeeper(keeper nftkeeper.Keeper, veKeeper func() Keeper) NftKeeper {
	return NftKeeper{Keeper: keeper, veKeeper: veKeeper}
}

func (k NftKeeper) Send(c context.Context, msg *nfttypes.MsgSend) (*nfttypes.MsgSendResponse, error) {
	if msg.ClassId == types.VeNftClass.Id {
		ctx := sdk.UnwrapSDKContext(c)
		veID := types.Uint64FromVeID(msg.Id)
		err := k.veKeeper().CheckVeAttached(ctx, veID)
		if err != nil {
			return nil, err
		}
	}
	return k.Keeper.Send(c, msg)
}

func (k Keeper) CheckVeAttached(ctx sdk.Context, veID uint64) error {
	if k.GetVeAttached(ctx, veID) != 0 || k.GetVeVoted(ctx, veID) {
		return types.ErrVeAttached
	}
	return nil
}

func (k Keeper) SaveNftClass(ctx sdk.Context) error {
	return k.nftKeeper.SaveClass(ctx, types.VeNftClass)
}

func (k Keeper) SetNextVeID(ctx sdk.Context, nextVeID uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := sdk.Uint64ToBigEndian(nextVeID)
	store.Set(types.NextVeIDKey(), bz)
}

func (k Keeper) GetNextVeID(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.NextVeIDKey())
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) IncVeAttached(ctx sdk.Context, veID uint64) {
	attached := k.GetVeAttached(ctx, veID)
	k.SetVeAttached(ctx, veID, attached+1)
}

func (k Keeper) DecVeAttached(ctx sdk.Context, veID uint64) {
	attached := k.GetVeAttached(ctx, veID)
	if attached == 0 {
		panic("invalid ve attached number")
	}
	k.SetVeAttached(ctx, veID, attached-1)
}

func (k Keeper) SetVeAttached(ctx sdk.Context, veID uint64, attached uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.AttachedKey(veID), sdk.Uint64ToBigEndian(attached))
}

func (k Keeper) GetVeAttached(ctx sdk.Context, veID uint64) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AttachedKey(veID))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

func (k Keeper) SetVeVoted(ctx sdk.Context, veID uint64, voted bool) {
	store := ctx.KVStore(k.storeKey)
	bz := []byte{0}
	if voted {
		bz[0] = 1
	}
	store.Set(types.VotedKey(veID), bz)
}

func (k Keeper) GetVeVoted(ctx sdk.Context, veID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.VotedKey(veID))
	if bz == nil || bz[0] == 0 {
		return false
	} else {
		return true
	}
}
