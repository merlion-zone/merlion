package ve

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/ve/keeper"
	"github.com/merlion-zone/merlion/x/ve/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	if err := k.SaveNftClass(ctx); err != nil {
		panic(err)
	}

	k.SetTotalLockedAmount(ctx, merlion.ZeroInt)
	k.SetNextVeID(ctx, 1)
	k.SetEpoch(ctx, 0)
	k.SetCheckpoint(ctx, 0, types.Checkpoint{
		Bias:      merlion.ZeroInt,
		Slope:     merlion.ZeroInt,
		Timestamp: 0,
		Block:     0,
	})
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
