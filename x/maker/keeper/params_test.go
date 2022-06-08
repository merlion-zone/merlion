package keeper_test

import "github.com/merlion-zone/merlion/x/maker/types"

func (suite *KeeperTestSuite) TestParams() {
	makerKeeper := suite.app.MakerKeeper
	params := makerKeeper.GetParams(suite.ctx)

	suite.Require().Equal(types.DefaultParams(), params)
	suite.Require().Equal(types.DefaultBackingRatioStep, makerKeeper.BackingRatioStep(suite.ctx))
	suite.Require().Equal(types.DefaultBackingRatioPriceBand, makerKeeper.BackingRatioPriceBand(suite.ctx))
	suite.Require().Equal(types.DefaultBackingRatioCooldownPeriod, makerKeeper.BackingRatioCooldownPeriod(suite.ctx))
	suite.Require().Equal(types.DefaultMintPriceBias, makerKeeper.MintPriceBias(suite.ctx))
	suite.Require().Equal(types.DefaultBurnPriceBias, makerKeeper.BurnPriceBias(suite.ctx))
	suite.Require().Equal(types.DefaultRebackBonus, makerKeeper.RebackBonus(suite.ctx))
	suite.Require().Equal(types.DefaultLiquidationCommissionFee, makerKeeper.LiquidationCommissionFee(suite.ctx))

	makerKeeper.SetParams(suite.ctx, params)
	newParams := makerKeeper.GetParams(suite.ctx)
	suite.Require().Equal(params, newParams)
}
