package keeper_test

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/keeper"
)

func (suite *KeeperTestSuite) TestEmitter_AddTotalEmission() {
	suite.SetupTest()
	emitter := keeper.NewEmitter(suite.app.VeKeeper)
	total := sdk.NewInt(100)
	emitter.AddTotalEmission(suite.ctx, total)

	emission := suite.app.VeKeeper.GetTotalEmission(suite.ctx)
	raw, _ := new(big.Int).SetString("450000000000000000000000100", 10)
	suite.Require().Equal(sdk.NewIntFromBigInt(raw), emission)

	emission = suite.app.VeKeeper.GetEmissionAtLastPeriod(suite.ctx)
	suite.Require().Equal("2979900176780079392645357", emission.BigInt().String())
}

func (suite *KeeperTestSuite) TestEmitter_CirculationSupply() {
	suite.SetupTest()
	emitter := keeper.NewEmitter((suite.app.VeKeeper))
	supply := emitter.CirculationSupply(suite.ctx)
	raw, _ := new(big.Int).SetString("500000000000000000000010000", 10)
	suite.Require().Equal(raw, supply.BigInt())
}

func (suite *KeeperTestSuite) TestEmitter_CirculationRate() {
	suite.SetupTest()
	emitter := keeper.NewEmitter((suite.app.VeKeeper))
	rate := emitter.CirculationRate(suite.ctx)
	suite.Require().Equal(sdk.NewDecFromInt(sdk.NewInt(1)), rate)
}

func (suite *KeeperTestSuite) TestEmitter_Emission() {
	suite.SetupTest()
	emitter := keeper.NewEmitter((suite.app.VeKeeper))
	emission := emitter.Emission(suite.ctx)
	suite.Require().Equal("2970033726709441703282758", emission.BigInt().String())
}

func (suite *KeeperTestSuite) TestEmitter_EmissionCompensation() {
	suite.SetupTest()
	emitter := keeper.NewEmitter((suite.app.VeKeeper))
	emission := emitter.EmissionCompensation(suite.ctx, sdk.NewInt(100))
	suite.Require().Equal(sdk.ZeroInt(), emission)
}

func (suite *KeeperTestSuite) TestEmitter_Emit() {
	// TODO
}

func (suite *KeeperTestSuite) TestKeeper_AddTotalEmission() {
	suite.SetupTest()
	total := sdk.NewInt(100)
	suite.app.VeKeeper.AddTotalEmission(suite.ctx, total)

	emission := suite.app.VeKeeper.GetTotalEmission(suite.ctx)
	raw, _ := new(big.Int).SetString("450000000000000000000000100", 10)
	suite.Require().Equal(sdk.NewIntFromBigInt(raw), emission)

	emission = suite.app.VeKeeper.GetEmissionAtLastPeriod(suite.ctx)
	suite.Require().Equal("2979900176780079392645357", emission.BigInt().String())
}

func (suite *KeeperTestSuite) TestKeeper_GetTotalEmission_SetTotalEmission() {
	suite.SetupTest()
	emission := suite.app.VeKeeper.GetTotalEmission(suite.ctx)
	raw, _ := new(big.Int).SetString("450000000000000000000000000", 10)
	suite.Require().Equal(sdk.NewIntFromBigInt(raw), emission)

	total := sdk.NewInt(100)
	suite.app.VeKeeper.SetTotalEmission(suite.ctx, total)
	emission = suite.app.VeKeeper.GetTotalEmission(suite.ctx)
	suite.Require().Equal(total, emission)
}

func (suite *KeeperTestSuite) TestKeeper_SetEmissionAtLastPeriod_GetEmissionAtLastPeriod() {
	suite.SetupTest()
	emission := suite.app.VeKeeper.GetEmissionAtLastPeriod(suite.ctx)
	raw, _ := new(big.Int).SetString("2979900176780079392645357", 10)
	suite.Require().Equal(sdk.NewIntFromBigInt(raw), emission)

	val := sdk.NewInt(100)
	suite.app.VeKeeper.SetEmissionAtLastPeriod(suite.ctx, val)
	emission = suite.app.VeKeeper.GetEmissionAtLastPeriod(suite.ctx)
	suite.Require().Equal(val, emission)
}

func (suite *KeeperTestSuite) TestKeeper_SetEmissionLastTimestamp_GetEmissionLastTimestamp() {
	suite.SetupTest()
	timestamp := suite.app.VeKeeper.GetEmissionLastTimestamp(suite.ctx)
	suite.Require().Equal(uint64(0), timestamp)

	timestamp = uint64(1754379718)
	suite.app.VeKeeper.SetEmissionLastTimestamp(suite.ctx, timestamp)
	stmp := suite.app.VeKeeper.GetEmissionLastTimestamp(suite.ctx)
	suite.Require().Equal(timestamp, stmp)
}
