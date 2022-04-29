package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

func (k Keeper) SaveNftClass(ctx sdk.Context) error {
	return k.nftKeeper.SaveClass(ctx, types.VeNftClass)
}
