package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/app"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/vesting/types"
	"github.com/stretchr/testify/require"
)

func TestMsgAddAirdrops_ValidateBasic(t *testing.T) {
	app.Setup(false)
	for _, tc := range []struct {
		desc     string
		sender   string
		airdrops []types.Airdrop
		valid    bool
	}{
		{
			desc:   "invalid sender address",
			sender: "",
		},
		{
			desc:   "invalid airdrop target address",
			sender: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			airdrops: []types.Airdrop{
				{},
			},
		},
		{
			desc:   "valid",
			sender: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			airdrops: []types.Airdrop{
				{
					TargetAddr: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
					Amount:     sdk.NewCoin(mertypes.AttoLionDenom, sdk.NewInt(1)),
				},
			},
			valid: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msg := &types.MsgAddAirdrops{
				Sender:   tc.sender,
				Airdrops: tc.airdrops,
			}
			err := msg.ValidateBasic()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgAddAirdrops_GetSigners(t *testing.T) {
	app.Setup(false)
	addr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"
	msg := &types.MsgAddAirdrops{
		Sender:   addr,
		Airdrops: []types.Airdrop{},
	}
	signers := msg.GetSigners()
	require.Equal(t, addr, signers[0].String())
}

func TestMsgExecuteAirdrops_ValidateBasic(t *testing.T) {
	app.Setup(false)
	for _, tc := range []struct {
		desc     string
		sender   string
		maxCount uint64
		valid    bool
	}{
		{
			desc:   "invalid sender address",
			sender: "",
		},
		{
			desc:     "max count must be > 0",
			sender:   "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			maxCount: 0,
		},
		{
			desc:     "valid",
			sender:   "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			maxCount: 1,
			valid:    true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msg := &types.MsgExecuteAirdrops{
				Sender:   tc.sender,
				MaxCount: tc.maxCount,
			}
			err := msg.ValidateBasic()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgExecuteAirdrops_GetSigners(t *testing.T) {
	app.Setup(false)
	addr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"
	msg := &types.MsgExecuteAirdrops{
		Sender: addr,
	}
	signers := msg.GetSigners()
	require.Equal(t, addr, signers[0].String())
}

func TestMsgSetAllocationAddress_ValidateBasic(t *testing.T) {
	app.Setup(false)
	for _, tc := range []struct {
		desc                          string
		sender                        string
		teamVestingAddr               string
		strategicReserveCustodianAddr string
		valid                         bool
	}{
		{
			desc:   "invalid sender address",
			sender: "",
		},
		{
			desc:            "invalid team vesting address",
			sender:          "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			teamVestingAddr: "xxx",
		},
		{
			desc:                          "invalid strategic reserve custodian address",
			sender:                        "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			teamVestingAddr:               "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg1",
			strategicReserveCustodianAddr: "xxx",
		},
		{
			desc:                          "valid",
			sender:                        "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			teamVestingAddr:               "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			strategicReserveCustodianAddr: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			valid:                         true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msg := &types.MsgSetAllocationAddress{
				Sender:                        tc.sender,
				TeamVestingAddr:               tc.teamVestingAddr,
				StrategicReserveCustodianAddr: tc.strategicReserveCustodianAddr,
			}
			err := msg.ValidateBasic()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestMsgSetAllocationAddress_GetSigners(t *testing.T) {
	app.Setup(false)
	addr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"
	msg := &types.MsgSetAllocationAddress{
		Sender: addr,
	}
	signers := msg.GetSigners()
	require.Equal(t, addr, signers[0].String())
}
