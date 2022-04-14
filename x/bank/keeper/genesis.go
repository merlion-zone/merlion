package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState *banktypes.GenesisState) {
	k.BaseKeeper.InitGenesis(ctx, genState)

	for _, balance := range genState.Balances {
		addr, err := sdk.AccAddressFromBech32(balance.Address)
		if err != nil {
			panic(err)
		}

		// Check erc20 balance
		for _, coin := range balance.Coins {
			erc20Balance := k.erc20Keeper().GetBalance(ctx, addr, coin.Denom)
			if !erc20Balance.IsEqual(coin) {
				panic(fmt.Sprintf("unequal balance of coin and erc20 representations; %s, %s", coin, erc20Balance))
			}
		}
	}
}
