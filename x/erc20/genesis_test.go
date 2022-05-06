package erc20

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/app"
	"github.com/merlion-zone/merlion/x/erc20/types"
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

func (suite *GenesisTestSuite) TestERC20InitGenesis() {
	testCases := []struct {
		name         string
		genesisState types.GenesisState
	}{
		{
			"empty genesis",
			types.GenesisState{},
		},
		{
			"default genesis",
			*types.DefaultGenesis(),
		},
		{
			"custom genesis",
			types.GenesisState{
				types.DefaultParams(),
				[]types.TokenPair{
					{
						Erc20Address:  "0x5dCA2483280D9727c80b5518faC4556617fb19ZZ",
						Denom:         "coin",
						ContractOwner: types.OWNER_MODULE,
					},
				},
			},
		},
	}

	for _, tc := range testCases {

		suite.Require().NotPanics(func() {
			InitGenesis(suite.ctx, suite.app.Erc20Keeper, suite.app.AccountKeeper, tc.genesisState)
		})
		params := suite.app.Erc20Keeper.GetParams(suite.ctx)

		tokenPairs := suite.app.Erc20Keeper.GetAllTokenPairs(suite.ctx)
		suite.Require().Equal(tc.genesisState.Params, params)
		if len(tokenPairs) > 0 {
			suite.Require().Equal(tc.genesisState.TokenPairs, tokenPairs)
		} else {
			suite.Require().Len(tc.genesisState.TokenPairs, 0)
		}
	}
}

func (suite *GenesisTestSuite) TestErc20ExportGenesis() {
	testGenCases := []struct {
		name         string
		genesisState types.GenesisState
	}{
		{
			"empty genesis",
			types.GenesisState{},
		},
		{
			"default genesis",
			*types.DefaultGenesis(),
		},
		{
			"custom genesis",
			types.GenesisState{
				types.DefaultParams(),
				[]types.TokenPair{
					{
						Erc20Address:  "0x5dCA2483280D9727c80b5518faC4556617fb19ZZ",
						Denom:         "coin",
						ContractOwner: types.OWNER_MODULE,
					},
				},
			},
		},
	}

	for _, tc := range testGenCases {
		InitGenesis(suite.ctx, suite.app.Erc20Keeper, suite.app.AccountKeeper, tc.genesisState)
		suite.Require().NotPanics(func() {
			genesisExported := ExportGenesis(suite.ctx, suite.app.Erc20Keeper)
			params := suite.app.Erc20Keeper.GetParams(suite.ctx)
			suite.Require().Equal(genesisExported.Params, params)

			tokenPairs := suite.app.Erc20Keeper.GetAllTokenPairs(suite.ctx)
			if len(tokenPairs) > 0 {
				suite.Require().Equal(genesisExported.TokenPairs, tokenPairs)
			} else {
				suite.Require().Len(genesisExported.TokenPairs, 0)
			}
		})
		// }
	}
}
