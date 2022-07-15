package keeper_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	testkeeper "github.com/merlion-zone/merlion/testutil/keeper"
	"github.com/merlion-zone/merlion/x/erc20/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetAllTokenPairs(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		pairs  = []types.TokenPair{
			types.NewTokenPair(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), "USDT", 0),
			types.NewTokenPair(common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), "USDC", 0)}
	)
	for _, p := range pairs {
		k.SetTokenPair(ctx, p)
	}
	allPairs := k.GetAllTokenPairs(ctx)
	require.Equal(t, pairs, allPairs)
}

func TestKeeper_SetTokenPair_GetTokenPair(t *testing.T) {
	var (
		k, ctx    = testkeeper.Erc20Keeper(t)
		goodPair  = types.NewTokenPair(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), "USDT", 0)
		emptyPair = types.TokenPair{}
		id        = goodPair.GetID()
	)
	k.SetTokenPair(ctx, goodPair)

	for _, tc := range []struct {
		desc string
		id   []byte
		pair types.TokenPair
		ok   bool
	}{
		{
			desc: "id is nil",
			id:   nil,
			pair: emptyPair,
			ok:   false,
		},
		{
			desc: "id is bad",
			id:   []byte("bad"),
			pair: emptyPair,
			ok:   false,
		},
		{
			desc: "id is good",
			id:   id,
			pair: goodPair,
			ok:   true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			res, ok := k.GetTokenPair(ctx, tc.id)
			require.Equal(t, tc.pair, res)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestKeeper_SetERC20Map_GetERC20Map(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		addr   = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		pair   = types.NewTokenPair(addr, "USDT", 0)
		id     = pair.GetID()
	)
	k.SetERC20Map(ctx, addr, id)
	res := k.GetERC20Map(ctx, addr)
	require.Equal(t, id, res)
}

func TestKeeper_SetDenomMap_GetDenomMap(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		addr   = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		pair   = types.NewTokenPair(addr, "USDT", 0)
		id     = pair.GetID()
	)
	k.SetDenomMap(ctx, pair.Denom, id)
	res := k.GetDenomMap(ctx, pair.Denom)
	require.Equal(t, id, res)
}

func TestKeeper_DeleteTokenPair(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		addr   = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		pair   = types.NewTokenPair(addr, "USDT", 0)
		id     = pair.GetID()
	)
	k.SetTokenPair(ctx, pair)
	k.DeleteTokenPair(ctx, pair)
	res, ok := k.GetTokenPair(ctx, id)
	require.Equal(t, res, types.TokenPair{})
	require.Equal(t, ok, false)
	require.Nil(t, k.GetERC20Map(ctx, addr))
	require.Nil(t, k.GetDenomMap(ctx, pair.Denom))
}

func TestKeeper_GetTokenPairID(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		addr   = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		pair   = types.NewTokenPair(addr, "USDT", 0)
		id     = pair.GetID()
	)

	for _, tc := range []struct {
		desc  string
		token string
		id    []byte
	}{
		{
			desc:  "token is address",
			token: addr.String(),
			id:    id,
		},
		{
			desc:  "token is denom",
			token: pair.Denom,
			id:    id,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			if common.IsHexAddress(tc.token) {
				k.SetERC20Map(ctx, common.HexToAddress(tc.token), id)
			} else {
				k.SetDenomMap(ctx, tc.token, id)
			}
			res := k.GetTokenPairID(ctx, tc.token)
			require.Equal(t, tc.id, res)
		})
	}
}

func TestKeeper_IsTokenPairRegistered(t *testing.T) {
	var (
		k, ctx   = testkeeper.Erc20Keeper(t)
		goodPair = types.NewTokenPair(common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), "USDT", 0)
		id       = goodPair.GetID()
	)
	k.SetTokenPair(ctx, goodPair)

	for _, tc := range []struct {
		desc string
		id   []byte
		ok   bool
	}{
		{
			desc: "token pair registered",
			id:   id,
			ok:   true,
		},
		{
			desc: "token pair not registered",
			id:   []byte("bad"),
			ok:   false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ok := k.IsTokenPairRegistered(ctx, tc.id)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestKeeper_IsERC20Registered(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		addr   = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		pair   = types.NewTokenPair(addr, "USDT", 0)
		id     = pair.GetID()
	)
	k.SetERC20Map(ctx, addr, id)

	for _, tc := range []struct {
		desc  string
		erc20 common.Address
		ok    bool
	}{
		{
			desc:  "erc20 registered",
			erc20: addr,
			ok:    true,
		},
		{
			desc:  "erc20 not registered",
			erc20: common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
			ok:    false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ok := k.IsERC20Registered(ctx, tc.erc20)
			require.Equal(t, tc.ok, ok)
		})
	}
}

func TestKeeper_IsDenomRegistered(t *testing.T) {
	var (
		k, ctx = testkeeper.Erc20Keeper(t)
		addr   = common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
		pair   = types.NewTokenPair(addr, "USDT", 0)
		id     = pair.GetID()
	)
	k.SetDenomMap(ctx, pair.Denom, id)

	for _, tc := range []struct {
		desc  string
		denom string
		ok    bool
	}{
		{
			desc:  "Denom registered",
			denom: pair.Denom,
			ok:    true,
		},
		{
			desc:  "Denom not registered",
			denom: "USDC",
			ok:    false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			ok := k.IsDenomRegistered(ctx, tc.denom)
			require.Equal(t, tc.ok, ok)
		})
	}
}
