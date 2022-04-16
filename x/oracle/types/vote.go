package types

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewAggregateExchangeRatePrevote creates an AggregateExchangeRatePrevote instance.
func NewAggregateExchangeRatePrevote(hash AggregateVoteHash, voter sdk.ValAddress, submitBlock uint64) AggregateExchangeRatePrevote {
	return AggregateExchangeRatePrevote{
		Hash:        hash.String(),
		Voter:       voter.String(),
		SubmitBlock: submitBlock,
	}
}

// String implement fmt.Stringer.
func (v AggregateExchangeRatePrevote) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

// NewAggregateExchangeRateVote creates an AggregateExchangeRateVote instance.
func NewAggregateExchangeRateVote(exchangeRateTuples ExchangeRateTuples, voter sdk.ValAddress) AggregateExchangeRateVote {
	return AggregateExchangeRateVote{
		ExchangeRateTuples: exchangeRateTuples,
		Voter:              voter.String(),
	}
}

// String implement fmt.Stringer.
func (v AggregateExchangeRateVote) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

// NewExchangeRateTuple creates an ExchangeRateTuple instance.
func NewExchangeRateTuple(denom string, exchangeRate sdk.Dec) ExchangeRateTuple {
	return ExchangeRateTuple{
		denom,
		exchangeRate,
	}
}

// String implement fmt.Stringer.
func (v ExchangeRateTuple) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

// ExchangeRateTuples defines array of ExchangeRateTuple.
type ExchangeRateTuples []ExchangeRateTuple

// String implements fmt.Stringer.
func (tuples ExchangeRateTuples) String() string {
	out, _ := yaml.Marshal(tuples)
	return string(out)
}

// ParseExchangeRateTuples parses ExchangeRateTuple from a string.
func ParseExchangeRateTuples(tuplesStr string) (ExchangeRateTuples, error) {
	tuplesStr = strings.TrimSpace(tuplesStr)
	if len(tuplesStr) == 0 {
		return nil, nil
	}

	tupleStrs := strings.Split(tuplesStr, ",")
	tuples := make(ExchangeRateTuples, len(tupleStrs))
	duplicateCheckMap := make(map[string]bool)
	for i, tupleStr := range tupleStrs {
		splits := strings.Split(strings.TrimSpace(tupleStr), ":")

		if len(splits) != 2 {
			return nil, fmt.Errorf("invalid exchange rate notation: %s", tupleStr)
		}

		denom := strings.TrimSpace(splits[0])
		err := sdk.ValidateDenom(denom)
		if err != nil {
			return nil, err
		}

		amountStr := strings.TrimSpace(splits[1])
		amount, err := sdk.NewDecFromStr(splits[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse decimal amount '%s': %w", amountStr, err)
		}

		tuples[i] = ExchangeRateTuple{
			Denom:        denom,
			ExchangeRate: amount,
		}

		if _, ok := duplicateCheckMap[denom]; ok {
			return nil, fmt.Errorf("duplicated denom %s", denom)
		}

		duplicateCheckMap[denom] = true
	}

	return tuples, nil
}
