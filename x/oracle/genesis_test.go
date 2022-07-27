package oracle_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/oracle"
	"github.com/merlion-zone/merlion/x/oracle/keeper"
	"github.com/merlion-zone/merlion/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestExportInitGenesis(t *testing.T) {
	input, _ := setup(t)

	input.OracleKeeper.SetFeederDelegation(input.Ctx, keeper.ValAddrs[0], keeper.Addrs[1])
	input.OracleKeeper.SetExchangeRate(input.Ctx, "denom", sdk.NewDec(123))
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, keeper.ValAddrs[0], types.NewAggregateExchangeRatePrevote(types.AggregateVoteHash{123}, keeper.ValAddrs[0], uint64(2)))
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, keeper.ValAddrs[0], types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{{Denom: "foo", ExchangeRate: sdk.NewDec(123)}}, keeper.ValAddrs[0]))
	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], 10)
	genesis := oracle.ExportGenesis(input.Ctx, input.OracleKeeper)

	newInput := keeper.CreateTestInput(t)
	oracle.InitGenesis(newInput.Ctx, newInput.OracleKeeper, *genesis)
	newGenesis := oracle.ExportGenesis(newInput.Ctx, newInput.OracleKeeper)

	require.Equal(t, genesis, newGenesis)
}

func TestInitGenesis(t *testing.T) {
	input, _ := setup(t)
	genesis := types.DefaultGenesis()
	require.NotPanics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, *genesis)
	})

	genesis.FeederDelegations = []types.FeederDelegation{{
		FeederAddress:    keeper.Addrs[0].String(),
		ValidatorAddress: "invalid",
	}}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, *genesis)
	})

	genesis.FeederDelegations = []types.FeederDelegation{{
		FeederAddress:    "invalid",
		ValidatorAddress: keeper.ValAddrs[0].String(),
	}}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, *genesis)
	})

	genesis.FeederDelegations = []types.FeederDelegation{{
		FeederAddress:    keeper.Addrs[0].String(),
		ValidatorAddress: keeper.ValAddrs[0].String(),
	}}

	genesis.MissCounters = []types.MissCounter{
		{
			ValidatorAddress: "invalid",
			MissCounter:      10,
		},
	}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, *genesis)
	})

	genesis.MissCounters = []types.MissCounter{
		{
			ValidatorAddress: keeper.ValAddrs[0].String(),
			MissCounter:      10,
		},
	}

	genesis.AggregateExchangeRatePrevotes = []types.AggregateExchangeRatePrevote{
		{
			Hash:        "hash",
			Voter:       "invalid",
			SubmitBlock: 100,
		},
	}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, *genesis)
	})

	genesis.AggregateExchangeRatePrevotes = []types.AggregateExchangeRatePrevote{
		{
			Hash:        "hash",
			Voter:       keeper.ValAddrs[0].String(),
			SubmitBlock: 100,
		},
	}

	genesis.AggregateExchangeRateVotes = []types.AggregateExchangeRateVote{
		{
			ExchangeRateTuples: []types.ExchangeRateTuple{
				{
					Denom:        "foo",
					ExchangeRate: sdk.NewDec(10),
				},
			},
			Voter: "invalid",
		},
	}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, *genesis)
	})

	genesis.AggregateExchangeRateVotes = []types.AggregateExchangeRateVote{
		{
			ExchangeRateTuples: []types.ExchangeRateTuple{
				{
					Denom:        "foo",
					ExchangeRate: sdk.NewDec(10),
				},
			},
			Voter: keeper.ValAddrs[0].String(),
		},
	}

	require.NotPanics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, *genesis)
	})
}
