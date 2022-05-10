package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m *VeDelegation) GetSharesByVeID(veID uint64) (VeShares, bool) {
	for _, shares := range m.VeShares {
		if shares.VeId == veID {
			return shares, true
		}
	}
	return VeShares{
		VeId:               veID,
		TokensMayUnsettled: sdk.ZeroInt(),
		Shares:             sdk.ZeroDec(),
	}, false
}

func (m *VeDelegation) SetSharesByVeID(veShares VeShares) (new bool) {
	for i, shares := range m.VeShares {
		if shares.VeId == veShares.VeId {
			m.VeShares[i] = veShares
			return false
		}
	}
	m.VeShares = append(m.VeShares, veShares)
	return true
}

func (m *VeDelegation) RemoveSharesByIndex(i int) {
	m.VeShares = append(m.VeShares[:i], m.VeShares[i+1:]...)
}

func (m *VeDelegation) Shares() sdk.Dec {
	total := sdk.ZeroDec()
	for _, shares := range m.VeShares {
		total = total.Add(shares.Shares)
	}
	return total
}

func (m *VeDelegation) Tokens() sdk.Int {
	total := sdk.ZeroInt()
	for _, shares := range m.VeShares {
		total = total.Add(shares.Tokens())
	}
	return total
}

func (m *VeShares) Tokens() sdk.Int {
	return m.TokensMayUnsettled
}

func (m *VeShares) AddTokensAndShares(tokens sdk.Int, shares sdk.Dec) {
	m.TokensMayUnsettled = m.TokensMayUnsettled.Add(tokens)
	m.Shares = m.Shares.Add(shares)
}

func NewVeUnbondingDelegation(
	delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
) VeUnbondingDelegation {
	return VeUnbondingDelegation{
		DelegatorAddress: delegatorAddr.String(),
		ValidatorAddress: validatorAddr.String(),
	}
}

func (ubd *VeUnbondingDelegation) AddEntry(veTokens VeTokensSlice) {
	veBalances := make([]VeUnbondingDelegationEntryBalances, 0, len(veTokens))
	for _, vt := range veTokens {
		veBalances = append(veBalances, VeUnbondingDelegationEntryBalances{
			VeId:           vt.VeId,
			InitialBalance: vt.Tokens,
			Balance:        vt.Tokens,
		})
	}
	ubd.Entries = append(ubd.Entries, VeUnbondingDelegationEntry{
		VeBalances: veBalances,
	})
}

func (ubd *VeUnbondingDelegation) RemoveEntry(i int) {
	ubd.Entries = append(ubd.Entries[:i], ubd.Entries[i+1:]...)
}

func (entry *VeUnbondingDelegationEntry) Balance() sdk.Int {
	total := sdk.ZeroInt()
	for _, b := range entry.VeBalances {
		total = total.Add(b.Balance)
	}
	return total
}

// CONTRACT: 0 < totalSlashAmt <= totalBalance
func (entry *VeUnbondingDelegationEntry) Slash(totalSlashAmt, totalBalance sdk.Int, veBurnedAmounts map[uint64]sdk.Int) (sdk.Int, map[uint64]sdk.Int) {
	totalVeSlashAmt := sdk.ZeroInt()
	for _, b := range entry.VeBalances {
		veSlash := totalSlashAmt.Mul(b.Balance).Quo(totalBalance)
		b.Balance = b.Balance.Sub(veSlash)
		totalVeSlashAmt = totalVeSlashAmt.Add(veSlash)

		if veBurnedAmt, ok := veBurnedAmounts[b.VeId]; ok {
			veBurnedAmounts[b.VeId] = veBurnedAmt.Add(veSlash)
		} else {
			veBurnedAmounts[b.VeId] = veSlash
		}
	}
	return totalVeSlashAmt, veBurnedAmounts
}

func NewVeRedelegation(delAddr sdk.AccAddress, valSrcAddr sdk.ValAddress, valDstAddr sdk.ValAddress) VeRedelegation {
	return VeRedelegation{
		DelegatorAddress:    delAddr.String(),
		ValidatorSrcAddress: valSrcAddr.String(),
		ValidatorDstAddress: valDstAddr.String(),
	}
}

func (red *VeRedelegation) AddEntry(veTokens VeTokensSlice, totalAmt sdk.Int, totalShares sdk.Dec) {
	veShares := make([]VeRedelegationEntryShares, 0, len(veTokens))
	for _, vt := range veTokens {
		veShares = append(veShares, VeRedelegationEntryShares{
			VeId:           vt.VeId,
			InitialBalance: vt.Tokens,
			SharesDst:      totalShares.MulInt(vt.Tokens).QuoInt(totalAmt),
		})
	}
	red.Entries = append(red.Entries, VeRedelegationEntry{
		VeShares: veShares,
	})
}

func (red *VeRedelegation) RemoveEntry(i int) {
	red.Entries = append(red.Entries[:i], red.Entries[i+1:]...)
}

func (entry *VeRedelegationEntry) InitialBalance() sdk.Int {
	total := sdk.ZeroInt()
	for _, b := range entry.VeShares {
		total = total.Add(b.InitialBalance)
	}
	return total
}

type VeTokensSlice []VeTokens

func (s VeTokensSlice) Tokens() sdk.Int {
	total := sdk.ZeroInt()
	for _, vt := range s {
		total = total.Add(vt.Tokens)
	}
	return total
}

func (s VeTokensSlice) AddToMap(veAmounts map[uint64]sdk.Int) map[uint64]sdk.Int {
	for _, vt := range s {
		if veAmt, ok := veAmounts[vt.VeId]; ok {
			veAmounts[vt.VeId] = veAmt.Add(vt.Tokens)
		} else {
			veAmounts[vt.VeId] = vt.Tokens
		}
	}
	return veAmounts
}
