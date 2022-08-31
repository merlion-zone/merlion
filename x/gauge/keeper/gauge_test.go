package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
)

func (suite *KeeperTestSuite) TestKeeper_CreateGauge() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	k.CreateGauge(suite.ctx, depoistDenom)

	suite.Require().Equal(true, k.HasGauge(suite.ctx, depoistDenom))
}

func (suite *KeeperTestSuite) TestKeeper_Gauge() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	k.CreateGauge(suite.ctx, depoistDenom)

	g := k.Gauge(suite.ctx, depoistDenom)
	suite.Require().Equal(depoistDenom, g.Base.PoolDenom())
	suite.Require().Equal(fmt.Sprintf("%s_%s", types.GaugePoolName, depoistDenom), g.Base.PoolName())
}

// TODO:
func (suite *KeeperTestSuite) TestKeeper_Deposit() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	k.CreateGauge(suite.ctx, depoistDenom)
	g := k.Gauge(suite.ctx, depoistDenom)
	veID := uint64(100)
	amount := sdk.NewInt(100)

	err := g.Deposit(suite.ctx, veID, amount)
	suite.Require().NoError(err)

	// suite.Require().Equal(depoistDenom, g.Base.PoolDenom())
	// suite.Require().Equal(fmt.Sprintf("%s_%s", types.GaugePoolName, depoistDenom), g.Base.PoolName())
}

// TODO:
func (suite *KeeperTestSuite) TestKeeper_Withdraw() {
}

// TODO:
func (suite *KeeperTestSuite) TestKeeper_DepositReward() {
}

// TODO:
func (suite *KeeperTestSuite) TestKeeper_DepositFees() {
}
