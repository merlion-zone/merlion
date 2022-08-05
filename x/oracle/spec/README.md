<!--
order: 0
title: "Oracle Overview"
parent:
  title: "oracle"
-->

# `oracle`

## Abstract

This document specifies the oracle module of the Merlion blockchain.

The oracle module provides the Merlion blockchain with an up-to-date and accurate price feed of exchange rates of various coins against USD so that 

As price information is extrinsic to the blockchain, the Merlion network relies on validators to periodically vote on the current coins exchange rate, with the protocol tallying up the results once per `VotePeriod` and updating the on-chain exchange rate as the weighted median of the ballot.

> Since the Oracle service is powered by validators, you may find it interesting to look at the [Staking](https://github.com/cosmos/cosmos-sdk/tree/master/x/staking/spec/README.md) module, which covers the logic for staking and validators.
