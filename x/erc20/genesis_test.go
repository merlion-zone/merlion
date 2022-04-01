package erc20_test

import (
	"testing"

	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	keepertest "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/testutil/nullify"
	"github.com/merlion-zone/merlion/x/erc20"
	"github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.Erc20Keeper(t)
	ak := authkeeper.AccountKeeper{} // TODO
	erc20.InitGenesis(ctx, *k, ak, genesisState)
	got := erc20.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
