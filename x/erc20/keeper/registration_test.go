package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
	erc20types "github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestKeeper_RegisterCoin() {
	suite.SetupTest()
	t := suite.T()

	_, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, merlion.DisplayDenom)
	require.Error(t, err, sdkerrors.Wrapf(erc20types.ErrEVMDenom, "cannot register the EVM denomination %s", merlion.DisplayDenom))

	_, err = suite.app.Erc20Keeper.RegisterCoin(suite.ctx, "USDT")
	require.Error(t, err, sdkerrors.Wrapf(erc20types.ErrEVMDenom, "cannot get metadata of denom %s", "USDT"))

	suite.app.BankKeeper.SetDenomMetaData(suite.ctx, suite.coinMetadata)
	_, err = suite.app.Erc20Keeper.RegisterCoin(suite.ctx, "uusd")
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "base denomination 'uusd' cannot have a supply of 0"))

	totalSupply := sdk.NewCoins(sdk.NewCoin("uusd", sdk.NewInt(10000000)))
	suite.app.BankKeeper.MintCoins(suite.ctx, erc20types.ModuleName, totalSupply)

	_, err = suite.app.Erc20Keeper.RegisterCoin(suite.ctx, "uusd")
	require.Error(t, err, sdkerrors.Wrapf(erc20types.ErrTokenPairAlreadyExists, "coin denomination already registered: %s", suite.coinMetadata.Base))

	addr, err := suite.app.Erc20Keeper.DeployERC20Contract(suite.ctx, suite.coinMetadata)
	require.NoError(t, err)
	require.Equal(t, "0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", addr.String())

	tokenPair, err := suite.app.Erc20Keeper.RegisterERC20(suite.ctx, addr)
	require.NoError(t, err)

	require.Equal(t, "0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", tokenPair.Erc20Address)
	require.Equal(t, "erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", tokenPair.Denom)
	require.Equal(t, erc20types.Owner(2), tokenPair.ContractOwner)

	// QueryERC20
	erc20Data, err := suite.app.Erc20Keeper.QueryERC20(suite.ctx, addr)
	require.NoError(t, err)
	require.Equal(t, "USDT", erc20Data.Name)
	require.Equal(t, "USDT", erc20Data.Symbol)
	require.Equal(t, uint8(0), erc20Data.Decimals)
}
