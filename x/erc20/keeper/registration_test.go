package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	testkeeper "github.com/merlion-zone/merlion/testutil/keeper"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
)

const faucetAccountName = "faucet"

func TestKeeper_DeployERC20Contract(t *testing.T) {
	var (
		k, bankKeeper, ctx = testkeeper.Erc20KeeperWithExtra(t)
		coinMetadata       = banktypes.Metadata{
			Description: "USDT",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "u" + "usd", Exponent: uint32(0), Aliases: []string{"micro" + "usd"}}, // e.g., uusd
				{Denom: "m" + "usd", Exponent: uint32(3), Aliases: []string{"milli" + "usd"}}, // e.g., musd
				{Denom: "usd", Exponent: uint32(6), Aliases: []string{}},                      // e.g., usd
			},
			Base:    "uusd",
			Display: "usd",
			Name:    "USDT",
			Symbol:  "USDT",
		}
		err error
	)

	_, err = k.RegisterCoin(ctx, merlion.DisplayDenom)
	require.Error(t, err, sdkerrors.Wrapf(types.ErrEVMDenom, "cannot register the EVM denomination %s", merlion.DisplayDenom))

	_, err = k.RegisterCoin(ctx, "USDT")
	require.Error(t, err, sdkerrors.Wrapf(types.ErrEVMDenom, "cannot get metadata of denom %s", "USDT"))

	bankKeeper.SetDenomMetaData(ctx, coinMetadata)
	_, err = k.RegisterCoin(ctx, "uusd")
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "base denomination 'uusd' cannot have a supply of 0"))

	totalSupply := sdk.NewCoins(sdk.NewCoin("uusd", sdk.NewInt(10000000)))
	bankKeeper.MintCoins(ctx, faucetAccountName, totalSupply)

	_, err = k.RegisterCoin(ctx, "uusd")
	require.NoError(t, err)
}

// func TestKeeper_RegisterERC20(t *testing.T) {
// 	var (
// 		k, ctx = testkeeper.Erc20Keeper(t)
// 		addr   = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
// 		err    error
// 	)

// 	_, err = k.RegisterERC20(ctx, addr)
// 	require.NoError(t, err)
// }
