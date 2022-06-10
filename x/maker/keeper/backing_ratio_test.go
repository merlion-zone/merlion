package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (suite *KeeperTestSuite) TestAdjustBackingRatio() {
	// default backing ratio is 1, set at block height 0, now height is 1
	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.queryClient.BackingRatio(ctx, &types.QueryBackingRatioRequest{})
	orgRes := &types.QueryBackingRatioResponse{
		BackingRatio:    sdk.OneDec(),
		LastUpdateBlock: 0,
	}
	suite.Require().NoError(err)
	suite.Require().Equal(orgRes, res)
	suite.Require().Equal(int64(1), suite.ctx.BlockHeight())

	suite.shortenBackingRatioCooldownPeriod()
	shortCooldownPeriod := suite.app.MakerKeeper.BackingRatioCooldownPeriod(suite.ctx)

	testCases := []struct {
		name     string
		malleate func()
		expPanic bool
		expRes   *types.QueryBackingRatioResponse
	}{
		{
			name: "within cool down period",
			malleate: func() {
				for i := int64(0); i < shortCooldownPeriod-2; i++ {
					suite.Commit()
				}
				suite.Require().Equal(shortCooldownPeriod-1, suite.ctx.BlockHeight())
			},
			expPanic: false,
			expRes:   orgRes,
		},
		{
			name: "mer price not set",
			malleate: func() {
				for i := int64(0); i < shortCooldownPeriod-1; i++ {
					suite.Commit()
				}
				suite.Require().Equal(shortCooldownPeriod, suite.ctx.BlockHeight())
			},
			expPanic: true,
			expRes:   nil,
		},
		{
			name: "mer price too high",
			malleate: func() {
				for i := int64(0); i < shortCooldownPeriod-1; i++ {
					suite.Commit()
				}
				suite.Require().Equal(shortCooldownPeriod, suite.ctx.BlockHeight())
				suite.app.OracleKeeper.SetExchangeRate(suite.ctx, merlion.MicroUSDDenom, sdk.NewDecWithPrec(101, 2))
			},
			expPanic: false,
			expRes: &types.QueryBackingRatioResponse{
				BackingRatio:    sdk.NewDecWithPrec(9975, 4),
				LastUpdateBlock: shortCooldownPeriod,
			},
		},
		{
			name: "min br is zero",
			malleate: func() {
				for i := int64(0); i < shortCooldownPeriod-1; i++ {
					suite.Commit()
				}
				suite.Require().Equal(shortCooldownPeriod, suite.ctx.BlockHeight())
				suite.app.OracleKeeper.SetExchangeRate(suite.ctx, merlion.MicroUSDDenom, sdk.NewDecWithPrec(101, 2))
				suite.app.MakerKeeper.SetBackingRatio(suite.ctx, types.DefaultBackingRatioStep.Sub(sdk.NewDecWithPrec(1, 4)))
			},
			expPanic: false,
			expRes: &types.QueryBackingRatioResponse{
				BackingRatio:    sdk.ZeroDec(),
				LastUpdateBlock: shortCooldownPeriod,
			},
		},
		{
			name: "max br is one",
			malleate: func() {
				for i := int64(0); i < shortCooldownPeriod-1; i++ {
					suite.Commit()
				}
				suite.Require().Equal(shortCooldownPeriod, suite.ctx.BlockHeight())
				suite.app.OracleKeeper.SetExchangeRate(suite.ctx, merlion.MicroUSDDenom, sdk.NewDecWithPrec(99, 2))
			},
			expPanic: false,
			expRes: &types.QueryBackingRatioResponse{
				BackingRatio:    sdk.OneDec(),
				LastUpdateBlock: shortCooldownPeriod,
			},
		},
		{
			name: "mer price too low",
			malleate: func() {
				for i := int64(0); i < shortCooldownPeriod-1; i++ {
					suite.Commit()
				}
				suite.app.MakerKeeper.SetBackingRatio(suite.ctx, sdk.NewDecWithPrec(9, 1))
				suite.Require().Equal(shortCooldownPeriod, suite.ctx.BlockHeight())
				suite.app.OracleKeeper.SetExchangeRate(suite.ctx, merlion.MicroUSDDenom, sdk.NewDecWithPrec(99, 2))
			},
			expPanic: false,
			expRes: &types.QueryBackingRatioResponse{
				BackingRatio:    sdk.NewDecWithPrec(9025, 4),
				LastUpdateBlock: shortCooldownPeriod,
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			suite.shortenBackingRatioCooldownPeriod()

			tc.malleate()

			if tc.expPanic {
				suite.Require().Panics(func() {
					suite.app.MakerKeeper.AdjustBackingRatio(suite.ctx)
				})
			} else {
				suite.app.MakerKeeper.AdjustBackingRatio(suite.ctx)
				suite.Commit()
				res, err := suite.queryClient.BackingRatio(ctx, &types.QueryBackingRatioRequest{})
				suite.Require().NoError(err, tc.name)
				suite.Require().Equal(tc.expRes, res, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) shortenBackingRatioCooldownPeriod() {
	// make cool down period shorter to speed up testing
	params := suite.app.MakerKeeper.GetParams(suite.ctx)
	params.BackingRatioCooldownPeriod = 10
	suite.app.MakerKeeper.SetParams(suite.ctx, params)
}
