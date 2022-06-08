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
	KeyBackingRatioStep           = []byte("BackingRatioStep")
	KeyBackingRatioPriceBand      = []byte("BackingRatioPriceBand")
	KeyBackingRatioCooldownPeriod = []byte("BackingRatioCooldownPeriod")
	KeyMintPriceBias              = []byte("MintPriceBias")
	KeyBurnPriceBias              = []byte("BurnPriceBias")
	KeyRebackBonus                = []byte("RebackBonus")
	KeyLiquidationCommissionFee   = []byte("LiquidationCommissionFee")
)

// Default parameter values
var (
	DefaultBackingRatioStep           = sdk.NewDecWithPrec(25, 4)    // 0.25%
	DefaultBackingRatioPriceBand      = sdk.NewDecWithPrec(5, 3)     // 0.5%
	DefaultBackingRatioCooldownPeriod = int64(merlion.BlocksPerHour) // 600
	DefaultMintPriceBias              = sdk.NewDecWithPrec(1, 2)     // 1%
	DefaultBurnPriceBias              = sdk.NewDecWithPrec(1, 2)     // 1%
	DefaultRebackBonus                = sdk.NewDecWithPrec(75, 4)    // 0.75%
	DefaultLiquidationCommissionFee   = sdk.NewDecWithPrec(10, 2)    // 10%
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		BackingRatioStep:           DefaultBackingRatioStep,
		BackingRatioPriceBand:      DefaultBackingRatioPriceBand,
		BackingRatioCooldownPeriod: DefaultBackingRatioCooldownPeriod,
		MintPriceBias:              DefaultMintPriceBias,
		BurnPriceBias:              DefaultBurnPriceBias,
		RebackBonus:                DefaultRebackBonus,
		LiquidationCommissionFee:   DefaultLiquidationCommissionFee,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBackingRatioStep, &p.BackingRatioStep, validateBackingRatioStep),
		paramtypes.NewParamSetPair(KeyBackingRatioPriceBand, &p.BackingRatioPriceBand, validateBackingRatioPriceBand),
		paramtypes.NewParamSetPair(KeyBackingRatioCooldownPeriod, &p.BackingRatioCooldownPeriod, validateBackingRatioCooldownPeriod),
		paramtypes.NewParamSetPair(KeyMintPriceBias, &p.MintPriceBias, validateMintBurnPriceBias),
		paramtypes.NewParamSetPair(KeyBurnPriceBias, &p.BurnPriceBias, validateMintBurnPriceBias),
		paramtypes.NewParamSetPair(KeyRebackBonus, &p.RebackBonus, validateRebackBonus),
		paramtypes.NewParamSetPair(KeyLiquidationCommissionFee, &p.LiquidationCommissionFee, validateLiquidationCommissionFee),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if !p.BackingRatioStep.IsPositive() || p.BackingRatioStep.GT(sdk.OneDec()) {
		return fmt.Errorf("backing ratio adjusting step should be a value between (0,1], is %s", p.BackingRatioStep)
	}
	if !p.BackingRatioPriceBand.IsPositive() || p.BackingRatioPriceBand.GT(sdk.OneDec()) {
		return fmt.Errorf("price band for adjusting backing ratio should be a value between (0,1], is %s", p.BackingRatioPriceBand)
	}
	if p.BackingRatioCooldownPeriod <= 0 {
		return fmt.Errorf("cooldown period for adjusting backing ratio should be positive, is %d", p.BackingRatioCooldownPeriod)
	}
	if !p.MintPriceBias.IsPositive() || p.MintPriceBias.GT(sdk.OneDec()) {
		return fmt.Errorf("mint price bias ratio should be a value between (0,1], is %s", p.MintPriceBias)
	}
	if !p.BurnPriceBias.IsPositive() || p.BurnPriceBias.GT(sdk.OneDec()) {
		return fmt.Errorf("burn price bias ratio should be a value between (0,1], is %s", p.MintPriceBias)
	}
	if p.RebackBonus.IsNegative() || p.RebackBonus.GT(sdk.OneDec()) {
		return fmt.Errorf("reback bonus ratio should be a value between [0,1], is %s", p.RebackBonus)
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

func validateBackingRatioStep(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("backing ratio adjusting step must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("backing ratio adjusting step is too large: %s", v)
	}

	return nil
}

func validateBackingRatioPriceBand(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsPositive() {
		return fmt.Errorf("price band for adjusting backing ratio must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("price band for adjusting backing ratio is too large: %s", v)
	}

	return nil
}

func validateBackingRatioCooldownPeriod(i interface{}) error {
	v, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("cooldown period for adjusting backing ratio must be positive: %d", v)
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

func validateRebackBonus(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("reback bonus ratio must be positive or zero: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("reback bonus ratio is too large: %s", v)
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
