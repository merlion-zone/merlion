package types_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/app"
	"github.com/merlion-zone/merlion/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestGetVeValidatorKey(t *testing.T) {
	app.Setup(false)
	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	val, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)
	key := types.GetVeValidatorKey(val)
	require.Equal(t, "a114dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea1", hex.EncodeToString(key))
}

func TestGetVeDelegationKey(t *testing.T) {
	app.Setup(false)
	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	key := types.GetVeDelegationKey(acc, val)
	require.Equal(t, "a214dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea114dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea1", hex.EncodeToString(key))
}

func TestGetVeDelegationsKey(t *testing.T) {
	app.Setup(false)
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	key := types.GetVeDelegationsKey(acc)
	require.Equal(t, "a214dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea1", hex.EncodeToString(key))
}

func TestGetVeUBDKey(t *testing.T) {
	app.Setup(false)
	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	key := types.GetVeUBDKey(acc, val)
	require.Equal(t, "a314dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea114dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea1", hex.EncodeToString(key))
}

func TestGetVeUBDsKey(t *testing.T) {
	app.Setup(false)
	valSrcAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	valDestAddr := "mervaloper1353a4uac03etdylz86tyq9ssm3x2704j6g0p55"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	valSrc, err := sdk.ValAddressFromBech32(valSrcAddr)
	require.NoError(t, err)

	valDest, err := sdk.ValAddressFromBech32(valDestAddr)
	require.NoError(t, err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	key := types.GetVeREDKey(acc, valSrc, valDest)
	require.Equal(t, "a414dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea114dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea1148d23daf3b87c72b693e23e96401610dc4caf3eb2", hex.EncodeToString(key))
}

func TestGetVeREDsKey(t *testing.T) {
	app.Setup(false)
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	key := types.GetVeREDsKey(acc)
	require.Equal(t, "a414dcd3b2e3d86a013b5b5a823b30f8fb791bbc0ea1", hex.EncodeToString(key))
}

func TestGetVeTokensKey(t *testing.T) {
	app.Setup(false)
	veID := uint64(100)

	key := types.GetVeTokensKey(veID)
	require.Equal(t, "a50000000000000064", hex.EncodeToString(key))
}
