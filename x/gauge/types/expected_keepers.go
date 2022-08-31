package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	SetAccount(ctx sdk.Context, acc authtypes.AccountI)
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type NftKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
}

type VeKeeper interface {
	GetTotalVotingPower(ctx sdk.Context, atTime uint64, atBlock int64) sdk.Int
	GetVotingPower(ctx sdk.Context, veID uint64, atTime uint64, atBlock int64) sdk.Int
	IncVeAttached(ctx sdk.Context, veID uint64)
	DecVeAttached(ctx sdk.Context, veID uint64)
}

type VoterKeeper interface {
	DistributeReward(ctx sdk.Context, poolDenom string)
}
