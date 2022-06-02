package maker_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/merlion-zone/merlion/app"
	"github.com/merlion-zone/merlion/x/maker"
	"github.com/merlion-zone/merlion/x/maker/types"
)

type GenesisTestSuite struct {
	suite.Suite
	ctx sdk.Context
	app *app.Merlion
}

func (suite *GenesisTestSuite) SetupTest() {
	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestMakerInitGenesis() {
	app := suite.app
	makerKeeper := app.MakerKeeper

	suite.Require().NotPanics(func() {
		maker.InitGenesis(suite.ctx, makerKeeper, *types.DefaultGenesis())
	})

	backingRatio := makerKeeper.GetBackingRatio(suite.ctx)
	params := makerKeeper.GetParams(suite.ctx)

	suite.Require().Equal(backingRatio, sdk.OneDec())
	suite.Require().Equal(params, types.DefaultParams())
}

func (suite *GenesisTestSuite) TestMakerExportGenesis() {
	app := suite.app
	makerKeeper := app.MakerKeeper

	suite.Require().NotPanics(func() {
		maker.InitGenesis(suite.ctx, makerKeeper, *types.DefaultGenesis())
	})

	genesisExported := maker.ExportGenesis(suite.ctx, makerKeeper)
	suite.Require().Equal(genesisExported.BackingRatio, sdk.OneDec())
	suite.Require().Equal(genesisExported.Params, types.DefaultParams())
}
