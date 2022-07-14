package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateDenomDescription(t *testing.T) {
	for _, tc := range []struct {
		desc      string
		address   string
		denomDesc string
	}{
		{
			desc:      "CreateDenomDescription is ok",
			address:   "address",
			denomDesc: "Merlion coin token representation of address",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			desc := CreateDenomDescription(tc.address)
			require.Equal(t, tc.denomDesc, desc)
		})
	}
}

func TestCreateDenom(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		address string
		denom   string
	}{
		{
			desc:    "CreateDenom is ok",
			address: "address",
			denom:   "erc20/address",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			denom := CreateDenom(tc.address)
			require.Equal(t, tc.denom, denom)
		})
	}
}

func TestSanitizeERC20Name(t *testing.T) {
	for _, tc := range []struct {
		desc string
		name string
		res  string
	}{
		{
			desc: "ToLower",
			name: "ABC",
			res:  "abc",
		},
		{
			desc: "remove token",
			name: "ABC token",
			res:  "abc",
		},
		{
			desc: "remove coin",
			name: "ABC coin",
			res:  "abc",
		},
		{
			desc: "TrimSpace",
			name: " ABC coin  ",
			res:  "abc",
		},
		{
			desc: "_",
			name: " ABC D coin  ",
			res:  "abc_d",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			res := SanitizeERC20Name(tc.name)
			require.Equal(t, tc.res, res)
		})
	}
}
