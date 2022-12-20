package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/gauge/keeper"
	"github.com/merlion-zone/merlion/x/gauge/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

func (suite *KeeperTestSuite) TestKeeper_HasGauge_SetGauge_GetGauges() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	k.SetGauge(suite.ctx, depoistDenom)
	suite.Require().Equal(false, k.HasGauge(suite.ctx, "xxx"))
	suite.Require().Equal(true, k.HasGauge(suite.ctx, depoistDenom))
	denoms := k.GetGauges(suite.ctx)
	suite.Require().Equal(depoistDenom, denoms[0])
}

func (suite *KeeperTestSuite) TestBase_SetTotalDepositedAmount_GetTotalDepositedAmount() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	amount := base.GetTotalDepositedAmount(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), amount)

	amt := sdk.NewInt(100)
	base.SetTotalDepositedAmount(suite.ctx, amt)

	amount = base.GetTotalDepositedAmount(suite.ctx)
	suite.Require().Equal(amt, amount)
}

func (suite *KeeperTestSuite) TestBase_SetDepositedAmountByUser_GetDepositedAmountByUser() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	veID := uint64(100)
	amount := base.GetDepositedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), amount)

	amt := sdk.NewInt(100)
	base.SetDepositedAmountByUser(suite.ctx, veID, amt)

	amount = base.GetDepositedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(amt, amount)
}

func (suite *KeeperTestSuite) TestBase_DeleteDepositedAmountByUser() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	veID := uint64(100)
	amt := sdk.NewInt(100)
	base.SetDepositedAmountByUser(suite.ctx, veID, amt)
	base.DeleteDepositedAmountByUser(suite.ctx, veID)

	amount := base.GetDepositedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), amount)
}

func (suite *KeeperTestSuite) TestBase_SetTotalDerivedAmount_GetTotalDerivedAmount() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	amount := base.GetTotalDerivedAmount(suite.ctx)
	suite.Require().Equal(sdk.ZeroInt(), amount)

	amt := sdk.NewInt(100)
	base.SetTotalDerivedAmount(suite.ctx, amt)

	amount = base.GetTotalDerivedAmount(suite.ctx)
	suite.Require().Equal(amt, amount)
}

func (suite *KeeperTestSuite) TestBase_SetDerivedAmountByUser_GetDerivedAmountByUser() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	veID := uint64(100)
	amount := base.GetDerivedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), amount)

	amt := sdk.NewInt(100)
	base.SetDerivedAmountByUser(suite.ctx, veID, amt)

	amount = base.GetDerivedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(amt, amount)
}

func (suite *KeeperTestSuite) TestBase_SetReward_GetReward() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	rewardDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	reward := base.GetReward(suite.ctx, rewardDenom)
	suite.Require().Equal(types.Reward{
		Denom:               rewardDenom,
		Rate:                sdk.ZeroInt(),
		FinishTime:          0,
		LastUpdateTime:      0,
		CumulativePerTicket: sdk.ZeroInt(),
		AccruedAmount:       sdk.ZeroInt(),
	}, reward)

	reward = types.Reward{
		Denom:               rewardDenom,
		Rate:                sdk.NewInt(1),
		FinishTime:          1,
		LastUpdateTime:      2,
		CumulativePerTicket: sdk.NewInt(100),
		AccruedAmount:       sdk.NewInt(100),
	}
	base.SetReward(suite.ctx, rewardDenom, reward)

	res := base.GetReward(suite.ctx, rewardDenom)
	suite.Require().Equal(reward, res)
}

func (suite *KeeperTestSuite) TestBase_SetUserReward_GetUserReward() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	rewardDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	veID := uint64(100)
	reward := base.GetUserReward(suite.ctx, rewardDenom, veID)
	suite.Require().Equal(types.UserReward{
		Denom:               rewardDenom,
		VeId:                veID,
		LastClaimTime:       0,
		CumulativePerTicket: sdk.ZeroInt(),
	}, reward)

	reward = types.UserReward{
		Denom:               rewardDenom,
		VeId:                veID,
		LastClaimTime:       1,
		CumulativePerTicket: sdk.NewInt(1),
	}
	base.SetUserReward(suite.ctx, rewardDenom, veID, reward)

	res := base.GetUserReward(suite.ctx, rewardDenom, veID)
	suite.Require().Equal(reward, res)
}

func (suite *KeeperTestSuite) TestBase_SetUserVeIDByAddress_GetUserVeIDByAddress() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	acc := suite.address.Bytes()
	veID := base.GetUserVeIDByAddress(suite.ctx, acc)
	suite.Require().Equal(uint64(vetypes.EmptyVeID), veID)

	veID = uint64(100)
	base.SetUserVeIDByAddress(suite.ctx, acc, veID)

	res := base.GetUserVeIDByAddress(suite.ctx, acc)
	suite.Require().Equal(veID, res)
}

func (suite *KeeperTestSuite) TestBase_DeleteUserVeIDByAddress() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	acc := suite.address.Bytes()
	veID := uint64(100)
	base.SetUserVeIDByAddress(suite.ctx, acc, veID)
	base.DeleteUserVeIDByAddress(suite.ctx, acc)
	veID = base.GetUserVeIDByAddress(suite.ctx, acc)
	suite.Require().Equal(uint64(vetypes.EmptyVeID), veID)
}

func (suite *KeeperTestSuite) TestBase_GetEpoch_SetEpoch() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	epoch := base.GetEpoch(suite.ctx)
	suite.Require().Equal(uint64(keeper.EmptyEpoch), epoch)

	epoch = uint64(1)
	base.SetEpoch(suite.ctx, epoch)

	res := base.GetEpoch(suite.ctx)
	suite.Require().Equal(epoch, res)
}

func (suite *KeeperTestSuite) TestBase_SetCheckpoint_GetCheckpoint() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	epoch := uint64(1)
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	point := base.GetCheckpoint(suite.ctx, epoch)
	suite.Require().Equal(types.Checkpoint{}, point)

	point = types.Checkpoint{
		Timestamp: uint64(1661248381),
		Amount:    sdk.NewInt(100),
	}
	base.SetCheckpoint(suite.ctx, epoch, point)

	res := base.GetCheckpoint(suite.ctx, epoch)
	suite.Require().Equal(point, res)
}

func (suite *KeeperTestSuite) TestBase_SetUserEpoch_GetUserEpoch() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	veID := uint64(100)
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	epoch := base.GetUserEpoch(suite.ctx, veID)
	suite.Require().Equal(uint64(keeper.EmptyEpoch), epoch)

	epoch = uint64(1)
	base.SetUserEpoch(suite.ctx, veID, epoch)

	res := base.GetUserEpoch(suite.ctx, veID)
	suite.Require().Equal(epoch, res)
}

func (suite *KeeperTestSuite) TestBase_SetUserCheckpoint_GetUserCheckpoint() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	veID := uint64(100)
	epoch := uint64(1)
	depoistDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	point := base.GetUserCheckpoint(suite.ctx, veID, epoch)
	suite.Require().Equal(types.Checkpoint{}, point)

	point = types.Checkpoint{
		Timestamp: uint64(1661248381),
		Amount:    sdk.NewInt(100),
	}
	base.SetUserCheckpoint(suite.ctx, veID, epoch, point)

	res := base.GetUserCheckpoint(suite.ctx, veID, epoch)
	suite.Require().Equal(point, res)
}

func (suite *KeeperTestSuite) TestBase_SetRewardEpoch_GetRewardEpoch() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	rewardDenom := mertypes.BaseDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	epoch := base.GetRewardEpoch(suite.ctx, rewardDenom)
	suite.Require().Equal(uint64(keeper.EmptyEpoch), epoch)

	epoch = uint64(1)
	base.SetRewardEpoch(suite.ctx, rewardDenom, epoch)

	res := base.GetRewardEpoch(suite.ctx, rewardDenom)
	suite.Require().Equal(epoch, res)
}

func (suite *KeeperTestSuite) TestBase_SetRewardCheckpoint_GetRewardCheckpoint() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	rewardDenom := depoistDenom
	base := keeper.NewBase(k, depoistDenom, types.GaugeKey(depoistDenom), true)

	point := base.GetRewardCheckpoint(suite.ctx, rewardDenom, uint64(100))
	suite.Require().Equal(types.Checkpoint{}, point)

	epoch := uint64(1)
	point = types.Checkpoint{
		Timestamp: uint64(1661248381),
		Amount:    sdk.NewInt(100),
	}
	base.SetRewardCheckpoint(suite.ctx, depoistDenom, epoch, point)

	res := base.GetRewardCheckpoint(suite.ctx, rewardDenom, epoch)
	suite.Require().Equal(point, res)
}
