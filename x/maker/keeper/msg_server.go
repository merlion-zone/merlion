package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	merlion "github.com/merlion-zone/merlion/types"
	"github.com/merlion-zone/merlion/x/maker/types"
	oracletypes "github.com/merlion-zone/merlion/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) MintBySwap(c context.Context, msg *types.MsgMintBySwap) (*types.MsgMintBySwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	receiver := sender
	if len(msg.To) > 0 {
		// user specifies receiver
		receiver, err = sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			return nil, err
		}
	}

	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(m.Keeper.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralRatio := m.Keeper.GetCollateralRatio(ctx)

	backingParams, found := m.Keeper.GetBackingRiskParams(ctx, msg.BackingInMax.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", msg.BackingInMax.Denom)
	}

	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", msg.BackingInMax.Denom)
	}

	mintOut := msg.MintOut
	mintFee := computeFee(mintOut, backingParams.MintFee)
	mintTotal := mintOut.AddAmount(mintFee.Amount)
	mintTotalInUSD := mintTotal.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	if backingParams.MaxMerMint != nil {
		// TODO
		totalMintedMer := m.Keeper.bankKeeper.GetSupply(ctx, merlion.MicroUSDDenom)
		if totalMintedMer.Amount.Add(mintTotal.Amount).GT(*backingParams.MaxMerMint) {
			return nil, sdkerrors.Wrapf(types.ErrMerCeiling, "existing(%s) + new(%s) > ceiling(%s)", totalMintedMer, mintTotal, backingParams.MaxMerMint)
		}
	}

	// get backing and lion price in usd
	backingPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, msg.BackingInMax.Denom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroLionDenom)
	if err != nil {
		return nil, err
	}

	backingNeeded := sdk.NewCoin(msg.BackingInMax.Denom, sdk.ZeroInt())
	lionNeeded := sdk.NewCoin(merlion.MicroLionDenom, sdk.ZeroInt())
	if collateralRatio.GTE(sdk.OneDec()) || msg.LionInMax.IsZero() {
		// full/over collateralized, or user selects full collateralization
		backingNeeded.Amount = mintTotalInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if collateralRatio.IsZero() {
		// algorithmic
		lionNeeded.Amount = mintTotalInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingNeeded.Amount = mintTotalInUSD.Mul(collateralRatio).QuoRoundUp(backingPrice).RoundInt()
		lionNeeded.Amount = mintTotalInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	if msg.BackingInMax.IsLT(backingNeeded) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinSlippage, "backing coin needed: %s", backingNeeded)
	}
	if msg.LionInMax.IsLT(lionNeeded) {
		return nil, sdkerrors.Wrapf(types.ErrLionCoinSlippage, "lion coin needed: %s", lionNeeded)
	}

	if backingParams.MaxBacking != nil {
		// get backing balance owned by module
		moduleOwnedBacking := m.Keeper.bankKeeper.GetBalance(ctx, m.Keeper.accountKeeper.GetModuleAddress(types.ModuleName), msg.BackingInMax.Denom)
		if moduleOwnedBacking.Amount.Add(backingNeeded.Amount).GT(*backingParams.MaxBacking) {
			return nil, sdkerrors.Wrapf(types.ErrBackingCeiling, "existing(%s) + new(%s) > ceiling(%s)", moduleOwnedBacking, backingNeeded, backingParams.MaxBacking)
		}
	}

	// take backing coin and lion coin
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(backingNeeded, lionNeeded))
	if err != nil {
		return nil, err
	}
	// burn lion
	if lionNeeded.IsPositive() {
		err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lionNeeded))
		if err != nil {
			return nil, err
		}
	}

	// mint mer stablecoin
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintTotal))
	if err != nil {
		return nil, err
	}
	// send mer to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(mintOut))
	if err != nil {
		return nil, err
	}
	// send mer fee to oracle
	if mintFee.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(mintFee))
		if err != nil {
			return nil, err
		}
	}

	// TODO: event

	return &types.MsgMintBySwapResponse{
		MintFee:   mintFee,
		BackingIn: backingNeeded,
		LionIn:    lionNeeded,
	}, nil
}

func (m msgServer) BurnBySwap(c context.Context, msg *types.MsgBurnBySwap) (*types.MsgBurnBySwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	receiver := sender
	if len(msg.To) > 0 {
		// user specifies receiver
		receiver, err = sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			return nil, err
		}
	}

	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price upper bound
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(m.Keeper.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralRatio := m.Keeper.GetCollateralRatio(ctx)

	backingParams, found := m.Keeper.GetBackingRiskParams(ctx, msg.BackingOutMin.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", msg.BackingOutMin.Denom)
	}

	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", msg.BackingOutMin.Denom)
	}

	burnIn := msg.BurnIn
	burnFee := computeFee(burnIn, backingParams.BurnFee)
	burn := burnIn.SubAmount(burnFee.Amount)
	burnInUSD := burn.Amount.ToDec().Mul(merlion.MicroUSDTarget)

	// get backing and lion price in usd
	backingPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, msg.BackingOutMin.Denom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroLionDenom)
	if err != nil {
		return nil, err
	}

	backingOut := sdk.NewCoin(msg.BackingOutMin.Denom, sdk.ZeroInt())
	lionOut := sdk.NewCoin(merlion.MicroLionDenom, sdk.ZeroInt())
	if collateralRatio.GTE(sdk.OneDec()) {
		// full/over collateralized
		backingOut.Amount = burnInUSD.QuoRoundUp(backingPrice).RoundInt()
	} else if collateralRatio.IsZero() {
		// algorithmic
		lionOut.Amount = burnInUSD.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingOut.Amount = burnInUSD.Mul(collateralRatio).QuoRoundUp(backingPrice).RoundInt()
		lionOut.Amount = burnInUSD.Mul(sdk.OneDec().Sub(collateralRatio)).QuoRoundUp(lionPrice).RoundInt()
	}

	moduleOwnedBacking := m.Keeper.bankKeeper.GetBalance(ctx, m.Keeper.accountKeeper.GetModuleAddress(types.ModuleName), msg.BackingOutMin.Denom)
	if moduleOwnedBacking.IsLT(backingOut) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinInsufficient, "backing coin out(%s) < balance(%s)", backingOut, moduleOwnedBacking)
	}

	if backingOut.IsLT(msg.BackingOutMin) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinSlippage, "backing coin out: %s", backingOut)
	}
	if lionOut.IsLT(msg.LionOutMin) {
		return nil, sdkerrors.Wrapf(types.ErrLionCoinSlippage, "lion coin out: %s", lionOut)
	}

	// take mer stablecoin
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(burnIn))
	if err != nil {
		return nil, err
	}
	// burn mer
	err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burn))
	if err != nil {
		return nil, err
	}
	// send mer fee to oracle
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(burnFee))
	if err != nil {
		return nil, err
	}

	// mint lion
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(lionOut))
	if err != nil {
		return nil, err
	}
	// send backing and lion to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(backingOut, lionOut))
	if err != nil {
		return nil, err
	}

	// TODO: event

	return &types.MsgBurnBySwapResponse{
		BurnFee:    burnFee,
		BackingOut: backingOut,
		LionOut:    lionOut,
	}, nil
}

func (m msgServer) MintByCollateral(c context.Context, msg *types.MsgMintByCollateral) (*types.MsgMintByCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	receiver := sender
	if len(msg.To) > 0 {
		// user specifies receiver
		receiver, err = sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			return nil, err
		}
	}

	collateralPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, msg.CollateralDenom)
	if err != nil {
		return nil, err
	}
	lionPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroLionDenom)
	if err != nil {
		return nil, err
	}
	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price lower bound
	merPriceLowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(m.Keeper.MintPriceBias(ctx)))
	if merPrice.LT(merPriceLowerBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, msg.CollateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", msg.CollateralDenom)
	}
	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", msg.CollateralDenom)
	}

	pool, found := m.Keeper.GetPoolCollateral(ctx, msg.CollateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoCollateral, "account has no collateral: %s", msg.CollateralDenom)
	}
	acc, found := m.Keeper.GetAccountCollateral(ctx, sender, msg.CollateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoCollateral, "account has no collateral: %s", msg.CollateralDenom)
	}

	// settle interestFee fee
	settleInterestFee(ctx, &acc, &pool, collateralParams.InterestFee)

	// compute mint amount
	mintFee := computeFee(msg.MintOut, collateralParams.MintFee)
	mint := msg.MintOut.Add(mintFee)

	// update debt
	acc.MerDebt = acc.MerDebt.Add(mint)
	pool.MerDebt = pool.MerDebt.Add(mint)

	// compute actual catalytic lion
	merTotal := acc.MerDebt.Add(acc.MerByLion)
	bestCatalyticLionInUSD := merTotal.Amount.ToDec().Mul(*collateralParams.CatalyticLionRatio)
	lionInMaxInUSD := msg.LionInMax.Amount.ToDec().Mul(lionPrice).TruncateInt()
	catalyticLionInUSD := sdk.MinDec(bestCatalyticLionInUSD, acc.MerByLion.Amount.Add(lionInMaxInUSD).ToDec()).TruncateInt()

	// compute actual lion-in
	lionInInUSD := catalyticLionInUSD.Sub(acc.MerByLion.Amount)
	if lionInInUSD.IsNegative() {
		lionInInUSD = sdk.ZeroInt()
	} else {
		acc.MerByLion = acc.MerByLion.AddAmount(lionInInUSD)
		pool.MerByLion = pool.MerByLion.AddAmount(lionInInUSD)
		acc.MerDebt = acc.MerDebt.SubAmount(lionInInUSD)
		pool.MerDebt = pool.MerDebt.SubAmount(lionInInUSD)
	}
	lionIn := sdk.NewCoin(merlion.MicroLionDenom, lionInInUSD.ToDec().Quo(lionPrice).TruncateInt())

	// compute actual catalytic ratio and max loan-to-value
	catalyticRatio := catalyticLionInUSD.ToDec().Quo(bestCatalyticLionInUSD)
	maxLoanToValue := collateralParams.BasicLoanToValue.Add(collateralParams.LoanToValue.Sub(*collateralParams.BasicLoanToValue).Mul(catalyticRatio))

	// check max mintable mer
	collateralInUSD := acc.Collateral.Amount.ToDec().Mul(collateralPrice)
	maxMerMintable := collateralInUSD.Mul(maxLoanToValue)
	if maxMerMintable.LT(acc.MerDebt.Amount.ToDec()) {
		return nil, sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "account has insufficient collateral: %s", msg.CollateralDenom)
	}

	// eventually update collateral
	m.Keeper.SetPoolCollateral(ctx, pool)
	m.Keeper.SetAccountCollateral(ctx, sender, acc)

	// take lion and burn it
	if lionIn.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(lionIn))
		if err != nil {
			return nil, err
		}
		err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lionIn))
		if err != nil {
			return nil, err
		}
	}

	// mint mer
	err = m.Keeper.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mint))
	if err != nil {
		return nil, err
	}
	// send mer to receiver
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(msg.MintOut))
	if err != nil {
		return nil, err
	}
	// send mint fee to oracle
	if mintFee.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(mintFee))
		if err != nil {
			return nil, err
		}
	}

	// TODO: event

	return &types.MsgMintByCollateralResponse{
		MintFee: mintFee,
		LionIn:  lionIn,
	}, nil
}

func (m msgServer) BurnByCollateral(c context.Context, msg *types.MsgBurnByCollateral) (*types.MsgBurnByCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}

	merPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, merlion.MicroUSDDenom)
	if err != nil {
		return nil, err
	}

	// check price upper bound
	merPriceUpperBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Add(m.Keeper.BurnPriceBias(ctx)))
	if merPrice.GT(merPriceUpperBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooHigh, "%s price too high: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, msg.CollateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", msg.CollateralDenom)
	}
	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", msg.CollateralDenom)
	}

	pool, found := m.Keeper.GetPoolCollateral(ctx, msg.CollateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoCollateral, "account has no collateral: %s", msg.CollateralDenom)
	}
	acc, found := m.Keeper.GetAccountCollateral(ctx, sender, msg.CollateralDenom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoCollateral, "account has no collateral: %s", msg.CollateralDenom)
	}

	// settle interestFee fee
	settleInterestFee(ctx, &acc, &pool, collateralParams.InterestFee)

	// compute burn-in
	if !acc.MerDebt.IsPositive() {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoDebt, "account has no debt for %s collateral", msg.CollateralDenom)
	}
	repayIn := sdk.NewCoin(msg.RepayInMax.Denom, sdk.MinInt(acc.MerDebt.Amount, msg.RepayInMax.Amount))
	interestFee := sdk.NewCoin(msg.RepayInMax.Denom, sdk.MinInt(acc.LastInterest.Amount, repayIn.Amount))
	burn := repayIn.Sub(interestFee)

	// update debt
	acc.LastInterest = acc.LastInterest.Sub(interestFee)
	acc.MerDebt = acc.MerDebt.Sub(repayIn)
	pool.MerDebt = pool.MerDebt.Sub(repayIn)

	// eventually update collateral
	m.Keeper.SetAccountCollateral(ctx, sender, acc)
	m.Keeper.SetPoolCollateral(ctx, pool)

	// take mer
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(repayIn))
	if err != nil {
		return nil, err
	}
	// burn mer
	if burn.IsPositive() {
		err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(burn))
		if err != nil {
			return nil, err
		}
	}
	// send fee to oracle
	if interestFee.IsPositive() {
		err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(interestFee))
		if err != nil {
			return nil, err
		}
	}

	// TODO: event

	return &types.MsgBurnByCollateralResponse{
		RepayIn: repayIn,
	}, nil
}

func (m msgServer) DepositCollateral(c context.Context, msg *types.MsgDepositCollateral) (*types.MsgDepositCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	receiver := sender
	if len(msg.To) > 0 {
		// user specifies receiver
		receiver, err = sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			return nil, err
		}
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, msg.Collateral.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", msg.Collateral.Denom)
	}

	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", msg.Collateral.Denom)
	}

	pool, found := m.Keeper.GetPoolCollateral(ctx, msg.Collateral.Denom)
	if !found {
		pool = types.PoolCollateral{
			Collateral: sdk.NewCoin(msg.Collateral.Denom, sdk.ZeroInt()),
			MerDebt:    sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
			MerByLion:  sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
		}
	}

	if pool.Collateral.Amount.Add(msg.Collateral.Amount).GT(*collateralParams.MaxCollateral) {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCeiling, "existing(%s) + new(%s) > ceiling(%s)", pool.Collateral.Amount, msg.Collateral.Amount, collateralParams.MaxCollateral)
	}

	acc, found := m.Keeper.GetAccountCollateral(ctx, receiver, msg.Collateral.Denom)
	if !found {
		acc = types.AccountCollateral{
			Account:             receiver.String(),
			Collateral:          sdk.NewCoin(msg.Collateral.Denom, sdk.ZeroInt()),
			MerDebt:             sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
			MerByLion:           sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
			LastInterest:        sdk.NewCoin(merlion.MicroUSDDenom, sdk.ZeroInt()),
			LastSettlementBlock: ctx.BlockHeight(),
		}
	} else {
		settleInterestFee(ctx, &acc, &pool, collateralParams.InterestFee)
	}

	acc.Collateral = acc.Collateral.Add(msg.Collateral)
	pool.Collateral = pool.Collateral.Add(msg.Collateral)
	m.Keeper.SetAccountCollateral(ctx, receiver, acc)
	m.Keeper.SetPoolCollateral(ctx, pool)

	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(msg.Collateral))
	if err != nil {
		return nil, err
	}

	// TODO: event

	return &types.MsgDepositCollateralResponse{}, nil
}

func (m msgServer) RedeemCollateral(c context.Context, msg *types.MsgRedeemCollateral) (*types.MsgRedeemCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, err
	}
	receiver := sender
	if len(msg.To) > 0 {
		// user specifies receiver
		receiver, err = sdk.AccAddressFromBech32(msg.To)
		if err != nil {
			return nil, err
		}
	}

	collateralParams, found := m.Keeper.GetCollateralRiskParams(ctx, msg.Collateral.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinNotFound, "collateral coin denomination not found: %s", msg.Collateral.Denom)
	}

	if !collateralParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrCollateralCoinDisabled, "collateral coin disabled: %s", msg.Collateral.Denom)
	}

	pool, found := m.Keeper.GetPoolCollateral(ctx, msg.Collateral.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoCollateral, "account has no collateral: %s", msg.Collateral.Denom)
	}

	acc, found := m.Keeper.GetAccountCollateral(ctx, sender, msg.Collateral.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrAccountNoCollateral, "account has no collateral: %s", msg.Collateral.Denom)
	}

	// update collateral
	pool.Collateral = pool.Collateral.Sub(msg.Collateral)
	acc.Collateral = acc.Collateral.Sub(msg.Collateral)

	settleInterestFee(ctx, &acc, &pool, collateralParams.InterestFee)

	collateralPrice, err := m.Keeper.oracleKeeper.GetExchangeRate(ctx, msg.Collateral.Denom)
	if err != nil {
		return nil, err
	}

	merTotal := acc.MerDebt.Add(acc.MerByLion)
	catalyticRatio := acc.MerByLion.Amount.ToDec().QuoInt(merTotal.Amount).Quo(*collateralParams.CatalyticLionRatio)
	if catalyticRatio.GT(sdk.OneDec()) {
		catalyticRatio = sdk.OneDec()
	}
	maxLoanToValue := collateralParams.BasicLoanToValue.Add(collateralParams.LoanToValue.Sub(*collateralParams.BasicLoanToValue).Mul(catalyticRatio))

	collateralInUSD := acc.Collateral.Amount.ToDec().Mul(collateralPrice)
	if acc.MerDebt.Amount.ToDec().LT(collateralInUSD.Mul(maxLoanToValue)) {
		return nil, sdkerrors.Wrapf(types.ErrAccountInsufficientCollateral, "account has insufficient collateral: %s", msg.Collateral.Denom)
	}

	// eventually persist collateral
	m.Keeper.SetPoolCollateral(ctx, pool)
	m.Keeper.SetAccountCollateral(ctx, sender, acc)

	err = m.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(msg.Collateral))
	if err != nil {
		return nil, err
	}

	// TODO: event

	return &types.MsgRedeemCollateralResponse{}, nil
}

func (m msgServer) LiquidateCollateral(c context.Context, msg *types.MsgLiquidateCollateral) (*types.MsgLiquidateCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) BuyBacking(c context.Context, msg *types.MsgBuyBacking) (*types.MsgBuyBackingResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) SellBacking(c context.Context, msg *types.MsgSellBacking) (*types.MsgSellBackingResponse, error) {
	// TODO implement me
	panic("implement me")
}

func settleInterestFee(ctx sdk.Context, acc *types.AccountCollateral, pool *types.PoolCollateral, apr *sdk.Dec) {
	if apr != nil {
		period := ctx.BlockHeight() - acc.LastSettlementBlock
		// principal debt, excluding interest debt
		principalDebt := acc.MerDebt.Sub(acc.LastInterest)
		interestOfPeriod := principalDebt.Amount.ToDec().Mul(*apr).MulInt64(period).QuoInt64(int64(merlion.BlocksPerYear)).RoundInt()
		// update remaining interest accumulation
		acc.LastInterest = acc.LastInterest.AddAmount(interestOfPeriod)
		// update debt
		acc.MerDebt = acc.MerDebt.AddAmount(interestOfPeriod)
		pool.MerDebt = pool.MerDebt.AddAmount(interestOfPeriod)
	}
	// update settlement block
	acc.LastSettlementBlock = ctx.BlockHeight()
}

func computeFee(coin sdk.Coin, rate *sdk.Dec) sdk.Coin {
	amt := sdk.ZeroInt()
	if rate != nil {
		amt = coin.Amount.ToDec().Mul(*rate).TruncateInt()
	}
	return sdk.NewCoin(coin.Denom, amt)
}
