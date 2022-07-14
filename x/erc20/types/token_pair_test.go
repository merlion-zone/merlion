package types

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestNewTokenPair(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		erc20Address common.Address
		denom        string
		owner        Owner
	}{
		{
			desc:         "NewTokenPair is ok",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
			owner:        0,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pair := NewTokenPair(tc.erc20Address, tc.denom, tc.owner)
			require.Equal(t, tc.erc20Address.String(), pair.Erc20Address)
			require.Equal(t, tc.denom, pair.Denom)
			require.Equal(t, tc.owner, pair.ContractOwner)
		})
	}
}

func TestTokenPair_GetID(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		erc20Address common.Address
		denom        string
		id           string
	}{
		{
			desc:         "GetID is ok",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
			id:           "1535ea4ac631edbb712eb516efeb1d653e1fe30667a8ca555ebf0a7d2441e75c",
		},
		{
			desc:         "GetID is ok",
			erc20Address: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			denom:        "USDC",
			id:           "5327245738d3ae7eef42b83b7a8fa596b632f296edfde23f7c3250171d8d0241",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pair := NewTokenPair(tc.erc20Address, tc.denom, 2)
			id := pair.GetID()
			require.Equal(t, tc.id, hex.EncodeToString(id))
		})
	}
}

func TestTokenPair_GetERC20Contract(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		erc20Address common.Address
		denom        string
	}{
		{
			desc:         "GetERC20Contract is ok",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
		},
		{
			desc:         "GetERC20Contract is ok",
			erc20Address: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			denom:        "USDC",
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pair := NewTokenPair(tc.erc20Address, tc.denom, 2)
			addr := pair.GetERC20Contract()
			require.Equal(t, tc.erc20Address, addr)
		})
	}
}

func TestTokenPair_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		erc20Address common.Address
		denom        string
		valid        bool
	}{
		{
			desc:         "TokenPair denom invalid",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "1X",
			valid:        false,
		},
		{
			desc:         "TokenPair valid",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
			valid:        true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pair := NewTokenPair(tc.erc20Address, tc.denom, 2)
			err := pair.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestTokenPair_IsNativeCoin(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		erc20Address common.Address
		denom        string
		owner        Owner
		native       bool
	}{
		{
			desc:         "NativeCoin",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
			owner:        OWNER_MODULE,
			native:       true,
		},
		{
			desc:         "Not NativeCoin",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
			owner:        OWNER_EXTERNAL,
			native:       false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pair := NewTokenPair(tc.erc20Address, tc.denom, tc.owner)
			native := pair.IsNativeCoin()
			require.Equal(t, tc.native, native)
		})
	}
}

func TestTokenPair_IsNativeERC20(t *testing.T) {
	for _, tc := range []struct {
		desc         string
		erc20Address common.Address
		denom        string
		owner        Owner
		native       bool
	}{
		{
			desc:         "NativeERC20",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
			owner:        OWNER_EXTERNAL,
			native:       true,
		},
		{
			desc:         "Not NativeERC20",
			erc20Address: common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
			denom:        "USDT",
			owner:        OWNER_MODULE,
			native:       false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			pair := NewTokenPair(tc.erc20Address, tc.denom, tc.owner)
			native := pair.IsNativeERC20()
			require.Equal(t, tc.native, native)
		})
	}
}
