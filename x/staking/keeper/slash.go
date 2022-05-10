package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (k Keeper) Slash(ctx sdk.Context, consAddr sdk.ConsAddress, infractionHeight int64, power int64, slashFactor sdk.Dec) {
	logger := k.Logger(ctx)

	if slashFactor.IsNegative() {
		panic(fmt.Errorf("attempted to slash with a negative slash factor: %v", slashFactor))
	}

	// Amount of slashing = slash slashFactor * power at time of infraction
	amount := k.TokensFromConsensusPower(ctx, power)
	slashAmountDec := amount.ToDec().Mul(slashFactor)
	slashAmount := slashAmountDec.TruncateInt()

	// ref https://github.com/cosmos/cosmos-sdk/issues/1348

	validator, found := k.GetValidatorByConsAddr(ctx, consAddr)
	if !found {
		// If not found, the validator must have been overslashed and removed - so we don't need to do anything
		// NOTE:  Correctness dependent on invariant that unbonding delegations / redelegations must also have been completely
		//        slashed in this case - which we don't explicitly check, but should be true.
		// Log the slash attempt for future reference (maybe we should tag it too)
		logger.Error(
			"WARNING: ignored attempt to slash a nonexistent validator; we recommend you investigate immediately",
			"validator", consAddr.String(),
		)
		return
	}

	// should not be slashing an unbonded validator
	if validator.IsUnbonded() {
		panic(fmt.Sprintf("should not be slashing unbonded validator: %s", validator.GetOperator()))
	}

	operatorAddress := validator.GetOperator()

	// call the before-modification hook
	k.BeforeValidatorModified(ctx, operatorAddress)

	// Track remaining slash amount for the validator
	// This will decrease when we slash unbondings and
	// redelegations, as that stake has since unbonded
	remainingSlashAmount := slashAmount

	switch {
	case infractionHeight > ctx.BlockHeight():
		// Can't slash infractions in the future
		panic(fmt.Sprintf(
			"impossible attempt to slash future infraction at height %d but we are at height %d",
			infractionHeight, ctx.BlockHeight()))

	case infractionHeight == ctx.BlockHeight():
		// Special-case slash at current height for efficiency - we don't need to
		// look through unbonding delegations or redelegations.
		logger.Info(
			"slashing at current height; not scanning unbonding delegations & redelegations",
			"height", infractionHeight,
		)

	case infractionHeight < ctx.BlockHeight():
		// Iterate through unbonding delegations from slashed validator
		unbondingDelegations := k.GetUnbondingDelegationsFromValidator(ctx, operatorAddress)
		for _, unbondingDelegation := range unbondingDelegations {
			amountSlashed := k.SlashUnbondingDelegation(ctx, unbondingDelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}

			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}

		// Iterate through redelegations from slashed source validator
		redelegations := k.GetRedelegationsFromSrcValidator(ctx, operatorAddress)
		for _, redelegation := range redelegations {
			amountSlashed := k.SlashRedelegation(ctx, validator, redelegation, infractionHeight, slashFactor)
			if amountSlashed.IsZero() {
				continue
			}

			remainingSlashAmount = remainingSlashAmount.Sub(amountSlashed)
		}
	}

	// cannot decrease balance below zero
	tokensToBurn := sdk.MinInt(remainingSlashAmount, validator.Tokens)
	tokensToBurn = sdk.MaxInt(tokensToBurn, sdk.ZeroInt()) // defensive.

	// we need to calculate the *effective* slash fraction for distribution
	if validator.Tokens.IsPositive() {
		effectiveFraction := tokensToBurn.ToDec().QuoRoundUp(validator.Tokens.ToDec())
		// possible if power has changed
		if effectiveFraction.GT(sdk.OneDec()) {
			effectiveFraction = sdk.OneDec()
		}
		// call the before-slashed hook
		k.BeforeValidatorSlashed(ctx, operatorAddress, effectiveFraction)
	}

	// Deduct from validator's bonded tokens and update the validator.
	// Burn the slashed tokens from the pool account and decrease the total supply.
	validator = k.RemoveValidatorTokens(ctx, validator, tokensToBurn)

	veValidator, found := k.GetVeValidator(ctx, validator.GetOperator())
	if found && validator.DelegatorShares.IsPositive() {
		tokensToBurn = tokensToBurn.ToDec().Mul(validator.DelegatorShares.Sub(veValidator.VeDelegatorShares)).Quo(validator.DelegatorShares).TruncateInt()
	}

	switch validator.GetStatus() {
	case stakingtypes.Bonded:
		if err := k.BurnBondedTokens(ctx, tokensToBurn); err != nil {
			panic(err)
		}
	case stakingtypes.Unbonding, stakingtypes.Unbonded:
		if err := k.BurnNotBondedTokens(ctx, tokensToBurn); err != nil {
			panic(err)
		}
	default:
		panic("invalid validator status")
	}

	logger.Info(
		"validator slashed by slash factor",
		"validator", validator.GetOperator().String(),
		"slash_factor", slashFactor.String(),
		"burned", tokensToBurn,
	)
}

func (k Keeper) SlashUnbondingDelegation(
	ctx sdk.Context,
	unbondingDelegation stakingtypes.UnbondingDelegation,
	infractionHeight int64,
	slashFactor sdk.Dec,
) (totalSlashAmount sdk.Int) {

	delAddr, err := sdk.AccAddressFromBech32(unbondingDelegation.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valAddr, err := sdk.ValAddressFromBech32(unbondingDelegation.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	veUbd, found := k.GetVeUnbondingDelegation(ctx, delAddr, valAddr)
	if !found || len(veUbd.Entries) != len(unbondingDelegation.Entries) {
		panic("inconsistent ve unbonding delegation")
	}

	now := ctx.BlockHeader().Time
	totalSlashAmount = sdk.ZeroInt()
	burnedAmount := sdk.ZeroInt()
	veBurnedAmounts := make(map[uint64]sdk.Int)

	// perform slashing on all entries within the unbonding delegation
	for i, entry := range unbondingDelegation.Entries {
		veEntry := veUbd.Entries[i]

		// If unbonding started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Unbonding delegation no longer eligible for slashing, skip it
			continue
		}

		// Calculate slash amount proportional to stake contributing to infraction
		slashAmountDec := slashFactor.MulInt(entry.InitialBalance)
		slashAmount := slashAmountDec.TruncateInt()
		totalSlashAmount = totalSlashAmount.Add(slashAmount)

		// Don't slash more tokens than held
		// Possible since the unbonding delegation may already
		// have been slashed, and slash amounts are calculated
		// according to stake held at time of infraction
		unbondingSlashAmount := sdk.MinInt(slashAmount, entry.Balance)

		// Update unbonding delegation if necessary
		if unbondingSlashAmount.IsZero() {
			continue
		}

		var veSlashAmt sdk.Int
		veSlashAmt, veBurnedAmounts = veEntry.Slash(unbondingSlashAmount, entry.Balance, veBurnedAmounts)
		veUbd.Entries[i] = veEntry
		k.SetVeUnbondingDelegation(ctx, veUbd)

		burnedAmount = burnedAmount.Add(unbondingSlashAmount.Sub(veSlashAmt))
		entry.Balance = entry.Balance.Sub(unbondingSlashAmount)
		unbondingDelegation.Entries[i] = entry
		k.SetUnbondingDelegation(ctx, unbondingDelegation)
	}

	if err := k.BurnNotBondedTokens(ctx, burnedAmount); err != nil {
		panic(err)
	}

	for veID, burnedAmt := range veBurnedAmounts {
		if burnedAmt.IsPositive() {
			k.SlashVeDelegatedAmount(ctx, veID, burnedAmt)
		}
	}

	return totalSlashAmount
}

func (k Keeper) SlashRedelegation(
	ctx sdk.Context,
	srcValidator stakingtypes.Validator,
	redelegation stakingtypes.Redelegation,
	infractionHeight int64,
	slashFactor sdk.Dec,
) (totalSlashAmount sdk.Int) {

	delAddr, err := sdk.AccAddressFromBech32(redelegation.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	valSrcAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorSrcAddress)
	if err != nil {
		panic(err)
	}
	valDstAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorDstAddress)
	if err != nil {
		panic(err)
	}
	veRed, found := k.GetVeRedelegation(ctx, delAddr, valSrcAddr, valDstAddr)
	if !found || len(veRed.Entries) != len(redelegation.Entries) {
		panic("inconsistent ve redelegation")
	}

	now := ctx.BlockHeader().Time
	totalSlashAmount = sdk.ZeroInt()
	bondedBurnedAmount, notBondedBurnedAmount := sdk.ZeroInt(), sdk.ZeroInt()
	veBurnedAmounts := make(map[uint64]sdk.Int)

	// perform slashing on all entries within the redelegation
	for i, entry := range redelegation.Entries {
		// If redelegation started before this height, stake didn't contribute to infraction
		if entry.CreationHeight < infractionHeight {
			continue
		}

		if entry.IsMature(now) {
			// Redelegation no longer eligible for slashing, skip it
			continue
		}

		// Calculate slash amount proportional to stake contributing to infraction
		slashAmountDec := slashFactor.MulInt(entry.InitialBalance)
		slashAmount := slashAmountDec.TruncateInt()
		totalSlashAmount = totalSlashAmount.Add(slashAmount)

		// Unbond from target validator
		sharesToUnbond := slashFactor.Mul(entry.SharesDst)
		if sharesToUnbond.IsZero() {
			continue
		}

		valDstAddr, err := sdk.ValAddressFromBech32(redelegation.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}

		delegatorAddress, err := sdk.AccAddressFromBech32(redelegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		delegation, found := k.GetDelegation(ctx, delegatorAddress, valDstAddr)
		if !found {
			// If deleted, delegation has zero shares, and we can't unbond any more
			continue
		}

		if sharesToUnbond.GT(delegation.Shares) {
			sharesToUnbond = delegation.Shares
		}

		returnAmount, veTokens, err := k.Unbond(ctx, delegatorAddress, valDstAddr, sharesToUnbond, false)
		if err != nil {
			panic(fmt.Errorf("error unbonding delegator: %v", err))
		}
		tokensToBurn := returnAmount.Sub(veTokens.Tokens())
		veBurnedAmounts = veTokens.AddToMap(veBurnedAmounts)

		dstValidator, found := k.GetValidator(ctx, valDstAddr)
		if !found {
			panic("destination validator not found")
		}

		// tokens of a redelegation currently live in the destination validator
		// therefor we must burn tokens from the destination-validator's bonding status
		switch {
		case dstValidator.IsBonded():
			bondedBurnedAmount = bondedBurnedAmount.Add(tokensToBurn)
		case dstValidator.IsUnbonded() || dstValidator.IsUnbonding():
			notBondedBurnedAmount = notBondedBurnedAmount.Add(tokensToBurn)
		default:
			panic("unknown validator status")
		}

		// TODO: should we update SharesDst?
		veEntry := veRed.Entries[i]
		_ = veEntry
	}

	if err := k.BurnBondedTokens(ctx, bondedBurnedAmount); err != nil {
		panic(err)
	}

	if err := k.BurnNotBondedTokens(ctx, notBondedBurnedAmount); err != nil {
		panic(err)
	}

	for veID, burnedAmt := range veBurnedAmounts {
		if burnedAmt.IsPositive() {
			k.SlashVeDelegatedAmount(ctx, veID, burnedAmt)
		}
	}

	return totalSlashAmount
}
