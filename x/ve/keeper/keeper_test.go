package keeper_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/merlion-zone/merlion/app"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/ve/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
	"github.com/tendermint/tendermint/version"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
	"github.com/tharsis/ethermint/tests"
)

type KeeperTestSuite struct {
	suite.Suite
	ctx sdk.Context
	app *app.Merlion

	address common.Address
	signer  keyring.Signer

	queryClient types.QueryClient
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) SetupTest() {
	require := suite.Require()

	// account key
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	suite.address = common.BytesToAddress(priv.PubKey().Address().Bytes())
	suite.signer = tests.NewSigner(priv)

	// consensus key
	privCons, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	consAddress := sdk.ConsAddress(privCons.PubKey().Address())

	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{
		Version: tmversion.Consensus{
			Block: version.BlockProtocol,
		},
		ChainID:         "merlion_5000-101",
		Height:          1,
		Time:            time.Now().UTC(),
		ProposerAddress: consAddress.Bytes(),
	})

	// set validator
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, privCons.PubKey(), stakingtypes.Description{})
	require.NoError(err)
	validator = stakingkeeper.TestingUpdateValidator(suite.app.StakingKeeper.Keeper, suite.ctx, validator, true)
	suite.app.StakingKeeper.AfterValidatorCreated(suite.ctx, validator.GetOperator())
	err = suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
	require.NoError(err)

	amount := sdk.NewInt64Coin(merlion.BaseDenom, 10000)
	err = app.FundAccount(suite.app.BankKeeper, suite.ctx, sdk.AccAddress(suite.address.Bytes()), sdk.NewCoins(amount))
	require.NoError(err)
}
