package keeper_test

import (
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	nfttypes "github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/app"
	"github.com/merlion-zone/merlion/types"
	mertypes "github.com/merlion-zone/merlion/types"
	erc20types "github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx   sdk.Context
	app   *app.Merlion
	addrs []sdk.AccAddress
}

var s *KeeperTestSuite

func TestKeeperTestSuite(t *testing.T) {
	s = new(KeeperTestSuite)
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

	suite.addrs = addrs
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

func (suite *KeeperTestSuite) TestKeeper_GetPaginatedTotalSupply() {
	suite.SetupTest()
	var (
		t                = suite.T()
		k                = suite.app.BankKeeper
		req              = &query.PageRequest{Key: nil, Limit: 2, CountTotal: true}
		supply, res, err = k.GetPaginatedTotalSupply(suite.ctx, req)
	)

	require.NoError(t, err)
	amt, _ := new(big.Int).SetString("500000000000000000000006170", 10)
	expectedSupply := sdk.NewCoins(sdk.NewCoin(suite.app.StakingKeeper.BondDenom(suite.ctx), sdk.NewIntFromBigInt(amt)))
	require.Equal(t, expectedSupply, supply)
	require.Equal(t, uint64(1), res.Total)
}

func (suite *KeeperTestSuite) TestKeeper_DelegateCoins() {
	suite.SetupTest()
	var (
		t             = suite.T()
		k             = suite.app.BankKeeper
		delegatorAddr = suite.addrs[0]
		moduleAccAddr = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, erc20types.ModuleName).GetAddress()
		denom         = suite.app.StakingKeeper.BondDenom(suite.ctx)
	)

	err := k.DelegateCoins(suite.ctx, delegatorAddr, moduleAccAddr, sdk.NewCoins(sdk.NewCoin("erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", sdk.NewInt(10000))))
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for delegation"))

	err = k.DelegateCoins(suite.ctx, delegatorAddr, moduleAccAddr, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000))))
	require.NoError(t, err)

	delegatorBalance := k.GetBalance(suite.ctx, delegatorAddr, denom)
	moduleAccBalance := k.GetBalance(suite.ctx, moduleAccAddr, denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(134)), delegatorBalance)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1000)), moduleAccBalance)
}

func (suite *KeeperTestSuite) TestKeeper_UndelegateCoins() {
	suite.SetupTest()
	var (
		t             = suite.T()
		k             = suite.app.BankKeeper
		delegatorAddr = suite.addrs[0]
		moduleAccAddr = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, erc20types.ModuleName).GetAddress()
		denom         = suite.app.StakingKeeper.BondDenom(suite.ctx)
	)

	err := k.UndelegateCoins(suite.ctx, delegatorAddr, moduleAccAddr, sdk.NewCoins(sdk.NewCoin("erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75", sdk.NewInt(10000))))
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for delegation"))

	err = k.UndelegateCoins(suite.ctx, moduleAccAddr, delegatorAddr, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000))))
	require.Error(t, err, "0alion is smaller than 1000alion: insufficient funds")

	k.DelegateCoins(suite.ctx, delegatorAddr, moduleAccAddr, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000))))
	err = k.UndelegateCoins(suite.ctx, moduleAccAddr, delegatorAddr, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100))))
	require.NoError(t, err)

	delegatorBalance := k.GetBalance(suite.ctx, delegatorAddr, denom)
	moduleAccBalance := k.GetBalance(suite.ctx, moduleAccAddr, denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(234)), delegatorBalance)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(900)), moduleAccBalance)
}

func (suite *KeeperTestSuite) TestKeeper_GetSupply() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = suite.app.StakingKeeper.BondDenom(suite.ctx)
	)
	supply := k.GetSupply(suite.ctx, "uusd")
	require.Equal(t, sdk.NewCoin("uusd", sdk.NewInt(0)), supply)

	erc20Denom := "erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75"
	supply = k.GetSupply(suite.ctx, erc20Denom)
	require.Equal(t, sdk.Coin{}, supply)

	supply = k.GetSupply(suite.ctx, denom)
	amt, _ := new(big.Int).SetString("500000000000000000000006170", 10)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewIntFromBigInt(amt)), supply)
}

func (suite *KeeperTestSuite) TestKeeper_HasSupply() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = suite.app.StakingKeeper.BondDenom(suite.ctx)
	)
	hasSupply := k.HasSupply(suite.ctx, "uusd")
	require.Equal(t, false, hasSupply)

	erc20Denom := "erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75"
	hasSupply = k.HasSupply(suite.ctx, erc20Denom)
	require.Equal(t, false, hasSupply)

	hasSupply = k.HasSupply(suite.ctx, denom)
	require.Equal(t, true, hasSupply)
}

func (suite *KeeperTestSuite) TestKeeper_GetDenomMetaData() {
	suite.SetupTest()
	var (
		t        = suite.T()
		k        = suite.app.BankKeeper
		denom    = types.MicroUSMDenom
		base     = denom
		display  = base[1:]
		uusmMeta = banktypes.Metadata{
			Description: "The native stable token of the Merlion.",
			DenomUnits: []*banktypes.DenomUnit{
				{Denom: "u" + display, Exponent: uint32(0), Aliases: []string{"micro" + display}}, // e.g., uusm
				{Denom: "m" + display, Exponent: uint32(3), Aliases: []string{"milli" + display}}, // e.g., musm
				{Denom: display, Exponent: uint32(6), Aliases: []string{""}},                      // e.g., usm
			},
			Base:    base,
			Display: display,
			Name:    strings.ToUpper(display), // e.g., USM
			Symbol:  strings.ToUpper(display), // e.g., USM
		}
	)
	meta, ok := k.GetDenomMetaData(suite.ctx, "uusd")
	require.Equal(t, false, ok)
	require.Equal(t, banktypes.Metadata{}, meta)

	meta, ok = k.GetDenomMetaData(suite.ctx, denom)
	require.Equal(t, true, ok)
	require.Equal(t, uusmMeta.Description, meta.Description)
	require.Equal(t, uusmMeta.Base, meta.Base)
	require.Equal(t, uusmMeta.Display, meta.Display)
	require.Equal(t, uusmMeta.Name, meta.Name)
	require.Equal(t, uusmMeta.Symbol, meta.Symbol)
	require.Equal(t, uusmMeta.DenomUnits[0], meta.DenomUnits[0])
	require.Equal(t, uusmMeta.DenomUnits[1], meta.DenomUnits[1])
}

func (suite *KeeperTestSuite) TestKeeper_SendCoinsFromModuleToAccount() {
	suite.SetupTest()
	var (
		t            = suite.T()
		k            = suite.app.BankKeeper
		senderModule = erc20types.ModuleName
		denom        = types.AttoLionDenom
		amt1         = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(200)))
		amt2         = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100)))
	)
	recipientAddr := authtypes.NewModuleAddress(banktypes.ModuleName)
	err := k.SendCoinsFromModuleToAccount(suite.ctx, senderModule, recipientAddr, amt1)
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr))
	// Raw balance check
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), k.GetBalance(suite.ctx, suite.addrs[0], denom))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(0)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(erc20types.ModuleName), denom))

	// Send coins
	err = k.SendCoins(suite.ctx, suite.addrs[0], authtypes.NewModuleAddress(erc20types.ModuleName), amt1)
	require.NoError(t, err)

	err = k.SendCoinsFromModuleToAccount(suite.ctx, senderModule, suite.addrs[0], amt2)
	require.NoError(t, err)

	receiverBalance := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	moduleAccBalance := k.GetBalance(suite.ctx, authtypes.NewModuleAddress(erc20types.ModuleName), denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1034)), receiverBalance)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(100)), moduleAccBalance)
}

func (suite *KeeperTestSuite) TestKeeper_SendCoinsFromModuleToModule() {
	suite.SetupTest()
	var (
		t            = suite.T()
		k            = suite.app.BankKeeper
		senderModule = erc20types.ModuleName
		recvModule   = nfttypes.ModuleName
		denom        = types.AttoLionDenom
		amt1         = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(200)))
		amt2         = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100)))
	)
	recipientAddr := authtypes.NewModuleAddress(banktypes.ModuleName)
	err := k.SendCoinsFromModuleToAccount(suite.ctx, senderModule, recipientAddr, amt1)
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr))
	// Raw balance check
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), k.GetBalance(suite.ctx, suite.addrs[0], denom))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(0)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(senderModule), denom))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(0)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(recvModule), denom))

	// Send coins
	err = k.SendCoins(suite.ctx, suite.addrs[0], authtypes.NewModuleAddress(senderModule), amt1)
	require.NoError(t, err)

	err = k.SendCoinsFromModuleToModule(suite.ctx, senderModule, recvModule, amt2)
	require.NoError(t, err)

	senderModuleBalance := k.GetBalance(suite.ctx, authtypes.NewModuleAddress(senderModule), denom)
	recvModuleBalance := k.GetBalance(suite.ctx, authtypes.NewModuleAddress(recvModule), denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(100)), senderModuleBalance)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(100)), recvModuleBalance)
}

func (suite *KeeperTestSuite) TestKeeper_SendCoinsFromAccountToModule() {
	suite.SetupTest()
	var (
		t     = suite.T()
		k     = suite.app.BankKeeper
		denom = types.AttoLionDenom
		amt   = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100)))
	)
	// Raw balance check
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), k.GetBalance(suite.ctx, suite.addrs[0], denom))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(0)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(erc20types.ModuleName), denom))

	// Send coins
	err := k.SendCoinsFromAccountToModule(suite.ctx, suite.addrs[0], erc20types.ModuleName, amt)
	require.NoError(t, err)

	senderBalance := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	recvBalance := k.GetBalance(suite.ctx, authtypes.NewModuleAddress(erc20types.ModuleName), denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1034)), senderBalance)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(100)), recvBalance)
}

func (suite *KeeperTestSuite) TestKeeper_DelegateCoinsFromAccountToModule() {
	suite.SetupTest()
	var (
		t        = suite.T()
		k        = suite.app.BankKeeper
		denom    = types.AttoLionDenom
		recvAddr = stakingtypes.BondedPoolName
		amt      = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100)))
	)
	// Raw balance check
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), k.GetBalance(suite.ctx, suite.addrs[0], denom))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(0)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(recvAddr), denom))

	err := k.DelegateCoinsFromAccountToModule(suite.ctx, suite.addrs[0], recvAddr, amt)
	require.NoError(t, err)

	senderBalance := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	recvBalance := k.GetBalance(suite.ctx, authtypes.NewModuleAddress(recvAddr), denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1034)), senderBalance)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(100)), recvBalance)
}

func (suite *KeeperTestSuite) TestKeeper_UndelegateCoinsFromModuleToAccount() {
	suite.SetupTest()
	var (
		t        = suite.T()
		k        = suite.app.BankKeeper
		denom    = types.AttoLionDenom
		recvAddr = stakingtypes.BondedPoolName
		amt1     = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(200)))
		amt2     = sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100)))
	)
	// Raw balance check
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1134)), k.GetBalance(suite.ctx, suite.addrs[0], denom))
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(0)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(recvAddr), denom))

	err := k.DelegateCoinsFromAccountToModule(suite.ctx, suite.addrs[0], recvAddr, amt1)
	require.NoError(t, err)

	err = k.UndelegateCoinsFromModuleToAccount(suite.ctx, recvAddr, suite.addrs[0], amt2)
	require.NoError(t, err)

	senderBalance := k.GetBalance(suite.ctx, suite.addrs[0], denom)
	recvBalance := k.GetBalance(suite.ctx, authtypes.NewModuleAddress(recvAddr), denom)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(1034)), senderBalance)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(100)), recvBalance)
}

func (suite *KeeperTestSuite) TestKeeper_MintCoins() {
	suite.SetupTest()
	var (
		t          = suite.T()
		k          = suite.app.BankKeeper
		denom      = types.AttoLionDenom
		moduleName = erc20types.ModuleName
		erc20Denom = "erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75"
	)
	err := k.MintCoins(suite.ctx, moduleName, sdk.NewCoins(sdk.NewCoin(erc20Denom, sdk.NewInt(100))))
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for mint"))

	err = k.MintCoins(suite.ctx, moduleName, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100))))
	require.NoError(t, err)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(100)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(moduleName), denom))
}

func (suite *KeeperTestSuite) TestKeeper_BurnCoins() {
	suite.SetupTest()
	var (
		t          = suite.T()
		k          = suite.app.BankKeeper
		denom      = types.AttoLionDenom
		moduleName = erc20types.ModuleName
		erc20Denom = "erc20/0xd567B3d7B8FE3C79a1AD8dA978812cfC4Fa05e75"
	)
	err := k.BurnCoins(suite.ctx, moduleName, sdk.NewCoins(sdk.NewCoin(erc20Denom, sdk.NewInt(100))))
	require.Error(t, err, sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for burn"))

	err = k.MintCoins(suite.ctx, moduleName, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100))))
	require.NoError(t, err)
	err = k.BurnCoins(suite.ctx, moduleName, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(50))))
	require.NoError(t, err)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(50)), k.GetBalance(suite.ctx, authtypes.NewModuleAddress(moduleName), denom))
}
