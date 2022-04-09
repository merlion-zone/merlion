package types

import govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

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
	// TODO implement me
	panic("implement me")
}

func (m *RegisterCollateralProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RegisterCollateralProposal) ProposalType() string {
	return ProposalTypeRegisterCollateral
}

func (m *RegisterCollateralProposal) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *SetBackingRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *SetBackingRiskParamsProposal) ProposalType() string {
	return ProposalTypeSetBackingRiskParams
}

func (m *SetBackingRiskParamsProposal) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *SetCollateralRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *SetCollateralRiskParamsProposal) ProposalType() string {
	return ProposalTypeSetCollateralRiskParams
}

func (m *SetCollateralRiskParamsProposal) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *BatchSetBackingRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *BatchSetBackingRiskParamsProposal) ProposalType() string {
	return ProposalTypeBatchSetBackingRiskParams
}

func (m *BatchSetBackingRiskParamsProposal) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}

func (m *BatchSetCollateralRiskParamsProposal) ProposalRoute() string {
	return RouterKey
}

func (m *BatchSetCollateralRiskParamsProposal) ProposalType() string {
	return ProposalTypeBatchSetCollateralRiskParams
}

func (m *BatchSetCollateralRiskParamsProposal) ValidateBasic() error {
	// TODO implement me
	panic("implement me")
}
