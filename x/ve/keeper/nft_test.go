package keeper_test

func (suite *KeeperTestSuite) TestKeeper_CheckVeAttached() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	veID := uint64(1)
	err := k.CheckVeAttached(suite.ctx, veID)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestKeeper_HasNftClass() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	has := k.HasNftClass(suite.ctx)
	suite.Require().Equal(true, has)
}

func (suite *KeeperTestSuite) TestKeeper_SetNextVeID_GetNextVeID() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	veID := k.GetNextVeID(suite.ctx)
	suite.Require().Equal(uint64(1), veID)

	val := uint64(3)
	k.SetNextVeID(suite.ctx, val)
	suite.Require().Equal(val, k.GetNextVeID(suite.ctx))
}

func (suite *KeeperTestSuite) TestKeeper_IncVeAttached_DecVeAttached() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	veID := uint64(1)
	val := uint64(2)
	k.SetVeAttached(suite.ctx, veID, val)
	attached := k.GetVeAttached(suite.ctx, veID)
	suite.Require().Equal(val, attached)

	k.IncVeAttached(suite.ctx, veID)
	suite.Require().Equal(val+1, k.GetVeAttached(suite.ctx, veID))

	k.DecVeAttached(suite.ctx, veID)
	suite.Require().Equal(val, k.GetVeAttached(suite.ctx, veID))
}

func (suite *KeeperTestSuite) TestKeeper_SetVeAttached_GetVeAttached() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	veID := uint64(1)
	attached := k.GetVeAttached(suite.ctx, veID)
	suite.Require().Equal(uint64(0), attached)

	val := uint64(2)
	k.SetVeAttached(suite.ctx, veID, val)
	attached = k.GetVeAttached(suite.ctx, veID)
	suite.Require().Equal(val, attached)
}

func (suite *KeeperTestSuite) TestKeeper_SetVeVoted_GetVeVoted() {
	suite.SetupTest()
	k := suite.app.VeKeeper
	voted := k.GetVeVoted(suite.ctx, uint64(1))
	suite.Require().Equal(false, voted)

	k.SetVeVoted(suite.ctx, uint64(1), true)
	voted = k.GetVeVoted(suite.ctx, uint64(1))
	suite.Require().Equal(true, voted)
}
