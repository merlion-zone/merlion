package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/types"
)

type Keeper struct {
	stakingkeeper.Keeper
	storeKey   sdk.StoreKey
	cdc        codec.BinaryCodec
	bankKeeper stakingtypes.BankKeeper
	nftKeeper  types.NftKeeper
	veKeeper   types.VeKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	ps paramtypes.Subspace,
	ak stakingtypes.AccountKeeper,
	bk stakingtypes.BankKeeper,
	nk types.NftKeeper,
	vk types.VeKeeper,
) Keeper {
	keeper := Keeper{
		Keeper:     stakingkeeper.NewKeeper(cdc, key, ak, bk, ps),
		storeKey:   key,
		cdc:        cdc,
		bankKeeper: bk,
		nftKeeper:  nk,
		veKeeper:   vk,
	}

	keeper.veKeeper.SetGetDelegatedAmountByUser(func(ctx sdk.Context, veID uint64) sdk.Int {
		return keeper.GetVeDelegatedAmount(ctx, veID)
	})

	return keeper
}

func (k *Keeper) SetHooks(sh stakingtypes.StakingHooks) *Keeper {
	k.Keeper = *k.Keeper.SetHooks(sh)
	return k
}

func (k *Keeper) CheckDenom(ctx sdk.Context) {
	if k.veKeeper.LockDenom(ctx) != k.BondDenom(ctx) {
		panic("bond denom is different from ve lock denom")
	}
}
