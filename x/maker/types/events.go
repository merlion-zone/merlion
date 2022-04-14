package types

// maker module event types
const (
	EventTypeMintBySwap          = "mint_by_swap"
	EventTypeBurnBySwap          = "burn_by_swap"
	EventTypeBuyBacking          = "buy_backing"
	EventTypeSellBacking         = "sell_backing"
	EventTypeMintByCollateral    = "mint_by_collateral"
	EventTypeBurnByCollateral    = "burn_by_collateral"
	EventTypeDepositCollateral   = "deposit_collateral"
	EventTypeRedeemCollateral    = "redeem_collateral"
	EventTypeLiquidateCollateral = "liquidate_collateral"

	AttributeKeySender   = "sender"
	AttributeKeyReceiver = "receiver"
	AttributeKeyCoinIn   = "coin_in"
	AttributeKeyCoinOut  = "coin_out"
	AttributeKeyFee      = "fee"

	EventTypeRegisterBacking         = "register_backing"
	EventTypeRegisterCollateral      = "register_collateral"
	EventTypeSetBackingRiskParams    = "set_backing_risk_params"
	EventTypeSetCollateralRiskParams = "set_collateral_risk_params"

	AttributeKeyRiskParams = "risk_params"

	AttributeValueCategory = ModuleName
)
