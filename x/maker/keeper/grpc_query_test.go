package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

func (suite *KeeperTestSuite) TestAllBackingRiskParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all backing risk params is empty
	res, err := suite.queryClient.AllBackingRiskParams(ctx, &types.QueryAllBackingRiskParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllBackingRiskParamsResponse{}, res)

	brp, brp2 := suite.dummyBackingRiskParams()
	suite.Require().NotEqual(brp.BackingDenom, brp2.BackingDenom)
	// add backing risk params
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, brp)
	res, err = suite.queryClient.AllBackingRiskParams(ctx, &types.QueryAllBackingRiskParamsRequest{})
	expRes := &types.QueryAllBackingRiskParamsResponse{
		RiskParams: []types.BackingRiskParams{brp},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// add another risk params
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, brp2)
	res, err = suite.queryClient.AllBackingRiskParams(ctx, &types.QueryAllBackingRiskParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(res.RiskParams, 2)
	suite.Require().Contains(res.RiskParams, brp)
	suite.Require().Contains(res.RiskParams, brp2)
}

func (suite *KeeperTestSuite) TestAllCollateralRiskParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all collateral risk params is empty
	res, err := suite.queryClient.AllCollateralRiskParams(ctx, &types.QueryAllCollateralRiskParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllCollateralRiskParamsResponse{}, res)

	crp, crp2 := suite.dummyCollateralRiskParams()
	suite.Require().NotEqual(crp.CollateralDenom, crp2.CollateralDenom)
	// add collateral risk params
	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, crp)
	res, err = suite.queryClient.AllCollateralRiskParams(ctx, &types.QueryAllCollateralRiskParamsRequest{})
	expRes := &types.QueryAllCollateralRiskParamsResponse{
		RiskParams: []types.CollateralRiskParams{crp},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// add another collateral risk params
	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, crp2)
	res, err = suite.queryClient.AllCollateralRiskParams(ctx, &types.QueryAllCollateralRiskParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(res.RiskParams, 2)
	suite.Require().Contains(res.RiskParams, crp)
	suite.Require().Contains(res.RiskParams, crp2)
}

func (suite *KeeperTestSuite) TestAllBackingPools() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all backing pools is empty
	res, err := suite.queryClient.AllBackingPools(ctx, &types.QueryAllBackingPoolsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllBackingPoolsResponse{}, res)

	// add backing pool
	poolBacking := types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(100)),
		Backing:    sdk.NewCoin(suite.bcDenom, sdk.NewInt(10)),
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
		MerMinted:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(200)),
		Backing:    sdk.NewCoin("eth", sdk.NewInt(20)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, poolBacking2)
	res, err = suite.queryClient.AllBackingPools(ctx, &types.QueryAllBackingPoolsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(res.BackingPools, 2)
	suite.Require().Contains(res.BackingPools, poolBacking)
	suite.Require().Contains(res.BackingPools, poolBacking2)
}

func (suite *KeeperTestSuite) TestAllCollateralPools() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default all collateral pools is empty
	res, err := suite.queryClient.AllCollateralPools(ctx, &types.QueryAllCollateralPoolsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(&types.QueryAllCollateralPoolsResponse{}, res)

	// add collateral pool
	poolColl := types.PoolCollateral{
		Collateral: sdk.NewCoin(suite.bcDenom, sdk.NewInt(10)),
		MerDebt:    sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(100)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(200)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	}
	suite.app.MakerKeeper.SetPoolCollateral(suite.ctx, poolColl)
	res, err = suite.queryClient.AllCollateralPools(ctx, &types.QueryAllCollateralPoolsRequest{})
	expRes := &types.QueryAllCollateralPoolsResponse{
		CollateralPools: []types.PoolCollateral{poolColl},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// add another collateral pool
	poolColl2 := types.PoolCollateral{
		Collateral: sdk.NewCoin("eth", sdk.NewInt(10)),
		MerDebt:    sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(200)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(300)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolCollateral(suite.ctx, poolColl2)
	res, err = suite.queryClient.AllCollateralPools(ctx, &types.QueryAllCollateralPoolsRequest{})
	suite.Require().NoError(err)
	suite.Require().Len(res.CollateralPools, 2)
	suite.Require().Contains(res.CollateralPools, poolColl)
	suite.Require().Contains(res.CollateralPools, poolColl2)
}

func (suite *KeeperTestSuite) TestBackingPool() {
	poolBacking := types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(200)),
		Backing:    sdk.NewCoin(suite.bcDenom, sdk.NewInt(20)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, poolBacking)

	// backing pool not found
	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.queryClient.BackingPool(ctx, &types.QueryBackingPoolRequest{BackingDenom: "eth"})
	suite.Require().Error(err)

	// correct
	res, err := suite.queryClient.BackingPool(ctx, &types.QueryBackingPoolRequest{BackingDenom: suite.bcDenom})
	expRes := &types.QueryBackingPoolResponse{BackingPool: poolBacking}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestCollateralPool() {
	poolCollateral := types.PoolCollateral{
		Collateral: sdk.NewCoin(suite.bcDenom, sdk.NewInt(10)),
		MerDebt:    sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(200)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(300)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	}
	suite.app.MakerKeeper.SetPoolCollateral(suite.ctx, poolCollateral)

	// collateral pool not found
	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err := suite.queryClient.CollateralPool(ctx, &types.QueryCollateralPoolRequest{CollateralDenom: "eth"})
	suite.Require().Error(err)

	// correct
	res, err := suite.queryClient.CollateralPool(ctx, &types.QueryCollateralPoolRequest{CollateralDenom: suite.bcDenom})
	expRes := &types.QueryCollateralPoolResponse{CollateralPool: poolCollateral}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestCollateralOfAccount() {
	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	accAddress := sdk.AccAddress(priv.PubKey().Address())

	priv2, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err)
	accAddress2 := sdk.AccAddress(priv2.PubKey().Address())

	accColl := types.AccountCollateral{
		Account:             accAddress.String(),
		Collateral:          sdk.NewCoin(suite.bcDenom, sdk.NewInt(100)),
		MerDebt:             sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(200)),
		MerByLion:           sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(50)),
		LionBurned:          sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
		LastInterest:        sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(10)),
		LastSettlementBlock: 666,
	}
	suite.app.MakerKeeper.SetAccountCollateral(suite.ctx, accAddress, accColl)

	// account not found
	ctx := sdk.WrapSDKContext(suite.ctx)
	_, err = suite.queryClient.CollateralOfAccount(ctx, &types.QueryCollateralOfAccountRequest{
		Account:         accAddress2.String(),
		CollateralDenom: suite.bcDenom,
	})
	suite.Require().Error(err)

	// collateral denom not found
	_, err = suite.queryClient.CollateralOfAccount(ctx, &types.QueryCollateralOfAccountRequest{
		Account:         accAddress.String(),
		CollateralDenom: "eth",
	})
	suite.Require().Error(err)

	// correct
	res, err := suite.queryClient.CollateralOfAccount(ctx, &types.QueryCollateralOfAccountRequest{
		Account:         accAddress.String(),
		CollateralDenom: suite.bcDenom,
	})
	expRes := &types.QueryCollateralOfAccountResponse{AccountCollateral: accColl}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestTotalBacking() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default total backing is all zero
	res, err := suite.queryClient.TotalBacking(ctx, &types.QueryTotalBackingRequest{})
	expRes := &types.QueryTotalBackingResponse{
		TotalBacking: types.TotalBacking{
			BackingValue: sdk.ZeroInt(),
			MerMinted:    sdk.Coin{"", sdk.ZeroInt()},
			LionBurned:   sdk.Coin{"", sdk.ZeroInt()},
		},
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// set prices, total backing & backings pools
	suite.setupEstimationTest()
	totalBacking := types.TotalBacking{
		BackingValue: sdk.OneInt(),
		MerMinted:    sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(100)),
		LionBurned:   sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	}
	suite.app.MakerKeeper.SetTotalBacking(suite.ctx, totalBacking)

	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(100)),
		Backing:    sdk.NewCoin(suite.bcDenom, sdk.NewInt(10_000000)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	})

	suite.app.MakerKeeper.SetPoolBacking(suite.ctx, types.PoolBacking{
		MerMinted:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(200)),
		Backing:    sdk.NewCoin("eth", sdk.NewInt(2)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(2000)),
	})

	totalBacking.BackingValue = sdk.NewInt(2009_900000) // 10_000000 * 0.99 + 2 * 1000_000000

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
		MerDebt:    sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(100)),
		MerByLion:  sdk.NewCoin(merlion.MicroUSMDenom, sdk.NewInt(50)),
		LionBurned: sdk.NewCoin(merlion.AttoLionDenom, sdk.NewInt(1000)),
	}
	suite.app.MakerKeeper.SetTotalCollateral(suite.ctx, totalCollateral)
	res, err = suite.queryClient.TotalCollateral(ctx, &types.QueryTotalCollateralRequest{})
	expRes = &types.QueryTotalCollateralResponse{TotalCollateral: totalCollateral}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestBackingRatio() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	// default backing ratio is 1, set at block height 0
	res, err := suite.queryClient.BackingRatio(ctx, &types.QueryBackingRatioRequest{})
	expRes := &types.QueryBackingRatioResponse{
		BackingRatio:    sdk.OneDec(),
		LastUpdateBlock: 0,
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)

	// set backing ratio manually
	newBR := sdk.NewDecWithPrec(80, 2)
	suite.app.MakerKeeper.SetBackingRatio(suite.ctx, newBR)
	suite.app.MakerKeeper.SetBackingRatioLastBlock(suite.ctx, 888)
	res, err = suite.queryClient.BackingRatio(ctx, &types.QueryBackingRatioRequest{})
	expRes = &types.QueryBackingRatioResponse{
		BackingRatio:    newBR,
		LastUpdateBlock: 888,
	}
	suite.Require().NoError(err)
	suite.Require().Equal(expRes, res)
}

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := sdk.WrapSDKContext(suite.ctx)
	res, err := suite.queryClient.Params(ctx, &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(types.DefaultParams(), res.Params)
}
