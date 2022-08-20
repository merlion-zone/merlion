package keeper_test

import "github.com/merlion-zone/merlion/x/vesting/types"

func (suite *KeeperTestSuite) TestKeeper_GetParams() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	params := k.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
}

func (suite *KeeperTestSuite) TestKeeper_SetParams() {
	suite.SetupTest()
	k := suite.app.VestingKeeper
	p := types.DefaultParams()
	k.SetParams(suite.ctx, p)

	params := k.GetParams(suite.ctx)
	suite.Require().Equal(p, params)
}
