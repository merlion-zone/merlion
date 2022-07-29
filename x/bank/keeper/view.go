package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k Keeper) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	if !k.erc20Keeper().IsDenomForErc20(amt.Denom) {
		return k.BaseKeeper.HasBalance(ctx, addr, amt)
	} else {
		bal := k.erc20Keeper().GetBalance(ctx, addr, amt.Denom)
		return !bal.Amount.IsNil() && bal.IsGTE(amt)
	}
}

func (k Keeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	// NOTE: not including balances of all erc20 native tokens
	return k.BaseKeeper.GetAllBalances(ctx, addr)
}

func (k Keeper) GetAccountsBalances(ctx sdk.Context) []banktypes.Balance {
	// NOTE: not including balances of all erc20 native tokens
	return k.BaseKeeper.GetAccountsBalances(ctx)
}

func (k Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	if !k.erc20Keeper().IsDenomForErc20(denom) {
		return k.BaseKeeper.GetBalance(ctx, addr, denom)
	} else {
		return k.erc20Keeper().GetBalance(ctx, addr, denom)
	}
}

func (k Keeper) IterateAccountBalances(ctx sdk.Context, addr sdk.AccAddress, cb func(sdk.Coin) bool) {
	// NOTE: not including balances of all erc20 native tokens
	k.BaseKeeper.IterateAccountBalances(ctx, addr, cb)
}

func (k Keeper) IterateAllBalances(ctx sdk.Context, cb func(sdk.AccAddress, sdk.Coin) bool) {
	// NOTE: not including balances of all erc20 native tokens
	k.BaseKeeper.IterateAllBalances(ctx, cb)
}

func (k Keeper) SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	// NOTE: not including balances of all erc20 native tokens
	return k.BaseKeeper.SpendableCoins(ctx, addr)
}

func (k Keeper) ValidateBalance(ctx sdk.Context, addr sdk.AccAddress) error {
	// NOTE: not including validating balances of all erc20 native tokens
	return k.BaseKeeper.ValidateBalance(ctx, addr)
}
