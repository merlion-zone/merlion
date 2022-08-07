package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetSlopeChange_GetSlopeChange() {
	suite.SetupTest()
	timestamp := uint64(1754379718)
	change := suite.app.VeKeeper.GetSlopeChange(suite.ctx, timestamp)
	suite.Require().Equal(sdk.ZeroInt(), change)

	suite.app.VeKeeper.SetSlopeChange(suite.ctx, timestamp, sdk.NewInt(10000))
	change = suite.app.VeKeeper.GetSlopeChange(suite.ctx, timestamp)
	suite.Require().Equal(sdk.NewInt(10000), change)
}

func (suite *KeeperTestSuite) TestKeeper_SetUserCheckpoint_GetUserCheckpoint() {
	suite.SetupTest()
	veID := uint64(10000)
	epoch := uint64(100)
	checkpoint := suite.app.VeKeeper.GetUserCheckpoint(suite.ctx, veID, epoch)
	suite.Require().Equal(sdk.ZeroInt(), checkpoint.Bias)
	suite.Require().Equal(sdk.ZeroInt(), checkpoint.Slope)

	point := types.Checkpoint{
		Bias:  sdk.NewInt(1000),
		Slope: sdk.NewInt(1000),
	}
	suite.app.VeKeeper.SetUserCheckpoint(suite.ctx, veID, epoch, point)
	checkpoint = suite.app.VeKeeper.GetUserCheckpoint(suite.ctx, veID, epoch)
	suite.Require().Equal(point.Bias, checkpoint.Bias)
	suite.Require().Equal(point.Slope, checkpoint.Slope)
}

func (suite *KeeperTestSuite) TestKeeper_SetUserEpoch_GetUserEpoch() {
	suite.SetupTest()
	veID := uint64(10000)
	epoch := suite.app.VeKeeper.GetUserEpoch(suite.ctx, veID)
	suite.Require().Equal(uint64(0), epoch)

	suite.app.VeKeeper.SetUserEpoch(suite.ctx, veID, uint64(1000))
	epoch = suite.app.VeKeeper.GetUserEpoch(suite.ctx, veID)
	suite.Require().Equal(uint64(1000), epoch)
}

func (suite *KeeperTestSuite) TestKeeper_SetCheckpoint_GetCheckpoint() {
	suite.SetupTest()
	epoch := uint64(10000)
	point := suite.app.VeKeeper.GetCheckpoint(suite.ctx, epoch)
	suite.Require().Equal(sdk.ZeroInt(), point.Bias)
	suite.Require().Equal(sdk.ZeroInt(), point.Slope)

	checkpoint := types.Checkpoint{
		Bias:  sdk.NewInt(1000),
		Slope: sdk.NewInt(1000),
	}
	suite.app.VeKeeper.SetCheckpoint(suite.ctx, epoch, checkpoint)
	point = suite.app.VeKeeper.GetCheckpoint(suite.ctx, epoch)
	suite.Require().Equal(checkpoint.Bias, point.Bias)
	suite.Require().Equal(checkpoint.Slope, point.Slope)
}

func (suite *KeeperTestSuite) TestKeeper_GetEpoch_SetEpoch() {
	suite.SetupTest()
	epoch := suite.app.VeKeeper.GetEpoch(suite.ctx)
	suite.Require().Equal(uint64(0), epoch)

	suite.app.VeKeeper.SetEpoch(suite.ctx, uint64(1000))
	epoch = suite.app.VeKeeper.GetEpoch(suite.ctx)
	suite.Require().Equal(uint64(1000), epoch)
}

func (suite *KeeperTestSuite) TestKeeper_RegulateCheckpoint() {
	suite.SetupTest()
	suite.app.VeKeeper.RegulateCheckpoint(suite.ctx)
	epoch := suite.app.VeKeeper.GetEpoch(suite.ctx)
	point := suite.app.VeKeeper.GetCheckpoint(suite.ctx, epoch)

	suite.Require().Equal(uint64(1), epoch)
	suite.Require().Equal(sdk.ZeroInt(), point.Bias)
	suite.Require().Equal(sdk.ZeroInt(), point.Slope)
}

func (suite *KeeperTestSuite) TestKeeper_RegulateUserCheckpoint() {
	suite.SetupTest()
	veID := uint64(100)
	lockedOld := types.LockedBalance{Amount: sdk.ZeroInt(), End: uint64(2537136000)}
	lockedNew := types.LockedBalance{Amount: sdk.NewInt(1000000), End: uint64(2543184000)}
	suite.app.VeKeeper.RegulateUserCheckpoint(suite.ctx, veID, lockedOld, lockedNew)

	epoch := suite.app.VeKeeper.GetEpoch(suite.ctx)
	point := suite.app.VeKeeper.GetCheckpoint(suite.ctx, epoch)

	suite.Require().Equal(uint64(1), epoch)
	suite.Require().Equal(sdk.ZeroInt(), point.Bias)
	suite.Require().Equal(sdk.ZeroInt(), point.Slope)

	userEpoch := suite.app.VeKeeper.GetUserEpoch(suite.ctx, veID)
	userPoint := suite.app.VeKeeper.GetUserCheckpoint(suite.ctx, veID, epoch)
	suite.Require().Equal(uint64(1), userEpoch)
	suite.Require().Equal(sdk.ZeroInt(), userPoint.Bias)
	suite.Require().Equal(sdk.ZeroInt(), userPoint.Slope)
}
