package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetTotalLockedAmount_GetTotalLockedAmount() {
	suite.SetupTest()
	amt := suite.app.VeKeeper.GetTotalLockedAmount(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), amt)

	suite.app.VeKeeper.SetTotalLockedAmount(suite.ctx, sdk.NewInt(10000))
	amt = suite.app.VeKeeper.GetTotalLockedAmount(suite.ctx)
	suite.Require().Equal(sdk.NewInt(10000), amt)
}

func (suite *KeeperTestSuite) TestKeeper_SetLockedAmountByUser_GetLockedAmountByUser() {
	suite.SetupTest()
	veID := uint64(10000)
	amt := suite.app.VeKeeper.GetLockedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(types.NewLockedBalance(), amt)

	locked := types.LockedBalance{
		Amount: sdk.NewInt(1000),
		End:    uint64(1754379718),
	}
	suite.app.VeKeeper.SetLockedAmountByUser(suite.ctx, uint64(10000), locked)
	amt = suite.app.VeKeeper.GetLockedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(locked, amt)
}

func (suite *KeeperTestSuite) TestKeeper_DeleteLockedAmountByUser() {
	suite.SetupTest()
	veID := uint64(10000)
	locked := types.LockedBalance{
		Amount: sdk.NewInt(1000),
		End:    uint64(1754379718),
	}
	suite.app.VeKeeper.SetLockedAmountByUser(suite.ctx, uint64(10000), locked)
	suite.app.VeKeeper.DeleteLockedAmountByUser(suite.ctx, veID)
	amt := suite.app.VeKeeper.GetLockedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(types.NewLockedBalance(), amt)
}
