package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/types"
	vekeeper "github.com/merlion-zone/merlion/x/ve/keeper"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

// TODO:
func (suite *KeeperTestSuite) TestKeeper_Slash() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	delAcct := sdk.AccAddress(suite.address.Bytes())
	ctx := sdk.WrapSDKContext(suite.ctx)
	denom := k.BondDenom(suite.ctx)
	veServer := vekeeper.NewMsgServerImpl(suite.app.VeKeeper)

	// VeDelegate in advance
	veTokens := types.VeTokensSlice{
		types.VeTokens{
			VeId:   uint64(1),
			Tokens: sdk.NewInt(100),
		},
	}
	amount := sdk.NewCoin(denom, veTokens[0].Tokens)
	boundAmt := sdk.NewInt(100)
	veID := fmt.Sprintf("ve-%d", veTokens[0].VeId)
	_, err := veServer.Create(ctx, &vetypes.MsgCreate{
		Sender:       delAcct.String(),
		To:           delAcct.String(),
		Amount:       amount,
		LockDuration: vetypes.RegulatedPeriod,
	})
	suite.Require().NoError(err)

	_, err = veServer.Deposit(ctx, &vetypes.MsgDeposit{
		Sender: delAcct.String(),
		VeId:   veID,
		Amount: amount,
	})
	suite.Require().NoError(err)

	_, err = k.VeDelegate(suite.ctx, delAcct, boundAmt, veTokens, stakingtypes.Bonded, suite.validator, false)
	suite.Require().NoError(err)

	valNoExistAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	valNoExist, err := sdk.ValAddressFromBech32(valNoExistAddr)
	suite.Require().NoError(err)

	// Try ErrNoDelegatorForAddress
	_, err = k.Undelegate(suite.ctx, delAcct, valNoExist, sdk.NewDec(1))
	suite.Require().Error(err, stakingtypes.ErrNoDelegatorForAddress)

	_, err = k.Undelegate(suite.ctx, delAcct, suite.validator.GetOperator(), sdk.NewDec(90))
	suite.Require().NoError(err)

	notBounded := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, stakingtypes.NotBondedPoolName).GetAddress()
	notBoundBalance := suite.app.BankKeeper.GetBalance(suite.ctx, notBounded, denom)
	suite.Require().Equal(sdk.ZeroInt(), notBoundBalance.Amount)
}
