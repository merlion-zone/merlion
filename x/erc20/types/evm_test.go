package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewERC20Data(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		name     string
		symbol   string
		decimals uint8
	}{
		{
			desc:     "NewERC20Data is ok",
			name:     "Moon Token",
			symbol:   "Moon",
			decimals: 8,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			d := NewERC20Data(tc.name, tc.symbol, tc.decimals)
			require.Equal(t, tc.name, d.Name)
			require.Equal(t, tc.symbol, d.Symbol)
			require.Equal(t, tc.decimals, d.Decimals)
		})
	}
}
