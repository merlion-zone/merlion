package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	erc20types "github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestKeeper_IsDenomForErc20() {
	suite.SetupTest()
	t := suite.T()
	k := suite.app.Erc20Keeper
	for _, tc := range []struct {
		desc  string
		denom string
		valid bool
	}{
		{
			desc:  "Denom is for erc20",
			denom: "erc20/0xdAC17F958D2ee523a2206206994597C13D831ec7",
			valid: true,
		},
		{
			desc:  "Denom is not for erc20",
			denom: "address",
			valid: false,
		},
		{
			desc:  "Denom is not for erc20",
			denom: "aaa/0xdAC17F958D2ee523a2206206994597C13D831ec7",
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			valid := k.IsDenomForErc20(tc.denom)
			require.Equal(t, tc.valid, valid)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetSupply() {
	suite.SetupTest()
	t := suite.T()
	k := suite.app.Erc20Keeper

	suite.app.BankKeeper.SetDenomMetaData(suite.ctx, suite.coinMetadata)
	uusdSupply := sdk.NewCoin("uusd", sdk.NewInt(10000000))
	totalSupply := sdk.NewCoins(uusdSupply)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, erc20types.ModuleName, totalSupply)
	require.NoError(t, err)

	for _, tc := range []struct {
		desc   string
		denom  string
		supply sdk.Coin
	}{
		{
			desc:   "denom invalid",
			denom:  "erc20/0xdAC17F958D2ee523a2206206994597C13D831ec7",
			supply: sdk.Coin{},
		},
		{
			desc:   "denom valid",
			denom:  uusdSupply.Denom,
			supply: uusdSupply,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			supply := k.GetSupply(suite.ctx, tc.denom)
			require.Equal(t, tc.supply, supply)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_GetBalance() {
	suite.SetupTest()
	t := suite.T()
	k := suite.app.Erc20Keeper

	suite.app.BankKeeper.SetDenomMetaData(suite.ctx, suite.coinMetadata)
	uusdSupply := sdk.NewCoin("uusd", sdk.NewInt(10000000))
	totalSupply := sdk.NewCoins(uusdSupply)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, erc20types.ModuleName, totalSupply)
	require.NoError(t, err)
	acc := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, erc20types.ModuleName)
	for _, tc := range []struct {
		desc   string
		denom  string
		supply sdk.Coin
	}{
		{
			desc:   "denom invalid",
			denom:  "erc20/0xdAC17F958D2ee523a2206206994597C13D831ec7",
			supply: sdk.Coin{},
		},
		{
			desc:   "denom valid",
			denom:  uusdSupply.Denom,
			supply: uusdSupply,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			balance := k.GetBalance(suite.ctx, acc.GetAddress(), tc.denom)
			require.Equal(t, tc.supply, balance)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_SplitCoinsByErc20() {
	suite.SetupTest()
	var (
		t           = suite.T()
		k           = suite.app.Erc20Keeper
		native      = sdk.NewCoin("uusd", sdk.NewInt(10000000))
		nativeErc20 = sdk.NewCoin("erc20/0xdAC17F958D2ee523a2206206994597C13D831ec7", sdk.NewInt(10000000))
		coins       = sdk.NewCoins(native, nativeErc20)
	)
	nativeCoins, nativeErc20Tokens := k.SplitCoinsByErc20(coins)
	require.Equal(t, sdk.NewCoins(native), nativeCoins)
	require.Equal(t, sdk.NewCoins(nativeErc20), nativeErc20Tokens)
}

func (suite *KeeperTestSuite) TestKeeper_SendCoins() {
	suite.SetupTest()
	t := suite.T()
	k := suite.app.Erc20Keeper

	suite.app.BankKeeper.SetDenomMetaData(suite.ctx, suite.coinMetadata)
	uusdSupply := sdk.NewCoin("uusd", sdk.NewInt(10000000))
	totalSupply := sdk.NewCoins(uusdSupply)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, erc20types.ModuleName, totalSupply)
	require.NoError(t, err)
	from := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, erc20types.ModuleName).GetAddress()
	to := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, authtypes.FeeCollectorName).GetAddress()
	amt := sdk.NewCoin("uusd", sdk.NewInt(100000))
	k.SendCoins(suite.ctx, from, to, sdk.Coins{amt}, sdk.Coins{})

	fromBalance := k.GetBalance(suite.ctx, from, "uusd")
	toBalance := k.GetBalance(suite.ctx, to, "uusd")
	require.Equal(t, uusdSupply.Sub(amt), fromBalance)
	require.Equal(t, amt, toBalance)
}
