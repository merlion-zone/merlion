package keeper

import (
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/merlion-zone/merlion/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestOrganizeAggregate(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	fooBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(17), "foo", ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(10), "foo", ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(6), "foo", ValAddrs[2], power),
	}
	barBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(1000), "bar", ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(1300), "bar", ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(2000), "bar", ValAddrs[2], power),
	}

	for i := range fooBallot {
		input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[i],
			types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
				{Denom: fooBallot[i].Denom, ExchangeRate: fooBallot[i].ExchangeRate},
				{Denom: barBallot[i].Denom, ExchangeRate: barBallot[i].ExchangeRate},
			}, ValAddrs[i]))
	}

	// organize votes by denom
	ballotMap := input.OracleKeeper.OrganizeBallotByDenom(input.Ctx, map[string]types.Claim{
		ValAddrs[0].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[0],
		},
		ValAddrs[1].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[1],
		},
		ValAddrs[2].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[2],
		},
	})

	// sort each ballot for comparison
	sort.Sort(fooBallot)
	sort.Sort(barBallot)
	sort.Sort(ballotMap["foo"])
	sort.Sort(ballotMap["bar"])

	require.Equal(t, fooBallot, ballotMap["foo"])
	require.Equal(t, barBallot, ballotMap["bar"])
}

func TestClearBallots(t *testing.T) {
	input := CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	fooBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(17), "foo", ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(10), "foo", ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(6), "foo", ValAddrs[2], power),
	}
	barBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(sdk.NewDec(1000), "bar", ValAddrs[0], power),
		types.NewVoteForTally(sdk.NewDec(1300), "bar", ValAddrs[1], power),
		types.NewVoteForTally(sdk.NewDec(2000), "bar", ValAddrs[2], power),
	}

	for i := range fooBallot {
		input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, ValAddrs[i], types.AggregateExchangeRatePrevote{
			Hash:        "",
			Voter:       ValAddrs[i].String(),
			SubmitBlock: uint64(input.Ctx.BlockHeight()),
		})

		input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, ValAddrs[i],
			types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
				{Denom: fooBallot[i].Denom, ExchangeRate: fooBallot[i].ExchangeRate},
				{Denom: barBallot[i].Denom, ExchangeRate: barBallot[i].ExchangeRate},
			}, ValAddrs[i]))
	}

	input.OracleKeeper.ClearBallots(input.Ctx, 5)

	prevoteCounter := 0
	voteCounter := 0
	input.OracleKeeper.IterateAggregateExchangeRatePrevotes(input.Ctx, func(_ sdk.ValAddress, _ types.AggregateExchangeRatePrevote) bool {
		prevoteCounter++
		return false
	})
	input.OracleKeeper.IterateAggregateExchangeRateVotes(input.Ctx, func(_ sdk.ValAddress, _ types.AggregateExchangeRateVote) bool {
		voteCounter++
		return false
	})

	require.Equal(t, prevoteCounter, 3)
	require.Equal(t, voteCounter, 0)

	input.OracleKeeper.ClearBallots(input.Ctx.WithBlockHeight(input.Ctx.BlockHeight()+6), 5)

	prevoteCounter = 0
	input.OracleKeeper.IterateAggregateExchangeRatePrevotes(input.Ctx, func(_ sdk.ValAddress, _ types.AggregateExchangeRatePrevote) bool {
		prevoteCounter++
		return false
	})
	require.Equal(t, prevoteCounter, 0)
}
