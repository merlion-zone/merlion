package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/x/ve/keeper"
	"github.com/merlion-zone/merlion/x/ve/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

func (suite *KeeperTestSuite) TestVeCreate() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	receiver := sdk.AccAddress(priv.PubKey().Address())
	sender := sdk.AccAddress(suite.address.Bytes())
	denom := "1alion"
	impl := keeper.NewMsgServerImpl(suite.app.VeKeeper)
	testCases := []struct {
		name         string
		pass         bool
		sender       sdk.AccAddress
		to           string
		amount       string
		lockDuration uint64
	}{
		{"receiver is different from sender", true, sender, receiver.String(), denom, types.RegulatedPeriod},
		{"invalid denom", false, sender, receiver.String(), "1lion", types.RegulatedPeriod},
		{"invalid sender", false, []byte("xxx"), receiver.String(), denom, types.RegulatedPeriod},
		{"invalid receiver", false, sender, "xxx", "1alion", types.RegulatedPeriod},
		{"unlock time is past time", false, sender, receiver.String(), denom, 0},
		{"unlock time is future time", false, sender, receiver.String(), denom, 10000 * types.MaxLockTime},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			sender := tc.sender
			amount, _ := sdk.ParseCoinNormalized(tc.amount)

			res, err := impl.Create(ctx, &types.MsgCreate{
				Sender:       sender.String(),
				To:           tc.to,
				Amount:       amount,
				LockDuration: tc.lockDuration,
			})
			if tc.pass {
				require.NoError(err, tc.name)

				expectedUnlockTime := types.RegulatedUnixTimeFromNow(suite.ctx, tc.lockDuration)
				require.Equal(expectedUnlockTime, res.UnlockTime)
			} else {
				require.Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVeDeposit() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	sender := sdk.AccAddress(suite.address.Bytes())
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	receiver := sdk.AccAddress(priv.PubKey().Address())
	impl := keeper.NewMsgServerImpl(suite.app.VeKeeper)
	denom := "alion"
	res, err := impl.Create(ctx, &types.MsgCreate{
		Sender:       sender.String(),
		To:           receiver.String(),
		Amount:       sdk.NewCoin(denom, sdk.NewInt(1)),
		LockDuration: types.RegulatedPeriod,
	})
	require.NoError(err)
	require.Equal("ve-1", res.VeId)

	testCases := []struct {
		name   string
		pass   bool
		sender sdk.AccAddress
		amount sdk.Coin
		veID   string
	}{
		{"invalid denom", false, sender, sdk.NewCoin("lion", sdk.NewInt(1)), "ve-100"},
		{"invalid sender", false, []byte("xxx"), sdk.NewCoin(denom, sdk.NewInt(1)), "ve-100"},
		{"invalid ve id", false, sender, sdk.NewCoin(denom, sdk.NewInt(1)), "ve-100"},
		{"nothing is locked for ve", false, sender, sdk.NewCoin(denom, sdk.NewInt(1)), "ve-200"},
		{"ok", true, sender, sdk.NewCoin(denom, sdk.NewInt(1)), "ve-1"},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			sender := tc.sender
			res, err := impl.Deposit(ctx, &types.MsgDeposit{
				Sender: sender.String(),
				VeId:   tc.veID,
				Amount: tc.amount,
			})
			if tc.pass {
				require.NoError(err, tc.name)
				require.NotNil(res)
			} else {
				require.Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVeExtendTime() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	impl := keeper.NewMsgServerImpl(suite.app.VeKeeper)
	sender := sdk.AccAddress(suite.address.Bytes())
	denom := "alion"
	// Create Valid VeID
	for i := 1; i <= 3; i++ {
		res, err := impl.Create(ctx, &types.MsgCreate{
			Sender:       sender.String(),
			To:           sender.String(),
			Amount:       sdk.NewCoin(denom, sdk.NewInt(1)),
			LockDuration: types.RegulatedPeriod,
		})
		require.NoError(err)
		require.Equal(fmt.Sprintf("ve-%d", i), res.VeId)
	}

	// Deposit for ve-3
	res, err := impl.Deposit(ctx, &types.MsgDeposit{
		Sender: sender.String(),
		VeId:   "ve-3",
		Amount: sdk.NewCoin(denom, sdk.NewInt(1)),
	})
	require.NoError(err)
	require.NotNil(res)

	testCases := []struct {
		name         string
		pass         bool
		sender       sdk.AccAddress
		veID         string
		lockDuration uint64
	}{
		{"invalid sender", false, []byte("xxx"), "ve-100", types.RegulatedPeriod},
		{"invalid ve id", false, sender, "ve-100", types.RegulatedPeriod},
		{"nothing is locked for ve", false, sender, "ve-2", types.RegulatedPeriod},
		{"not increased lock time", false, sender, "ve-3", 0},
		{"too long lock time", false, sender, "ve-3", 10000 * types.RegulatedPeriod},
		{"valid", true, sender, "ve-3", 2 * types.RegulatedPeriod},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			sender := tc.sender
			res, err := impl.ExtendTime(ctx, &types.MsgExtendTime{
				Sender:       sender.String(),
				VeId:         tc.veID,
				LockDuration: tc.lockDuration,
			})
			if tc.pass {
				require.NoError(err, tc.name)
				require.NotNil(res)
			} else {
				require.Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVeMerge() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	impl := keeper.NewMsgServerImpl(suite.app.VeKeeper)
	sender := sdk.AccAddress(suite.address.Bytes())
	denom := "alion"
	// Create Valid VeID
	for i := 1; i <= 2; i++ {
		res, err := impl.Create(ctx, &types.MsgCreate{
			Sender:       sender.String(),
			To:           sender.String(),
			Amount:       sdk.NewCoin(denom, sdk.NewInt(1)),
			LockDuration: types.RegulatedPeriod,
		})
		require.NoError(err)
		require.Equal(fmt.Sprintf("ve-%d", i), res.VeId)

		// Deposit
		_, err = impl.Deposit(ctx, &types.MsgDeposit{
			Sender: sender.String(),
			VeId:   res.VeId,
			Amount: sdk.NewCoin(denom, sdk.NewInt(1)),
		})
		require.NoError(err)
	}

	// Another NFT
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	receiver := sdk.AccAddress(priv.PubKey().Address())
	_, err = impl.Create(ctx, &types.MsgCreate{
		Sender:       sender.String(),
		To:           receiver.String(),
		Amount:       sdk.NewCoin(denom, sdk.NewInt(1)),
		LockDuration: types.RegulatedPeriod,
	})
	require.NoError(err)

	testCases := []struct {
		name     string
		pass     bool
		sender   sdk.AccAddress
		fromVeId string
		toVeId   string
	}{
		{"invalid sender", false, []byte("xxx"), "ve-1", "ve-2"},
		{"user doesn't own fromVeId", false, sender, "ve-3", "ve-2"},
		{"user doesn't own toVeId", false, sender, "ve-1", "ve-3"},
		{"ok", true, sender, "ve-1", "ve-2"},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			sender := tc.sender
			res, err := impl.Merge(ctx, &types.MsgMerge{
				Sender:   sender.String(),
				FromVeId: tc.fromVeId,
				ToVeId:   tc.toVeId,
			})
			if tc.pass {
				require.NoError(err, tc.name)
				require.NotNil(res)
			} else {
				require.Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestVeWithdraw() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	impl := keeper.NewMsgServerImpl(suite.app.VeKeeper)
	sender := sdk.AccAddress(suite.address.Bytes())
	denom := "alion"
	// Create Valid VeID
	for i := 1; i <= 2; i++ {
		res, err := impl.Create(ctx, &types.MsgCreate{
			Sender:       sender.String(),
			To:           sender.String(),
			Amount:       sdk.NewCoin(denom, sdk.NewInt(1)),
			LockDuration: types.RegulatedPeriod,
		})
		require.NoError(err)
		require.Equal(fmt.Sprintf("ve-%d", i), res.VeId)

		// Deposit
		_, err = impl.Deposit(ctx, &types.MsgDeposit{
			Sender: sender.String(),
			VeId:   res.VeId,
			Amount: sdk.NewCoin(denom, sdk.NewInt(1)),
		})
		require.NoError(err)
	}

	// Another NFT
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	receiver := sdk.AccAddress(priv.PubKey().Address())
	_, err = impl.Create(ctx, &types.MsgCreate{
		Sender:       sender.String(),
		To:           receiver.String(),
		Amount:       sdk.NewCoin(denom, sdk.NewInt(1)),
		LockDuration: types.RegulatedPeriod,
	})
	require.NoError(err)

	testCases := []struct {
		name   string
		pass   bool
		sender sdk.AccAddress
		veId   string
	}{
		{"invalid sender", false, []byte("xxx"), "ve-1"},
		{"user doesn't own veId", false, sender, "ve-3"},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			sender := tc.sender
			res, err := impl.Withdraw(ctx, &types.MsgWithdraw{
				Sender: sender.String(),
				VeId:   tc.veId,
			})
			if tc.pass {
				require.NoError(err, tc.name)
				require.NotNil(res)
			} else {
				require.Error(err, tc.name)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestKeeper_DepositFor() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	impl := keeper.NewMsgServerImpl(suite.app.VeKeeper)
	sender := sdk.AccAddress(suite.address.Bytes())
	denom := "alion"
	// Create Valid VeID
	for i := 1; i <= 2; i++ {
		res, err := impl.Create(ctx, &types.MsgCreate{
			Sender:       sender.String(),
			To:           sender.String(),
			Amount:       sdk.NewCoin(denom, sdk.NewInt(1)),
			LockDuration: types.RegulatedPeriod,
		})
		require.NoError(err)
		require.Equal(fmt.Sprintf("ve-%d", i), res.VeId)

		// Deposit
		_, err = impl.Deposit(ctx, &types.MsgDeposit{
			Sender: sender.String(),
			VeId:   res.VeId,
			Amount: sdk.NewCoin(denom, sdk.NewInt(1)),
		})
		require.NoError(err)
	}

	testCases := []struct {
		name      string
		sendCoins bool
		amount    sdk.Int
	}{
		{"does not send coins", false, sdk.NewInt(1)},
		{"do send coins", true, sdk.NewInt(1)},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			err := suite.app.VeKeeper.DepositFor(suite.ctx, sender, uint64(1),
				tc.amount, 0, types.NewLockedBalance(), tc.sendCoins)
			require.NoError(err, tc.name)
		})
	}
}
