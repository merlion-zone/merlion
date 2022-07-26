package keeper

import (
	"bytes"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/oracle/types"
	"github.com/stretchr/testify/require"
)

const (
	fooDenom1 = "bar1"
	fooDenom2 = "bar2"
	fooDenom3 = "bar3"
	fooDenom4 = "bar4"
	fooDenom5 = "bar5"
	fooDenom6 = "bar6"

	MicroUnit = int64(1e6)
)

func TestExchangeRate(t *testing.T) {
	input := CreateTestInput(t)

	bar4ExchangeRate := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(MicroUnit)
	bar5ExchangeRate := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(MicroUnit)
	bar2ExchangeRate := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetExchangeRate(input.Ctx, fooDenom4, bar4ExchangeRate)
	rate, err := input.OracleKeeper.GetExchangeRate(input.Ctx, fooDenom4)
	require.NoError(t, err)
	require.Equal(t, bar4ExchangeRate, rate)

	input.OracleKeeper.SetExchangeRate(input.Ctx, fooDenom5, bar5ExchangeRate)
	rate, err = input.OracleKeeper.GetExchangeRate(input.Ctx, fooDenom5)
	require.NoError(t, err)
	require.Equal(t, bar5ExchangeRate, rate)

	input.OracleKeeper.SetExchangeRate(input.Ctx, fooDenom2, bar2ExchangeRate)
	rate, err = input.OracleKeeper.GetExchangeRate(input.Ctx, fooDenom2)
	require.NoError(t, err)
	require.Equal(t, bar2ExchangeRate, rate)

	input.OracleKeeper.DeleteExchangeRate(input.Ctx, fooDenom2)
	_, err = input.OracleKeeper.GetExchangeRate(input.Ctx, fooDenom2)
	require.Error(t, err)

	numExchangeRates := 0
	handler := func(denom string, exchangeRate sdk.Dec) (stop bool) {
		numExchangeRates = numExchangeRates + 1
		return false
	}
	input.OracleKeeper.IterateExchangeRates(input.Ctx, handler)

	require.True(t, numExchangeRates == 2)
}

func TestIterateExchangeRates(t *testing.T) {
	input := CreateTestInput(t)

	bar4ExchangeRate := sdk.NewDecWithPrec(839, int64(OracleDecPrecision)).MulInt64(MicroUnit)
	bar5ExchangeRate := sdk.NewDecWithPrec(4995, int64(OracleDecPrecision)).MulInt64(MicroUnit)
	bar2ExchangeRate := sdk.NewDecWithPrec(2838, int64(OracleDecPrecision)).MulInt64(MicroUnit)
	bar1ExchangeRate := sdk.NewDecWithPrec(3282384, int64(OracleDecPrecision)).MulInt64(MicroUnit)

	// Set & get rates
	input.OracleKeeper.SetExchangeRate(input.Ctx, fooDenom4, bar4ExchangeRate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, fooDenom5, bar5ExchangeRate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, fooDenom2, bar2ExchangeRate)
	input.OracleKeeper.SetExchangeRate(input.Ctx, fooDenom1, bar1ExchangeRate)

	input.OracleKeeper.IterateExchangeRates(input.Ctx, func(denom string, rate sdk.Dec) (stop bool) {
		switch denom {
		case fooDenom4:
			require.Equal(t, bar4ExchangeRate, rate)
		case fooDenom5:
			require.Equal(t, bar5ExchangeRate, rate)
		case fooDenom2:
			require.Equal(t, bar2ExchangeRate, rate)
		case fooDenom1:
			require.Equal(t, bar1ExchangeRate, rate)
		}
		return false
	})

}

func TestRewardPool(t *testing.T) {
	input := CreateTestInput(t)

	fees := sdk.NewCoins(sdk.NewCoin(fooDenom3, sdk.NewInt(1000)))
	acc := input.AccountKeeper.GetModuleAccount(input.Ctx, types.ModuleName)
	err := FundAccount(input, acc.GetAddress(), fees)
	if err != nil {
		panic(err) // never occurs
	}

	KFees := input.OracleKeeper.GetRewardPool(input.Ctx, fooDenom3)
	require.Equal(t, fees[0], KFees)
}

func TestParams(t *testing.T) {
	input := CreateTestInput(t)

	// Test default params setting
	input.OracleKeeper.SetParams(input.Ctx, types.DefaultParams())
	params := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, params)

	// Test custom params setting
	votePeriod := uint64(10)
	voteThreshold := sdk.NewDecWithPrec(33, 2)
	oracleRewardBand := sdk.NewDecWithPrec(1, 2)
	rewardDistributionWindow := uint64(10000000000000)
	slashFraction := sdk.NewDecWithPrec(1, 2)
	slashWindow := uint64(1000)
	minValidPerWindow := sdk.NewDecWithPrec(1, 4)

	// Should really test validateParams, but skipping because obvious
	newParams := types.Params{
		VotePeriod:               votePeriod,
		VoteThreshold:            voteThreshold,
		RewardBand:               oracleRewardBand,
		RewardDistributionWindow: rewardDistributionWindow,
		SlashFraction:            slashFraction,
		SlashWindow:              slashWindow,
		MinValidPerWindow:        minValidPerWindow,
	}
	input.OracleKeeper.SetParams(input.Ctx, newParams)

	storedParams := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, storedParams)
	require.Equal(t, storedParams, newParams)
}

func TestFeederDelegation(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetFeederDelegation(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[0], delegate)

	input.OracleKeeper.SetFeederDelegation(input.Ctx, ValAddrs[0], Addrs[1])
	delegate = input.OracleKeeper.GetFeederDelegation(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[1], delegate)
}

func TestIterateFeederDelegations(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	delegate := input.OracleKeeper.GetFeederDelegation(input.Ctx, ValAddrs[0])
	require.Equal(t, Addrs[0], delegate)

	input.OracleKeeper.SetFeederDelegation(input.Ctx, ValAddrs[0], Addrs[1])

	var delegators []sdk.ValAddress
	var delegates []sdk.AccAddress
	input.OracleKeeper.IterateFeederDelegations(input.Ctx, func(delegator sdk.ValAddress, delegate sdk.AccAddress) (stop bool) {
		delegators = append(delegators, delegator)
		delegates = append(delegates, delegate)
		return false
	})

	require.Equal(t, 1, len(delegators))
	require.Equal(t, 1, len(delegates))
	require.Equal(t, ValAddrs[0], delegators[0])
	require.Equal(t, Addrs[1], delegates[0])
}

func TestMissCounter(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, uint64(0), counter)

	missCounter := uint64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[0], missCounter)
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, missCounter, counter)

	input.OracleKeeper.DeleteMissCounter(input.Ctx, ValAddrs[0])
	counter = input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, uint64(0), counter)
}

func TestIterateMissCounters(t *testing.T) {
	input := CreateTestInput(t)

	// Test default getters and setters
	counter := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, uint64(0), counter)

	missCounter := uint64(10)
	input.OracleKeeper.SetMissCounter(input.Ctx, ValAddrs[1], missCounter)

	var operators []sdk.ValAddress
	var missCounters []uint64
	input.OracleKeeper.IterateMissCounters(input.Ctx, func(delegator sdk.ValAddress, missCounter uint64) (stop bool) {
		operators = append(operators, delegator)
		missCounters = append(missCounters, missCounter)
		return false
	})

	require.Equal(t, 1, len(operators))
	require.Equal(t, 1, len(missCounters))
	require.Equal(t, ValAddrs[1], operators[0])
	require.Equal(t, missCounter, missCounters[0])
}

func TestAggregatePrevoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	hash := types.GetAggregateVoteHash("salt", "foo:100,bar:1000", sdk.ValAddress(Addrs[0]))
	aggregatePrevote := types.NewAggregateExchangeRatePrevote(hash, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[0]), aggregatePrevote)

	KPrevote, err := input.OracleKeeper.GetAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, aggregatePrevote, KPrevote)

	input.OracleKeeper.DeleteAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[0]))
	_, err = input.OracleKeeper.GetAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestAggregatePrevoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	hash := types.GetAggregateVoteHash("salt", "foo:100,bar:1000", sdk.ValAddress(Addrs[0]))
	aggregatePrevote1 := types.NewAggregateExchangeRatePrevote(hash, sdk.ValAddress(Addrs[0]), 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[0]), aggregatePrevote1)

	hash2 := types.GetAggregateVoteHash("salt", "foo:100,bar:1000", sdk.ValAddress(Addrs[1]))
	aggregatePrevote2 := types.NewAggregateExchangeRatePrevote(hash2, sdk.ValAddress(Addrs[1]), 0)
	input.OracleKeeper.SetAggregateExchangeRatePrevote(input.Ctx, sdk.ValAddress(Addrs[1]), aggregatePrevote2)

	i := 0
	bigger := bytes.Compare(Addrs[0], Addrs[1])
	input.OracleKeeper.IterateAggregateExchangeRatePrevotes(input.Ctx, func(voter sdk.ValAddress, p types.AggregateExchangeRatePrevote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, aggregatePrevote1, p)
			require.Equal(t, voter.String(), p.Voter)
		} else {
			require.Equal(t, aggregatePrevote2, p)
			require.Equal(t, voter.String(), p.Voter)
		}

		i++
		return false
	})
}

func TestAggregateVoteAddDelete(t *testing.T) {
	input := CreateTestInput(t)

	aggregateVote := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
		{Denom: "foo", ExchangeRate: sdk.NewDec(-1)},
		{Denom: "foo", ExchangeRate: sdk.NewDec(0)},
		{Denom: "foo", ExchangeRate: sdk.NewDec(1)},
	}, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[0]), aggregateVote)

	KVote, err := input.OracleKeeper.GetAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.NoError(t, err)
	require.Equal(t, aggregateVote, KVote)

	input.OracleKeeper.DeleteAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[0]))
	_, err = input.OracleKeeper.GetAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[0]))
	require.Error(t, err)
}

func TestAggregateVoteIterate(t *testing.T) {
	input := CreateTestInput(t)

	aggregateVote1 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
		{Denom: "foo", ExchangeRate: sdk.NewDec(-1)},
		{Denom: "foo", ExchangeRate: sdk.NewDec(0)},
		{Denom: "foo", ExchangeRate: sdk.NewDec(1)},
	}, sdk.ValAddress(Addrs[0]))
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[0]), aggregateVote1)

	aggregateVote2 := types.NewAggregateExchangeRateVote(types.ExchangeRateTuples{
		{Denom: "foo", ExchangeRate: sdk.NewDec(-1)},
		{Denom: "foo", ExchangeRate: sdk.NewDec(0)},
		{Denom: "foo", ExchangeRate: sdk.NewDec(1)},
	}, sdk.ValAddress(Addrs[1]))
	input.OracleKeeper.SetAggregateExchangeRateVote(input.Ctx, sdk.ValAddress(Addrs[1]), aggregateVote2)

	i := 0
	bigger := bytes.Compare(address.MustLengthPrefix(Addrs[0]), address.MustLengthPrefix(Addrs[1]))
	input.OracleKeeper.IterateAggregateExchangeRateVotes(input.Ctx, func(voter sdk.ValAddress, p types.AggregateExchangeRateVote) (stop bool) {
		if (i == 0 && bigger == -1) || (i == 1 && bigger == 1) {
			require.Equal(t, aggregateVote1, p)
			require.Equal(t, voter.String(), p.Voter)
		} else {
			require.Equal(t, aggregateVote2, p)
			require.Equal(t, voter.String(), p.Voter)
		}

		i++
		return false
	})
}

func TestVoteTargets(t *testing.T) {
	input := CreateTestInput(t)

	input.OracleKeeper.ClearVoteTargets(input.Ctx)

	expectedVoteTargets := []string{"bar", "foo", "whoowhoo"}
	for _, voteTarget := range expectedVoteTargets {
		input.OracleKeeper.SetVoteTarget(input.Ctx, voteTarget)
	}

	for _, voteTarget := range expectedVoteTargets {
		require.True(t, input.OracleKeeper.IsVoteTarget(input.Ctx, voteTarget))
	}

	i := 0
	input.OracleKeeper.IterateVoteTargets(input.Ctx, func(denom string) (stop bool) {
		require.Equal(t, expectedVoteTargets[i], denom)
		i++
		return false
	})

	voteTargets := input.OracleKeeper.GetVoteTargets(input.Ctx)
	require.Equal(t, expectedVoteTargets, voteTargets)
}

func TestTargets(t *testing.T) {
	input := CreateTestInput(t)

	expectedTargets := []string{"bar", "foo", "whoowhoo"}
	for _, target := range expectedTargets {
		input.OracleKeeper.SetTarget(input.Ctx, target)
	}

	for _, target := range expectedTargets {
		require.True(t, input.OracleKeeper.IsTarget(input.Ctx, target))
	}

	i := 0
	input.OracleKeeper.IterateTargets(input.Ctx, func(denom string) (stop bool) {
		require.Equal(t, expectedTargets[i], denom)
		i++
		return false
	})

	targets := input.OracleKeeper.GetTargets(input.Ctx)
	require.Equal(t, expectedTargets, targets)
}

func TestValidateFeeder(t *testing.T) {
	// initial setup
	input := CreateTestInput(t)
	addr, val := ValAddrs[0], ValPubKeys[0]
	addr1, val1 := ValAddrs[1], ValPubKeys[1]
	amt := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	sh := staking.NewHandler(input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh(ctx, NewTestMsgCreateValidator(addr, val, amt))
	require.NoError(t, err)
	_, err = sh(ctx, NewTestMsgCreateValidator(addr1, val1, amt))
	require.NoError(t, err)
	staking.EndBlocker(ctx, input.StakingKeeper)

	require.Equal(
		t, input.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr)),
		sdk.NewCoins(sdk.NewCoin(input.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, addr).GetBondedTokens())
	require.Equal(
		t, input.BankKeeper.GetAllBalances(ctx, sdk.AccAddress(addr1)),
		sdk.NewCoins(sdk.NewCoin(input.StakingKeeper.GetParams(ctx).BondDenom, InitTokens.Sub(amt))),
	)
	require.Equal(t, amt, input.StakingKeeper.Validator(ctx, addr1).GetBondedTokens())

	require.NoError(t, input.OracleKeeper.ValidateFeeder(input.Ctx, sdk.AccAddress(addr), sdk.ValAddress(addr)))
	require.NoError(t, input.OracleKeeper.ValidateFeeder(input.Ctx, sdk.AccAddress(addr1), sdk.ValAddress(addr1)))

	// delegate works
	input.OracleKeeper.SetFeederDelegation(input.Ctx, sdk.ValAddress(addr), sdk.AccAddress(addr1))
	require.NoError(t, input.OracleKeeper.ValidateFeeder(input.Ctx, sdk.AccAddress(addr1), sdk.ValAddress(addr)))
	require.Error(t, input.OracleKeeper.ValidateFeeder(input.Ctx, sdk.AccAddress(Addrs[2]), sdk.ValAddress(addr)))

	// only active validators can do oracle votes
	validator, found := input.StakingKeeper.GetValidator(input.Ctx, sdk.ValAddress(addr))
	require.True(t, found)
	validator.Status = stakingtypes.Unbonded
	input.StakingKeeper.SetValidator(input.Ctx, validator)
	require.Error(t, input.OracleKeeper.ValidateFeeder(input.Ctx, sdk.AccAddress(addr1), sdk.ValAddress(addr)))
}
