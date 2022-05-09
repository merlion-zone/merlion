package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

// GetVeDelegation returns a specific ve delegation.
func (k Keeper) GetVeDelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation types.VeDelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetDelegationKey(delAddr, valAddr)
	bz := store.Get(key)
	if bz == nil {
		return delegation, false
	}

	k.cdc.MustUnmarshal(bz, &delegation)
	return delegation, true
}

// SetVeDelegation sets a ve delegation.
func (k Keeper) SetVeDelegation(ctx sdk.Context, delegation types.VeDelegation) {
	delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&delegation)
	store.Set(types.GetDelegationKey(delAddr, valAddr), bz)
}

func (k Keeper) VeDelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, veID uint64, bondAmt sdk.Int, tokenSrc stakingtypes.BondStatus,
	validator stakingtypes.Validator,
) (newShares sdk.Dec, err error) {
	if veID == vetypes.EmptyVeID {
		panic("invalid ve id")
	}

	if validator.InvalidExRate() {
		return sdk.ZeroDec(), stakingtypes.ErrDelegatorShareExRateInvalid
	}

	delegation, found := k.GetDelegation(ctx, delAddr, validator.GetOperator())
	if !found {
		delegation = stakingtypes.NewDelegation(delAddr, validator.GetOperator(), sdk.ZeroDec())
	}

	veDelegation, got := k.GetVeDelegation(ctx, delAddr, validator.GetOperator())
	if !got {
		veDelegation = types.VeDelegation{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: validator.String(),
		}
	}

	if found {
		k.BeforeDelegationSharesModified(ctx, delAddr, validator.GetOperator())
	} else {
		k.BeforeDelegationCreated(ctx, delAddr, validator.GetOperator())
	}

	locked := k.veKeeper.GetLockedAmountByUser(ctx, veID)
	if locked.End <= uint64(ctx.BlockTime().Unix()) {
		return sdk.ZeroDec(), sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ve expired due to unlocking time %s", time.Unix(int64(locked.End), 0))
	}

	// TODO: settle ve shares and tokens

	veShares, _ := veDelegation.GetSharesByVeID(veID)

	if bondAmt.Add(veShares.Tokens()).GT(locked.Amount) {
		return sdk.ZeroDec(), sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "insufficient ve locked amount")
	}

	_, newShares = k.AddValidatorTokensAndShares(ctx, validator, bondAmt)

	delegation.Shares = delegation.Shares.Add(newShares)
	k.SetDelegation(ctx, delegation)

	veShares.AddTokensAndShares(bondAmt, newShares)
	veDelegation.SetSharesByVeID(veShares)
	k.SetVeDelegation(ctx, veDelegation)

	k.AfterDelegationModified(ctx, delAddr, delegation.GetValidatorAddr())

	return newShares, nil
}
