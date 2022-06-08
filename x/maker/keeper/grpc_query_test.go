package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (suite *KeeperTestSuite) TestAllBackingRiskParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all backing risk params is empty
	res, err := suite.queryClient.AllBackingRiskParams(ctx, &types.QueryAllBackingRiskParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllBackingRiskParamsResponse{}, res)

	// add backing risk params
	maxBacking := sdk.NewInt(100)
	maxMerMint := sdk.NewInt(10)
	mintFee := sdk.NewDecWithPrec(5, 3)
	burnFee := sdk.NewDecWithPrec(6, 3)
	buybackFee := sdk.NewDecWithPrec(7, 3)
	rebackFee := sdk.NewDecWithPrec(8, 3)
	backingRiskParams := types.BackingRiskParams{
		BackingDenom: "btc",
		Enabled:      false,
		MaxBacking:   &maxBacking,
		MaxMerMint:   &maxMerMint,
		MintFee:      &mintFee,
		BurnFee:      &burnFee,
		BuybackFee:   &buybackFee,
		RebackFee:    &rebackFee,
	}
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, backingRiskParams)
	res, err = suite.queryClient.AllBackingRiskParams(ctx, &types.QueryAllBackingRiskParamsRequest{})
	expRes := &types.QueryAllBackingRiskParamsResponse{
		RiskParams: []types.BackingRiskParams{backingRiskParams},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// add another risk params
	maxBacking2 := sdk.NewInt(200)
	maxMerMint2 := sdk.NewInt(20)
	mintFee2 := sdk.NewDecWithPrec(6, 3)
	burnFee2 := sdk.NewDecWithPrec(7, 3)
	buybackFee2 := sdk.NewDecWithPrec(8, 3)
	rebackFee2 := sdk.NewDecWithPrec(9, 3)
	backingRiskParams2 := types.BackingRiskParams{
		BackingDenom: "eth",
		Enabled:      true,
		MaxBacking:   &maxBacking2,
		MaxMerMint:   &maxMerMint2,
		MintFee:      &mintFee2,
		BurnFee:      &burnFee2,
		BuybackFee:   &buybackFee2,
		RebackFee:    &rebackFee2,
	}
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, backingRiskParams2)
	res, err = suite.queryClient.AllBackingRiskParams(ctx, &types.QueryAllBackingRiskParamsRequest{})
	expRes = &types.QueryAllBackingRiskParamsResponse{
		RiskParams: []types.BackingRiskParams{backingRiskParams, backingRiskParams2},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestAllCollateralRiskParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all collateral risk params is empty
	res, err := suite.queryClient.AllCollateralRiskParams(ctx, &types.QueryAllCollateralRiskParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllCollateralRiskParamsResponse{}, res)

	// add collateral risk params
	maxCollateral := sdk.NewInt(100)
	maxMerMint := sdk.NewInt(10)
	liquidationThreshold := sdk.NewDecWithPrec(75, 2)
	loanToValue := sdk.NewDecWithPrec(70, 2)
	basicLoanToValue := sdk.NewDecWithPrec(50, 2)
	catalyticLionRation := sdk.NewDecWithPrec(5, 2)
	liquidationFee := sdk.NewDecWithPrec(10, 2)
	mintFee := sdk.NewDecWithPrec(1, 2)
	InterestFee := sdk.NewDecWithPrec(3, 2)
	collateralRiskParams := types.CollateralRiskParams{
		CollateralDenom:      "btc",
		Enabled:              true,
		MaxCollateral:        &maxCollateral,
		MaxMerMint:           &maxMerMint,
		LiquidationThreshold: &liquidationThreshold,
		LoanToValue:          &loanToValue,
		BasicLoanToValue:     &basicLoanToValue,
		CatalyticLionRatio:   &catalyticLionRation,
		LiquidationFee:       &liquidationFee,
		MintFee:              &mintFee,
		InterestFee:          &InterestFee,
	}
	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, collateralRiskParams)
	res, err = suite.queryClient.AllCollateralRiskParams(ctx, &types.QueryAllCollateralRiskParamsRequest{})
	expRes := &types.QueryAllCollateralRiskParamsResponse{
		RiskParams: []types.CollateralRiskParams{collateralRiskParams},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// add another collateral risk params
	maxCollateral2 := sdk.NewInt(200)
	maxMerMint2 := sdk.NewInt(20)
	liquidationThreshold2 := sdk.NewDecWithPrec(76, 2)
	loanToValue2 := sdk.NewDecWithPrec(71, 2)
	basicLoanToValue2 := sdk.NewDecWithPrec(51, 2)
	catalyticLionRation2 := sdk.NewDecWithPrec(6, 2)
	liquidationFee2 := sdk.NewDecWithPrec(11, 2)
	mintFee2 := sdk.NewDecWithPrec(2, 2)
	InterestFee2 := sdk.NewDecWithPrec(4, 2)
	collateralRiskParams2 := types.CollateralRiskParams{
		CollateralDenom:      "eth",
		Enabled:              false,
		MaxCollateral:        &maxCollateral2,
		MaxMerMint:           &maxMerMint2,
		LiquidationThreshold: &liquidationThreshold2,
		LoanToValue:          &loanToValue2,
		BasicLoanToValue:     &basicLoanToValue2,
		CatalyticLionRatio:   &catalyticLionRation2,
		LiquidationFee:       &liquidationFee2,
		MintFee:              &mintFee2,
		InterestFee:          &InterestFee2,
	}

	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, collateralRiskParams2)
	res, err = suite.queryClient.AllCollateralRiskParams(ctx, &types.QueryAllCollateralRiskParamsRequest{})
	expRes = &types.QueryAllCollateralRiskParamsResponse{
		RiskParams: []types.CollateralRiskParams{collateralRiskParams, collateralRiskParams2},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestAllBackingPools() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all backing pools is empty
	res, err := suite.queryClient.AllBackingPools(ctx, &types.QueryAllBackingPoolsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllBackingPoolsResponse{}, res)

	// add backing pool
	poolBacking := types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(100)),
		Backing:    sdk.NewCoin("btc", sdk.NewInt(10)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	}
	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, poolBacking)
	res, err = suite.queryClient.AllBackingPools(ctx, &types.QueryAllBackingPoolsRequest{})
	expRes := &types.QueryAllBackingPoolsResponse{
		BackingPools: []types.PoolBacking{poolBacking},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// add another backing pool
	poolBacking2 := types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(200)),
		Backing:    sdk.NewCoin("eth", sdk.NewInt(20)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, poolBacking2)
	res, err = suite.queryClient.AllBackingPools(ctx, &types.QueryAllBackingPoolsRequest{})
	expRes = &types.QueryAllBackingPoolsResponse{
		BackingPools: []types.PoolBacking{poolBacking, poolBacking2},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestAllCollateralPools() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all collateral pools is empty
	res, err := suite.queryClient.AllCollateralPools(ctx, &types.QueryAllCollateralPoolsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllCollateralPoolsResponse{}, res)

	// add collateral pool
	poolCollateral := types.PoolCollateral{
		Collateral: sdk.NewCoin("btc", sdk.NewInt(10)),
		MerDebt:    sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(100)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(200)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	}
	suite.app.MakerKeeper.SetPoolCollateral(suite.ctx, poolCollateral)
	res, err = suite.queryClient.AllCollateralPools(ctx, &types.QueryAllCollateralPoolsRequest{})
	expRes := &types.QueryAllCollateralPoolsResponse{
		CollateralPools: []types.PoolCollateral{poolCollateral},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// add another collateral pool
	poolCollateral2 := types.PoolCollateral{
		Collateral: sdk.NewCoin("eth", sdk.NewInt(10)),
		MerDebt:    sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(200)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(300)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolCollateral(suite.ctx, poolCollateral2)
	res, err = suite.queryClient.AllCollateralPools(ctx, &types.QueryAllCollateralPoolsRequest{})
	expRes = &types.QueryAllCollateralPoolsResponse{
		CollateralPools: []types.PoolCollateral{poolCollateral, poolCollateral2},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestBackingPool() {
	poolBacking := types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(200)),
		Backing:    sdk.NewCoin("eth", sdk.NewInt(20)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, poolBacking)

	// backing pool not found
	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.queryClient.BackingPool(ctx, &types.QueryBackingPoolRequest{BackingDenom: "btc"})
	suite.Require().Error(err)

	// backing pool found
	res, err := suite.queryClient.BackingPool(ctx, &types.QueryBackingPoolRequest{BackingDenom: "eth"})
	expRes := &types.QueryBackingPoolResponse{BackingPool: poolBacking}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestCollateralPool() {
	poolCollateral := types.PoolCollateral{
		Collateral: sdk.NewCoin("eth", sdk.NewInt(10)),
		MerDebt:    sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(200)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(300)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolCollateral(suite.ctx, poolCollateral)

	// collateral pool not found
	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.queryClient.CollateralPool(ctx, &types.QueryCollateralPoolRequest{CollateralDenom: "btc"})
	suite.Require().Error(err)

	// collateral pool found
	res, err := suite.queryClient.CollateralPool(ctx, &types.QueryCollateralPoolRequest{CollateralDenom: "eth"})
	expRes := &types.QueryCollateralPoolResponse{CollateralPool: poolCollateral}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestCollateralOfAccount() {

}

func (suite *KeeperTestSuite) TestTotalBacking() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default total backing is all zero
	res, err := suite.queryClient.TotalBacking(ctx, &types.QueryTotalBackingRequest{})
	expRes := &types.QueryTotalBackingResponse{
		TotalBacking: types.TotalBacking{
			MerMinted:  sdk.Coin{"", sdk.ZeroInt()},
			LionBurned: sdk.Coin{"", sdk.ZeroInt()},
		},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// set total backing
	totalBacking := types.TotalBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(100)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	}
	suite.app.MakerKeeper.SetTotalBacking(suite.ctx, totalBacking)
	res, err = suite.queryClient.TotalBacking(ctx, &types.QueryTotalBackingRequest{})
	expRes = &types.QueryTotalBackingResponse{TotalBacking: totalBacking}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestTotalCollateral() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default total collateral is all zero
	res, err := suite.queryClient.TotalCollateral(ctx, &types.QueryTotalCollateralRequest{})
	expRes := &types.QueryTotalCollateralResponse{
		TotalCollateral: types.TotalCollateral{
			MerDebt:    sdk.Coin{"", sdk.ZeroInt()},
			MerByLion:  sdk.Coin{"", sdk.ZeroInt()},
			LionBurned: sdk.Coin{"", sdk.ZeroInt()},
		},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// set total collateral
	totalCollateral := types.TotalCollateral{
		MerDebt:    sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(100)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.NewInt(50)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	}
	suite.app.MakerKeeper.SetTotalCollateral(suite.ctx, totalCollateral)
	res, err = suite.queryClient.TotalCollateral(ctx, &types.QueryTotalCollateralRequest{})
	expRes = &types.QueryTotalCollateralResponse{TotalCollateral: totalCollateral}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestBackingRatio() {

}

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.queryClient.Params(ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(types.DefaultParams(), res.Params)
}
