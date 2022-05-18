package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	ethermint "github.com/tharsis/ethermint/types"
	"gopkg.in/yaml.v2"
)

// Parameter keys
var (
	KeyAllocationAmounts = []byte("AllocationAmounts")
)

// Default parameter values
var (
	DefaultTotalAmount                = sdk.NewInt(1000_000_000).Mul(ethermint.PowerReduction)
	DefaultAirdropAmountRate          = sdk.NewDecWithPrec(5, 2)  // 5%
	DefaultVeVestingAmountRate        = sdk.NewDecWithPrec(45, 2) // 45%
	DefaultStakingRewardAmountRate    = sdk.NewDecWithPrec(5, 2)  // 5%
	DefaultCommunityPoolAmountRate    = sdk.NewDecWithPrec(5, 2)  // 5%
	DefaultStrategicReserveAmountRate = sdk.NewDecWithPrec(20, 2) // 20%
	DefaultTeamVestingAmountRate      = sdk.NewDecWithPrec(20, 2) // 20%
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	airdropAmount := DefaultTotalAmount.ToDec().Mul(DefaultAirdropAmountRate).TruncateInt()
	veVestingAmount := DefaultTotalAmount.ToDec().Mul(DefaultVeVestingAmountRate).TruncateInt()
	stakingRewardAmount := DefaultTotalAmount.ToDec().Mul(DefaultStakingRewardAmountRate).TruncateInt()
	communityPoolAmount := DefaultTotalAmount.ToDec().Mul(DefaultCommunityPoolAmountRate).TruncateInt()
	strategicReserveAmount := DefaultTotalAmount.ToDec().Mul(DefaultStrategicReserveAmountRate).TruncateInt()
	teamVestingAmount := DefaultTotalAmount.ToDec().Mul(DefaultTeamVestingAmountRate).TruncateInt()

	if !DefaultTotalAmount.Equal(airdropAmount.Add(veVestingAmount).Add(stakingRewardAmount).Add(communityPoolAmount).Add(strategicReserveAmount).Add(teamVestingAmount)) {
		panic("sum of all allocation amounts must be equal to total amount")
	}

	return Params{
		Allocation: AllocationAmounts{
			TotalAmount:            DefaultTotalAmount,
			AirdropAmount:          airdropAmount,
			VeVestingAmount:        veVestingAmount,
			StakingRewardAmount:    stakingRewardAmount,
			CommunityPoolAmount:    communityPoolAmount,
			StrategicReserveAmount: strategicReserveAmount,
			TeamVestingAmount:      teamVestingAmount,
		},
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAllocationAmounts, &p.Allocation, validateAllocationAmounts),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return validateAllocationAmounts(p.Allocation)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateAllocationAmounts(i interface{}) error {
	v, ok := i.(AllocationAmounts)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.TotalAmount.Equal(v.AirdropAmount.Add(v.VeVestingAmount).Add(v.StakingRewardAmount).Add(v.CommunityPoolAmount).Add(v.StrategicReserveAmount).Add(v.TeamVestingAmount)) {
		return fmt.Errorf("sum of all allocation amounts must be equal to total amount")
	}

	return nil
}
