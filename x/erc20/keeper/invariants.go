package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/erc20/types"
)

// RegisterInvariants registers the erc20 module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "consistent-balance", ConsistentBalanceInvariant(k))
}

// ConsistentBalanceInvariant checks that all accounts have consistent balances in bank and erc20
func ConsistentBalanceInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg   string
			count int
		)

		// Only iterate on native coins
		k.bankKeeper.IterateAllBalances(ctx, func(addr sdk.AccAddress, balance sdk.Coin) bool {
			if strings.Contains(balance.Denom, merlion.DisplayDenom) {
				// skip gas token
				return false
			}

			erc20Balance := k.GetBalance(ctx, addr, balance.Denom)
			if !erc20Balance.IsEqual(balance) {
				count++
				msg += fmt.Sprintf("\t%s has unequal balances: %s != %s\n", addr, balance, erc20Balance)
			}

			return false
		})

		broken := count != 0

		return sdk.FormatInvariant(
			types.ModuleName, "consistent-balance",
			fmt.Sprintf("inconsistent balances found %d\n%s", count, msg),
		), broken
	}
}
