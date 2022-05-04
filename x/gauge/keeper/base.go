package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

type Base struct {
	keeper Keeper

	depoistDenom string
	prefixKey    []byte
	isGauge      bool
}

func (b *Base) PoolDenom() string {
	return b.depoistDenom
}

func (b *Base) PoolName() string {
	name := types.BribePoolName
	if b.isGauge {
		name = types.GaugePoolName
	}
	return fmt.Sprintf("%s_%s", name, b.depoistDenom)
}

func (b *Base) EscrowPool(ctx sdk.Context) authtypes.ModuleAccountI {
	poolName := b.PoolName()
	acc := b.keeper.accountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(poolName))
	if acc != nil {
		macc, ok := acc.(authtypes.ModuleAccountI)
		if !ok {
			panic("account is not a module account")
		}
		return macc
	}

	// create a new module account
	macc := authtypes.NewEmptyModuleAccount(poolName)
	maccI := (b.keeper.accountKeeper.NewAccount(ctx, macc)).(authtypes.ModuleAccountI) // set the account number
	b.keeper.accountKeeper.SetModuleAccount(ctx, maccI)

	return maccI
}

func (b *Base) isRewardDenom(ctx sdk.Context, denom string) bool {
	var ok bool
	b.IterateRewards(ctx, func(reward types.Reward) (stop bool) {
		if reward.Denom == denom {
			ok = true
			return true
		}
		return false
	})
	return ok
}

func (b *Base) getRewardDenoms(ctx sdk.Context) []string {
	var denoms []string
	b.IterateRewards(ctx, func(reward types.Reward) (stop bool) {
		denoms = append(denoms, reward.Denom)
		return false
	})
	return denoms
}

func (b *Base) findPriorEpoch(ctx sdk.Context, veID uint64, rewardDenom string, timestamp uint64) uint64 {
	var epoch uint64
	if veID != vetypes.EmptyVeID {
		if len(rewardDenom) != 0 {
			panic("either veID or rewardDenom")
		}
		epoch = b.GetUserEpoch(ctx, veID)
	} else if len(rewardDenom) != 0 {
		epoch = b.GetRewardEpoch(ctx, rewardDenom)
	} else {
		epoch = b.GetEpoch(ctx)
	}
	if epoch == EmptyEpoch {
		return EmptyEpoch
	}

	getCheckpoint := func(epoch uint64) types.Checkpoint {
		if veID != vetypes.EmptyVeID {
			return b.GetUserCheckpoint(ctx, veID, epoch)
		} else if len(rewardDenom) != 0 {
			return b.GetRewardCheckpoint(ctx, rewardDenom, epoch)
		} else {
			return b.GetCheckpoint(ctx, epoch)
		}
	}

	// check most recent
	if getCheckpoint(epoch).Timestamp <= timestamp {
		return epoch
	}

	// check initial
	if getCheckpoint(FirstEpoch).Timestamp > timestamp {
		return EmptyEpoch
	}

	// binary search
	lower := uint64(0)
	upper := epoch
	for upper > lower {
		mid := (lower + upper + 1) / 2 // ceil
		point := getCheckpoint(mid)
		if point.Timestamp == timestamp {
			return mid
		} else if point.Timestamp < timestamp {
			lower = mid
		} else {
			upper = mid - 1
		}
	}
	return lower
}

func (b *Base) findPriorRewardPerTicket(ctx sdk.Context, rewardDenom string, timestamp uint64) sdk.Int {
	epoch := b.findPriorEpoch(ctx, vetypes.EmptyVeID, rewardDenom, timestamp)
	return b.GetRewardCheckpoint(ctx, rewardDenom, epoch).Amount
}

func (b *Base) lastTimeRewardAvailable(ctx sdk.Context, rewardDenom string) uint64 {
	return merlion.Min(uint64(ctx.BlockTime().Unix()), b.GetReward(ctx, rewardDenom).FinishTime)
}

func (b *Base) rewardPerTicket(ctx sdk.Context, rewardDenom string) sdk.Int {
	totalAmount := b.GetTotalDepositedAmount(ctx)
	if b.isGauge {
		totalAmount = b.GetTotalDerivedAmount(ctx)
	}

	reward := b.GetReward(ctx, rewardDenom)
	if totalAmount.IsZero() {
		return reward.CumulativePerTicket
	}

	timeFromLast := b.lastTimeRewardAvailable(ctx, rewardDenom) - merlion.Min(reward.LastUpdateTime, reward.FinishTime)
	rewardFromLast := reward.Rate.MulRaw(int64(timeFromLast)).Quo(totalAmount)
	return reward.CumulativePerTicket.Add(rewardFromLast)
}

func (b *Base) derivedTickets(ctx sdk.Context, veID uint64) sdk.Int {
	if !b.isGauge {
		return sdk.ZeroInt()
	}

	totalDeposited := b.GetTotalDepositedAmount(ctx)
	deposited := b.GetDepositedAmountByUser(ctx, veID)

	derived := deposited.MulRaw(40).QuoRaw(100)

	adjusted := sdk.ZeroInt()
	totalPower := b.keeper.veKeeper.GetTotalVotingPower(ctx, uint64(ctx.BlockTime().Unix()), 0)
	if totalPower.IsPositive() {
		power := b.keeper.veKeeper.GetVotingPower(ctx, veID, uint64(ctx.BlockTime().Unix()), 0)
		adjusted = totalDeposited.Mul(power).Quo(totalPower).MulRaw(60).QuoRaw(100)
	}

	tickets := derived.Add(adjusted)
	if tickets.LT(deposited) {
		return tickets
	} else {
		return deposited
	}
}

func (b *Base) calculateRewardPerTicket(ctx sdk.Context, rewardDenom string, timestamp0, timestamp1, startTimestamp uint64, totalAmount sdk.Int) (rewardPerTicket sdk.Int, endTimestamp uint64) {
	reward := b.GetReward(ctx, rewardDenom)
	endTimestamp = merlion.Max(timestamp1, startTimestamp)
	duration := merlion.Min(endTimestamp, reward.FinishTime) - merlion.Min(merlion.Max(timestamp0, startTimestamp), reward.FinishTime)
	rewardPerTicket = reward.Rate.MulRaw(int64(duration)).Quo(totalAmount)
	return rewardPerTicket, endTimestamp
}

func (b *Base) updateRewardPerTicket(ctx sdk.Context, rewardDenom string) {
	epoch := b.GetEpoch(ctx)
	if epoch == EmptyEpoch {
		return
	}

	now := uint64(ctx.BlockTime().Unix())

	reward := b.GetReward(ctx, rewardDenom)
	startTimestamp := reward.LastUpdateTime

	if reward.Rate.IsZero() {
		reward.LastUpdateTime = now
		b.SetReward(ctx, rewardDenom, reward)
		return
	}

	rewardCumulativePerTicket := reward.CumulativePerTicket
	startEpoch := b.findPriorEpoch(ctx, 0, "", startTimestamp)

	for i := startEpoch; i < epoch-1; i++ {
		point := b.GetCheckpoint(ctx, i)
		if point.Amount.IsPositive() {
			nextPoint := b.GetCheckpoint(ctx, i+1)
			rewardPerTicket, endTimestamp := b.calculateRewardPerTicket(ctx, rewardDenom, point.Timestamp, nextPoint.Timestamp, startTimestamp, point.Amount)
			rewardCumulativePerTicket = rewardCumulativePerTicket.Add(rewardPerTicket)
			b.writeRewardPerTicketCheckpoint(ctx, rewardDenom, rewardCumulativePerTicket, endTimestamp)
			startTimestamp = endTimestamp
		}
	}

	point := b.GetCheckpoint(ctx, epoch)
	if point.Amount.IsPositive() {
		rewardPerTicket, _ := b.calculateRewardPerTicket(ctx, rewardDenom, merlion.Max(point.Timestamp, startTimestamp), b.lastTimeRewardAvailable(ctx, rewardDenom), startTimestamp, point.Amount)
		rewardCumulativePerTicket = rewardCumulativePerTicket.Add(rewardPerTicket)
		b.writeRewardPerTicketCheckpoint(ctx, rewardDenom, rewardCumulativePerTicket, now)
		startTimestamp = now
	}

	reward.CumulativePerTicket = rewardCumulativePerTicket
	reward.LastUpdateTime = startTimestamp
	b.SetReward(ctx, rewardDenom, reward)
}

func (b *Base) userReward(ctx sdk.Context, rewardDenom string, veID uint64) sdk.Int {
	epoch := b.GetUserEpoch(ctx, veID)
	if epoch == EmptyEpoch {
		return sdk.ZeroInt()
	}

	userReward := b.GetUserReward(ctx, rewardDenom, veID)
	startTimestamp := merlion.Max(userReward.LastClaimTime, b.GetRewardCheckpoint(ctx, rewardDenom, FirstEpoch).Timestamp)

	reward := sdk.ZeroInt()
	startEpoch := b.findPriorEpoch(ctx, veID, "", startTimestamp)

	for i := startEpoch; i < epoch-1; i++ {
		point := b.GetUserCheckpoint(ctx, veID, i)
		nextPoint := b.GetUserCheckpoint(ctx, veID, i+1)
		rewardPerTicket := b.findPriorRewardPerTicket(ctx, rewardDenom, point.Timestamp)
		nextRewardPerTicket := b.findPriorRewardPerTicket(ctx, rewardDenom, nextPoint.Timestamp)
		reward = reward.Add(point.Amount.Mul(nextRewardPerTicket.Sub(rewardPerTicket)))
	}

	point := b.GetUserCheckpoint(ctx, veID, epoch)
	rewardPerTicket := b.findPriorRewardPerTicket(ctx, rewardDenom, point.Timestamp)
	reward = reward.Add(point.Amount.Mul(b.rewardPerTicket(ctx, rewardDenom).Sub(sdk.MaxInt(rewardPerTicket, userReward.CumulativePerTicket))))

	return reward
}

func (b *Base) claimReward(ctx sdk.Context, veID uint64) (err error) {
	owner := b.keeper.nftKeeper.GetOwner(ctx, vetypes.VeNftClass.Id, vetypes.VeID(veID))
	pool := b.EscrowPool(ctx)

	denoms := b.getRewardDenoms(ctx)
	for _, rewardDenom := range denoms {
		b.updateRewardPerTicket(ctx, rewardDenom)

		rewardAmount := b.userReward(ctx, rewardDenom, veID)

		reward := b.GetReward(ctx, rewardDenom)
		userReward := b.GetUserReward(ctx, rewardDenom, veID)
		userReward.LastClaimTime = uint64(ctx.BlockTime().Unix())
		userReward.CumulativePerTicket = reward.CumulativePerTicket

		if rewardAmount.IsPositive() {
			coin := sdk.NewCoin(rewardDenom, rewardAmount)
			err = b.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, pool.GetName(), owner, sdk.NewCoins(coin))
			if err != nil {
				return err
			}
		}
	}

	b.deriveAmountForUser(ctx, veID)

	return nil
}

func (b *Base) deriveAmountForUser(ctx sdk.Context, veID uint64) {
	if b.isGauge {
		derived := b.GetDerivedAmountByUser(ctx, veID)
		totalDerived := b.GetTotalDerivedAmount(ctx)
		totalDerived = totalDerived.Sub(derived)

		derived = b.derivedTickets(ctx, veID)
		b.SetDerivedAmountByUser(ctx, veID, derived)
		totalDerived = totalDerived.Add(derived)
		b.SetTotalDerivedAmount(ctx, totalDerived)

		b.writeUserCheckpoint(ctx, veID, derived)
		b.writeCheckpoint(ctx)
	}
}

func (b *Base) RemainingReward(ctx sdk.Context, rewardDenom string) sdk.Int {
	now := uint64(ctx.BlockTime().Unix())
	reward := b.GetReward(ctx, rewardDenom)
	if reward.FinishTime <= now {
		return sdk.ZeroInt()
	}
	return reward.Rate.MulRaw(int64(reward.FinishTime - now))
}

func (b *Base) depositReward(ctx sdk.Context, sender sdk.AccAddress, rewardDenom string, amount sdk.Int) error {
	if rewardDenom == b.depoistDenom {
		// TODO: error
	}
	if !amount.IsPositive() {
		// TODO: error
	}

	now := uint64(ctx.BlockTime().Unix())

	reward := b.GetReward(ctx, rewardDenom)
	if reward.Rate.IsZero() {
		b.writeRewardPerTicketCheckpoint(ctx, rewardDenom, sdk.ZeroInt(), now)
	}
	b.updateRewardPerTicket(ctx, rewardDenom)
	reward = b.GetReward(ctx, rewardDenom)

	coin := sdk.NewCoin(rewardDenom, amount)
	err := b.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, b.EscrowPool(ctx).GetName(), sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	rewardAmount := amount
	if reward.FinishTime > now {
		remaining := reward.Rate.MulRaw(int64(reward.FinishTime - now))
		if rewardAmount.LTE(remaining) {
			// TODO: error
		}
		rewardAmount = rewardAmount.Add(remaining)
	}
	reward.Rate = rewardAmount.QuoRaw(vetypes.RegulatedPeriod)
	reward.FinishTime = now + vetypes.RegulatedPeriod

	if !reward.Rate.IsPositive() {
		// TODO: error
	}

	balance := b.keeper.bankKeeper.GetBalance(ctx, b.EscrowPool(ctx).GetAddress(), rewardDenom)
	if balance.Amount.LT(rewardAmount) {
		// TODO: error
	}

	b.SetReward(ctx, rewardDenom, reward)
	return nil
}

func (b *Base) writeUserCheckpoint(ctx sdk.Context, veID uint64, amount sdk.Int) {
	now := uint64(ctx.BlockTime().Unix())
	epoch := b.GetUserEpoch(ctx, veID)
	point := b.GetUserCheckpoint(ctx, veID, epoch)

	point.Amount = amount
	if epoch > 0 && point.Timestamp == now {
		// nothing
	} else {
		epoch += 1
		b.SetUserEpoch(ctx, veID, epoch)
		point.Timestamp = now
	}
	b.SetUserCheckpoint(ctx, veID, epoch, point)
}

func (b *Base) writeRewardPerTicketCheckpoint(ctx sdk.Context, rewardDenom string, reward sdk.Int, timestamp uint64) {
	epoch := b.GetRewardEpoch(ctx, rewardDenom)
	point := b.GetRewardCheckpoint(ctx, rewardDenom, epoch)

	point.Amount = reward
	if epoch > 0 && point.Timestamp == timestamp {
		// nothing
	} else {
		epoch += 1
		b.SetRewardEpoch(ctx, rewardDenom, epoch)
		point.Timestamp = timestamp
	}
	b.SetRewardCheckpoint(ctx, rewardDenom, epoch, point)
}

func (b *Base) writeCheckpoint(ctx sdk.Context) {
	now := uint64(ctx.BlockTime().Unix())
	epoch := b.GetEpoch(ctx)
	point := b.GetCheckpoint(ctx, epoch)

	if b.isGauge {
		point.Amount = b.GetTotalDerivedAmount(ctx)
	} else {
		point.Amount = b.GetTotalDepositedAmount(ctx)
	}

	if epoch > 0 && point.Timestamp == now {
		// nothing
	} else {
		epoch += 1
		b.SetEpoch(ctx, epoch)
		point.Timestamp = now
	}

	b.SetCheckpoint(ctx, epoch, point)
}
