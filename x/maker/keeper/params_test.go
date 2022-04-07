package keeper_test

import (
	"testing"

	testkeeper "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/x/maker/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.MakerKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
