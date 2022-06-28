package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/maker/types"
)

func (suite *KeeperTestSuite) TestSetGetBackingRiskParams() {
	brp, brp2 := suite.dummyBackingRiskParams()
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, brp)

	// backing denom not found
	suite.Require().NotEqual(brp.BackingDenom, brp2.BackingDenom)
	isBackingRegistered := suite.app.MakerKeeper.IsBackingRegistered(suite.ctx, brp2.BackingDenom)
	suite.Require().False(isBackingRegistered)
	gotBrp, found := suite.app.MakerKeeper.GetBackingRiskParams(suite.ctx, brp2.BackingDenom)
	suite.Require().False(found)
	suite.Require().Equal(types.BackingRiskParams{}, gotBrp)

	// backing denom found
	isBackingRegistered = suite.app.MakerKeeper.IsBackingRegistered(suite.ctx, brp.BackingDenom)
	suite.Require().True(isBackingRegistered)
	gotBrp, found = suite.app.MakerKeeper.GetBackingRiskParams(suite.ctx, brp.BackingDenom)
	suite.Require().True(found)
	suite.Require().Equal(brp, gotBrp)

	// update backing risk params
	newMaxBacking := sdk.NewInt(300)
	newBrp := brp
	newBrp.MaxBacking = &newMaxBacking
	suite.Require().NotEqual(brp.MaxBacking, newBrp.MaxBacking)
	suite.app.MakerKeeper.SetBackingRiskParams(suite.ctx, newBrp)
	gotBrp, found = suite.app.MakerKeeper.GetBackingRiskParams(suite.ctx, brp.BackingDenom)
	suite.Require().True(found)
	suite.Require().Equal(newBrp.MaxBacking, gotBrp.MaxBacking)
}

func (suite *KeeperTestSuite) TestSetGetCollateralRiskParams() {
	crp, crp2 := suite.dummyCollateralRiskParams()
	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, crp)

	// collateral denom not found
	suite.Require().NotEqual(crp.CollateralDenom, crp2.CollateralDenom)
	isCollateralRegistered := suite.app.MakerKeeper.IsCollateralRegistered(suite.ctx, crp2.CollateralDenom)
	suite.Require().False(isCollateralRegistered)
	gotCrp, found := suite.app.MakerKeeper.GetCollateralRiskParams(suite.ctx, crp2.CollateralDenom)
	suite.Require().False(found)
	suite.Require().Equal(types.CollateralRiskParams{}, gotCrp)

	// collateral denom found
	isCollateralRegistered = suite.app.MakerKeeper.IsCollateralRegistered(suite.ctx, crp.CollateralDenom)
	suite.Require().True(isCollateralRegistered)
	gotCrp, found = suite.app.MakerKeeper.GetCollateralRiskParams(suite.ctx, crp.CollateralDenom)
	suite.Require().True(found)
	suite.Require().Equal(crp, gotCrp)

	// update Collateral risk params
	newMaxCollateral := sdk.NewInt(300)
	newCrp := crp
	newCrp.MaxCollateral = &newMaxCollateral
	suite.Require().NotEqual(crp.MaxCollateral, newCrp.MaxCollateral)
	suite.app.MakerKeeper.SetCollateralRiskParams(suite.ctx, newCrp)
	gotCrp, found = suite.app.MakerKeeper.GetCollateralRiskParams(suite.ctx, crp.CollateralDenom)
	suite.Require().True(found)
	suite.Require().Equal(newCrp.MaxCollateral, gotCrp.MaxCollateral)

}

func (suite *KeeperTestSuite) dummyBackingRiskParams() (brp, brp2 types.BackingRiskParams) {
	maxBacking := sdk.NewInt(10_100000)
	maxMerMint := sdk.NewInt(10_000000)
	mintFee := sdk.NewDecWithPrec(5, 3)
	burnFee := sdk.NewDecWithPrec(6, 3)
	buybackFee := sdk.NewDecWithPrec(7, 3)
	rebackFee := sdk.NewDecWithPrec(4, 3)
	brp = types.BackingRiskParams{
		BackingDenom: suite.bcDenom,
		Enabled:      true,
		MaxBacking:   &maxBacking,
		MaxMerMint:   &maxMerMint,
		MintFee:      &mintFee,
		BurnFee:      &burnFee,
		BuybackFee:   &buybackFee,
		RebackFee:    &rebackFee,
	}

	maxBacking2 := sdk.NewInt(200)
	maxMerMint2 := sdk.NewInt(2000_000000)
	mintFee2 := sdk.NewDecWithPrec(6, 3)
	burnFee2 := sdk.NewDecWithPrec(7, 3)
	buybackFee2 := sdk.NewDecWithPrec(8, 3)
	rebackFee2 := sdk.NewDecWithPrec(9, 3)
	brp2 = types.BackingRiskParams{
		BackingDenom: "eth",
		Enabled:      false,
		MaxBacking:   &maxBacking2,
		MaxMerMint:   &maxMerMint2,
		MintFee:      &mintFee2,
		BurnFee:      &burnFee2,
		BuybackFee:   &buybackFee2,
		RebackFee:    &rebackFee2,
	}

	return
}

func (suite *KeeperTestSuite) dummyCollateralRiskParams() (crp, crp2 types.CollateralRiskParams) {
	maxCollateral := sdk.NewInt(16_000000)
	maxMerMint := sdk.NewInt(10_000000)
	liquidationThreshold := sdk.NewDecWithPrec(90, 2)
	loanToValue := sdk.NewDecWithPrec(80, 2)
	basicLoanToValue := sdk.NewDecWithPrec(50, 2)
	catalyticLionRatio := sdk.NewDecWithPrec(5, 2)
	liquidationFee := sdk.NewDecWithPrec(10, 2)
	mintFee := sdk.NewDecWithPrec(1, 2)
	interestFee := sdk.NewDec(4)
	crp = types.CollateralRiskParams{
		CollateralDenom:      suite.bcDenom,
		Enabled:              true,
		MaxCollateral:        &maxCollateral,
		MaxMerMint:           &maxMerMint,
		LiquidationThreshold: &liquidationThreshold,
		LoanToValue:          &loanToValue,
		BasicLoanToValue:     &basicLoanToValue,
		CatalyticLionRatio:   &catalyticLionRatio,
		LiquidationFee:       &liquidationFee,
		MintFee:              &mintFee,
		InterestFee:          &interestFee,
	}

	maxCollateral2 := sdk.NewInt(200)
	maxMerMint2 := sdk.NewInt(20)
	liquidationThreshold2 := sdk.NewDecWithPrec(91, 2)
	loanToValue2 := sdk.NewDecWithPrec(71, 2)
	basicLoanToValue2 := sdk.NewDecWithPrec(51, 2)
	catalyticLionRatio2 := sdk.NewDecWithPrec(6, 2)
	liquidationFee2 := sdk.NewDecWithPrec(11, 2)
	mintFee2 := sdk.NewDecWithPrec(2, 2)
	interestFee2 := sdk.NewDecWithPrec(4, 2)
	crp2 = types.CollateralRiskParams{
		CollateralDenom:      "eth",
		Enabled:              false,
		MaxCollateral:        &maxCollateral2,
		MaxMerMint:           &maxMerMint2,
		LiquidationThreshold: &liquidationThreshold2,
		LoanToValue:          &loanToValue2,
		BasicLoanToValue:     &basicLoanToValue2,
		CatalyticLionRatio:   &catalyticLionRatio2,
		LiquidationFee:       &liquidationFee2,
		MintFee:              &mintFee2,
		InterestFee:          &interestFee2,
	}

	return
}
