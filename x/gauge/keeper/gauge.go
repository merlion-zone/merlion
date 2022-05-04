package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

type FeeClaimee interface {
	ClaimFees(ctx sdk.Context, claimer sdk.AccAddress) sdk.Coins
}

type Gauge struct {
	Base
	bribe      Bribe
	feeClaimee FeeClaimee
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
			prefixKey:    types.GaugeKey(depoistDenom),
			depoistDenom: depoistDenom,
			keeper:       k,
		},
	}
}

func (g Gauge) ClaimReward(ctx sdk.Context, veID uint64) (err error) {
	// TODO: voter distribute for gauge

	return g.Base.claimReward(ctx, veID)
}

func (g Gauge) Deposit(ctx sdk.Context, veID uint64, amount sdk.Int) (err error) {
	owner := g.Base.keeper.nftKeeper.GetOwner(ctx, vetypes.VeNftClass.Id, vetypes.VeID(veID))
	coin := sdk.NewCoin(g.Base.depoistDenom, amount)
	err = g.Base.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, g.Base.EscrowPool(ctx).GetName(), sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	totalDeposited := g.Base.GetTotalDepositedAmount(ctx)
	deposited := g.Base.GetDepositedAmountByUser(ctx, veID)
	totalDeposited = totalDeposited.Add(amount)
	deposited = deposited.Add(amount)
	g.Base.SetTotalDepositedAmount(ctx, totalDeposited)
	g.Base.SetDepositedAmountByUser(ctx, veID, deposited)

	// if first-time deposit
	if g.Base.GetUserVeIDByAddress(ctx, owner) == vetypes.EmptyVeID {
		g.Base.SetUserVeIDByAddress(ctx, owner, veID)
		// TODO: voter attach
	}

	g.Base.deriveAmountForUser(ctx, veID)

	// TODO: voter emit deposit

	return nil
}

func (g Gauge) Withdraw(ctx sdk.Context, veID uint64, amount sdk.Int) (err error) {
	owner := g.Base.keeper.nftKeeper.GetOwner(ctx, vetypes.VeNftClass.Id, vetypes.VeID(veID))

	totalDeposited := g.Base.GetTotalDepositedAmount(ctx)
	deposited := g.Base.GetDepositedAmountByUser(ctx, veID)
	if amount.GT(deposited) {
		// TODO: error
	}
	totalDeposited = totalDeposited.Sub(amount)
	deposited = deposited.Sub(amount)
	g.Base.SetTotalDepositedAmount(ctx, totalDeposited)
	if deposited.IsPositive() {
		g.Base.SetDepositedAmountByUser(ctx, veID, deposited)
	} else {
		g.Base.DeleteDepositedAmountByUser(ctx, veID)
		g.Base.DeleteUserVeIDByAddress(ctx, owner)
		// TODO: voter detach
	}

	coin := sdk.NewCoin(g.Base.depoistDenom, amount)
	err = g.Base.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, g.Base.EscrowPool(ctx).GetName(), owner, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	g.Base.deriveAmountForUser(ctx, veID)

	// TODO: voter emit withdraw

	return nil
}

func (g Gauge) DepositReward(ctx sdk.Context, sender sdk.AccAddress, rewardDenom string, amount sdk.Int) error {
	err := g.DepositFees(ctx)
	if err != nil {
		return err
	}

	return g.Base.depositReward(ctx, sender, rewardDenom, amount)
}

func (g Gauge) DepositFees(ctx sdk.Context) (err error) {
	claimed := g.feeClaimee.ClaimFees(ctx, g.Base.EscrowPool(ctx).GetAddress())

	acc := g.Base.EscrowPool(ctx).GetAddress()
	for _, fee := range claimed {
		reward := g.Base.GetReward(ctx, fee.Denom)
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

		g.Base.SetReward(ctx, fee.Denom, reward)
	}
	return nil
}
