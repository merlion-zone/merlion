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
	lowerBound := merlion.MicroUSDTarget.Mul(sdk.OneDec().Sub(m.Keeper.CollateralRatioPriceBand(ctx)))
	if merPrice.LT(lowerBound) {
		return nil, sdkerrors.Wrapf(types.ErrMerPriceTooLow, "%s price too low: %s", merlion.MicroUSDDenom, merPrice)
	}

	collateralRaio := m.Keeper.GetCollateralRatio(ctx)
	backingParams, found := m.Keeper.GetBackingRiskParams(ctx, msg.BackingInMax.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinNotFound, "backing coin denomination not found: %s", msg.BackingInMax.Denom)
	}
	if !backingParams.Enabled {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinDisabled, "backing coin disabled: %s", msg.BackingInMax.Denom)
	}

	mintOut := msg.MintOut
	mintFee := sdk.NewCoin(mintOut.Denom, mintOut.Amount.ToDec().Mul(*backingParams.MintFee).RoundInt())
	mintTotal := mintOut.AddAmount(mintFee.Amount)
	mintTotalDec := mintTotal.Amount.ToDec()

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
	if collateralRaio.GTE(sdk.OneDec()) || msg.LionInMax.IsZero() {
		// full/over collateralized, or user selects full collateralization
		backingNeeded.Amount = mintTotalDec.QuoRoundUp(backingPrice).RoundInt()
	} else if collateralRaio.IsZero() {
		// algorithmic
		lionNeeded.Amount = mintTotalDec.QuoRoundUp(lionPrice).RoundInt()
	} else {
		// fractional
		backingNeeded.Amount = mintTotalDec.Mul(collateralRaio).QuoRoundUp(backingPrice).RoundInt()
		lionNeeded.Amount = mintTotalDec.Mul(sdk.OneDec().Sub(collateralRaio)).QuoRoundUp(lionPrice).RoundInt()
	}

	if msg.BackingInMax.IsLT(backingNeeded) {
		return nil, sdkerrors.Wrapf(types.ErrBackingCoinSlippage, "backing coin needed: %s", backingNeeded)
	}
	if msg.LionInMax.IsLT(lionNeeded) {
		return nil, sdkerrors.Wrapf(types.ErrLionCoinSlippage, "lion coin needed: %s", lionNeeded)
	}

	// take backing coin and lion coin
	err = m.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.NewCoins(backingNeeded, lionNeeded))
	if err != nil {
		return nil, err
	}
	// burn lion
	err = m.Keeper.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(lionNeeded))
	if err != nil {
		return nil, err
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
	err = m.Keeper.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, oracletypes.ModuleName, sdk.NewCoins(mintFee))
	if err != nil {
		return nil, err
	}

	// TODO: event

	return &types.MsgMintBySwapResponse{
		MintFee:   mintFee,
		BackingIn: backingNeeded,
		LionIn:    lionNeeded,
	}, nil
}

func (m msgServer) BurnBySwap(ctx context.Context, swap *types.MsgBurnBySwap) (*types.MsgBurnBySwapResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) MintByCollateral(ctx context.Context, collateral *types.MsgMintByCollateral) (*types.MsgMintByCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) BurnByCollateral(ctx context.Context, collateral *types.MsgBurnByCollateral) (*types.MsgBurnByCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) DepositCollateral(ctx context.Context, collateral *types.MsgDepositCollateral) (*types.MsgDepositCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) RedeemCollateral(ctx context.Context, collateral *types.MsgRedeemCollateral) (*types.MsgRedeemCollateralResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) BuyBack(ctx context.Context, back *types.MsgBuyBack) (*types.MsgBuyBackResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (m msgServer) ReCollateralize(ctx context.Context, collateralize *types.MsgReCollateralize) (*types.MsgReCollateralizeResponse, error) {
	// TODO implement me
	panic("implement me")
}
