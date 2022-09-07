package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/app"
	mertypes "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestMsgVeDelegate_ValidateBasic(t *testing.T) {
	app.Setup(false)
	for _, tc := range []struct {
		desc             string
		delegatorAddress string
		validatorAddress string
		veID             string
		amount           sdk.Coin
		valid            bool
	}{
		{
			desc:             "ErrEmptyDelegatorAddr",
			delegatorAddress: "",
		},
		{
			desc:             "ErrEmptyValidatorAddr",
			delegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			validatorAddress: "",
		},
		{
			desc:             "invalid ve id",
			delegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			validatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
			veID:             "",
		},
		{
			desc:             "invalid delegation amount",
			delegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			validatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
			veID:             "ve-100",
			amount:           sdk.NewCoin(mertypes.AttoLionDenom, sdk.NewInt(0)),
		},
		{
			desc:             "valid",
			delegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
			validatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
			veID:             "ve-100",
			amount:           sdk.NewCoin(mertypes.AttoLionDenom, sdk.NewInt(1)),
			valid:            true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			msg := &types.MsgVeDelegate{
				DelegatorAddress: tc.delegatorAddress,
				ValidatorAddress: tc.validatorAddress,
				VeId:             tc.veID,
				Amount:           tc.amount,
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

func TestMsgVeDelegate_GetSigners(t *testing.T) {
	app.Setup(false)
	msg := &types.MsgVeDelegate{
		DelegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
		ValidatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
		Amount:           sdk.NewCoin(mertypes.AttoLionDenom, sdk.NewInt(1)),
		VeId:             "ve-100",
	}
	signers := msg.GetSigners()
	sender, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	require.NoError(t, err)
	require.Equal(t, sender, signers[0])
}
