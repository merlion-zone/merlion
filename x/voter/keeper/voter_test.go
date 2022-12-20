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

func (suite *KeeperTestSuite) TestKeeper_Abstain() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	depoistDenom := mertypes.BaseDenom
	k.CreateGauge(suite.ctx, depoistDenom)

	veID := uint64(100)
	k.Abstain(suite.ctx, veID)

	poolWeightedVotesByUser := k.GetPoolWeightedVotesByUser(suite.ctx, veID, depoistDenom)
	suite.Require().Equal(sdk.ZeroInt(), poolWeightedVotesByUser)

	poolWeightedVotes := k.GetPoolWeightedVotes(suite.ctx, depoistDenom)
	suite.Require().Equal(sdk.ZeroInt(), poolWeightedVotes)

	userVotes := k.GetTotalVotesByUser(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), userVotes)

	totalVotes := k.GetTotalVotes(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), totalVotes)

	veVoted := suite.app.VeKeeper.GetVeVoted(suite.ctx, veID)
	suite.Require().Equal(false, veVoted)
}

func (suite *KeeperTestSuite) TestKeeper_DepositReward() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	amount := sdk.NewInt(100)
	sender := sdk.AccAddress(suite.address.Bytes())
	k.DepositReward(suite.ctx, sender, amount)

	denom := suite.app.VeKeeper.LockDenom(suite.ctx)
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, sender, denom)
	suite.Require().Equal(sdk.NewInt(9900), balance.Amount)

	index := k.GetIndex(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), index)
}
