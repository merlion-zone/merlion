package keeper_test

import (
	"testing"

	testkeeper "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestKeeper_IsDenomForErc20(t *testing.T) {
	k, _ := testkeeper.Erc20Keeper(t)
	for _, tc := range []struct {
		desc  string
		denom string
		valid bool
	}{
		{
			desc:  "Denom is for erc20",
			denom: "erc20/0xdAC17F958D2ee523a2206206994597C13D831ec7",
			valid: true,
		},
		{
			desc:  "Denom is not for erc20",
			denom: "address",
			valid: false,
		},
		{
			desc:  "Denom is not for erc20",
			denom: "aaa/0xdAC17F958D2ee523a2206206994597C13D831ec7",
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			valid := k.IsDenomForErc20(tc.denom)
			require.Equal(t, tc.valid, valid)
		})
	}
}
