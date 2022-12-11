package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/keeper"
	"github.com/merlion-zone/merlion/x/staking/types"
	vekeeper "github.com/merlion-zone/merlion/x/ve/keeper"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

func (suite *KeeperTestSuite) TestVeDelegate() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	ctx := sdk.WrapSDKContext(suite.ctx)
	delAcct := sdk.AccAddress(suite.address.Bytes())
	denom := "alion"
	impl := keeper.NewMsgServerImpl(k)

	// Prepare valid case
	amount := sdk.NewCoin(denom, sdk.NewInt(100))
	veServer := vekeeper.NewMsgServerImpl(suite.app.VeKeeper)
	_, err := veServer.Create(ctx, &vetypes.MsgCreate{
		Sender:       delAcct.String(),
		To:           delAcct.String(),
		Amount:       amount,
		LockDuration: vetypes.RegulatedPeriod,
	})
	suite.Require().NoError(err)

	_, err = veServer.Deposit(ctx, &vetypes.MsgDeposit{
		Sender: delAcct.String(),
		VeId:   "ve-1",
		Amount: amount,
	})
	suite.Require().NoError(err)

	testCases := []struct {
		name      string
		pass      bool
		delegator string
		validator string
		veID      string
		amount    sdk.Coin
	}{
		{"invalid validator address", false, delAcct.String(), "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjrx", "ve-100", sdk.NewCoin(denom, sdk.NewInt(100))},
		{"validator not found", false, delAcct.String(), "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3", "ve-100", sdk.NewCoin(denom, sdk.NewInt(100))},
		{"invalid delegator address", false, "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg2", suite.validator.GetOperator().String(), "ve-100", sdk.NewCoin(denom, sdk.NewInt(100))},
		{"invalid coin denomination", false, delAcct.String(), suite.validator.GetOperator().String(), "ve-100", sdk.NewCoin("denom", sdk.NewInt(100))},
		{"ve not owned by delegator", false, delAcct.String(), suite.validator.GetOperator().String(), "ve-100", sdk.NewCoin(denom, sdk.NewInt(100))},
		{"ok", true, delAcct.String(), suite.validator.GetOperator().String(), "ve-1", sdk.NewCoin(denom, sdk.NewInt(100))},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			msg := &types.MsgVeDelegate{
				DelegatorAddress: tc.delegator,
				ValidatorAddress: tc.validator,
				VeId:             tc.veID,
				Amount:           tc.amount,
			}
			_, err := impl.VeDelegate(ctx, msg)
			if tc.pass {
				suite.Require().NoError(err, tc.name)

			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestBeginRedelegate() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	ctx := sdk.WrapSDKContext(suite.ctx)
	delAcct := sdk.AccAddress(suite.address.Bytes())
	denom := "alion"
	impl := keeper.NewMsgServerImpl(k)

	// VeDelegate in advance
	veServer := vekeeper.NewMsgServerImpl(suite.app.VeKeeper)
	veID := "ve-1"
	veTokens := types.VeTokensSlice{
		types.VeTokens{
			VeId:   uint64(1),
			Tokens: sdk.NewInt(100),
		},
	}
	amount := sdk.NewCoin(denom, veTokens[0].Tokens)
	boundAmt := sdk.NewInt(100)
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

	valSrcAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	valDestAddr := "mervaloper1353a4uac03etdylz86tyq9ssm3x2704j6g0p55"
	valSrc, err := sdk.ValAddressFromBech32(valSrcAddr)
	suite.Require().NoError(err)
	valDest, err := sdk.ValAddressFromBech32(valDestAddr)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		pass         bool
		delegator    string
		validatorSrc string
		validatorDst string
		amount       sdk.Coin
	}{
		{"invalid validator src address", false, delAcct.String(), "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjrx", suite.validator.GetOperator().String(), sdk.NewCoin(denom, sdk.NewInt(100))},
		{"invalid delegator address", false, "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg2", suite.validator.GetOperator().String(), suite.validator.GetOperator().String(), sdk.NewCoin(denom, sdk.NewInt(100))},
		{"invalid coin denomination", false, delAcct.String(), valSrc.String(), valDest.String(), sdk.NewCoin("denom", sdk.NewInt(100))},
		{"invalid validator dst address", false, delAcct.String(), valSrc.String(), "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjrx", sdk.NewCoin(denom, sdk.NewInt(100))},
		{"ok", true, delAcct.String(), suite.validator.GetOperator().String(), suite.validatorRedelegate.GetOperator().String(), sdk.NewCoin(denom, sdk.NewInt(100))},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			msg := &stakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    tc.delegator,
				ValidatorSrcAddress: tc.validatorSrc,
				ValidatorDstAddress: tc.validatorDst,
				Amount:              tc.amount,
			}
			_, err := impl.BeginRedelegate(ctx, msg)
			if tc.pass {
				suite.Require().NoError(err, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestUndelegate() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	ctx := sdk.WrapSDKContext(suite.ctx)
	delAcct := sdk.AccAddress(suite.address.Bytes())
	denom := "alion"
	impl := keeper.NewMsgServerImpl(k)

	// VeDelegate in advance
	veServer := vekeeper.NewMsgServerImpl(suite.app.VeKeeper)
	veID := "ve-1"
	veTokens := types.VeTokensSlice{
		types.VeTokens{
			VeId:   uint64(1),
			Tokens: sdk.NewInt(100),
		},
	}
	amount := sdk.NewCoin(denom, veTokens[0].Tokens)
	boundAmt := sdk.NewInt(100)
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

	testCases := []struct {
		name      string
		pass      bool
		delegator string
		validator string
		amount    sdk.Coin
	}{
		{"invalid validator address", false, delAcct.String(), "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjrx", sdk.NewCoin(denom, sdk.NewInt(100))},
		{"invalid delegator address", false, "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg2", suite.validator.GetOperator().String(), sdk.NewCoin(denom, sdk.NewInt(100))},
		{"invalid coin denomination", false, delAcct.String(), suite.validator.GetOperator().String(), sdk.NewCoin("denom", sdk.NewInt(100))},
		{"ok", true, delAcct.String(), suite.validator.GetOperator().String(), sdk.NewCoin(denom, sdk.NewInt(100))},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			msg := &stakingtypes.MsgUndelegate{
				DelegatorAddress: tc.delegator,
				ValidatorAddress: tc.validator,
				Amount:           tc.amount,
			}
			_, err := impl.Undelegate(ctx, msg)
			if tc.pass {
				suite.Require().NoError(err, tc.name)
			} else {
				suite.Require().Error(err, tc.name)
			}
		})
	}
}
