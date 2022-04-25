---
order: 5
---

# Tokens

Learn about the different types of native tokens/coins available in Merlion. {synopsis}

## Native Coins

The two native coin types of Merlion are **Mer** and **Lion**:

- **Mer**: Stablecoins that track the price of fiat currencies, and they are named for their fiat counterparts. In the
  early stage of the mainnet launch, it will mainly issue **MerUSD**, or **USM**, which tracks/pegs the price of $USD.
- **Lion**: Native staking coin that partially absorbs the price volatility of Mer. Users stake Lion to validators to
  add blocks of transactions to the blockchain, and earn various fees and rewards. Holders of Lion also can vote on
  proposals and participate in on-chain governance. And Lion is also used for gas consumption for running smart
  contracts on the EVM.

Merlion uses [Atto](https://en.wikipedia.org/wiki/Atto-) Lion or `alion` as the base denomination to maintain parity
with Ethereum.

```
1 lion = 1 * 1e18 alion
```

And the base denomination of Mer is `uusd`.

```
1 MerUSD = 1 USM = 1 * 1e6 uusd
```

## Other Cosmos Coins

Accounts can own Cosmos SDK coins in their balance, which are used for operations with other Cosmos modules and
transactions. Example of these are IBC vouchers.

## ERC-20 Tokens

Any ERC-20 tokens are natively supported by the EVM of Merlion.
