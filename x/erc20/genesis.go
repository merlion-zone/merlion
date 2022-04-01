package erc20

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/erc20/keeper"
	"github.com/merlion-zone/merlion/x/erc20/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	accountKeeper types.AccountKeeper,
	genState types.GenesisState,
) {
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)

	// ensure erc20 module account is set on genesis
	if acc := accountKeeper.GetModuleAccount(ctx, types.ModuleName); acc == nil {
		// NOTE: shouldn't occur
		panic("the erc20 module account has not been set")
	}

	for _, pair := range genState.TokenPairs {
		id := pair.GetID()
		k.SetTokenPair(ctx, pair)
		k.SetDenomMap(ctx, pair.Denom, id)
		k.SetERC20Map(ctx, pair.GetERC20Contract(), id)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)
	genesis.TokenPairs = k.GetAllTokenPairs(ctx)

	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
