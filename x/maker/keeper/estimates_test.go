package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (suite *KeeperTestSuite) TestEstimateMintBySwapIn() {
	testCases := []struct {
		name     string
		malleate func()
		req      *types.EstimateMintBySwapInRequest
		expPass  bool
		expErr   error
		expRes   *types.EstimateMintBySwapInResponse
	}{
		{
			name: "mer price too low",
			malleate: func() {
				suite.app.OracleKeeper.SetExchangeRate(suite.ctx, merlion.MicroUSDDenom, sdk.NewDecWithPrec(989, 3))
			},
			req:     &types.EstimateMintBySwapInRequest{BackingDenom: suite.bcDenom},
			expPass: false,
			expErr:  types.ErrMerPriceTooLow,
		},
		{
			name:    "backing denom not found",
			req:     &types.EstimateMintBySwapInRequest{BackingDenom: "fil"},
			expPass: false,
			expErr:  types.ErrBackingCoinNotFound,
		},
		{
			name:    "backing denom disabled",
			req:     &types.EstimateMintBySwapInRequest{BackingDenom: "eth"},
			expPass: false,
			expErr:  types.ErrBackingCoinDisabled,
		},
		{
			name: "mer over ceiling",
			req: &types.EstimateMintBySwapInRequest{
				MintOut:      sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(2_000000)),
				BackingDenom: suite.bcDenom,
			},
			expPass: false,
			expErr:  types.ErrMerCeiling,
		},
		{
			name: "default full backing",
			req: &types.EstimateMintBySwapInRequest{
				MintOut:      sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1_000000)),
				BackingDenom: suite.bcDenom,
			},
			expPass: true,
			expRes: &types.EstimateMintBySwapInResponse{
				BackingIn: sdk.NewCoin(suite.bcDenom, sdk.NewInt(1015152)), // 1_000000 * (1+0.005) / 0.99
				LionIn:    sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt()),
				MintFee:   sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(5000)),
			},
		},
		{
			name: "user asked full backing",
			malleate: func() {
				suite.app.MakerKeeper.SetBackingRatio(suite.ctx, sdk.NewDecWithPrec(80, 2))
			},
			req: &types.EstimateMintBySwapInRequest{
				MintOut:      sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1_000000)),
				BackingDenom: suite.bcDenom,
				FullBacking:  true,
			},
			expPass: true,
			expRes: &types.EstimateMintBySwapInResponse{
				BackingIn: sdk.NewCoin(suite.bcDenom, sdk.NewInt(1015152)), // 1_000000 * (1+0.005) / 0.99
				LionIn:    sdk.NewCoin(merlion.AttoLionDenom, sdk.ZeroInt()),
				MintFee:   sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(5000)),
			},
		},
		{
			name: "full algorithmic",
			malleate: func() {
				suite.app.MakerKeeper.SetBackingRatio(suite.ctx, sdk.ZeroDec())
			},
			req: &types.EstimateMintBySwapInRequest{
				MintOut:      sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1_000000)),
				BackingDenom: suite.bcDenom,
			},
			expPass: true,
			expRes: &types.EstimateMintBySwapInResponse{
				BackingIn: sdk.NewCoin(suite.bcDenom, sdk.ZeroInt()),
				LionIn:    sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(10050000000000000)), // 1_000000 * (1+0.005) / (100*10^-12)
				MintFee:   sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(5000)),
			},
		},
		{
			name: "backing over ceiling",
			req: &types.EstimateMintBySwapInRequest{
				MintOut:      sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1_500000)),
				BackingDenom: suite.bcDenom,
			},
			expPass: false,
			expErr:  types.ErrBackingCeiling,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			suite.SetupTest() // reset
			suite.setupEstimation()
			if tc.malleate != nil {
				tc.malleate()
			}

			ctx := sdk.WrapSDKContext(suite.ctx)
			res, err := suite.queryClient.EstimateMintBySwapIn(ctx, tc.req)
			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expRes, res)
			} else {
				suite.Require().Error(err)
				suite.Require().ErrorIs(err, tc.expErr)
			}

		})
	}
}

func (suite *KeeperTestSuite) setupEstimation() {
	// set prices
	suite.app.OracleKeeper.SetExchangeRate(suite.ctx, suite.bcDenom, sdk.NewDecWithPrec(99, 2))
	suite.app.OracleKeeper.SetExchangeRate(suite.ctx, "eth", sdk.NewDec(1000_000000))
	suite.app.OracleKeeper.SetExchangeRate(suite.ctx, "fil", sdk.NewDec(5_000000))
	suite.app.OracleKeeper.SetExchangeRate(suite.ctx, merlion.AttoLionDenom, sdk.NewDecWithPrec(100, 12))
	suite.app.OracleKeeper.SetExchangeRate(suite.ctx, merlion.MicroUSDDenom, sdk.NewDecWithPrec(101, 2))

	// set risk params
	brp, brp2 := suite.dummyBackingRiskParams()
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, brp)
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, brp2)

	crp, crp2 := suite.dummyCollateralRiskParams()
	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, crp)
	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, crp2)

	// set pool and total backing/collateral
	lionBurned, ok := sdk.NewIntFromString("100_000000_000000_000000")
	suite.Require().True(ok)

	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(8_000000)),
		Backing:    sdk.NewCoin(suite.bcDenom, sdk.NewInt(9_000000)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, lionBurned),
	})
	suite.app.MakerKeeper.SetTotalBacking(suite.ctx, types.TotalBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1)),
	})
	suite.app.MakerKeeper.SetPoolCollateral(suite.ctx, types.PoolCollateral{
		Collateral: sdk.NewCoin(suite.bcDenom, sdk.NewInt(1)),
		MerDebt:    sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1)),
	})
	suite.app.MakerKeeper.SetTotalCollateral(suite.ctx, types.TotalCollateral{
		MerDebt:    sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(1)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1)),
	})
}
