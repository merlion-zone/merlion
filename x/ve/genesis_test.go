package ve

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/app"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/ve/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type GenesisTestSuite struct {
	suite.Suite
	ctx sdk.Context
	app *app.Merlion
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
}

func (suite *GenesisTestSuite) TestVeInitGenesis() {
	app := suite.app
	veKeeper := app.VeKeeper

	suite.Require().NotPanics(func() {
		InitGenesis(suite.ctx, veKeeper, *types.DefaultGenesis())
	})

	params := veKeeper.GetParams(suite.ctx)
	suite.Require().Equal(params.GetLockDenom(), merlion.BaseDenom)

	suite.Require().Equal(sdk.ZeroInt(), veKeeper.GetTotalLockedAmount(suite.ctx))
	suite.Require().EqualValues(types.FirstVeID, veKeeper.GetNextVeID(suite.ctx))
	suite.Require().EqualValues(types.EmptyEpoch, veKeeper.GetEpoch(suite.ctx))
	suite.Require().Equal(types.Checkpoint{
		Bias:      sdk.ZeroInt(),
		Slope:     sdk.ZeroInt(),
		Timestamp: 0,
		Block:     0,
	}, veKeeper.GetCheckpoint(suite.ctx, types.EmptyEpoch))
}

func (suite *GenesisTestSuite) TestVeExportGenesis() {
	app := suite.app
	veKeeper := app.VeKeeper

	suite.Require().NotPanics(func() {
		InitGenesis(suite.ctx, veKeeper, *types.DefaultGenesis())
	})

	genesisExported := ExportGenesis(suite.ctx, app.VeKeeper)
	suite.Require().Equal(genesisExported.Params.GetLockDenom(), merlion.BaseDenom)
}
