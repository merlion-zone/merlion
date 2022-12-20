package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/merlion-zone/merlion/app"
	"github.com/merlion-zone/merlion/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestVeDelegation_GetSharesByVeID(t *testing.T) {
	app.Setup(false)

	veDelegation := types.VeDelegation{
		DelegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
		ValidatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
		VeShares: []types.VeShares{
			{
				VeId:               uint64(100),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(9),
			},
			{
				VeId:               uint64(101),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(1),
			},
		},
	}

	s1, ok := veDelegation.GetSharesByVeID(uint64(100))
	require.Equal(t, veDelegation.VeShares[0], s1)
	require.Equal(t, true, ok)

	s2, ok := veDelegation.GetSharesByVeID(uint64(102))
	require.Equal(t, uint64(102), s2.VeId)
	require.Equal(t, sdk.ZeroInt(), s2.TokensMayUnsettled)
	require.Equal(t, sdk.ZeroDec(), s2.Shares)
	require.Equal(t, false, ok)
}

func TestVeDelegation_SetSharesByVeID(t *testing.T) {
	app.Setup(false)

	veDelegation := types.VeDelegation{
		DelegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
		ValidatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
		VeShares: []types.VeShares{
			{
				VeId:               uint64(100),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(9),
			},
			{
				VeId:               uint64(101),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(1),
			},
		},
	}

	s1 := types.VeShares{
		VeId:               uint64(100),
		TokensMayUnsettled: sdk.NewInt(100),
		Shares:             sdk.NewDec(7),
	}
	new := veDelegation.SetSharesByVeID(s1)
	require.Equal(t, veDelegation.VeShares[0], s1)
	require.Equal(t, false, new)

	s2 := types.VeShares{
		VeId:               uint64(102),
		TokensMayUnsettled: sdk.NewInt(100),
		Shares:             sdk.NewDec(7),
	}
	new = veDelegation.SetSharesByVeID(s2)
	require.Equal(t, veDelegation.VeShares[2], s2)
	require.Equal(t, true, new)
}

func TestVeDelegation_RemoveSharesByIndex(t *testing.T) {
	app.Setup(false)

	veDelegation := types.VeDelegation{
		DelegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
		ValidatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
		VeShares: []types.VeShares{
			{
				VeId:               uint64(100),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(9),
			},
			{
				VeId:               uint64(101),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(1),
			},
		},
	}

	veDelegation.RemoveSharesByIndex(1)
	require.Equal(t, len(veDelegation.VeShares), 1)
}

func TestVeDelegation_Shares(t *testing.T) {
	app.Setup(false)

	veDelegation := types.VeDelegation{
		DelegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
		ValidatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
		VeShares: []types.VeShares{
			{
				VeId:               uint64(100),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(9),
			},
			{
				VeId:               uint64(101),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(1),
			},
		},
	}

	shares := veDelegation.Shares()
	require.Equal(t, shares, sdk.NewDec(10))
}

func TestVeDelegation_Tokens(t *testing.T) {
	app.Setup(false)

	veDelegation := types.VeDelegation{
		DelegatorAddress: "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5",
		ValidatorAddress: "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3",
		VeShares: []types.VeShares{
			{
				VeId:               uint64(100),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(9),
			},
			{
				VeId:               uint64(101),
				TokensMayUnsettled: sdk.NewInt(100),
				Shares:             sdk.NewDec(1),
			},
		},
	}

	tokens := veDelegation.Tokens()
	require.Equal(t, tokens, sdk.NewInt(200))
}

func TestVeShares_Tokens(t *testing.T) {
	app.Setup(false)

	share := types.VeShares{
		VeId:               uint64(100),
		TokensMayUnsettled: sdk.NewInt(100),
		Shares:             sdk.NewDec(9),
	}

	tokens := share.Tokens()
	require.Equal(t, sdk.NewInt(100), tokens)
}

func TestVeShares_AddTokensAndShares(t *testing.T) {
	app.Setup(false)

	share := types.VeShares{
		VeId:               uint64(100),
		TokensMayUnsettled: sdk.NewInt(100),
		Shares:             sdk.NewDec(9),
	}

	share.AddTokensAndShares(sdk.NewInt(100), sdk.NewDec(1))
	require.Equal(t, sdk.NewInt(200), share.Tokens())
	require.Equal(t, sdk.NewDec(10), share.Shares)
}

func TestNewVeUnbondingDelegation(t *testing.T) {
	app.Setup(false)

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	delegation := types.NewVeUnbondingDelegation(acc, val)
	require.Equal(t, delAddr, delegation.DelegatorAddress)
	require.Equal(t, valAddr, delegation.ValidatorAddress)
}

func TestVeUnbondingDelegation_AddEntry(t *testing.T) {
	app.Setup(false)

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	delegation := types.NewVeUnbondingDelegation(acc, val)
	veTokens := types.VeTokensSlice{
		{VeId: uint64(100), Tokens: sdk.NewInt(100)},
		{VeId: uint64(200), Tokens: sdk.NewInt(200)},
	}
	delegation.AddEntry(veTokens)

	entries := delegation.Entries
	l := len(entries)
	for i := 0; i < l; i++ {
		require.Equal(t, veTokens[i].VeId, entries[0].VeBalances[i].VeId)
		require.Equal(t, veTokens[i].Tokens, entries[0].VeBalances[i].InitialBalance)
		require.Equal(t, veTokens[i].Tokens, entries[0].VeBalances[i].Balance)
	}
}

func TestVeUnbondingDelegation_RemoveEntry(t *testing.T) {
	app.Setup(false)

	valAddr := "mervaloper1mnfm9c7cdgqnkk66sganp78m0ydmcr4pctrjr3"
	delAddr := "mer1mnfm9c7cdgqnkk66sganp78m0ydmcr4ppeaeg5"

	val, err := sdk.ValAddressFromBech32(valAddr)
	require.NoError(t, err)

	acc, err := sdk.AccAddressFromBech32(delAddr)
	require.NoError(t, err)

	delegation := types.NewVeUnbondingDelegation(acc, val)
	veTokens := types.VeTokensSlice{
		{VeId: uint64(100), Tokens: sdk.NewInt(100)},
		{VeId: uint64(200), Tokens: sdk.NewInt(200)},
	}
	delegation.AddEntry(veTokens)
	delegation.RemoveEntry(0)
	entries := delegation.Entries
	require.Equal(t, 0, len(entries))
}

func TestVeUnbondingDelegationEntry_Balance(t *testing.T) {
	app.Setup(false)

	veBalances := []types.VeUnbondingDelegationEntryBalances{
		{
			VeId:           uint64(100),
			InitialBalance: sdk.NewInt(100),
			Balance:        sdk.NewInt(100),
		},
		{
			VeId:           uint64(200),
			InitialBalance: sdk.NewInt(100),
			Balance:        sdk.NewInt(100),
		},
	}
	entry := types.VeUnbondingDelegationEntry{
		VeBalances: veBalances,
	}
	require.Equal(t, sdk.NewInt(200), entry.Balance())
}

func TestVeUnbondingDelegationEntry_Slash(t *testing.T) {
	app.Setup(false)

	veBalances := []types.VeUnbondingDelegationEntryBalances{
		{
			VeId:           uint64(100),
			InitialBalance: sdk.NewInt(100),
			Balance:        sdk.NewInt(100),
		},
		{
			VeId:           uint64(200),
			InitialBalance: sdk.NewInt(100),
			Balance:        sdk.NewInt(100),
		},
	}
	entry := types.VeUnbondingDelegationEntry{
		VeBalances: veBalances,
	}
	totalBalance := entry.Balance()
	totalSlashAmt := sdk.NewInt(50)
	totalVeSlashAmt, veBurnedAmounts := entry.Slash(totalSlashAmt, totalBalance, make(map[uint64]sdk.Int))
	require.Equal(t, sdk.NewInt(50), totalVeSlashAmt)
	require.Equal(t, sdk.NewInt(25), veBurnedAmounts[uint64(100)])
	require.Equal(t, sdk.NewInt(25), veBurnedAmounts[uint64(200)])
}

func TestNewVeRedelegation(t *testing.T) {
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

	delegation := types.NewVeRedelegation(acc, valSrc, valDest)
	require.Equal(t, delAddr, delegation.DelegatorAddress)
	require.Equal(t, valSrcAddr, delegation.ValidatorSrcAddress)
	require.Equal(t, valDestAddr, delegation.ValidatorDstAddress)
}

func TestVeRedelegation_AddEntry(t *testing.T) {
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

	delegation := types.NewVeRedelegation(acc, valSrc, valDest)
	veTokens := types.VeTokensSlice{
		{VeId: uint64(100), Tokens: sdk.NewInt(100)},
		{VeId: uint64(200), Tokens: sdk.NewInt(200)},
	}
	totalAmt := veTokens.Tokens()
	totalShares := sdk.NewDec(int64(100))
	delegation.AddEntry(veTokens, totalAmt, totalShares)

	entries := delegation.Entries
	l := len(entries)
	for i := 0; i < l; i++ {
		require.Equal(t, veTokens[i].VeId, entries[0].VeShares[i].VeId)
		require.Equal(t, veTokens[i].Tokens, entries[0].VeShares[i].InitialBalance)
	}
}

func TestVeRedelegation_RemoveEntry(t *testing.T) {
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

	delegation := types.NewVeRedelegation(acc, valSrc, valDest)
	veTokens := types.VeTokensSlice{
		{VeId: uint64(100), Tokens: sdk.NewInt(100)},
		{VeId: uint64(200), Tokens: sdk.NewInt(200)},
	}
	totalAmt := veTokens.Tokens()
	totalShares := sdk.NewDec(int64(100))
	delegation.AddEntry(veTokens, totalAmt, totalShares)
	delegation.RemoveEntry(0)
	entries := delegation.Entries
	require.Equal(t, 0, len(entries))
}

func TestVeRedelegationEntry_InitialBalance(t *testing.T) {
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

	delegation := types.NewVeRedelegation(acc, valSrc, valDest)
	veTokens := types.VeTokensSlice{
		{VeId: uint64(100), Tokens: sdk.NewInt(100)},
		{VeId: uint64(200), Tokens: sdk.NewInt(200)},
	}
	totalAmt := veTokens.Tokens()
	totalShares := sdk.NewDec(int64(100))
	delegation.AddEntry(veTokens, totalAmt, totalShares)

	entries := delegation.Entries
	balance := entries[0].InitialBalance()
	require.Equal(t, sdk.NewInt(300), balance)
}

func TestVeTokensSlice_Tokens(t *testing.T) {
	app.Setup(false)
	veTokens := types.VeTokensSlice{
		{VeId: uint64(100), Tokens: sdk.NewInt(100)},
		{VeId: uint64(200), Tokens: sdk.NewInt(200)},
	}
	totalAmt := veTokens.Tokens()
	require.Equal(t, sdk.NewInt(300), totalAmt)
}

func TestVeTokensSlice_AddToMap(t *testing.T) {
	app.Setup(false)
	veTokens := types.VeTokensSlice{
		{VeId: uint64(100), Tokens: sdk.NewInt(100)},
		{VeId: uint64(200), Tokens: sdk.NewInt(200)},
	}
	veAmounts := veTokens.AddToMap(map[uint64]sdk.Int{
		uint64(100): sdk.NewInt(100),
		uint64(300): sdk.NewInt(100),
	})
	require.Equal(t, sdk.NewInt(200), veAmounts[uint64(100)])
	require.Equal(t, sdk.NewInt(100), veAmounts[uint64(300)])
}
