package ve_test

import (
	"testing"

	keepertest "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/testutil/nullify"
	"github.com/merlion-zone/merlion/x/ve"
	"github.com/merlion-zone/merlion/x/ve/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.VeKeeper(t)
	ve.InitGenesis(ctx, *k, genesisState)
	got := ve.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
