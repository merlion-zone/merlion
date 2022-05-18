package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/merlion-zone/merlion/x/vesting/types"
)

type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   sdk.StoreKey
	paramstore paramtypes.Subspace

	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
}

// NewKeeper creates a new vesting Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ps paramtypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		paramstore:    ps,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
	}
}
