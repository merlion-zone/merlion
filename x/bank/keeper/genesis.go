package keeper

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	merlion "github.com/merlion-zone/merlion/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState *banktypes.GenesisState) {
	// 1. Init balances, set denom metadata, and set supply
	k.BaseKeeper.InitGenesis(ctx, genState)

	// 2. Set denom metadata for predefined stable coins
	merlion.SetDenomMetaDataForStableCoins(ctx, k)

	// 3. Register erc20 for coins, and set erc20 balance
	for _, balance := range genState.Balances {
		addr, err := sdk.AccAddressFromBech32(balance.Address)
		if err != nil {
			panic(err)
		}

		for _, coin := range balance.Coins {
			if strings.Contains(coin.Denom, merlion.DisplayDenom) {
				// skip gas token
				continue
			}

			// register erc20 for coins
			if !k.erc20Keeper().IsDenomRegistered(ctx, coin.Denom) {
				_, err := k.erc20Keeper().RegisterCoin(ctx, coin.Denom)
				if err != nil {
					panic(fmt.Sprintf("register native coin to erc20 module: %s", err))
				}
			}

			// balance alignment
			erc20Balance := k.erc20Keeper().GetBalance(ctx, addr, coin.Denom)
			if erc20Balance.IsLT(coin) {
				// mint missing
				k.erc20Keeper().SendCoins(ctx, nil, addr, sdk.NewCoins(coin.Sub(erc20Balance)), nil)
			} else if coin.IsLT(erc20Balance) {
				// burn excess
				k.erc20Keeper().SendCoins(ctx, addr, nil, sdk.NewCoins(erc20Balance.Sub(coin)), nil)
			}

			// check balance again
			erc20Balance = k.erc20Keeper().GetBalance(ctx, addr, coin.Denom)
			if !erc20Balance.IsEqual(coin) {
				panic(fmt.Sprintf("unequal balance of coin and erc20 representations: %s, %s", coin, erc20Balance))
			}
		}
	}
}
