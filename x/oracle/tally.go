package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/merlion-zone/merlion/x/oracle/keeper"
	"github.com/merlion-zone/merlion/x/oracle/types"
)

// Tally calculates the median and returns it. Sets the set of voters to be rewarded, i.e., voted within
// a reasonable spread from the weighted median to the store
// CONTRACT: pb must be sorted
func Tally(ctx sdk.Context, pb types.ExchangeRateBallot, rewardBand sdk.Dec, validatorClaimMap map[string]types.Claim) (weightedMedian sdk.Dec) {
	weightedMedian = pb.WeightedMedianWithAssertion()

	standardDeviation := pb.StandardDeviation(weightedMedian)
	rewardSpread := weightedMedian.Mul(rewardBand.QuoInt64(2))

	if standardDeviation.GT(rewardSpread) {
		rewardSpread = standardDeviation
	}

	for _, vote := range pb {
		// Filter ballot winners & abstain voters
		if (vote.ExchangeRate.GTE(weightedMedian.Sub(rewardSpread)) &&
			vote.ExchangeRate.LTE(weightedMedian.Add(rewardSpread))) ||
			!vote.ExchangeRate.IsPositive() {

			key := vote.Voter.String()
			claim := validatorClaimMap[key]
			claim.Weight += vote.Power
			claim.WinCount++
			validatorClaimMap[key] = claim
		}
	}

	return
}

// ballotIsPassing returns the ballot power and whether the power is passing the threshold amount of voting power.
func ballotIsPassing(ballot types.ExchangeRateBallot, thresholdVotes sdk.Int) (sdk.Int, bool) {
	ballotPower := sdk.NewInt(ballot.Power())
	return ballotPower, !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes)
}

// PickReferenceMer chooses Reference Mer with the highest voter turnout.
// If the voting power of the two denominations is the same,
// select reference Mer in alphabetical order.
func PickReferenceMer(ctx sdk.Context, k keeper.Keeper, voteTargets map[string]struct{}, voteMap map[string]types.ExchangeRateBallot) string {
	largestBallotPower := int64(0)
	referenceMer := ""

	stakingKeeper := k.StakingKeeper()
	totalBondedPower := sdk.TokensToConsensusPower(stakingKeeper.TotalBondedTokens(ctx), stakingKeeper.PowerReduction(ctx))
	voteThreshold := k.VoteThreshold(ctx)
	thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()

	for denom, ballot := range voteMap {
		// If denom is not in the voteTargets, or the ballot for it has failed, then skip
		// and remove it from voteMap for iteration efficiency.
		if _, exists := voteTargets[denom]; !exists {
			delete(voteMap, denom)
			continue
		}

		ballotPower := int64(0)

		// If the ballot is not passed, remove it from the voteTargets array
		// to prevent slashing validators who did valid vote.
		if power, ok := ballotIsPassing(ballot, thresholdVotes); ok {
			ballotPower = power.Int64()
		} else {
			delete(voteTargets, denom)
			delete(voteMap, denom)
			continue
		}

		if ballotPower > largestBallotPower || largestBallotPower == 0 {
			referenceMer = denom
			largestBallotPower = ballotPower
		} else if largestBallotPower == ballotPower && referenceMer > denom {
			referenceMer = denom
		}
	}

	return referenceMer
}
