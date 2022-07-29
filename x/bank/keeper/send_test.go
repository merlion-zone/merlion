package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/merlion-zone/merlion/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestKeeper_InputOutputCoins() {
	suite.SetupTest()
	var (
		t      = suite.T()
		k      = suite.app.BankKeeper
		denom  = types.AttoLionDenom
		amt    = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100)))
		amtOut = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(200)))
		inputs = []banktypes.Input{
			{Address: suite.addrs[0].String(), Coins: amt},
			{Address: suite.addrs[1].String(), Coins: amt},
		}
		outputs = []banktypes.Output{
			{Address: suite.addrs[2].String(), Coins: amtOut},
		}
	)
	// Raw balance check
	bal0 := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	bal1 := k.GetBalance(suite.ctx, suite.addrs[1], denom)
	bal2 := k.GetBalance(suite.ctx, suite.addrs[2], denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), bal0)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1234)), bal1)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1234)), bal2)

	// Send coins
	err := k.InputOutputCoins(suite.ctx, inputs, outputs)
	require.NoError(t, err)

	// Balance check
	bal0 = k.GetBalance(suite.ctx, suite.addrs[0], denom)
	bal1 = k.GetBalance(suite.ctx, suite.addrs[1], denom)
	bal2 = k.GetBalance(suite.ctx, suite.addrs[2], denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1034)), bal0)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), bal1)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1434)), bal2)
}

func (suite *KeeperTestSuite) TestKeeper_SendCoins() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = types.AttoLionDenom
		amt   = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100)))
	)
	// Raw balance check
	bal0 := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	bal1 := k.GetBalance(suite.ctx, suite.addrs[1], denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), bal0)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1234)), bal1)

	// Send coins
	err := k.SendCoins(suite.ctx, suite.addrs[0], suite.addrs[1], amt)
	require.NoError(t, err)

	// Balance check
	bal0 = k.GetBalance(suite.ctx, suite.addrs[0], denom)
	bal1 = k.GetBalance(suite.ctx, suite.addrs[1], denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1034)), bal0)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1334)), bal1)
}
