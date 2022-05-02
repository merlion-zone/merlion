package types

import (
	"math"

	merlion "github.com/merlion-zone/merlion/types"
)

const (
	EmptyVeID = 0
	FirstVeID = 1
	MaxVeID   = math.MaxUint64 - 1

	MaxLockTime = merlion.SecondsPer4Years

	MaxUnixTime = math.MaxInt64

	RegulatedPeriod = merlion.SecondsPerWeek

	EmptyEpoch = 0
	FirstEpoch = 1
)
