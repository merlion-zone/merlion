package oracle_test

import (
	"fmt"
	"math"
	"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/oracle"
	"github.com/merlion-zone/merlion/x/oracle/keeper"
	"github.com/merlion-zone/merlion/x/oracle/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/rand"
)

func TestOracleThreshold(t *testing.T) {
	input, h := setup(t)

	// Case 1.
	// Less than the threshold signs, exchange rate consensus fails
	salt := "1"
	hash := types.GetAggregateVoteHash(salt, denom1ExchangeRateStr, keeper.ValAddrs[0])
	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, denom1ExchangeRateStr, keeper.Addrs[0], keeper.ValAddrs[0])

	_, err1 := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 := h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx.WithBlockHeight(1), denom1)
	require.Error(t, err)

	// Case 2.
	// More than the threshold signs, exchange rate consensus succeeds
	salt = "1"
	hash = types.GetAggregateVoteHash(salt, denom1ExchangeRateStr, keeper.ValAddrs[0])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, denom1ExchangeRateStr, keeper.Addrs[0], keeper.ValAddrs[0])

	_, err1 = h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	salt = "2"
	hash = types.GetAggregateVoteHash(salt, denom1ExchangeRateStr, keeper.ValAddrs[1])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, denom1ExchangeRateStr, keeper.Addrs[1], keeper.ValAddrs[1])

	_, err1 = h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	salt = "3"
	hash = types.GetAggregateVoteHash(salt, denom1ExchangeRateStr, keeper.ValAddrs[2])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[2], keeper.ValAddrs[2])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, denom1ExchangeRateStr, keeper.Addrs[2], keeper.ValAddrs[2])

	_, err1 = h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rate, err := input.OracleKeeper.GetExchangeRate(input.Ctx.WithBlockHeight(1), denom1)
	require.NoError(t, err)
	require.Equal(t, randomExchangeRate, rate)

	// Case 3.
	// Increase voting power of absent validator, exchange rate consensus fails
	val, _ := input.StakingKeeper.GetValidator(input.Ctx, keeper.ValAddrs[2])
	input.StakingKeeper.Delegate(input.Ctx.WithBlockHeight(0), keeper.Addrs[2], stakingAmt.MulRaw(3), stakingtypes.Unbonded, val, false)

	salt = "1"
	hash = types.GetAggregateVoteHash(salt, denom1ExchangeRateStr, keeper.ValAddrs[0])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[0], keeper.ValAddrs[0])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, denom1ExchangeRateStr, keeper.Addrs[0], keeper.ValAddrs[0])

	_, err1 = h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	salt = "2"
	hash = types.GetAggregateVoteHash(salt, denom1ExchangeRateStr, keeper.ValAddrs[1])
	prevoteMsg = types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[1], keeper.ValAddrs[1])
	voteMsg = types.NewMsgAggregateExchangeRateVote(salt, denom1ExchangeRateStr, keeper.Addrs[1], keeper.ValAddrs[1])

	_, err1 = h(input.Ctx.WithBlockHeight(0), prevoteMsg)
	_, err2 = h(input.Ctx.WithBlockHeight(1), voteMsg)
	require.NoError(t, err1)
	require.NoError(t, err2)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	_, err = input.OracleKeeper.GetExchangeRate(input.Ctx.WithBlockHeight(1), denom1)
	require.Error(t, err)
}

func TestOracleDrop(t *testing.T) {
	input, h := setup(t)

	input.OracleKeeper.SetExchangeRate(input.Ctx, denom2, randomExchangeRate)

	// Account 1, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 0)

	// Immediately swap halt after an illiquid oracle vote
	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx, denom2)
	require.Error(t, err)
}

func TestOracleTally(t *testing.T) {
	input, _ := setup(t)

	ballot := types.ExchangeRateBallot{}
	rates, valAddrs, stakingKeeper := types.GenerateRandomTestCase()
	input.OracleKeeper.SetStakingKeeper(stakingKeeper)
	h := oracle.NewHandler(input.OracleKeeper)
	for i, rate := range rates {

		decExchangeRate := sdk.NewDecWithPrec(int64(rate*math.Pow10(keeper.OracleDecPrecision)), int64(keeper.OracleDecPrecision))

		salt := fmt.Sprintf("%d", i)
		hash := types.GetAggregateVoteHash(salt, denom1ExchangeRateStr, valAddrs[i])
		prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, sdk.AccAddress(valAddrs[i]), valAddrs[i])
		voteMsg := types.NewMsgAggregateExchangeRateVote(salt, denom1ExchangeRateStr, sdk.AccAddress(valAddrs[i]), valAddrs[i])

		_, err1 := h(input.Ctx.WithBlockHeight(0), prevoteMsg)
		_, err2 := h(input.Ctx.WithBlockHeight(1), voteMsg)
		require.NoError(t, err1)
		require.NoError(t, err2)

		power := stakingAmt.QuoRaw(microUnit).Int64()
		if decExchangeRate.IsZero() {
			power = int64(0)
		}

		vote := types.NewVoteForTally(
			decExchangeRate, denom1, valAddrs[i], power)
		ballot = append(ballot, vote)

		// change power of every three validator
		if i%3 == 0 {
			stakingKeeper.Validators()[i].SetConsensusPower(int64(i + 1))
		}
	}

	validatorClaimMap := make(map[string]types.Claim)
	for _, valAddr := range valAddrs {
		validatorClaimMap[valAddr.String()] = types.Claim{
			Power:     stakingKeeper.Validator(input.Ctx, valAddr).GetConsensusPower(sdk.DefaultPowerReduction),
			Weight:    int64(0),
			WinCount:  int64(0),
			Recipient: valAddr,
		}
	}
	sort.Sort(ballot)
	weightedMedian := ballot.WeightedMedianWithAssertion()
	standardDeviation := ballot.StandardDeviation(weightedMedian)
	maxSpread := weightedMedian.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	if standardDeviation.GT(maxSpread) {
		maxSpread = standardDeviation
	}

	expectedValidatorClaimMap := make(map[string]types.Claim)
	for _, valAddr := range valAddrs {
		expectedValidatorClaimMap[valAddr.String()] = types.Claim{
			Power:     stakingKeeper.Validator(input.Ctx, valAddr).GetConsensusPower(sdk.DefaultPowerReduction),
			Weight:    int64(0),
			WinCount:  int64(0),
			Recipient: valAddr,
		}
	}

	for _, vote := range ballot {
		if (vote.ExchangeRate.GTE(weightedMedian.Sub(maxSpread)) &&
			vote.ExchangeRate.LTE(weightedMedian.Add(maxSpread))) ||
			!vote.ExchangeRate.IsPositive() {
			key := vote.Voter.String()
			claim := expectedValidatorClaimMap[key]
			claim.Weight += vote.Power
			claim.WinCount++
			expectedValidatorClaimMap[key] = claim
		}
	}

	tallyMedian := oracle.Tally(input.Ctx, ballot, input.OracleKeeper.RewardBand(input.Ctx), validatorClaimMap)

	require.Equal(t, validatorClaimMap, expectedValidatorClaimMap)
	require.Equal(t, tallyMedian.MulInt64(100).TruncateInt(), weightedMedian.MulInt64(100).TruncateInt())
}

func TestOracleTallyTiming(t *testing.T) {
	input, h := setup(t)

	// all the keeper.Addrs vote for the block ... not last period block yet, so tally fails
	for i := range keeper.Addrs[:2] {
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom1, Amount: randomExchangeRate}}, i)
	}

	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 10 // set vote period to 10 for now, for convenience
	input.OracleKeeper.SetParams(input.Ctx, params)
	require.Equal(t, 0, int(input.Ctx.BlockHeight()))

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx, denom1)
	require.Error(t, err)

	input.Ctx = input.Ctx.WithBlockHeight(int64(params.VotePeriod - 1))

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	_, err = input.OracleKeeper.GetExchangeRate(input.Ctx, denom1)
	require.NoError(t, err)
}

func TestOracleRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// Account 1, denom1
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom1, Amount: randomExchangeRate}}, 0)

	// Account 2, denom1
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom1, Amount: randomExchangeRate}}, 1)

	rewardsAmt := sdk.NewInt(100000000)
	err := input.BankKeeper.MintCoins(input.Ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(merlion.AttoLionDenom, rewardsAmt)))
	require.NoError(t, err)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	votePeriodsPerWindow := uint64(sdk.NewDec(int64(input.OracleKeeper.RewardDistributionWindow(input.Ctx))).QuoInt64(int64(input.OracleKeeper.VotePeriod(input.Ctx))).TruncateInt64())
	expectedRewardAmt := sdk.NewDecFromInt(rewardsAmt.QuoRaw(2)).QuoInt64(int64(votePeriodsPerWindow)).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
}

func TestOracleRewardBand(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	input.OracleKeeper.SetParams(input.Ctx, params)

	// reset vote targets
	input.OracleKeeper.ClearVoteTargets(input.Ctx)
	input.OracleKeeper.SetVoteTarget(input.Ctx, denom2)

	rewardSpread := randomExchangeRate.Mul(input.OracleKeeper.RewardBand(input.Ctx).QuoInt64(2))

	// no one will miss the vote
	// Account 1, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate.Sub(rewardSpread)}}, 0)

	// Account 2, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 1)

	// Account 3, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate.Add(rewardSpread)}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// Account 1 will miss the vote due to raward band condition
	// Account 1, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate.Sub(rewardSpread.Add(sdk.OneDec()))}}, 0)

	// Account 2, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 1)

	// Account 3, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate.Add(rewardSpread)}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

}

func TestOracleMultiRewardDistribution(t *testing.T) {
	input, h := setup(t)

	// denom1 and denom2 have the same voting power, but denom2 has been chosen as referenceMer by alphabetical order.
	// Account 1, denom1, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom1, Amount: randomExchangeRate}, {Denom: denom2, Amount: randomExchangeRate}}, 0)

	// Account 2, denom1
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom1, Amount: randomExchangeRate}}, 1)

	// Account 3, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 2)

	rewardAmt := sdk.NewInt(100000000)
	err := input.BankKeeper.MintCoins(input.Ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(merlion.AttoLionDenom, rewardAmt)))
	require.NoError(t, err)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)

	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(3).MulRaw(2)).QuoInt64(int64(rewardDistributedWindow)).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(3)).QuoInt64(int64(rewardDistributedWindow)).TruncateInt()
	expectedRewardAmt3 := sdk.ZeroInt() // even vote power is same denom2 with denom1, denom1 chosen referenceMer because alphabetical order

	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt3, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
}

func TestOracleExchangeRate(t *testing.T) {
	input, h := setup(t)

	denom2RandomExchangeRate := sdk.NewDecWithPrec(1000000000, int64(6)).MulInt64(microUnit)
	usmRandomExchangeRate := sdk.NewDecWithPrec(1000000, int64(6)).MulInt64(microUnit)

	// denom2 has been chosen as referenceMer by highest voting power
	// Account 1, USM, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: merlion.MicroUSMDenom, Amount: usmRandomExchangeRate}, {Denom: denom2, Amount: denom2RandomExchangeRate}}, 0)

	// Account 2, USM, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: merlion.MicroUSMDenom, Amount: usmRandomExchangeRate}, {Denom: denom2, Amount: denom2RandomExchangeRate}}, 1)

	// Account 3, denom2, denom1
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: denom2RandomExchangeRate}, {Denom: denom1, Amount: randomExchangeRate}}, 2)

	rewardAmt := sdk.NewInt(100000000)
	err := input.BankKeeper.MintCoins(input.Ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(merlion.AttoLionDenom, rewardAmt)))
	require.NoError(t, err)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(5).MulRaw(2)).QuoInt64(int64(rewardDistributedWindow)).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(5).MulRaw(1)).QuoInt64(int64(rewardDistributedWindow)).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards = input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
}

func TestOracleEnsureSorted(t *testing.T) {
	input, h := setup(t)

	for i := 0; i < 100; i++ {
		denom2ExchangeRate1 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(microUnit)
		usmExchangeRate1 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(microUnit)

		denom2ExchangeRate2 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(microUnit)
		usmExchangeRate2 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(microUnit)

		denom2ExchangeRate3 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(microUnit)
		usmExchangeRate3 := sdk.NewDecWithPrec(int64(rand.Uint64()%100000000), 6).MulInt64(microUnit)

		// Account 1, USM, Denom2
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: merlion.MicroUSMDenom, Amount: usmExchangeRate1}, {Denom: denom2, Amount: denom2ExchangeRate1}}, 0)

		// Account 2, USM, Denom2
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: merlion.MicroUSMDenom, Amount: usmExchangeRate2}, {Denom: denom2, Amount: denom2ExchangeRate2}}, 1)

		// Account 3, USM, Denom2
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: merlion.MicroUSMDenom, Amount: denom2ExchangeRate3}, {Denom: denom2, Amount: usmExchangeRate3}}, 2)

		require.NotPanics(t, func() {
			oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)
		})
	}
}

func TestOracleExchangeRateVal5(t *testing.T) {
	input, h := setupVal5(t)

	denom2ExchangeRate := sdk.NewDecWithPrec(505000, int64(6)).MulInt64(microUnit)
	denom2ExchangeRateWithErr := sdk.NewDecWithPrec(500000, int64(6)).MulInt64(microUnit)
	usmExchangeRate := sdk.NewDecWithPrec(505, int64(6)).MulInt64(microUnit)
	usmExchangeRateWithErr := sdk.NewDecWithPrec(500, int64(6)).MulInt64(microUnit)

	// denom2 has been chosen as referenceMer by highest voting power
	// Account 1, denom2, USM
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: denom2ExchangeRate}, {Denom: merlion.MicroUSMDenom, Amount: usmExchangeRate}}, 0)

	// Account 2, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: denom2ExchangeRate}}, 1)

	// Account 3, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: denom2ExchangeRate}}, 2)

	// Account 4, denom2, USM
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: denom2ExchangeRateWithErr}, {Denom: merlion.MicroUSMDenom, Amount: usmExchangeRateWithErr}}, 3)

	// Account 5, denom2, USM
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: denom2ExchangeRateWithErr}, {Denom: merlion.MicroUSMDenom, Amount: usmExchangeRateWithErr}}, 4)

	rewardAmt := sdk.NewInt(100000000)
	err := input.BankKeeper.MintCoins(input.Ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(merlion.AttoLionDenom, rewardAmt)))
	require.NoError(t, err)

	oracle.EndBlocker(input.Ctx.WithBlockHeight(1), input.OracleKeeper)

	denom2Rate, err := input.OracleKeeper.GetExchangeRate(input.Ctx, denom2)
	require.NoError(t, err)
	usm, err := input.OracleKeeper.GetExchangeRate(input.Ctx, merlion.MicroUSMDenom)
	require.NoError(t, err)

	// legacy version case
	require.NotEqual(t, usmExchangeRateWithErr, usm)

	// new version case
	require.Equal(t, denom2ExchangeRate, denom2Rate)
	require.Equal(t, usmExchangeRate, usm)

	rewardDistributedWindow := input.OracleKeeper.RewardDistributionWindow(input.Ctx)
	expectedRewardAmt := sdk.NewDecFromInt(rewardAmt.QuoRaw(8).MulRaw(2)).QuoInt64(int64(rewardDistributedWindow)).TruncateInt()
	expectedRewardAmt2 := sdk.NewDecFromInt(rewardAmt.QuoRaw(8).MulRaw(1)).QuoInt64(int64(rewardDistributedWindow)).TruncateInt()
	rewards := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[0])
	require.Equal(t, expectedRewardAmt, rewards.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards1 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[1])
	require.Equal(t, expectedRewardAmt2, rewards1.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards2 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[2])
	require.Equal(t, expectedRewardAmt2, rewards2.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards3 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[3])
	require.Equal(t, expectedRewardAmt, rewards3.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
	rewards4 := input.DistrKeeper.GetValidatorOutstandingRewards(input.Ctx.WithBlockHeight(2), keeper.ValAddrs[4])
	require.Equal(t, expectedRewardAmt, rewards4.Rewards.AmountOf(merlion.AttoLionDenom).TruncateInt())
}

func TestInvalidVotesSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	input.OracleKeeper.SetParams(input.Ctx, params)
	input.OracleKeeper.SetVoteTarget(input.Ctx, denom2)

	votePeriodsPerWindow := sdk.NewDec(int64(input.OracleKeeper.SlashWindow(input.Ctx))).QuoInt64(int64(input.OracleKeeper.VotePeriod(input.Ctx))).TruncateInt64()
	slashFraction := input.OracleKeeper.SlashFraction(input.Ctx)
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := uint64(0); i < uint64(sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64()); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, denom2
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 0)

		// Account 2, denom2, miss vote
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate.Add(sdk.NewDec(100000000000000))}}, 1)

		// Account 3, denom2
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 2)

		oracle.EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, i+1, input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())

	// one more miss vote will inccur keeper.ValAddrs[1] slashing
	// Account 1, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 0)

	// Account 2, denom2, miss vote
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate.Add(sdk.NewDec(100000000000000))}}, 1)

	// Account 3, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 2)

	input.Ctx = input.Ctx.WithBlockHeight(votePeriodsPerWindow - 1)
	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	validator = input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, sdk.OneDec().Sub(slashFraction).MulInt(stakingAmt).TruncateInt(), validator.GetBondedTokens())
}

func TestNotPassedBallotSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	input.OracleKeeper.SetParams(input.Ctx, params)

	// reset vote targets
	input.OracleKeeper.ClearVoteTargets(input.Ctx)
	input.OracleKeeper.SetVoteTarget(input.Ctx, denom2)

	input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

	// Account 1, denom2
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 0)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestAbstainSlashing(t *testing.T) {
	input, h := setup(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	input.OracleKeeper.SetParams(input.Ctx, params)

	// clear tobin tax to reset vote targets
	input.OracleKeeper.ClearVoteTargets(input.Ctx)
	input.OracleKeeper.SetVoteTarget(input.Ctx, denom2)

	votePeriodsPerWindow := sdk.NewDec(int64(input.OracleKeeper.SlashWindow(input.Ctx))).QuoInt64(int64(input.OracleKeeper.VotePeriod(input.Ctx))).TruncateInt64()
	minValidPerWindow := input.OracleKeeper.MinValidPerWindow(input.Ctx)

	for i := uint64(0); i <= uint64(sdk.OneDec().Sub(minValidPerWindow).MulInt64(votePeriodsPerWindow).TruncateInt64()); i++ {
		input.Ctx = input.Ctx.WithBlockHeight(input.Ctx.BlockHeight() + 1)

		// Account 1, denom2
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 0)

		// Account 2, denom2, abstain vote
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: sdk.ZeroDec()}}, 1)

		// Account 3, denom2
		makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 2)

		oracle.EndBlocker(input.Ctx, input.OracleKeeper)
		require.Equal(t, uint64(0), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	}

	validator := input.StakingKeeper.Validator(input.Ctx, keeper.ValAddrs[1])
	require.Equal(t, stakingAmt, validator.GetBondedTokens())
}

func TestVoteTargets(t *testing.T) {
	input, h := setup(t)

	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	// missing
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))

	// reset vote targets
	input.OracleKeeper.ClearVoteTargets(input.Ctx)
	input.OracleKeeper.SetVoteTarget(input.Ctx, denom2)

	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 0)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 1)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: randomExchangeRate}}, 2)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)

	// no missing
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[0]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[1]))
	require.Equal(t, uint64(1), input.OracleKeeper.GetMissCounter(input.Ctx, keeper.ValAddrs[2]))
}

func TestAbstainWithSmallStakingPower(t *testing.T) {
	input, h := setupWithSmallVotingPower(t)

	// reset vote targets
	input.OracleKeeper.ClearVoteTargets(input.Ctx)
	input.OracleKeeper.SetVoteTarget(input.Ctx, denom2)
	makeAggregatePrevoteAndVote(t, input, h, 0, sdk.DecCoins{{Denom: denom2, Amount: sdk.ZeroDec()}}, 0)

	oracle.EndBlocker(input.Ctx, input.OracleKeeper)
	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx, denom2)
	require.Error(t, err)
}

func makeAggregatePrevoteAndVote(t *testing.T, input keeper.TestInput, h sdk.Handler, height int64, rates sdk.DecCoins, idx int) {
	// Account 1, denom1
	salt := "1"
	hash := types.GetAggregateVoteHash(salt, decCoins2ExchangeRates(rates), keeper.ValAddrs[idx])

	prevoteMsg := types.NewMsgAggregateExchangeRatePrevote(hash, keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err := h(input.Ctx.WithBlockHeight(height), prevoteMsg)
	require.NoError(t, err)

	voteMsg := types.NewMsgAggregateExchangeRateVote(salt, decCoins2ExchangeRates(rates), keeper.Addrs[idx], keeper.ValAddrs[idx])
	_, err = h(input.Ctx.WithBlockHeight(height+1), voteMsg)
	require.NoError(t, err)
}

func decCoins2ExchangeRates(coins sdk.DecCoins) string {
	if len(coins) == 0 {
		return ""
	}

	out := ""
	for _, coin := range coins {
		out += fmt.Sprintf("%v:%v,", coin.Denom, coin.Amount)
	}

	return out[:len(out)-1]
}
