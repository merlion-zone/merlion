package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	vekeeper "github.com/merlion-zone/merlion/x/ve/keeper"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
	"github.com/merlion-zone/merlion/x/voter/types"
)

func (k Keeper) CreateGauge(ctx sdk.Context, depoistDenom string) {
	if !k.gaugeKeeper.HasGauge(ctx, depoistDenom) {
		k.gaugeKeeper.CreateGauge(ctx, depoistDenom)
	}

	k.updateClaimableForGauge(ctx, depoistDenom)
}

func (k Keeper) Abstain(ctx sdk.Context, veID uint64) {
	poolDenoms := k.gaugeKeeper.GetGauges(ctx)

	totalVotes := k.GetTotalVotes(ctx)

	for _, poolDenom := range poolDenoms {
		k.updateClaimableForGauge(ctx, poolDenom)

		weightedVotes := k.GetPoolWeightedVotesByUser(ctx, veID, poolDenom)
		if weightedVotes.IsZero() {
			continue
		}

		k.DeletePoolWeightedVotesByUser(ctx, veID, poolDenom)
		k.SetPoolWeightedVotes(ctx, poolDenom, k.GetPoolWeightedVotes(ctx, poolDenom).Sub(weightedVotes))

		totalVotes = totalVotes.Sub(weightedVotes.Abs())

		// only concurring votes need canceling bribe deposit
		if weightedVotes.IsPositive() {
			bribe := k.gaugeKeeper.Bribe(ctx, poolDenom)
			err := bribe.Withdraw(ctx, veID, weightedVotes)
			if err != nil {
				panic(err)
			}
		}
	}

	k.DeleteTotalVotesByUser(ctx, veID)
	k.SetTotalVotes(ctx, totalVotes)

	k.veKeeper.SetVeVoted(ctx, veID, false)
}

func (k Keeper) Vote(ctx sdk.Context, veID uint64, poolWeights map[string]sdk.Dec) {
	// reset voting for user
	k.Abstain(ctx, veID)

	votingPower := k.veKeeper.GetVotingPower(ctx, veID, uint64(ctx.BlockTime().Unix()), 0)

	totalVotesByUser := k.GetTotalVotesByUser(ctx, veID)
	totalVotes := k.GetTotalVotes(ctx)

	totalWeights := sdk.ZeroDec()
	for poolDenom, weight := range poolWeights {
		if !k.gaugeKeeper.HasGauge(ctx, poolDenom) {
			panic("gauge not found")
		}
		k.updateClaimableForGauge(ctx, poolDenom)

		totalWeights = totalWeights.Add(weight.Abs())

		// <votes for gauge> = <voting power> * <weight for gauge>
		weightedVotes := votingPower.ToDec().Mul(weight).TruncateInt()
		if weightedVotes.IsZero() {
			panic("gauge weighted votes must be nonzero")
		}

		// total votes also accumulate negative votes
		totalVotesByUser = totalVotesByUser.Add(weightedVotes.Abs())
		totalVotes = totalVotes.Add(weightedVotes.Abs())

		k.SetPoolWeightedVotesByUser(ctx, veID, poolDenom, weightedVotes)
		poolWeightedVotes := k.GetPoolWeightedVotes(ctx, poolDenom).Add(weightedVotes)
		k.SetPoolWeightedVotes(ctx, poolDenom, poolWeightedVotes)

		// only concurring votes will receive bribe
		if weightedVotes.IsPositive() {
			bribe := k.gaugeKeeper.Bribe(ctx, poolDenom)
			bribe.Deposit(ctx, veID, weightedVotes)
		}
	}

	if !totalWeights.Equal(sdk.OneDec()) {
		panic("sum of pool weights must be one")
	}

	k.SetTotalVotesByUser(ctx, veID, totalVotesByUser)
	k.SetTotalVotes(ctx, totalVotes)

	k.veKeeper.SetVeVoted(ctx, veID, true)
}

// Poke adjusts votes due to updated voting power of user
func (k Keeper) Poke(ctx sdk.Context, veID uint64) {
	totalVotesByUser := k.GetTotalVotesByUser(ctx, veID)
	if !totalVotesByUser.IsPositive() {
		// no voting so no poke
		return
	}

	poolDenoms := k.gaugeKeeper.GetGauges(ctx)

	totalWeights := sdk.ZeroDec()
	poolWeights := make(map[string]sdk.Dec)

	var fineTuning string
	for _, poolDenom := range poolDenoms {
		weightedVotes := k.GetPoolWeightedVotesByUser(ctx, veID, poolDenom)
		if weightedVotes.IsZero() {
			continue
		}
		weight := weightedVotes.ToDec().QuoInt(totalVotesByUser)
		poolWeights[poolDenom] = weight
		totalWeights = totalWeights.Add(weight)
		fineTuning = poolDenom
	}
	if totalWeights.Equal(sdk.OneDec()) {
		// it's ok to compensate for accuracy loss
		poolWeights[fineTuning] = poolWeights[fineTuning].Add(sdk.OneDec().Sub(totalWeights))
	}

	k.Vote(ctx, veID, poolWeights)
}

func (k Keeper) DepositReward(ctx sdk.Context, sender sdk.AccAddress, amount sdk.Int) {
	coin := sdk.NewCoin(k.veKeeper.LockDenom(ctx), amount)
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(coin))
	if err != nil {
		panic(err)
	}

	index := k.GetIndex(ctx)
	totalVotes := k.GetTotalVotes(ctx)
	index = index.Add(amount.Quo(totalVotes))
	k.SetIndex(ctx, index)
}

func (k Keeper) EmitReward(ctx sdk.Context) {
	emitter := vekeeper.NewEmitter(k.veKeeper.(vekeeper.Keeper))
	emission := emitter.Emit(ctx)
	if emission.IsPositive() {
		k.DepositReward(ctx, k.accountKeeper.GetModuleAddress(vetypes.EmissionPoolName), emission)
	}
}

func (k Keeper) DistributeReward(ctx sdk.Context, poolDenom string) {
	gauge := k.gaugeKeeper.Gauge(ctx, poolDenom)

	k.EmitReward(ctx)

	k.updateClaimableForGauge(ctx, poolDenom)

	claimable := k.GetClaimableRewardByGauge(ctx, poolDenom)
	if claimable.GT(gauge.RemainingReward(ctx, poolDenom)) && claimable.QuoRaw(vetypes.RegulatedPeriod).IsPositive() {
		k.SetClaimableRewardByGauge(ctx, poolDenom, sdk.ZeroInt())

		err := gauge.DepositReward(ctx, k.accountKeeper.GetModuleAddress(types.ModuleName), poolDenom, claimable)
		if err != nil {
			panic(err)
		}
	}
}

func (k Keeper) updateClaimableForGauge(ctx sdk.Context, poolDenom string) {
	// votes owned by this gauge
	votes := k.GetPoolWeightedVotes(ctx, poolDenom)
	// cumulative reward per vote
	index := k.GetIndex(ctx)

	if votes.IsPositive() {
		// cumulative reward per vote which was recorded at last update for this gauge
		indexLast := k.GetIndexAtLastUpdatedByGauge(ctx, poolDenom)

		delta := index.Sub(indexLast)
		if delta.IsPositive() {
			// delta claimable reward = delta index * votes
			claimableDelta := delta.Mul(votes)

			claimable := k.GetClaimableRewardByGauge(ctx, poolDenom)
			claimable = claimable.Add(claimableDelta)
			k.SetClaimableRewardByGauge(ctx, poolDenom, claimable)
		}
	}

	// record cumulative reward per vote for this gauge
	k.SetIndexAtLastUpdatedByGauge(ctx, poolDenom, index)
}
