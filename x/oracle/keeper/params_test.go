package keeper_test

import (
	"testing"

	testkeeper "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.OracleKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
