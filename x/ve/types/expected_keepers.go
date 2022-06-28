package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetModuleAddress(name string) sdk.AccAddress
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetSupply(ctx sdk.Context, denom string) sdk.Coin
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	// Methods imported from bank should be defined here
}

// NftKeeper defines the expected interface needed to transfer NFT tokens.
type NftKeeper interface {
	Transfer(ctx sdk.Context, classID string, nftID string, receiver sdk.AccAddress) error
	SaveClass(ctx sdk.Context, class nft.Class) error
	HasClass(ctx sdk.Context, classID string) bool
	Mint(ctx sdk.Context, token nft.NFT, receiver sdk.AccAddress) error
	Burn(ctx sdk.Context, classID string, nftID string) error
	GetNFTsOfClassByOwner(ctx sdk.Context, classID string, owner sdk.AccAddress) (nfts []nft.NFT)
	GetOwner(ctx sdk.Context, classID string, nftID string) sdk.AccAddress
	HasNFT(ctx sdk.Context, classID, id string) bool
	NFTs(goCtx context.Context, r *nft.QueryNFTsRequest) (*nft.QueryNFTsResponse, error)
	NFT(goCtx context.Context, r *nft.QueryNFTRequest) (*nft.QueryNFTResponse, error)
	// Methods imported from nft should be defined here
}
