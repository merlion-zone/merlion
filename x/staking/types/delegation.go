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

func (m *VeShares) Tokens() sdk.Int {
	return m.TokensMayUnsettled
}

func (m *VeShares) AddTokensAndShares(tokens sdk.Int, shares sdk.Dec) {
	m.TokensMayUnsettled = m.TokensMayUnsettled.Add(tokens)
	m.Shares = m.Shares.Add(shares)
}
