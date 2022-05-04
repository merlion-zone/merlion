package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
)

type Bribe struct {
	Base
}

func (k Keeper) Bribe(ctx sdk.Context, depoistDenom string) Bribe {
	if !k.HasGauge(ctx, depoistDenom) {
		panic("gauge not found")
	}
	return Bribe{
		Base: Base{
			prefixKey:    types.BribeKey(depoistDenom),
			depoistDenom: depoistDenom,
			keeper:       k,
		},
	}
}

func (b Bribe) ClaimReward(ctx sdk.Context, veID uint64) (err error) {
	return b.Base.claimReward(ctx, veID)
}

func (b Bribe) Deposit(ctx sdk.Context, veID uint64, amount sdk.Int) {
	totalDeposited := b.Base.GetTotalDepositedAmount(ctx)
	deposited := b.Base.GetDepositedAmountByUser(ctx, veID)
	totalDeposited = totalDeposited.Add(amount)
	deposited = deposited.Add(amount)
	b.Base.SetTotalDepositedAmount(ctx, totalDeposited)
	b.Base.SetDepositedAmountByUser(ctx, veID, deposited)

	b.writeUserCheckpoint(ctx, veID, deposited)
	b.writeCheckpoint(ctx)
}

func (b Bribe) Withdraw(ctx sdk.Context, veID uint64, amount sdk.Int) (err error) {
	totalDeposited := b.Base.GetTotalDepositedAmount(ctx)
	deposited := b.Base.GetDepositedAmountByUser(ctx, veID)
	if amount.GT(deposited) {
		// TODO: error
	}

	totalDeposited = totalDeposited.Sub(amount)
	deposited = deposited.Sub(amount)
	b.Base.SetTotalDepositedAmount(ctx, totalDeposited)
	if deposited.IsPositive() {
		b.Base.SetDepositedAmountByUser(ctx, veID, deposited)
	} else {
		b.Base.DeleteDepositedAmountByUser(ctx, veID)
	}

	b.writeUserCheckpoint(ctx, veID, deposited)
	b.writeCheckpoint(ctx)

	return nil
}
