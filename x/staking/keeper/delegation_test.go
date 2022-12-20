package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/types"
	vekeeper "github.com/merlion-zone/merlion/x/ve/keeper"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

func (suite *KeeperTestSuite) TestKeeper_SetVeValidator_GetVeValidator() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	valAddr := sdk.ValAddress(suite.address.Bytes())
	validator, found := k.GetVeValidator(suite.ctx, valAddr)
	suite.Require().Equal(types.VeValidator{}, validator)
	suite.Require().Equal(false, found)

	v := types.VeValidator{
		OperatorAddress:   valAddr.String(),
		VeDelegatorShares: sdk.NewDec(5),
	}
	k.SetVeValidator(suite.ctx, v)

	validator, found = k.GetVeValidator(suite.ctx, valAddr)
	suite.Require().Equal(v, validator)
	suite.Require().Equal(true, found)
}

func (suite *KeeperTestSuite) TestKeeper_RemoveVeValidator() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	valAddr := sdk.ValAddress(suite.address.Bytes())
	v := types.VeValidator{
		OperatorAddress:   valAddr.String(),
		VeDelegatorShares: sdk.NewDec(5),
	}
	k.SetVeValidator(suite.ctx, v)
	k.RemoveVeValidator(suite.ctx, valAddr)

	validator, found := k.GetVeValidator(suite.ctx, valAddr)
	suite.Require().Equal(types.VeValidator{}, validator)
	suite.Require().Equal(false, found)
}

func (suite *KeeperTestSuite) TestKeeper_SetVeDelegation_GetVeDelegation() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	delegation, found := k.GetVeDelegation(suite.ctx, acc, val)
	suite.Require().Equal(types.VeDelegation{}, delegation)
	suite.Require().Equal(false, found)

	d := types.VeDelegation{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		VeShares:         []types.VeShares{},
	}
	k.SetVeDelegation(suite.ctx, d)

	delegation, found = k.GetVeDelegation(suite.ctx, acc, val)
	suite.Require().Equal(d.DelegatorAddress, delegation.DelegatorAddress)
	suite.Require().Equal(d.ValidatorAddress, delegation.ValidatorAddress)
	suite.Require().Equal(true, found)
}

func (suite *KeeperTestSuite) TestKeeper_RemoveVeDelegation() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	d := types.VeDelegation{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		VeShares:         []types.VeShares{},
	}
	k.SetVeDelegation(suite.ctx, d)
	k.RemoveVeDelegation(suite.ctx, d)
	delegation, found := k.GetVeDelegation(suite.ctx, acc, val)
	suite.Require().Equal(types.VeDelegation{}, delegation)
	suite.Require().Equal(false, found)
}

func (suite *KeeperTestSuite) TestKeeper_SetVeDelegatedAmount_GetVeDelegatedAmount() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	veID := uint64(100)
	amount := k.GetVeDelegatedAmount(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), amount)

	amt := sdk.NewInt(100)
	k.SetVeDelegatedAmount(suite.ctx, veID, amt)

	amount = k.GetVeDelegatedAmount(suite.ctx, veID)
	suite.Require().Equal(amt, amount)
}

func (suite *KeeperTestSuite) TestKeeper_RemoveVeDelegatedAmount() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	veID := uint64(100)
	amt := sdk.NewInt(100)
	k.SetVeDelegatedAmount(suite.ctx, veID, amt)
	k.RemoveVeDelegatedAmount(suite.ctx, veID)

	amount := k.GetVeDelegatedAmount(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), amount)
}

func (suite *KeeperTestSuite) TestKeeper_SubVeDelegatedAmount() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	veID := uint64(100)
	amt := sdk.NewInt(100)
	k.SubVeDelegatedAmount(suite.ctx, veID, amt)

	amount := k.GetVeDelegatedAmount(suite.ctx, veID)
	suite.Require().Equal(sdk.ZeroInt(), amount)

	k.SetVeDelegatedAmount(suite.ctx, veID, amt)
	subAmt := sdk.NewInt(20)
	k.SubVeDelegatedAmount(suite.ctx, veID, subAmt)

	amount = k.GetVeDelegatedAmount(suite.ctx, veID)
	suite.Require().Equal(amt.Sub(subAmt), amount)
}

func (suite *KeeperTestSuite) TestKeeper_SlashVeDelegatedAmount() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	coin := sdk.NewCoin(suite.app.VeKeeper.LockDenom(suite.ctx), sdk.NewInt(1000))
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, sdk.AccAddress(suite.address.Bytes()), vetypes.ModuleName, sdk.NewCoins(coin))
	suite.Require().NoError(err)

	veID := uint64(100)
	amt := sdk.NewInt(100)

	k.SetVeDelegatedAmount(suite.ctx, veID, amt)
	subAmt := sdk.NewInt(20)
	k.SlashVeDelegatedAmount(suite.ctx, veID, subAmt)

	amount := k.GetVeDelegatedAmount(suite.ctx, veID)
	suite.Require().Equal(amt.Sub(subAmt), amount)
}

func (suite *KeeperTestSuite) TestKeeper_SettleVeDelegation() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	coin := sdk.NewCoin(suite.app.VeKeeper.LockDenom(suite.ctx), sdk.NewInt(1000))
	err := suite.app.BankKeeper.SendCoinsFromAccountToModule(suite.ctx, sdk.AccAddress(suite.address.Bytes()), vetypes.ModuleName, sdk.NewCoins(coin))
	suite.Require().NoError(err)

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"
	delegation := types.VeDelegation{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
		VeShares: []types.VeShares{
			{
				VeId:               uint64(100),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(1),
			},
		},
	}
	validator := suite.validator
	validator.DelegatorShares = sdk.NewDec(10)
	delegation = k.SettleVeDelegation(suite.ctx, delegation, validator)
	suite.Require().Equal(sdk.ZeroInt(), delegation.VeShares[0].TokensMayUnsettled)

	amount := k.GetVeDelegatedAmount(suite.ctx, delegation.VeShares[0].VeId)
	suite.Require().Equal(sdk.ZeroInt(), amount)
}

func (suite *KeeperTestSuite) TestKeeper_SetVeUnbondingDelegation_GetVeUnbondingDelegation() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	ubd, found := k.GetVeUnbondingDelegation(suite.ctx, acc, val)
	suite.Require().Equal(types.VeUnbondingDelegation{}, ubd)
	suite.Require().Equal(false, found)

	unbond := types.VeUnbondingDelegation{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
	}
	k.SetVeUnbondingDelegation(suite.ctx, unbond)
	ubd, found = k.GetVeUnbondingDelegation(suite.ctx, acc, val)
	suite.Require().Equal(unbond, ubd)
	suite.Require().Equal(true, found)
}

func (suite *KeeperTestSuite) TestKeeper_RemoveVeUnbondingDelegation() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	unbond := types.VeUnbondingDelegation{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr,
	}
	k.SetVeUnbondingDelegation(suite.ctx, unbond)
	k.RemoveVeUnbondingDelegation(suite.ctx, unbond)

	ubd, found := k.GetVeUnbondingDelegation(suite.ctx, acc, val)
	suite.Require().Equal(types.VeUnbondingDelegation{}, ubd)
	suite.Require().Equal(false, found)
}

func (suite *KeeperTestSuite) TestKeeper_SetVeUnbondingDelegationEntry() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	veTokens := types.VeTokensSlice{
		types.VeTokens{
			VeId:   uint64(100),
			Tokens: sdk.NewInt(100),
		},
	}
	unbond := k.SetVeUnbondingDelegationEntry(suite.ctx, acc, val, veTokens)
	ubd, found := k.GetVeUnbondingDelegation(suite.ctx, acc, val)
	suite.Require().Equal(unbond, ubd)
	suite.Require().Equal(true, found)
}

func (suite *KeeperTestSuite) TestKeeper_SetVeRedelegation_GetVeRedelegation() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valSrcAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	valDestAddr := "mervaloper1353a4uac03etdylz86tyq9ssm3x2704j6g0p55"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	valSrc, err := sdk.ValAddressFromBech32(valSrcAddr)
	suite.Require().NoError(err)

	valDest, err := sdk.ValAddressFromBech32(valDestAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	red, found := k.GetVeRedelegation(suite.ctx, acc, valSrc, valDest)
	suite.Require().Equal(types.VeRedelegation{}, red)
	suite.Require().Equal(false, found)

	delegation := types.VeRedelegation{
		DelegatorAddress:    delAddr,
		ValidatorSrcAddress: valSrcAddr,
		ValidatorDstAddress: valDestAddr,
	}
	k.SetVeRedelegation(suite.ctx, delegation)
	red, found = k.GetVeRedelegation(suite.ctx, acc, valSrc, valDest)
	suite.Require().Equal(delegation, red)
	suite.Require().Equal(true, found)
}

func (suite *KeeperTestSuite) TestKeeper_RemoveVeRedelegation() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valSrcAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	valDestAddr := "mervaloper1353a4uac03etdylz86tyq9ssm3x2704j6g0p55"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	valSrc, err := sdk.ValAddressFromBech32(valSrcAddr)
	suite.Require().NoError(err)

	valDest, err := sdk.ValAddressFromBech32(valDestAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	delegation := types.VeRedelegation{
		DelegatorAddress:    delAddr,
		ValidatorSrcAddress: valSrcAddr,
		ValidatorDstAddress: valDestAddr,
	}
	k.SetVeRedelegation(suite.ctx, delegation)
	k.RemoveVeRedelegation(suite.ctx, delegation)

	red, found := k.GetVeRedelegation(suite.ctx, acc, valSrc, valDest)
	suite.Require().Equal(types.VeRedelegation{}, red)
	suite.Require().Equal(false, found)
}

func (suite *KeeperTestSuite) TestKeeper_SetVeRedelegationEntry() {
	suite.SetupTest()
	k := suite.app.StakingKeeper

	valSrcAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	valDestAddr := "mervaloper1353a4uac03etdylz86tyq9ssm3x2704j6g0p55"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	valSrc, err := sdk.ValAddressFromBech32(valSrcAddr)
	suite.Require().NoError(err)

	valDest, err := sdk.ValAddressFromBech32(valDestAddr)
	suite.Require().NoError(err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	suite.Require().NoError(err)

	totalAmount := sdk.NewInt(100)
	veTokens := types.VeTokensSlice{
		types.VeTokens{
			VeId:   uint64(100),
			Tokens: sdk.NewInt(100),
		},
	}
	totalShares := sdk.NewDec(100)
	k.SetVeRedelegationEntry(suite.ctx, acc, valSrc, valDest, totalAmount, veTokens, totalShares)

	red, found := k.GetVeRedelegation(suite.ctx, acc, valSrc, valDest)
	suite.Require().Equal(veTokens[0].Tokens, red.Entries[0].VeShares[0].InitialBalance)
	suite.Require().Equal(true, found)
}

func (suite *KeeperTestSuite) TestKeeper_VeDelegate() {
	suite.SetupTest()
	k := suite.app.StakingKeeper
	delAcct := sdk.AccAddress(suite.address.Bytes())
	ctx := sdk.WrapSDKContext(suite.ctx)
	denom := k.BondDenom(suite.ctx)
	veServer := vekeeper.NewMsgServerImpl(suite.app.VeKeeper)

	testCases := []struct {
		name            string
		delAddr         sdk.AccAddress
		bondAmt         sdk.Int
		veTokens        types.VeTokensSlice
		tokenSrc        stakingtypes.BondStatus
		validator       stakingtypes.Validator
		subtractAccount bool
		shares          sdk.Dec
	}{
		{"not subtract account", delAcct, sdk.NewInt(100), types.VeTokensSlice{
			types.VeTokens{
				VeId:   uint64(1),
				Tokens: sdk.NewInt(100),
			},
		}, stakingtypes.Bonded, suite.validator, false, sdk.NewDec(100)},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			amount := sdk.NewCoin(denom, tc.veTokens[0].Tokens)
			veID := fmt.Sprintf("ve-%d", tc.veTokens[0].VeId)
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

			newShares, err := k.VeDelegate(suite.ctx, tc.delAddr, tc.bondAmt, tc.veTokens, tc.tokenSrc, tc.validator, tc.subtractAccount)
			suite.Require().NoError(err)
			suite.Require().Equal(tc.shares, newShares)
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_BeginRedelegation() {
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

	valSrcAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	valDestAddr := "mervaloper1353a4uac03etdylz86tyq9ssm3x2704j6g0p55"
	valSrc, err := sdk.ValAddressFromBech32(valSrcAddr)
	suite.Require().NoError(err)
	valDest, err := sdk.ValAddressFromBech32(valDestAddr)
	suite.Require().NoError(err)

	testCases := []struct {
		name         string
		delAddr      sdk.AccAddress
		valSrcAddr   sdk.ValAddress
		valDstAddr   sdk.ValAddress
		sharesAmount sdk.Dec
		err          error
	}{
		{"ErrSelfRedelegation", delAcct, suite.validator.GetOperator(), suite.validator.GetOperator(), sdk.NewDec(10), stakingtypes.ErrSelfRedelegation},
		{"ErrBadRedelegationSrc", delAcct, valSrc, suite.validator.GetOperator(), sdk.NewDec(10), stakingtypes.ErrBadRedelegationDst},
		{"ErrBadRedelegationDst", delAcct, suite.validator.GetOperator(), valDest, sdk.NewDec(10), stakingtypes.ErrBadRedelegationDst},
		{"ErrTinyRedelegationAmount", delAcct, suite.validator.GetOperator(), suite.validatorRedelegate.GetOperator(), sdk.NewDec(0), stakingtypes.ErrTinyRedelegationAmount},
		{"ok", delAcct, suite.validator.GetOperator(), suite.validatorRedelegate.GetOperator(), sdk.NewDec(10), nil},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			_, err := k.BeginRedelegation(suite.ctx, tc.delAddr, tc.valSrcAddr, tc.valDstAddr, tc.sharesAmount)
			if err != nil {
				suite.Require().Equal(tc.err, err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_Undelegate() {
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
