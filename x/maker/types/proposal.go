package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterBacking              = "RegisterBacking"
	ProposalTypeRegisterCollateral           = "RegisterCollateral"
	ProposalTypeSetBackingRiskParams         = "SetBackingRiskParams"
	ProposalTypeSetCollateralRiskParams      = "SetCollateralRiskParams"
	ProposalTypeBatchSetBackingRiskParams    = "BatchSetBackingRiskParams"
	ProposalTypeBatchSetCollateralRiskParams = "BatchSetCollateralRiskParams"
)

var (
	_ govtypes.Content = &RegisterBackingProposal{}
	_ govtypes.Content = &RegisterCollateralProposal{}
	_ govtypes.Content = &SetBackingRiskParamsProposal{}
	_ govtypes.Content = &SetCollateralRiskParamsProposal{}
	_ govtypes.Content = &BatchSetBackingRiskParamsProposal{}
	_ govtypes.Content = &BatchSetCollateralRiskParamsProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeRegisterBacking)
	govtypes.RegisterProposalType(ProposalTypeRegisterCollateral)
	govtypes.RegisterProposalType(ProposalTypeSetBackingRiskParams)
	govtypes.RegisterProposalType(ProposalTypeSetCollateralRiskParams)
	govtypes.RegisterProposalType(ProposalTypeBatchSetBackingRiskParams)
	govtypes.RegisterProposalType(ProposalTypeBatchSetCollateralRiskParams)
	govtypes.RegisterProposalTypeCodec(&RegisterBackingProposal{}, "maker/RegisterBackingProposal")
	govtypes.RegisterProposalTypeCodec(&RegisterCollateralProposal{}, "maker/RegisterCollateralProposal")
	govtypes.RegisterProposalTypeCodec(&SetBackingRiskParamsProposal{}, "maker/SetBackingRiskParamsProposal")
	govtypes.RegisterProposalTypeCodec(&SetCollateralRiskParamsProposal{}, "maker/SetCollateralRiskParamsProposal")
	govtypes.RegisterProposalTypeCodec(&BatchSetBackingRiskParamsProposal{}, "maker/BatchSetBackingRiskParamsProposal")
	govtypes.RegisterProposalTypeCodec(&BatchSetCollateralRiskParamsProposal{}, "maker/BatchSetCollateralRiskParamsProposal")
}

func (m *RegisterBackingProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RegisterBackingProposal) ProposalType() string {
	return ProposalTypeRegisterBacking
}

func (m *RegisterBackingProposal) ValidateBasic() error {
	return validateBackingRiskParams(&m.RiskParams)
}

func (m *RegisterCollateralProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RegisterCollateralProposal) ProposalType() string {
	return ProposalTypeRegisterCollateral
}

func (m *RegisterCollateralProposal) ValidateBasic() error {
	return validateCollateralRiskParams(&m.RiskParams)
}

func (m *SetBackingRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *SetBackingRiskParamsProposal) ProposalType() string {
	return ProposalTypeSetBackingRiskParams
}

func (m *SetBackingRiskParamsProposal) ValidateBasic() error {
	return validateBackingRiskParams(&m.RiskParams)
}

func (m *SetCollateralRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *SetCollateralRiskParamsProposal) ProposalType() string {
	return ProposalTypeSetCollateralRiskParams
}

func (m *SetCollateralRiskParamsProposal) ValidateBasic() error {
	return validateCollateralRiskParams(&m.RiskParams)
}

func (m *BatchSetBackingRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *BatchSetBackingRiskParamsProposal) ProposalType() string {
	return ProposalTypeBatchSetBackingRiskParams
}

func (m *BatchSetBackingRiskParamsProposal) ValidateBasic() error {
	for _, params := range m.RiskParams {
		if err := validateBackingRiskParams(&params); err != nil {
			return err
		}
	}
	return nil
}

func (m *BatchSetCollateralRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *BatchSetCollateralRiskParamsProposal) ProposalType() string {
	return ProposalTypeBatchSetCollateralRiskParams
}

func (m *BatchSetCollateralRiskParamsProposal) ValidateBasic() error {
	for _, params := range m.RiskParams {
		if err := validateCollateralRiskParams(&params); err != nil {
			return err
		}
	}
	return nil
}

func validateBackingRiskParams(params *BackingRiskParams) error {
	if params.MaxBacking != nil && params.MaxBacking.IsNegative() {
		return fmt.Errorf("max backing value must be not negative")
	}
	if params.MaxMerMint != nil && params.MaxMerMint.IsNegative() {
		return fmt.Errorf("max mer mint value must be not negative")
	}
	if params.MintFee != nil && (params.MintFee.IsNegative() || params.MintFee.GT(sdk.OneDec())) {
		return fmt.Errorf("mint fee must be in [0, 1]")
	}
	if params.BurnFee != nil && (params.BurnFee.IsNegative() || params.BurnFee.GT(sdk.OneDec())) {
		return fmt.Errorf("burn fee must be in [0, 1]")
	}
	if params.BuybackFee != nil && (params.BuybackFee.IsNegative() || params.BuybackFee.GT(sdk.OneDec())) {
		return fmt.Errorf("buyback fee must be in [0, 1]")
	}
	if params.RecollateralizeFee != nil && (params.RecollateralizeFee.IsNegative() || params.RecollateralizeFee.GT(sdk.OneDec())) {
		return fmt.Errorf("recollateralize fee must be in [0, 1]")
	}
	return nil
}

func validateCollateralRiskParams(params *CollateralRiskParams) error {
	if params.MaxCollateral != nil && params.MaxCollateral.IsNegative() {
		return fmt.Errorf("max collateral value must be not negative")
	}
	if params.MaxMerMint != nil && params.MaxMerMint.IsNegative() {
		return fmt.Errorf("max mer mint value must be not negative")
	}
	if params.LiquidationThreshold != nil && (params.LiquidationThreshold.IsNegative() || params.LiquidationThreshold.GT(sdk.OneDec())) {
		return fmt.Errorf("liquidation threshold must be in [0, 1]")
	}
	if params.LoanToValue != nil && (params.LoanToValue.IsNegative() || params.LoanToValue.GT(sdk.OneDec())) {
		return fmt.Errorf("loan-to-value must be in [0, 1]")
	}
	if params.BasicLoanToValue != nil && (params.BasicLoanToValue.IsNegative() || params.BasicLoanToValue.GT(sdk.OneDec())) {
		return fmt.Errorf("basic loan-to-value must be in [0, 1]")
	}
	if params.CatalyticLionRatio != nil && (params.CatalyticLionRatio.IsNegative() || params.CatalyticLionRatio.GT(sdk.OneDec())) {
		return fmt.Errorf("catalytic lion ratio must be in [0, 1]")
	}
	if params.LiquidationFee != nil && (params.LiquidationFee.IsNegative() || params.LiquidationFee.GT(sdk.OneDec())) {
		return fmt.Errorf("liquidation fee must be in [0, 1]")
	}
	if params.MintFee != nil && (params.MintFee.IsNegative() || params.MintFee.GT(sdk.OneDec())) {
		return fmt.Errorf("mint fee must be in [0, 1]")
	}
	if params.InterestFee != nil && (params.InterestFee.IsNegative() || params.InterestFee.GT(sdk.OneDec())) {
		return fmt.Errorf("interest fee must be in [0, 1]")
	}
	return nil
}
