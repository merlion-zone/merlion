package types

import (
	"math"

	merlion "github.com/merlion-zone/merlion/types"
)

const (
	EmptyVeID = 0
	MaxVeID   = math.MaxUint64 - 1

	MaxLockTime = merlion.SecondsPer4Years

	MaxUnixTime = math.MaxInt64
)
