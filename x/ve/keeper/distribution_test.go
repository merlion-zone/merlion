package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) TestDistributor_DistributePerPeriod() {
	// TODO
}

func (suite *KeeperTestSuite) TestDistributor_Claim() {
	// TODO
}

func (suite *KeeperTestSuite) TestKeeper_SetDistributionAccruedLastTimestamp_GetDistributionAccruedLastTimestamp() {
	suite.SetupTest()
	timestamp := suite.app.VeKeeper.GetDistributionAccruedLastTimestamp(suite.ctx)
	suite.Require().Equal(uint64(0), timestamp)

	timestamp = uint64(1754379718)
	suite.app.VeKeeper.SetDistributionAccruedLastTimestamp(suite.ctx, timestamp)
	stmp := suite.app.VeKeeper.GetDistributionAccruedLastTimestamp(suite.ctx)
	suite.Require().Equal(timestamp, stmp)
}

func (suite *KeeperTestSuite) TestKeeper_SetDistributionTotalAmount_GetDistributionTotalAmount() {
	suite.SetupTest()
	amt := suite.app.VeKeeper.GetDistributionTotalAmount(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), amt)

	total := sdk.NewInt(10000)
	suite.app.VeKeeper.SetDistributionTotalAmount(suite.ctx, total)
	amt = suite.app.VeKeeper.GetDistributionTotalAmount(suite.ctx)
	suite.Require().Equal(total, amt)
}

func (suite *KeeperTestSuite) TestKeeper_SetDistributionPerPeriod_GetDistributionPerPeriod() {
	suite.SetupTest()
	timestamp := uint64(1754379718)
	period := suite.app.VeKeeper.GetDistributionPerPeriod(suite.ctx, timestamp)
	suite.Require().Equal(sdk.ZeroInt(), period)

	amt := sdk.NewInt(10000)
	suite.app.VeKeeper.SetDistributionPerPeriod(suite.ctx, timestamp, amt)
	period = suite.app.VeKeeper.GetDistributionPerPeriod(suite.ctx, timestamp)
	suite.Require().Equal(amt, period)
}

func (suite *KeeperTestSuite) TestKeeper_SetDistributionClaimLastTimestampByUser_GetDistributionClaimLastTimestampByUser() {
	suite.SetupTest()
	timestamp := uint64(1754379718)
	veID := uint64(100)
	stmp := suite.app.VeKeeper.GetDistributionClaimLastTimestampByUser(suite.ctx, veID)
	suite.Require().Equal(uint64(0), stmp)

	suite.app.VeKeeper.SetDistributionClaimLastTimestampByUser(suite.ctx, veID, timestamp)
	stmp = suite.app.VeKeeper.GetDistributionClaimLastTimestampByUser(suite.ctx, veID)
	suite.Require().Equal(timestamp, stmp)
}
