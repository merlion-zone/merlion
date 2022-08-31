package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/merlion-zone/merlion/x/gauge/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

type FeeClaimee interface {
	ClaimFees(ctx sdk.Context, claimant sdk.AccAddress) sdk.Coins
}

type Gauge struct {
	Base
	bribe Bribe
}

func (k Keeper) CreateGauge(ctx sdk.Context, depoistDenom string) {
	// TODO: whitelists
	if k.HasGauge(ctx, depoistDenom) {
		panic("gauge exists")
	}
	k.SetGauge(ctx, depoistDenom)
}

func (k Keeper) Gauge(ctx sdk.Context, depoistDenom string) Gauge {
	if !k.HasGauge(ctx, depoistDenom) {
		panic("gauge not found")
	}
	return Gauge{
		Base: Base{
			keeper:       k,
			depoistDenom: depoistDenom,
			prefixKey:    types.GaugeKey(depoistDenom),
			isGauge:      true,
		},
	}
}

func (g Gauge) ClaimReward(ctx sdk.Context, veID uint64, voterKeeper types.VoterKeeper) (err error) {
	voterKeeper.DistributeReward(ctx, g.PoolDenom())

	return g.claimReward(ctx, veID)
}

func (g Gauge) Deposit(ctx sdk.Context, veID uint64, amount sdk.Int) (err error) {
	owner := g.keeper.nftKeeper.GetOwner(ctx, vetypes.VeNftClass.Id, vetypes.VeIDFromUint64(veID))
	coin := sdk.NewCoin(g.depoistDenom, amount)
	err = g.keeper.bankKeeper.SendCoins(ctx, owner, g.EscrowPool(ctx).GetAddress(), sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	totalDeposited := g.GetTotalDepositedAmount(ctx)
	deposited := g.GetDepositedAmountByUser(ctx, veID)
	totalDeposited = totalDeposited.Add(amount)
	deposited = deposited.Add(amount)
	g.SetTotalDepositedAmount(ctx, totalDeposited)
	g.SetDepositedAmountByUser(ctx, veID, deposited)

	// if first-time deposit
	if g.GetUserVeIDByAddress(ctx, owner) == vetypes.EmptyVeID {
		g.SetUserVeIDByAddress(ctx, owner, veID)
		g.keeper.veKeeper.IncVeAttached(ctx, veID)
	}

	g.deriveAmountForUser(ctx, veID)

	// TODO: voter emit deposit

	return nil
}

func (g Gauge) Withdraw(ctx sdk.Context, veID uint64, amount sdk.Int) (err error) {
	owner := g.keeper.nftKeeper.GetOwner(ctx, vetypes.VeNftClass.Id, vetypes.VeIDFromUint64(veID))

	totalDeposited := g.GetTotalDepositedAmount(ctx)
	deposited := g.GetDepositedAmountByUser(ctx, veID)
	if amount.GT(deposited) {
		return types.ErrTooLargeAmount
	}
	totalDeposited = totalDeposited.Sub(amount)
	deposited = deposited.Sub(amount)
	g.SetTotalDepositedAmount(ctx, totalDeposited)
	if deposited.IsPositive() {
		g.SetDepositedAmountByUser(ctx, veID, deposited)
	} else {
		g.DeleteDepositedAmountByUser(ctx, veID)
		g.DeleteUserVeIDByAddress(ctx, owner)
		g.keeper.veKeeper.DecVeAttached(ctx, veID)
	}

	coin := sdk.NewCoin(g.depoistDenom, amount)
	err = g.keeper.bankKeeper.SendCoins(ctx, g.EscrowPool(ctx).GetAddress(), owner, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	g.deriveAmountForUser(ctx, veID)

	// TODO: voter emit withdraw

	return nil
}

func (g Gauge) DepositReward(ctx sdk.Context, sender sdk.AccAddress, rewardDenom string, amount sdk.Int) error {
	return g.depositReward(ctx, sender, rewardDenom, amount)
}

func (g Gauge) DepositFees(ctx sdk.Context, feeClaimee FeeClaimee) (err error) {
	claimed := feeClaimee.ClaimFees(ctx, g.EscrowPool(ctx).GetAddress())

	acc := g.EscrowPool(ctx).GetAddress()
	for _, fee := range claimed {
		reward := g.GetReward(ctx, fee.Denom)
		feeAmount := reward.AccruedAmount.Add(fee.Amount)

		if feeAmount.GT(g.bribe.RemainingReward(ctx, fee.Denom)) && feeAmount.QuoRaw(vetypes.RegulatedPeriod).IsPositive() {
			err = g.bribe.depositReward(ctx, acc, fee.Denom, feeAmount)
			if err != nil {
				return err
			}
			// clear accrued undistributed amount
			reward.AccruedAmount = sdk.ZeroInt()
		} else {
			// accrue undistributed amount
			reward.AccruedAmount = feeAmount
		}

		g.SetReward(ctx, fee.Denom, reward)
	}
	return nil
}
