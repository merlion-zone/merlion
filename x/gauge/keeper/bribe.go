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
		Base: NewBase(k, depoistDenom, types.BribeKey(depoistDenom), false),
	}
}

func (b Bribe) ClaimReward(ctx sdk.Context, veID uint64) (err error) {
	return b.claimReward(ctx, veID)
}

func (b Bribe) Deposit(ctx sdk.Context, veID uint64, amount sdk.Int) {
	totalDeposited := b.GetTotalDepositedAmount(ctx)
	deposited := b.GetDepositedAmountByUser(ctx, veID)
	totalDeposited = totalDeposited.Add(amount)
	deposited = deposited.Add(amount)
	b.SetTotalDepositedAmount(ctx, totalDeposited)
	b.SetDepositedAmountByUser(ctx, veID, deposited)

	b.writeUserCheckpoint(ctx, veID, deposited)
	b.writeCheckpoint(ctx)
}

func (b Bribe) Withdraw(ctx sdk.Context, veID uint64, amount sdk.Int) (err error) {
	totalDeposited := b.GetTotalDepositedAmount(ctx)
	deposited := b.GetDepositedAmountByUser(ctx, veID)
	if amount.GT(deposited) {
		return types.ErrTooLargeAmount
	}

	totalDeposited = totalDeposited.Sub(amount)
	deposited = deposited.Sub(amount)
	b.SetTotalDepositedAmount(ctx, totalDeposited)
	if deposited.IsPositive() {
		b.SetDepositedAmountByUser(ctx, veID, deposited)
	} else {
		b.DeleteDepositedAmountByUser(ctx, veID)
	}

	b.writeUserCheckpoint(ctx, veID, deposited)
	b.writeCheckpoint(ctx)

	return nil
}
