package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/x/maker/keeper"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.MakerKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
