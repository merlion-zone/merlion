<!--
order: 0
title: "Oracle Overview"
parent:
  title: "oracle"
-->

# `oracle`

## Abstract

This document specifies the Oracle module of the Merlion blockchain.

The Oracle module provides the Merlion blockchain with an up-to-date and accurate price feed of exchange rates of various coins against USD so that the maker may provide fair exchanges between different currency pairs.

As price information is extrinsic to the blockchain, the Merlion network relies on validators to periodically vote on the current coins' exchange rate, with the protocol tallying up the results once per `VotePeriod` and updating the on-chain exchange rate as the weighted median of the ballot.

> Since the Oracle service is powered by validators, you may find it interesting to look at the [Staking](https://github.com/cosmos/cosmos-sdk/tree/master/x/staking/spec/README.md) module, which covers the logic for staking and validators.

## Contents

1. **[Concepts](01_concepts.md)**
    - [Voting Procedure](01_concepts.md#voting-procedure)
    - [Reward Band](01_concepts.md#reward-band)
    - [Slashing](01_concepts.md#slashing)
    - [Abstaining from Voting](01_concepts.md#abstaining-from-voting)
    - [Transitions](01_concepts.md#transitions)
2. **[State](02_state.md)**
   - [AggregateExchangeRatePrevote](02_state.md#aggregateexchangerateprevote)
   - [AggregateExchangeRateVote](02_state.md#aggregateexchangeratevote)
   - [ExchangeRate](02_state.md#exchangerate)
   - [FeederDelegation](02_state.md#feederdelegation)
   - [MissCounter](02_state.md#misscounter)
3. **[EndBlock](03_end_block.md)**
   - [Tally Exchange Rate Votes](03_end_block.md#tally-exchange-rate-votes)