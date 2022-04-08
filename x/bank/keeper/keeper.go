package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/merlion-zone/merlion/x/bank/types"
)

// Keeper manages transfers between accounts.
// Actually it wraps the bankkeeper.BaseKeeper struct, to supply hooks around sending coins.
type Keeper struct {
	bankkeeper.BaseKeeper

	ak          banktypes.AccountKeeper
	erc20Keeper types.Erc20Keeper
}

func NewKeeper(cdc codec.BinaryCodec,
	storeKey sdk.StoreKey,
	ak banktypes.AccountKeeper,
	paramSpace paramtypes.Subspace,
	blockedAddrs map[string]bool) Keeper {
	return Keeper{
		BaseKeeper: bankkeeper.NewBaseKeeper(cdc, storeKey, ak, paramSpace, blockedAddrs),
		ak:         ak,
	}
}

func (k *Keeper) SetErc20Keeper(erc20Keeper types.Erc20Keeper) {
	k.erc20Keeper = erc20Keeper
}

func (k Keeper) GetPaginatedTotalSupply(ctx sdk.Context, pagination *query.PageRequest) (sdk.Coins, *query.PageResponse, error) {
	// NOTE: not including total supply of all erc20 native tokens
	return k.BaseKeeper.GetPaginatedTotalSupply(ctx, pagination)
}

func (k Keeper) DelegateCoins(ctx sdk.Context, delegatorAddr, moduleAccAddr sdk.AccAddress, amt sdk.Coins) error {
	nativeCoins, nativeErc20Tokens := k.erc20Keeper.SplitCoinsByErc20(amt)
	if !nativeErc20Tokens.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for delegation")
	}
	if err := k.BaseKeeper.DelegateCoins(ctx, delegatorAddr, moduleAccAddr, amt); err != nil {
		return err
	}
	return k.erc20Keeper.SendCoins(ctx, delegatorAddr, moduleAccAddr, nativeCoins, nil)
}

func (k Keeper) UndelegateCoins(ctx sdk.Context, moduleAccAddr, delegatorAddr sdk.AccAddress, amt sdk.Coins) error {
	nativeCoins, nativeErc20Tokens := k.erc20Keeper.SplitCoinsByErc20(amt)
	if !nativeErc20Tokens.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for delegation")
	}
	if err := k.BaseKeeper.UndelegateCoins(ctx, moduleAccAddr, delegatorAddr, amt); err != nil {
		return err
	}
	return k.erc20Keeper.SendCoins(ctx, moduleAccAddr, delegatorAddr, nativeCoins, nil)
}

func (k Keeper) GetSupply(ctx sdk.Context, denom string) sdk.Coin {
	if !k.erc20Keeper.IsDenomForErc20(denom) {
		return k.BaseKeeper.GetSupply(ctx, denom)
	} else {
		return k.erc20Keeper.GetSupply(ctx, denom)
	}
}

func (k Keeper) HasSupply(ctx sdk.Context, denom string) bool {
	if !k.erc20Keeper.IsDenomForErc20(denom) {
		return k.BaseKeeper.HasSupply(ctx, denom)
	} else {
		return !k.erc20Keeper.GetSupply(ctx, denom).IsZero()
	}
}

func (k Keeper) GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool) {
	if !k.erc20Keeper.IsDenomForErc20(denom) {
		return k.BaseKeeper.GetDenomMetaData(ctx, denom)
	} else {
		return k.erc20Keeper.GetDenomMetaData(ctx, denom)
	}
}

func (k Keeper) GetAllDenomMetaData(ctx sdk.Context) []banktypes.Metadata {
	// NOTE: not including denom metadata of all erc20 native tokens
	return k.BaseKeeper.GetAllDenomMetaData(ctx)
}

func (k Keeper) IterateAllDenomMetaData(ctx sdk.Context, cb func(banktypes.Metadata) bool) {
	// NOTE: not including denom metadata of all erc20 native tokens
	k.BaseKeeper.IterateAllDenomMetaData(ctx, cb)
}

func (k Keeper) SendCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if k.BlockedAddr(recipientAddr) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", recipientAddr)
	}

	return k.SendCoins(ctx, senderAddr, recipientAddr, amt)
}

func (k Keeper) SendCoinsFromModuleToModule(
	ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins,
) error {
	senderAddr := k.ak.GetModuleAddress(senderModule)
	if senderAddr == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

func (k Keeper) SendCoinsFromAccountToModule(
	ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins,
) error {
	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	return k.SendCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

func (k Keeper) DelegateCoinsFromAccountToModule(
	ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins,
) error {
	recipientAcc := k.ak.GetModuleAccount(ctx, recipientModule)
	if recipientAcc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", recipientModule))
	}

	if !recipientAcc.HasPermission(authtypes.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to receive delegated coins", recipientModule))
	}

	return k.DelegateCoins(ctx, senderAddr, recipientAcc.GetAddress(), amt)
}

func (k Keeper) UndelegateCoinsFromModuleToAccount(
	ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins,
) error {
	acc := k.ak.GetModuleAccount(ctx, senderModule)
	if acc == nil {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", senderModule))
	}

	if !acc.HasPermission(authtypes.Staking) {
		panic(sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "module account %s does not have permissions to undelegate coins", senderModule))
	}

	return k.UndelegateCoins(ctx, acc.GetAddress(), recipientAddr, amt)
}

func (k Keeper) MintCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error {
	nativeCoins, nativeErc20Tokens := k.erc20Keeper.SplitCoinsByErc20(amounts)
	if !nativeErc20Tokens.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for mint")
	}
	if err := k.BaseKeeper.MintCoins(ctx, moduleName, amounts); err != nil {
		return err
	}
	for _, coin := range nativeCoins {
		if !k.erc20Keeper.IsDenomRegistered(ctx, coin.Denom) {
			if _, err := k.erc20Keeper.RegisterCoin(ctx, coin.Denom); err != nil {
				return err
			}
		}
	}
	acc := k.ak.GetModuleAccount(ctx, moduleName) // had checked nil
	return k.erc20Keeper.SendCoins(ctx, nil, acc.GetAddress(), nativeCoins, nil)
}

func (k Keeper) BurnCoins(ctx sdk.Context, moduleName string, amounts sdk.Coins) error {
	nativeCoins, nativeErc20Tokens := k.erc20Keeper.SplitCoinsByErc20(amounts)
	if !nativeErc20Tokens.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "erc20 native tokens unqualified for burn")
	}
	if err := k.BaseKeeper.BurnCoins(ctx, moduleName, amounts); err != nil {
		return err
	}
	acc := k.ak.GetModuleAccount(ctx, moduleName) // had checked nil
	return k.erc20Keeper.SendCoins(ctx, acc.GetAddress(), nil, nativeCoins, nil)
}

func (k Keeper) IterateTotalSupply(ctx sdk.Context, cb func(sdk.Coin) bool) {
	// NOTE: not including total supply of all erc20 native tokens
	k.BaseKeeper.IterateTotalSupply(ctx, cb)
}
