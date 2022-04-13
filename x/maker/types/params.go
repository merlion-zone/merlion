package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	merlion "github.com/merlion-zone/merlion/types"
	"gopkg.in/yaml.v2"
)

// Parameter keys
var (
	KeyCollateralRatioStep           = []byte("CollateralRatioStep")
	KeyCollateralRatioPriceBand      = []byte("CollateralRatioPriceBand")
	KeyCollateralRatioCooldownPeriod = []byte("CollateralRatioCooldownPeriod")
	KeyMintPriceBias                 = []byte("MintPriceBias")
	KeyBurnPriceBias                 = []byte("BurnPriceBias")
	KeyRecollateralizeBonus          = []byte("RecollateralizeBonus")
	KeyLiquidationCommissionFee      = []byte("LiquidationCommissionFee")
)

// Default parameter values
var (
	DefaultCollateralRatioStep           = sdk.NewDecWithPrec(25, 4)    // 0.25%
	DefaultCollateralRatioPriceBand      = sdk.NewDecWithPrec(5, 3)     // 0.5%
	DefaultCollateralRatioCooldownPeriod = int64(merlion.BlocksPerHour) // 600
	DefaultMintPriceBias                 = sdk.NewDecWithPrec(1, 2)     // 1%
	DefaultBurnPriceBias                 = sdk.NewDecWithPrec(1, 2)     // 1%
	DefaultRecollateralizeBonus          = sdk.NewDecWithPrec(75, 4)    // 0.75%
	DefaultLiquidationCommissionFee      = sdk.NewDecWithPrec(10, 2)    // 10%
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		CollateralRatioStep:           DefaultCollateralRatioStep,
		CollateralRatioPriceBand:      DefaultCollateralRatioPriceBand,
		CollateralRatioCooldownPeriod: DefaultCollateralRatioCooldownPeriod,
		MintPriceBias:                 DefaultMintPriceBias,
		BurnPriceBias:                 DefaultBurnPriceBias,
		RecollateralizeBonus:          DefaultRecollateralizeBonus,
		LiquidationCommissionFee:      DefaultLiquidationCommissionFee,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCollateralRatioStep, &p.CollateralRatioStep, validateCollateralRatioStep),
		paramtypes.NewParamSetPair(KeyCollateralRatioPriceBand, &p.CollateralRatioPriceBand, validateCollateralRatioPriceBand),
		paramtypes.NewParamSetPair(KeyCollateralRatioCooldownPeriod, &p.CollateralRatioCooldownPeriod, validateCollateralRatioCooldownPeriod),
		paramtypes.NewParamSetPair(KeyMintPriceBias, &p.MintPriceBias, validateMintBurnPriceBias),
		paramtypes.NewParamSetPair(KeyBurnPriceBias, &p.BurnPriceBias, validateMintBurnPriceBias),
		paramtypes.NewParamSetPair(KeyRecollateralizeBonus, &p.RecollateralizeBonus, validateRecollateralizeBonus),
		paramtypes.NewParamSetPair(KeyLiquidationCommissionFee, &p.LiquidationCommissionFee, validateLiquidationCommissionFee),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if !p.CollateralRatioStep.IsPositive() || p.CollateralRatioStep.GT(sdk.OneDec()) {
		return fmt.Errorf("collateral ratio adjusting step should be a value between (0,1], is %s", p.CollateralRatioStep)
	}
	if !p.CollateralRatioPriceBand.IsPositive() || p.CollateralRatioPriceBand.GT(sdk.OneDec()) {
		return fmt.Errorf("price band for adjusting collateral ratio should be a value between (0,1], is %s", p.CollateralRatioPriceBand)
	}
	if p.CollateralRatioCooldownPeriod <= 0 {
		return fmt.Errorf("cooldown period for adjusting collateral ratio should be positive, is %d", p.CollateralRatioCooldownPeriod)
	}
	if !p.MintPriceBias.IsPositive() || p.MintPriceBias.GT(sdk.OneDec()) {
		return fmt.Errorf("mint price bias ratio should be a value between (0,1], is %s", p.MintPriceBias)
	}
	if !p.BurnPriceBias.IsPositive() || p.BurnPriceBias.GT(sdk.OneDec()) {
		return fmt.Errorf("burn price bias ratio should be a value between (0,1], is %s", p.MintPriceBias)
	}
	if p.RecollateralizeBonus.IsNegative() || p.RecollateralizeBonus.GT(sdk.OneDec()) {
		return fmt.Errorf("recollateralization bonus ratio should be a value between [0,1], is %s", p.RecollateralizeBonus)
	}
	if p.LiquidationCommissionFee.IsNegative() || p.LiquidationCommissionFee.GT(sdk.OneDec()) {
		return fmt.Errorf("liquidation commission fee ratio should be a value between [0,1], is %s", p.LiquidationCommissionFee)
	}
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateCollateralRatioStep(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("collateral ratio adjusting step must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("collateral ratio adjusting step is too large: %s", v)
	}

	return nil
}

func validateCollateralRatioPriceBand(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("price band for adjusting collateral ratio must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("price band for adjusting collateral ratio is too large: %s", v)
	}

	return nil
}

func validateCollateralRatioCooldownPeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("cooldown period for adjusting collateral ratio must be positive: %d", v)
	}

	return nil
}

func validateMintBurnPriceBias(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("mint/burn price bias ratio must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("mint/burn price bias ratio is too large: %s", v)
	}

	return nil
}

func validateRecollateralizeBonus(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("recollateralization bonus ratio must be positive or zero: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("recollateralization bonus ratio is too large: %s", v)
	}

	return nil
}

func validateLiquidationCommissionFee(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("liquidation commission fee ratio must be positive or zero: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("liquidation commission fee ratio is too large: %s", v)
	}

	return nil
}
