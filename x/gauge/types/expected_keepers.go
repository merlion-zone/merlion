package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	NewAccount(ctx sdk.Context, acc authtypes.AccountI) authtypes.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc authtypes.ModuleAccountI)
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

type NftKeeper interface {
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
}

type VeKeeper interface {
	GetTotalVotingPower(ctx sdk.Context, atTime uint64, atBlock int64) sdk.Int
	GetVotingPower(ctx sdk.Context, veID uint64, atTime uint64, atBlock int64) sdk.Int
}
