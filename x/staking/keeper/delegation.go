package keeper

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/merlion-zone/merlion/x/staking/types"
	vetypes "github.com/merlion-zone/merlion/x/ve/types"
)

func (k Keeper) GetVeValidator(ctx sdk.Context, addr sdk.ValAddress) (validator types.VeValidator, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVeValidatorKey(addr)
	bz := store.Get(key)
	if bz == nil {
		return validator, false
	}

	k.cdc.MustUnmarshal(bz, &validator)
	return validator, true
}

func (k Keeper) SetVeValidator(ctx sdk.Context, validator types.VeValidator) {
	valAddr, err := sdk.ValAddressFromBech32(validator.OperatorAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&validator)
	store.Set(types.GetVeValidatorKey(valAddr), bz)
}

func (k Keeper) RemoveVeValidator(ctx sdk.Context, addr sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVeValidatorKey(addr))
}

// GetVeDelegation returns a specific ve delegation.
func (k Keeper) GetVeDelegation(ctx sdk.Context,
	delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation types.VeDelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetVeDelegationKey(delAddr, valAddr)
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
	store.Set(types.GetVeDelegationKey(delAddr, valAddr), bz)
}

// RemoveVeDelegation deletes a ve delegation.
func (k Keeper) RemoveVeDelegation(ctx sdk.Context, delegation types.VeDelegation) {
	delAddr, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVeDelegationKey(delAddr, valAddr))
}

func (k Keeper) SetVeDelegatedAmount(ctx sdk.Context, veID uint64, amount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: amount})
	store.Set(types.GetVeTokensKey(veID), bz)
}

func (k Keeper) GetVeDelegatedAmount(ctx sdk.Context, veID uint64) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetVeTokensKey(veID))
	if bz == nil {
		return sdk.ZeroInt()
	}
	var amount sdk.IntProto
	k.cdc.MustUnmarshal(bz, &amount)
	return amount.Int
}

func (k Keeper) RemoveVeDelegatedAmount(ctx sdk.Context, veID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVeTokensKey(veID))
}

func (k Keeper) SubVeDelegatedAmount(ctx sdk.Context, veID uint64, subAmount sdk.Int) {
	if !subAmount.IsPositive() {
		return
	}
	veDelegatedAmt := k.GetVeDelegatedAmount(ctx, veID)
	veDelegatedAmt = veDelegatedAmt.Sub(subAmount)
	if !veDelegatedAmt.IsPositive() {
		k.RemoveVeDelegatedAmount(ctx, veID)
	} else {
		k.SetVeDelegatedAmount(ctx, veID, veDelegatedAmt)
	}
}

func (k Keeper) SlashVeDelegatedAmount(ctx sdk.Context, veID uint64, burnAmount sdk.Int) {
	if !burnAmount.IsPositive() {
		return
	}

	k.SubVeDelegatedAmount(ctx, veID, burnAmount)

	k.veKeeper.SlashLockedAmountByUser(ctx, veID, burnAmount)
}

func (k Keeper) SettleVeDelegation(ctx sdk.Context, delegation types.VeDelegation, validator stakingtypes.Validator) types.VeDelegation {
	for i, shares := range delegation.VeShares {
		initialTokens := shares.TokensMayUnsettled
		shares.TokensMayUnsettled = validator.TokensFromShares(shares.Shares).TruncateInt()
		delegation.VeShares[i] = shares

		burnAmount := initialTokens.Sub(shares.TokensMayUnsettled)
		k.SlashVeDelegatedAmount(ctx, shares.VeId, burnAmount)
	}

	if delegation.Shares().IsZero() {
		k.RemoveVeDelegation(ctx, delegation)
	} else {
		k.SetVeDelegation(ctx, delegation)
	}
	return delegation
}

func (k Keeper) GetVeUnbondingDelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress,
) (ubd types.VeUnbondingDelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetVeUBDKey(delAddr, valAddr))
	if bz == nil {
		return ubd, false
	}

	k.cdc.MustUnmarshal(bz, &ubd)
	return ubd, true
}

func (k Keeper) SetVeUnbondingDelegation(ctx sdk.Context, ubd types.VeUnbondingDelegation) {
	delAddr, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&ubd)
	store.Set(types.GetVeUBDKey(delAddr, valAddr), bz)
}

func (k Keeper) RemoveVeUnbondingDelegation(
	ctx sdk.Context, ubd types.VeUnbondingDelegation,
) {
	delAddr, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(ubd.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVeUBDKey(delAddr, valAddr))
}

func (k Keeper) SetVeUnbondingDelegationEntry(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
	veTokens types.VeTokensSlice,
) types.VeUnbondingDelegation {
	ubd, found := k.GetVeUnbondingDelegation(ctx, delegatorAddr, validatorAddr)
	if !found {
		ubd = types.NewVeUnbondingDelegation(delegatorAddr, validatorAddr)
	}

	ubd.AddEntry(veTokens)

	k.SetVeUnbondingDelegation(ctx, ubd)
	return ubd
}

func (k Keeper) GetVeRedelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr sdk.ValAddress, valDstAddr sdk.ValAddress,
) (red types.VeRedelegation, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetVeREDKey(delAddr, valSrcAddr, valDstAddr))
	if bz == nil {
		return red, false
	}

	k.cdc.MustUnmarshal(bz, &red)
	return red, true
}

func (k Keeper) SetVeRedelegation(ctx sdk.Context, red types.VeRedelegation) {
	delAddr, err := sdk.AccAddressFromBech32(red.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
	if err != nil {
		panic(err)
	}
	valDstAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&red)
	store.Set(types.GetVeREDKey(delAddr, valSrcAddr, valDstAddr), bz)
}

func (k Keeper) RemoveVeRedelegation(
	ctx sdk.Context, red types.VeRedelegation,
) {
	delAddr, err := sdk.AccAddressFromBech32(red.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valSrcAddr, err := sdk.ValAddressFromBech32(red.ValidatorSrcAddress)
	if err != nil {
		panic(err)
	}
	valDstAddr, err := sdk.ValAddressFromBech32(red.ValidatorDstAddress)
	if err != nil {
		panic(err)
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetVeREDKey(delAddr, valSrcAddr, valDstAddr))
}

func (k Keeper) SetVeRedelegationEntry(
	ctx sdk.Context, delegatorAddr sdk.AccAddress, valSrcAddr sdk.ValAddress, valDstAddr sdk.ValAddress,
	totalAmount sdk.Int, veTokens types.VeTokensSlice, totalShares sdk.Dec,
) types.VeRedelegation {
	red, found := k.GetVeRedelegation(ctx, delegatorAddr, valSrcAddr, valDstAddr)
	if !found {
		red = types.NewVeRedelegation(delegatorAddr, valSrcAddr, valDstAddr)
	}

	red.AddEntry(veTokens, totalAmount, totalShares)

	k.SetVeRedelegation(ctx, red)
	return red
}

func (k Keeper) VeDelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, bondAmt sdk.Int, veTokens types.VeTokensSlice,
	tokenSrc stakingtypes.BondStatus, validator stakingtypes.Validator, subtractAccount bool,
) (newShares sdk.Dec, err error) {
	if validator.InvalidExRate() {
		return sdk.ZeroDec(), stakingtypes.ErrDelegatorShareExRateInvalid
	}
	valAddr := validator.GetOperator()
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		delegation = stakingtypes.NewDelegation(delAddr, valAddr, sdk.ZeroDec())
	}

	veDelegation, got := k.GetVeDelegation(ctx, delAddr, valAddr)
	if !got {
		veDelegation = types.VeDelegation{
			DelegatorAddress: delAddr.String(),
			ValidatorAddress: validator.OperatorAddress,
		}
	} else {
		veDelegation = k.SettleVeDelegation(ctx, veDelegation, validator)
	}
	if found {
		k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	} else {
		k.BeforeDelegationCreated(ctx, delAddr, valAddr)
	}

	totalNewShares := sdk.ZeroDec()
	totalNewVeShares := sdk.ZeroDec()

	unconstrainedAmt := bondAmt.Sub(veTokens.Tokens())
	if unconstrainedAmt.IsPositive() {
		// if subtractAccount is true then we are
		// performing a delegation and not a redelegation, thus the source tokens are
		// all non bonded
		if subtractAccount {
			if tokenSrc == stakingtypes.Bonded {
				panic("delegation token source cannot be bonded")
			}

			var sendName string

			switch {
			case validator.IsBonded():
				sendName = stakingtypes.BondedPoolName
			case validator.IsUnbonding(), validator.IsUnbonded():
				sendName = stakingtypes.NotBondedPoolName
			default:
				panic("invalid validator status")
			}

			coins := sdk.NewCoins(sdk.NewCoin(k.BondDenom(ctx), unconstrainedAmt))
			if err := k.bankKeeper.DelegateCoinsFromAccountToModule(ctx, delAddr, sendName, coins); err != nil {
				return sdk.Dec{}, err
			}
		} else {
			// potentially transfer tokens between pools, if
			switch {
			case tokenSrc == stakingtypes.Bonded && validator.IsBonded():
				// do nothing
			case (tokenSrc == stakingtypes.Unbonded || tokenSrc == stakingtypes.Unbonding) && !validator.IsBonded():
				// do nothing
			case (tokenSrc == stakingtypes.Unbonded || tokenSrc == stakingtypes.Unbonding) && validator.IsBonded():
				// transfer pools
				k.NotBondedTokensToBonded(ctx, unconstrainedAmt)
			case tokenSrc == stakingtypes.Bonded && !validator.IsBonded():
				// transfer pools
				k.BondedTokensToNotBonded(ctx, unconstrainedAmt)
			default:
				panic("unknown token source bond status")
			}
		}

		_, newShares = k.AddValidatorTokensAndShares(ctx, validator, unconstrainedAmt)

		// Update delegation
		delegation.Shares = delegation.Shares.Add(newShares)
		k.SetDelegation(ctx, delegation)

		totalNewShares = totalNewShares.Add(newShares)
	}

	for _, vt := range veTokens {
		veID := vt.VeId
		veBondAmt := vt.Tokens

		if veID == vetypes.EmptyVeID {
			panic("invalid ve id")
		}

		locked := k.veKeeper.GetLockedAmountByUser(ctx, veID)
		if locked.End <= uint64(ctx.BlockTime().Unix()) {
			return sdk.ZeroDec(), sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "ve expired due to unlocking time %s", time.Unix(int64(locked.End), 0))
		}

		veShares, _ := veDelegation.GetSharesByVeID(veID)

		veDelegatedAmt := k.GetVeDelegatedAmount(ctx, veID)
		veDelegatedAmt = veDelegatedAmt.Add(veBondAmt)
		if veDelegatedAmt.GT(locked.Amount) {
			return sdk.ZeroDec(), sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "insufficient ve locked amount")
		}

		k.SetVeDelegatedAmount(ctx, veID, veDelegatedAmt)

		_, newShares = k.AddValidatorTokensAndShares(ctx, validator, veBondAmt)

		delegation.Shares = delegation.Shares.Add(newShares)
		k.SetDelegation(ctx, delegation)

		veShares.AddTokensAndShares(veBondAmt, newShares)
		veDelegation.SetSharesByVeID(veShares)
		k.SetVeDelegation(ctx, veDelegation)

		totalNewShares = totalNewShares.Add(newShares)
		totalNewVeShares = totalNewVeShares.Add(newShares)
	}

	if totalNewShares.IsPositive() {
		veValidator, found := k.GetVeValidator(ctx, valAddr)
		if !found {
			veValidator = types.VeValidator{
				OperatorAddress:   validator.OperatorAddress,
				VeDelegatorShares: sdk.ZeroDec(),
			}
		}
		veValidator.VeDelegatorShares = veValidator.VeDelegatorShares.Add(totalNewShares)
		k.SetVeValidator(ctx, veValidator)
	}

	k.AfterDelegationModified(ctx, delAddr, delegation.GetValidatorAddr())

	return totalNewShares, nil
}

func (k Keeper) BeginRedelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (completionTime time.Time, err error) {
	if bytes.Equal(valSrcAddr, valDstAddr) {
		return time.Time{}, stakingtypes.ErrSelfRedelegation
	}

	dstValidator, found := k.GetValidator(ctx, valDstAddr)
	if !found {
		return time.Time{}, stakingtypes.ErrBadRedelegationDst
	}

	srcValidator, found := k.GetValidator(ctx, valSrcAddr)
	if !found {
		return time.Time{}, stakingtypes.ErrBadRedelegationDst
	}

	// check if this is a transitive redelegation
	if k.HasReceivingRedelegation(ctx, delAddr, valSrcAddr) {
		return time.Time{}, stakingtypes.ErrTransitiveRedelegation
	}

	if k.HasMaxRedelegationEntries(ctx, delAddr, valSrcAddr, valDstAddr) {
		return time.Time{}, stakingtypes.ErrMaxRedelegationEntries
	}

	returnAmount, veTokens, err := k.Unbond(ctx, delAddr, valSrcAddr, sharesAmount, true)
	if err != nil {
		return time.Time{}, err
	}

	if returnAmount.IsZero() {
		return time.Time{}, stakingtypes.ErrTinyRedelegationAmount
	}

	sharesCreated, err := k.VeDelegate(ctx, delAddr, returnAmount, veTokens, srcValidator.GetStatus(), dstValidator, false)
	if err != nil {
		return time.Time{}, err
	}

	// create the unbonding delegation
	completionTime, height, completeNow := k.GetBeginInfo(ctx, valSrcAddr)

	if completeNow { // no need to create the redelegation object
		return completionTime, nil
	}

	red := k.SetRedelegationEntry(
		ctx, delAddr, valSrcAddr, valDstAddr,
		height, completionTime, returnAmount, sharesAmount, sharesCreated,
	)
	k.InsertRedelegationQueue(ctx, red, completionTime)

	veRed := k.SetVeRedelegationEntry(ctx, delAddr, valSrcAddr, valDstAddr, returnAmount, veTokens, sharesCreated)
	if len(red.Entries) != len(veRed.Entries) {
		panic("inconsistent ve redelegation entries")
	}

	return completionTime, nil
}

func (k Keeper) Undelegate(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec,
) (time.Time, error) {
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		return time.Time{}, stakingtypes.ErrNoDelegatorForAddress
	}

	if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
		return time.Time{}, stakingtypes.ErrMaxUnbondingDelegationEntries
	}

	returnAmount, veTokens, err := k.Unbond(ctx, delAddr, valAddr, sharesAmount, false)
	if err != nil {
		return time.Time{}, err
	}
	unconstrainedAmt := returnAmount.Sub(veTokens.Tokens())

	// transfer the validator tokens to the not bonded pool
	if validator.IsBonded() && unconstrainedAmt.IsPositive() {
		k.BondedTokensToNotBonded(ctx, unconstrainedAmt)
	}

	completionTime := ctx.BlockHeader().Time.Add(k.UnbondingTime(ctx))
	ubd := k.SetUnbondingDelegationEntry(ctx, delAddr, valAddr, ctx.BlockHeight(), completionTime, returnAmount)
	k.InsertUBDQueue(ctx, ubd, completionTime)

	veUbd := k.SetVeUnbondingDelegationEntry(ctx, delAddr, valAddr, veTokens)
	if len(ubd.Entries) != len(veUbd.Entries) {
		panic("inconsistent ve unbonding delegation entries")
	}

	return completionTime, nil
}

func (k Keeper) Unbond(
	ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec, updateVeAmt bool,
) (amount sdk.Int, veTokens types.VeTokensSlice, err error) {
	// check if a delegation object exists in the store
	delegation, found := k.GetDelegation(ctx, delAddr, valAddr)
	if !found {
		err = stakingtypes.ErrNoDelegatorForAddress
		return
	}

	veDelegation, found := k.GetVeDelegation(ctx, delAddr, valAddr)
	if !found {
		// no ve delegation, so still go to original Unbound method
		amount, err = k.Keeper.Unbond(ctx, delAddr, valAddr, shares)
		return
	}

	// call the before-delegation-modified hook
	k.BeforeDelegationSharesModified(ctx, delAddr, valAddr)

	// ensure that we have enough shares to remove
	if delegation.Shares.LT(shares) {
		err = sdkerrors.Wrap(stakingtypes.ErrNotEnoughDelegationShares, delegation.Shares.String())
		return
	}

	// get validator
	validator, found := k.GetValidator(ctx, valAddr)
	if !found {
		err = stakingtypes.ErrNoValidatorFound
		return
	}

	veDelegation = k.SettleVeDelegation(ctx, veDelegation, validator)

	unconstrainedShares := delegation.Shares.Sub(veDelegation.Shares())

	// subtract shares from delegation
	delegation.Shares = delegation.Shares.Sub(shares)

	delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
	if err != nil {
		return
	}

	isValidatorOperator := delegatorAddress.Equals(validator.GetOperator())

	// If the delegation is the operator of the validator and undelegating will decrease the validator's
	// self-delegation below their minimum, we jail the validator.
	if isValidatorOperator && !validator.Jailed &&
		validator.TokensFromShares(delegation.Shares).TruncateInt().LT(validator.MinSelfDelegation) {
		k.JailValidator(ctx, validator)
		validator = k.MustGetValidator(ctx, validator.GetOperator())
	}

	totalVeShares := sdk.MaxDec(shares.Sub(unconstrainedShares), sdk.ZeroDec())
	remainingShares := totalVeShares
	for i := 0; i < len(veDelegation.VeShares); i++ {
		if remainingShares.IsZero() {
			break
		}

		veShares := veDelegation.VeShares[i]

		veDelegatedAmt := k.GetVeDelegatedAmount(ctx, veShares.VeId)

		minShares := sdk.MinDec(remainingShares, veShares.Shares)
		veShares.Shares = veShares.Shares.Sub(minShares)
		remainingShares = remainingShares.Sub(minShares)

		if minShares.IsPositive() {
			amt := validator.TokensFromShares(minShares).TruncateInt()
			veTokens = append(veTokens, types.VeTokens{
				VeId:   veShares.VeId,
				Tokens: amt,
			})

			veShares.TokensMayUnsettled = veShares.TokensMayUnsettled.Sub(amt)
			if veShares.TokensMayUnsettled.IsNegative() {
				panic("inconsistent tokens and shares")
			}

			if updateVeAmt {
				veDelegatedAmt = veDelegatedAmt.Sub(amt)
			}
		}

		if veShares.Shares.IsZero() {
			veDelegation.RemoveSharesByIndex(i)
			i--
		} else {
			veDelegation.SetSharesByVeID(veShares)
		}

		if updateVeAmt {
			if veDelegatedAmt.IsZero() {
				k.RemoveVeDelegatedAmount(ctx, veShares.VeId)
			} else {
				k.SetVeDelegatedAmount(ctx, veShares.VeId, veDelegatedAmt)
			}
		}
	}

	if len(veDelegation.VeShares) != 0 {
		k.SetVeDelegation(ctx, veDelegation)
	} else {
		k.RemoveVeDelegation(ctx, veDelegation)
	}

	if !remainingShares.IsZero() {
		panic("inconsistent shares")
	}

	// remove the delegation
	if delegation.Shares.IsZero() {
		k.RemoveDelegation(ctx, delegation)
	} else {
		k.SetDelegation(ctx, delegation)
		// call the after delegation modification hook
		k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetValidatorAddr())
	}

	// remove the shares and coins from the validator
	// NOTE that the amount is later (in keeper.Delegation) moved between staking module pools
	validator, amount = k.RemoveValidatorTokensAndShares(ctx, validator, shares)

	if validator.DelegatorShares.IsZero() && validator.IsUnbonded() {
		// if not unbonded, we must instead remove validator in EndBlocker once it finishes its unbonding period
		k.RemoveValidator(ctx, validator.GetOperator())
		k.RemoveVeValidator(ctx, validator.GetOperator())
	} else if totalVeShares.IsPositive() {
		veValidator, found := k.GetVeValidator(ctx, validator.GetOperator())
		if !found {
			panic("inconsistent ve validator")
		}
		veValidator.VeDelegatorShares = veValidator.VeDelegatorShares.Sub(totalVeShares)
		if veValidator.VeDelegatorShares.IsNegative() {
			panic("inconsistent ve validator")
		}
		k.SetVeValidator(ctx, veValidator)
	}

	return
}

func (k Keeper) CompleteUnbonding(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (sdk.Coins, error) {
	ubd, found := k.GetUnbondingDelegation(ctx, delAddr, valAddr)
	if !found {
		return nil, stakingtypes.ErrNoUnbondingDelegation
	}

	veUbd, found := k.GetVeUnbondingDelegation(ctx, delAddr, valAddr)
	if !found || len(veUbd.Entries) != len(ubd.Entries) {
		panic("inconsistent ve unbonding delegation")
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	delegatorAddress, err := sdk.AccAddressFromBech32(ubd.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	// loop through all the entries and complete unbonding mature entries
	for i := 0; i < len(ubd.Entries); i++ {
		entry := ubd.Entries[i]
		veEntry := veUbd.Entries[i]
		if entry.IsMature(ctxTime) {
			ubd.RemoveEntry(int64(i))
			veUbd.RemoveEntry(i)
			i--

			for _, vb := range veEntry.VeBalances {
				k.SubVeDelegatedAmount(ctx, vb.VeId, vb.Balance)
			}

			unconstrainedBalance := entry.Balance.Sub(veEntry.Balance())
			if unconstrainedBalance.IsNegative() {
				panic("inconsistent ve unbonding delegation")
			}

			// track undelegation only when remaining or truncated shares are non-zero
			if !unconstrainedBalance.IsZero() {
				amt := sdk.NewCoin(bondDenom, unconstrainedBalance)
				if err := k.bankKeeper.UndelegateCoinsFromModuleToAccount(
					ctx, stakingtypes.NotBondedPoolName, delegatorAddress, sdk.NewCoins(amt),
				); err != nil {
					return nil, err
				}
			}

			balances = balances.Add(sdk.NewCoin(bondDenom, entry.Balance))
		}
	}

	// set the unbonding delegation or remove it if there are no more entries
	if len(ubd.Entries) == 0 {
		k.RemoveUnbondingDelegation(ctx, ubd)
		k.RemoveVeUnbondingDelegation(ctx, veUbd)
	} else {
		k.SetUnbondingDelegation(ctx, ubd)
		k.SetVeUnbondingDelegation(ctx, veUbd)
	}

	return balances, nil
}

func (k Keeper) CompleteRedelegation(
	ctx sdk.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress,
) (sdk.Coins, error) {
	red, found := k.GetRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	if !found {
		return nil, stakingtypes.ErrNoRedelegation
	}

	veRed, found := k.GetVeRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	if !found || len(veRed.Entries) != len(red.Entries) {
		panic("inconsistent ve redelegation")
	}

	bondDenom := k.GetParams(ctx).BondDenom
	balances := sdk.NewCoins()
	ctxTime := ctx.BlockHeader().Time

	// loop through all the entries and complete mature redelegation entries
	for i := 0; i < len(red.Entries); i++ {
		entry := red.Entries[i]
		if entry.IsMature(ctxTime) {
			red.RemoveEntry(int64(i))
			veRed.RemoveEntry(i)
			i--

			if !entry.InitialBalance.IsZero() {
				balances = balances.Add(sdk.NewCoin(bondDenom, entry.InitialBalance))
			}
		}
	}

	// set the redelegation or remove it if there are no more entries
	if len(red.Entries) == 0 {
		k.RemoveRedelegation(ctx, red)
		k.RemoveVeRedelegation(ctx, veRed)
	} else {
		k.SetRedelegation(ctx, red)
		k.SetVeRedelegation(ctx, veRed)
	}

	return balances, nil
}
