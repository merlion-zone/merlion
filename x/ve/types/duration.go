package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegulatedUnixTime(timestamp uint64) uint64 {
	// divide and round down to align to regulated periods
	return timestamp / RegulatedPeriod * RegulatedPeriod
}

func RegulatedUnixTimeFromNow(ctx sdk.Context, seconds uint64) uint64 {
	timestamp := ctx.BlockTime().Add(time.Duration(seconds) * time.Second).Unix()
	return RegulatedUnixTime(uint64(timestamp))
}

func NextRegulatedUnixTime(timestamp uint64) uint64 {
	CheckRegulatedUnixTime(timestamp)
	if timestamp > MaxUnixTime-RegulatedPeriod {
		panic("too large unix time")
	}
	return timestamp + RegulatedPeriod
}

func PreviousRegulatedUnixTime(timestamp uint64) uint64 {
	CheckRegulatedUnixTime(timestamp)
	if timestamp < RegulatedPeriod {
		panic("too small unix time")
	}
	return timestamp - RegulatedPeriod
}

func CheckRegulatedUnixTime(timestamp uint64) {
	if RegulatedUnixTime(timestamp) != timestamp {
		panic("invalid regulated unix time")
	}
}
