package keeper_test

import sdk "github.com/cosmos/cosmos-sdk/types"

func (suite *KeeperTestSuite) TestKeeper_SetTotalVotes_GetTotalVotes() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	votes := k.GetTotalVotes(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), votes)

	v := sdk.NewInt(100)
	k.SetTotalVotes(suite.ctx, v)
	votes = k.GetTotalVotes(suite.ctx)
	suite.Require().Equal(v, votes)
}

func (suite *KeeperTestSuite) TestKeeper_SetTotalVotesByUser_GetTotalVotesByUser() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	veID := uint64(100)
	votes := k.GetTotalVotesByUser(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), votes)

	v := sdk.NewInt(100)
	k.SetTotalVotesByUser(suite.ctx, veID, v)
	votes = k.GetTotalVotesByUser(suite.ctx, veID)
	suite.Require().Equal(v, votes)
}

func (suite *KeeperTestSuite) TestKeeper_DeleteTotalVotesByUser() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	veID := uint64(100)
	v := sdk.NewInt(100)
	k.SetTotalVotesByUser(suite.ctx, veID, v)
	k.DeleteTotalVotesByUser(suite.ctx, veID)
	votes := k.GetTotalVotesByUser(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), votes)
}

func (suite *KeeperTestSuite) TestKeeper_SetPoolWeightedVotes_GetPoolWeightedVotes() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	poolDenom := "alion"
	votes := k.GetPoolWeightedVotes(suite.ctx, poolDenom)
	suite.Require().Equal(sdk.ZeroInt(), votes)

	v := sdk.NewInt(100)
	k.SetPoolWeightedVotes(suite.ctx, poolDenom, v)
	votes = k.GetPoolWeightedVotes(suite.ctx, poolDenom)
	suite.Require().Equal(v, votes)
}

func (suite *KeeperTestSuite) TestKeeper_SetPoolWeightedVotesByUser_GetPoolWeightedVotesByUser() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	veID := uint64(100)
	poolDenom := "alion"
	votes := k.GetPoolWeightedVotesByUser(suite.ctx, veID, poolDenom)
	suite.Require().Equal(sdk.ZeroInt(), votes)

	v := sdk.NewInt(100)
	k.SetPoolWeightedVotesByUser(suite.ctx, veID, poolDenom, v)
	votes = k.GetPoolWeightedVotesByUser(suite.ctx, veID, poolDenom)
	suite.Require().Equal(v, votes)
}

func (suite *KeeperTestSuite) TestKeeper_DeletePoolWeightedVotesByUser() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	veID := uint64(100)
	poolDenom := "alion"
	v := sdk.NewInt(100)
	k.SetPoolWeightedVotesByUser(suite.ctx, veID, poolDenom, v)
	k.DeletePoolWeightedVotesByUser(suite.ctx, veID, poolDenom)
	votes := k.GetPoolWeightedVotesByUser(suite.ctx, veID, poolDenom)
	suite.Require().Equal(sdk.ZeroInt(), votes)
}

func (suite *KeeperTestSuite) TestKeeper_SetIndex_GetIndex() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	index := k.GetIndex(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), index)

	v := sdk.NewInt(100)
	k.SetIndex(suite.ctx, v)
	index = k.GetIndex(suite.ctx)
	suite.Require().Equal(v, index)
}

func (suite *KeeperTestSuite) TestKeeper_SetIndexAtLastUpdatedByGauge_GetIndexAtLastUpdatedByGauge() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	poolDenom := "alion"
	index := k.GetIndexAtLastUpdatedByGauge(suite.ctx, poolDenom)
	suite.Require().Equal(sdk.ZeroInt(), index)

	v := sdk.NewInt(100)
	k.SetIndexAtLastUpdatedByGauge(suite.ctx, poolDenom, v)
	index = k.GetIndexAtLastUpdatedByGauge(suite.ctx, poolDenom)
	suite.Require().Equal(v, index)
}

func (suite *KeeperTestSuite) TestKeeper_SetClaimableRewardByGauge_GetClaimableRewardByGauge() {
	suite.SetupTest()
	k := suite.app.VoterKeeper
	poolDenom := "alion"
	index := k.GetClaimableRewardByGauge(suite.ctx, poolDenom)
	suite.Require().Equal(sdk.ZeroInt(), index)

	v := sdk.NewInt(100)
	k.SetClaimableRewardByGauge(suite.ctx, poolDenom, v)
	index = k.GetClaimableRewardByGauge(suite.ctx, poolDenom)
	suite.Require().Equal(v, index)
}
