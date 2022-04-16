package types

import (
	"fmt"
	"strings"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
)

const (
	ProposalTypeSetDenomMetaData = "SetDenomMetaData"
)

var (
	_ govtypes.Content = &SetDenomMetadataProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeSetDenomMetaData)
	govtypes.RegisterProposalTypeCodec(&SetDenomMetadataProposal{}, "bank/SetDenomMetadataProposal")
}

func (m *SetDenomMetadataProposal) ProposalRoute() string {
	return banktypes.RouterKey
}

func (m *SetDenomMetadataProposal) ProposalType() string {
	return ProposalTypeSetDenomMetaData
}

func (m *SetDenomMetadataProposal) ValidateBasic() error {
	if err := m.Metadata.Validate(); err != nil {
		return err
	}

	if err := ibctransfertypes.ValidateIBCDenom(m.Metadata.Base); err != nil {
		return err
	}

	if err := validateIBC(m.Metadata); err != nil {
		return err
	}

	return govtypes.ValidateAbstract(m)
}

func validateIBC(metadata banktypes.Metadata) error {
	// Check ibc/ denom
	denomSplit := strings.SplitN(metadata.Base, "/", 2)

	if denomSplit[0] == metadata.Base && strings.TrimSpace(metadata.Base) != "" {
		// Not IBC
		return nil
	}

	if len(denomSplit) != 2 || denomSplit[0] != ibctransfertypes.DenomPrefix {
		// NOTE: should be unaccessible (covered on ValidateIBCDenom)
		return fmt.Errorf("invalid metadata. %s denomination should be prefixed with the format 'ibc/", metadata.Base)
	}

	if !strings.Contains(metadata.Name, "channel-") {
		return fmt.Errorf("invalid metadata (Name) for ibc. %s should include channel", metadata.Name)
	}

	if !strings.HasPrefix(metadata.Symbol, "ibc") {
		return fmt.Errorf("invalid metadata (Symbol) for ibc. %s should include \"ibc\" prefix", metadata.Symbol)
	}

	return nil
}
