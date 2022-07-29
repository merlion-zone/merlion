package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestKeeper_HasBalance() {
	suite.SetupTest()
	var (
		t          = suite.T()
		k          = suite.app.BankKeeper
		denom      = types.AttoLionDenom
		amt        = sdk.NewCoin(denom, sdk.NewInt(100))
		erc20Denom = "erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75"
	)
	// Raw balance check
	bal0 := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), bal0)

	has := k.HasBalance(suite.ctx, suite.addrs[0], amt)
	require.Equal(t, true, has)

	has = k.HasBalance(suite.ctx, suite.addrs[0], sdk.NewCoin(denom, sdk.NewInt(1135)))
	require.Equal(t, false, has)

	has = k.HasBalance(suite.ctx, suite.addrs[0], sdk.NewCoin(erc20Denom, sdk.NewInt(100)))
	require.Equal(t, false, has)
}

func (suite *KeeperTestSuite) TestKeeper_GetAllBalances() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = types.AttoLionDenom
	)
	bal := k.GetAllBalances(suite.ctx, suite.addrs[0])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1134))), bal)
}

func (suite *KeeperTestSuite) TestKeeper_GetBalance() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = types.AttoLionDenom
	)
	// Raw balance check
	bal := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), bal)
}

func (suite *KeeperTestSuite) TestKeeper_SpendableCoins() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = types.AttoLionDenom
	)
	bal := k.SpendableCoins(suite.ctx, suite.addrs[0])
	require.Equal(t, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1134))), bal)
}

func (suite *KeeperTestSuite) TestKeeper_ValidateBalance() {
	suite.SetupTest()
	var (
		t = suite.T()
		k = suite.app.BankKeeper
	)
	err := k.ValidateBalance(suite.ctx, suite.addrs[0])
	require.NoError(t, err)
}
