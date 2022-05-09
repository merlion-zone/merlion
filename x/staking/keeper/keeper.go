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
	storeKey  sdk.StoreKey
	cdc       codec.BinaryCodec
	nftKeeper types.NftKeeper
	veKeeper  types.VeKeeper
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key sdk.StoreKey,
	ak stakingtypes.AccountKeeper,
	bk stakingtypes.BankKeeper,
	ps paramtypes.Subspace,
) Keeper {
	return Keeper{
		Keeper:   stakingkeeper.NewKeeper(cdc, key, ak, bk, ps),
		storeKey: key,
		cdc:      cdc,
	}
}
