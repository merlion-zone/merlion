package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/gauge/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
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

	totalDeposited := g.GetTotalDepositedAmount(suite.ctx)
	deposited := g.GetDepositedAmountByUser(suite.ctx, veID)

	suite.Require().Equal(amount, totalDeposited)
	suite.Require().Equal(amount, deposited)

	owner := suite.app.NftKeeper.GetOwner(suite.ctx, vetypes.VeNftClass.Id, vetypes.VeIDFromUint64(veID))
	userVeID := g.GetUserVeIDByAddress(suite.ctx, owner)
	suite.Require().Equal(veID, userVeID)
}

func (suite *KeeperTestSuite) TestKeeper_Withdraw() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.BaseDenom
	k.CreateGauge(suite.ctx, depoistDenom)
	g := k.Gauge(suite.ctx, depoistDenom)

	// Deposit first
	veID := uint64(100)
	amount := sdk.NewInt(100)

	err := g.Deposit(suite.ctx, veID, amount)
	suite.Require().NoError(err)

	withdrawAmt := sdk.NewInt(50)
	err = g.Withdraw(suite.ctx, veID, amount.Add(sdk.NewInt(1)))
	suite.Require().Error(err, types.ErrTooLargeAmount)

	err = g.Withdraw(suite.ctx, veID, withdrawAmt)
	suite.Require().NoError(err)

	totalDeposited := g.GetTotalDepositedAmount(suite.ctx)
	deposited := g.GetDepositedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(amount.Sub(withdrawAmt), totalDeposited)
	suite.Require().Equal(amount.Sub(withdrawAmt), deposited)

	// Withdraw all
	err = g.Withdraw(suite.ctx, veID, withdrawAmt)
	suite.Require().NoError(err)

	totalDeposited = g.GetTotalDepositedAmount(suite.ctx)
	deposited = g.GetDepositedAmountByUser(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), totalDeposited)
	suite.Require().Equal(sdk.ZeroInt(), deposited)

	owner := suite.app.NftKeeper.GetOwner(suite.ctx, vetypes.VeNftClass.Id, vetypes.VeIDFromUint64(veID))
	userVeID := g.GetUserVeIDByAddress(suite.ctx, owner)
	suite.Require().Equal(uint64(vetypes.EmptyVeID), userVeID)
}

func (suite *KeeperTestSuite) TestKeeper_DepositReward() {
	suite.SetupTest()
	k := suite.app.GaugeKeeper
	depoistDenom := mertypes.MicroUSMDenom
	k.CreateGauge(suite.ctx, depoistDenom)
	g := k.Gauge(suite.ctx, depoistDenom)

	sender := sdk.AccAddress(suite.address.Bytes())
	rewardDenom := mertypes.BaseDenom
	amount := sdk.NewInt(5000000)

	err := g.DepositReward(suite.ctx, sender, mertypes.MicroUSMDenom, amount)
	suite.Require().Error(err, types.ErrInvalidDepositDenom)

	err = g.DepositReward(suite.ctx, sender, rewardDenom, sdk.ZeroInt())
	suite.Require().Error(err, types.ErrInvalidAmount)

	err = g.DepositReward(suite.ctx, sender, rewardDenom, amount)
	suite.Require().NoError(err)

	// Check results
	addr := g.EscrowPool(suite.ctx).GetAddress()
	balance := suite.app.BankKeeper.GetBalance(suite.ctx, addr, rewardDenom)
	suite.Require().Equal(amount, balance.Amount)

	reward := g.GetReward(suite.ctx, rewardDenom)
	suite.Require().Equal(rewardDenom, reward.Denom)
	suite.Require().Equal(sdk.NewInt(8), reward.Rate)
	suite.Require().Equal(sdk.NewInt(0), reward.CumulativePerTicket)
	suite.Require().Equal(sdk.NewInt(0), reward.AccruedAmount)
}
