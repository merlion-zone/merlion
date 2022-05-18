package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// DefaultGenesis returns the default vesting genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return gs.Params.Validate()
}

func (a AllocationAddresses) GetStrategicReserveCustodianAddr() sdk.AccAddress {
	if len(a.StrategicReserveCustodianAddr) == 0 {
		return sdk.AccAddress{}
	}
	srca, err := sdk.AccAddressFromBech32(a.StrategicReserveCustodianAddr)
	if err != nil {
		panic(err)
	}
	return srca
}

func (a AllocationAddresses) GetTeamVestingAddr() sdk.AccAddress {
	if len(a.TeamVestingAddr) == 0 {
		return sdk.AccAddress{}
	}
	tva, err := sdk.AccAddressFromBech32(a.TeamVestingAddr)
	if err != nil {
		panic(err)
	}
	return tva
}

func (a Airdrop) Empty() bool {
	return len(a.TargetAddr) == 0
}

func (a Airdrop) GetTargetAddr() sdk.AccAddress {
	if len(a.TargetAddr) == 0 {
		return sdk.AccAddress{}
	}
	ta, err := sdk.AccAddressFromBech32(a.TargetAddr)
	if err != nil {
		panic(err)
	}
	return ta
}
