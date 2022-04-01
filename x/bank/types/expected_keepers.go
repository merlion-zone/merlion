package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	erc20types "github.com/merlion-zone/merlion/x/erc20/types"
)

type Erc20Keeper interface {
	IsDenomRegistered(ctx sdk.Context, denom string) bool
	RegisterCoin(ctx sdk.Context, denom string) (*erc20types.TokenPair, error)
	IsDenomForErc20(denom string) bool
	SplitCoinsByErc20(amt sdk.Coins) (nativeCoins sdk.Coins, nativeErc20Tokens sdk.Coins)
	GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, nativeCoins sdk.Coins, nativeErc20Tokens sdk.Coins) error
}
