package ve

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/keeper"
	"github.com/merlion-zone/merlion/x/ve/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetParams(ctx, genState.Params)

	// we only allow upgrade by in-place store migrations
	// so check only once at block 1
	if !k.HasNftClass(ctx) {
		if err := k.SaveNftClass(ctx); err != nil {
			panic(err)
		}

		k.SetTotalLockedAmount(ctx, sdk.ZeroInt())
		k.SetNextVeID(ctx, types.FirstVeID)
		k.SetEpoch(ctx, types.EmptyEpoch)
		k.SetCheckpoint(ctx, types.EmptyEpoch, types.Checkpoint{
			Bias:      sdk.ZeroInt(),
			Slope:     sdk.ZeroInt(),
			Timestamp: 0,
			Block:     0,
		})
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	return genesis
}
