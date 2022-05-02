package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

func (k Keeper) DepositFees(ctx sdk.Context, stakedDenom string, amount sdk.Coins) {
}

type Gauge struct {
	Base
}

func (k Keeper) GetGauge(depoistDenom string) Gauge {
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
	err = g.Base.keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, owner, g.Base.GetEscrowPool(ctx).GetName(), sdk.NewCoins(coin))
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
	err = g.Base.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, g.Base.GetEscrowPool(ctx).GetName(), owner, sdk.NewCoins(coin))
	if err != nil {
		return err
	}

	g.Base.deriveAmountForUser(ctx, veID)

	// TODO: voter emit withdraw

	return nil
}
