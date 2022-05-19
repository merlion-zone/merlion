package vesting

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/vesting/keeper"
	"github.com/merlion-zone/merlion/x/vesting/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if merlion.IsPeriodLastBlock(ctx, types.ClaimVestedPeriod) {
		k.ClaimVested(ctx)
	}
}
