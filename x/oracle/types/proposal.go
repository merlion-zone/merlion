package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeRegisterTarget = "RegisterTarget"
)

var (
	_ govtypes.Content = &RegisterTargetProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeRegisterTarget)
	govtypes.RegisterProposalTypeCodec(&RegisterTargetProposal{}, "oracle/RegisterTargetProposal")
}

func (m *RegisterTargetProposal) ProposalRoute() string {
	return RouterKey
}

func (m *RegisterTargetProposal) ProposalType() string {
	return ProposalTypeRegisterTarget
}

func (m *RegisterTargetProposal) ValidateBasic() error {
	return validateTargetParams(&m.TargetParams)
}

func validateTargetParams(params *TargetParams) error {
	err := sdk.ValidateDenom(params.Denom)
	if err != nil {
		return err
	}
	if params.Source <= TARGET_SOURCE_UNSPECIFIED {
		return fmt.Errorf("target source must be specified")
	}
	// TODO
	return nil
}
