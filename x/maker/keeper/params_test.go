package keeper_test

import "github.com/merlion-zone/merlion/x/maker/types"

func (suite *KeeperTestSuite) TestParams() {
	makerKeeper := suite.app.MakerKeeper
	params := makerKeeper.GetParams(suite.ctx)

	suite.Require().Equal(params, types.DefaultParams())
	suite.Require().Equal(makerKeeper.BackingRatioStep(suite.ctx), types.DefaultBackingRatioStep)
	suite.Require().Equal(makerKeeper.BackingRatioPriceBand(suite.ctx), types.DefaultBackingRatioPriceBand)
	suite.Require().Equal(makerKeeper.BackingRatioCooldownPeriod(suite.ctx), types.DefaultBackingRatioCooldownPeriod)
	suite.Require().Equal(makerKeeper.MintPriceBias(suite.ctx), types.DefaultMintPriceBias)
	suite.Require().Equal(makerKeeper.BurnPriceBias(suite.ctx), types.DefaultBurnPriceBias)
	suite.Require().Equal(makerKeeper.RecollateralizeBonus(suite.ctx), types.DefaultRecollateralizeBonus)
	suite.Require().Equal(makerKeeper.LiquidationCommissionFee(suite.ctx), types.DefaultLiquidationCommissionFee)

	makerKeeper.SetParams(suite.ctx, params)
	newParams := makerKeeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
