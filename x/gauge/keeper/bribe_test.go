package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
)

func (suite *KeeperTestSuite) TestBribe_Deposit() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	k.SetGauge(suite.ctx, depoistDenom)
	bribe := k.Bribe(suite.ctx, depoistDenom)
	veID := uint64(100)
	amount := sdk.NewInt(100)
	bribe.Deposit(suite.ctx, veID, amount)

	suite.Require().Equal(amount, bribe.GetTotalDepositedAmount(suite.ctx))
	suite.Require().Equal(amount, bribe.GetDepositedAmountByUser(suite.ctx, veID))
}

func (suite *KeeperTestSuite) TestBribe_Withdraw() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	k.SetGauge(suite.ctx, depoistDenom)
	bribe := k.Bribe(suite.ctx, depoistDenom)
	veID := uint64(100)
	amount := sdk.NewInt(100)
	err := bribe.Withdraw(suite.ctx, veID, amount)
	suite.Require().Error(err, types.ErrTooLargeAmount)

	bribe.Deposit(suite.ctx, veID, amount)
	err = bribe.Withdraw(suite.ctx, veID, sdk.NewInt(50))
	suite.Require().Nil(err)

	suite.Require().Equal(sdk.NewInt(50), bribe.GetTotalDepositedAmount(suite.ctx))
	suite.Require().Equal(sdk.NewInt(50), bribe.GetDepositedAmountByUser(suite.ctx, veID))

	err = bribe.Withdraw(suite.ctx, veID, sdk.NewInt(50))
	suite.Require().Nil(err)

	suite.Require().Equal(sdk.ZeroInt(), bribe.GetTotalDepositedAmount(suite.ctx))
	suite.Require().Equal(sdk.ZeroInt(), bribe.GetDepositedAmountByUser(suite.ctx, veID))
}
