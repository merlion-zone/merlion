package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashAndResetMissCounters slashes any operator who over criteria and clears all operators' miss counter to zero.
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	stakingKeeper := k.StakingKeeper()
	height := ctx.BlockHeight()
	distributionHeight := height - sdk.ValidatorUpdateDelay - 1

	// slash_window / vote_period
	votePeriodsPerWindow := uint64(
		sdk.NewDec(int64(k.SlashWindow(ctx))).
			QuoInt64(int64(k.VotePeriod(ctx))).
			TruncateInt64(),
	)
	minValidPerWindow := k.MinValidPerWindow(ctx)
	slashFraction := k.SlashFraction(ctx)
	powerReduction := stakingKeeper.PowerReduction(ctx)

	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) bool {

		// Calculate valid vote rate: (SlashWindow - MissCounter) / SlashWindow
		validVoteRate := sdk.NewDecFromInt(
			sdk.NewInt(int64(votePeriodsPerWindow - missCounter))).
			QuoInt64(int64(votePeriodsPerWindow))

		// Penalize the validator whose valid vote rate is smaller than min threshold
		if validVoteRate.LT(minValidPerWindow) {
			validator := stakingKeeper.Validator(ctx, operator)
			if validator.IsBonded() && !validator.IsJailed() {
				consAddr, err := validator.GetConsAddr()
				if err != nil {
					panic(err)
				}

				stakingKeeper.Slash(
					ctx, consAddr,
					distributionHeight, validator.GetConsensusPower(powerReduction), slashFraction,
				)
				stakingKeeper.Jail(ctx, consAddr)
			}
		}

		k.DeleteMissCounter(ctx, operator)
		return false
	})
}
