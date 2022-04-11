package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// nolint
const (
	MicroLionDenom = "ulion"
	MicroUSDDenom  = "uusd"

	MicroUnit = int64(1e6)
)

var (
	MicroUSDTarget = sdk.OneDec()
)
