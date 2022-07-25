package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/app"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx          sdk.Context
	app          *app.Merlion
	coinMetadata banktypes.Metadata
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
	s.coinMetadata = banktypes.Metadata{
		Description: "USDT",
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: "u" + "usd", Exponent: uint32(0), Aliases: []string{"micro" + "usd"}}, // e.g., uusd
			{Denom: "m" + "usd", Exponent: uint32(3), Aliases: []string{"milli" + "usd"}}, // e.g., musd
			{Denom: "usd", Exponent: uint32(6), Aliases: []string{}},                      // e.g., usd
		},
		Base:    "uusd",
		Display: "usd",
		Name:    "USDT",
		Symbol:  "USDT",
	}
	suite.Run(t, s)
}

func (suite *KeeperTestSuite) SetupTest() {
	var (
		PKS   = simapp.CreateTestPubKeys(5)
		addrs = []sdk.AccAddress{
			sdk.AccAddress(PKS[0].Address()),
			sdk.AccAddress(PKS[1].Address()),
			sdk.AccAddress(PKS[2].Address()),
			sdk.AccAddress(PKS[3].Address()),
			sdk.AccAddress(PKS[4].Address()),
		}
		valConsPk1 = PKS[0]
	)

	// init app
	suite.app = app.Setup(false)
	ctx := suite.app.BaseApp.NewContext(false, tmproto.Header{})
	app.FundTestAddrs(suite.app, ctx, addrs, sdk.NewInt(1234))

	// setup context
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		ChainID:         "merlion_5000-101",
		Height:          1,
		Time:            time.Now().UTC(),
		ProposerAddress: addrs[0].Bytes(),
	})

	// setup validator
	tstaking := teststaking.NewHelper(suite.T(), suite.ctx, suite.app.StakingKeeper.Keeper)
	tstaking.Denom = mertypes.AttoLionDenom

	// create validator with 50% commission
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(sdk.ValAddress(valConsPk1.Address()), valConsPk1, sdk.NewInt(100), true)
}
