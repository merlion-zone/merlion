package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// NftKeeper defines the expected interface needed to transfer NFT tokens.
type NftKeeper interface {
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
	SaveClass(ctx sdk.Context, class nft.Class) error
	// Methods imported from nft should be defined here
}
