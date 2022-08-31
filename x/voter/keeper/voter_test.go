package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mertypes "github.com/merlion-zone/merlion/types"
)

func (suite *KeeperTestSuite) TestKeeper_CreateGauge() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	depoistDenom := mertypes.BaseDenom
	k.CreateGauge(suite.ctx, depoistDenom)
	hasGauge := suite.app.GaugeKeeper.HasGauge(suite.ctx, depoistDenom)
	suite.Require().Equal(true, hasGauge)

	claimable := k.GetClaimableRewardByGauge(suite.ctx, depoistDenom)
	index := k.GetIndexAtLastUpdatedByGauge(suite.ctx, depoistDenom)
	suite.Require().Equal(sdk.ZeroInt(), claimable)
	suite.Require().Equal(sdk.ZeroInt(), index)
}
