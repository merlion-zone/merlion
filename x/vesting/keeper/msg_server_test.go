package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/vesting/keeper"
	"github.com/merlion-zone/merlion/x/vesting/types"
	"github.com/tharsis/ethermint/crypto/ethsecp256k1"
)

func (suite *KeeperTestSuite) TestAddAirdrops() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	k := suite.app.VestingKeeper
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	sender := sdk.AccAddress(priv.PubKey().Address())
	receiver := sdk.AccAddress(suite.address.Bytes())
	denom := mertypes.AttoLionDenom
	k.SetAllocationAddresses(suite.ctx, types.AllocationAddresses{
		TeamVestingAddr:               receiver.String(),
		StrategicReserveCustodianAddr: receiver.String(),
	})
	teamAddr := receiver
	impl := keeper.NewMsgServerImpl(k)
	cap := k.GetParams(suite.ctx).Allocation.AirdropAmount.Add(sdk.NewInt(1))
	testCases := []struct {
		name     string
		pass     bool
		sender   string
		airdrops []types.Airdrop
	}{
		{"invalid sender", false, "xxx", []types.Airdrop{}},
		{"unauthorized sender", false, sender.String(), []types.Airdrop{}},
		{"invalid targetAddr", false, teamAddr.String(), []types.Airdrop{
			{TargetAddr: "", Amount: sdk.NewCoin(denom, sdk.NewInt(1))},
		}},
		{"total amount should not be greater than its cap", false, teamAddr.String(), []types.Airdrop{
			{TargetAddr: receiver.String(), Amount: sdk.NewCoin(denom, cap)},
		}},
		{"valid", true, teamAddr.String(), []types.Airdrop{
			{TargetAddr: receiver.String(), Amount: sdk.NewCoin(denom, sdk.NewInt(1))},
		}},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			res, err := impl.AddAirdrops(ctx, &types.MsgAddAirdrops{
				Sender:   tc.sender,
				Airdrops: tc.airdrops,
			})
			if tc.pass {
				require.NoError(err, tc.name)
				expected := sdk.NewInt(0)
				for _, airdrop := range tc.airdrops {
					amount, err := sdk.ParseCoinNormalized(airdrop.Amount.String())
					require.NoError(err)
					expected = expected.Add(amount.Amount)
				}
				actual := k.GetAirdropTotalAmount(suite.ctx)
				require.Equal(expected, actual)
			} else {
				require.Error(err, tc.name)
				require.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestExecuteAirdrops() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	k := suite.app.VestingKeeper
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	sender := sdk.AccAddress(priv.PubKey().Address())
	receiver := sdk.AccAddress(suite.address.Bytes())
	denom := mertypes.AttoLionDenom
	k.SetAllocationAddresses(suite.ctx, types.AllocationAddresses{
		TeamVestingAddr:               receiver.String(),
		StrategicReserveCustodianAddr: receiver.String(),
	})
	teamAddr := receiver
	impl := keeper.NewMsgServerImpl(k)
	testCases := []struct {
		name     string
		pass     bool
		sender   string
		airdrops []types.Airdrop
	}{
		{"invalid sender", false, "xxx", []types.Airdrop{}},
		{"unauthorized sender", false, sender.String(), []types.Airdrop{}},
		{"valid", true, teamAddr.String(), []types.Airdrop{
			{TargetAddr: sender.String(), Amount: sdk.NewCoin(denom, sdk.NewInt(100))},
		}},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			res, err := impl.AddAirdrops(ctx, &types.MsgAddAirdrops{
				Sender:   tc.sender,
				Airdrops: tc.airdrops,
			})
			if tc.pass {
				_, err := impl.ExecuteAirdrops(ctx, &types.MsgExecuteAirdrops{
					Sender:   tc.sender,
					MaxCount: 100,
				})
				require.NoError(err, tc.name)
			} else {
				require.Error(err, tc.name)
				require.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSetAllocationAddress() {
	require := suite.Require()
	ctx := sdk.WrapSDKContext(suite.ctx)
	k := suite.app.VestingKeeper
	priv, err := ethsecp256k1.GenerateKey()
	require.NoError(err)
	receiver := sdk.AccAddress(priv.PubKey().Address())
	sender := sdk.AccAddress(suite.address.Bytes())
	impl := keeper.NewMsgServerImpl(k)
	testCases := []struct {
		name                          string
		pass                          bool
		sender                        string
		teamVestingAddr               string
		strategicReserveCustodianAddr string
	}{
		{"invalid sender", false, "xxx", sender.String(), sender.String()},
		{"only one allocation address must be given", false, sender.String(), sender.String(), sender.String()},
		{"only one allocation address must be given", false, sender.String(), "", ""},
		{"invalid StrategicReserveCustodianAddr", false, receiver.String(), "", "xxx"},
		{"unauthorized sender", false, receiver.String(), "", sender.String()},
		{"invalid TeamVestingAddr", false, receiver.String(), "xxx", ""},
		{"unauthorized sender", false, receiver.String(), sender.String(), ""},
		{"ok", true, receiver.String(), sender.String(), ""},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.name), func() {
			if tc.pass {
				k.SetAllocationAddresses(suite.ctx, types.AllocationAddresses{
					TeamVestingAddr: receiver.String(),
				})
			}
			res, err := impl.SetAllocationAddress(ctx, &types.MsgSetAllocationAddress{
				Sender:                        tc.sender,
				TeamVestingAddr:               tc.teamVestingAddr,
				StrategicReserveCustodianAddr: tc.strategicReserveCustodianAddr,
			})
			if tc.pass {
				require.NoError(err, tc.name)
			} else {
				require.Error(err, tc.name)
				require.Nil(res)
			}
		})
	}
}
