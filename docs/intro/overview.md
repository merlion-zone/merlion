---
order: 1
---

# High-level Overview

## What is Merlion

The Merlion project is the most decentralized, scalable, and high-throughput blockchain for fractional-algorithmic
stablecoin and various vanward DeFi-specific innovations. It is built
using [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and [Tendermint Core](https://github.com/tendermint/tendermint)
. It also enables fully compatibility and interchangeability with Ethereum/EVM based DApps, by embedding the evm module
from [Ethermint](https://github.com/tharsis/ethermint).

Absorbed the creations and lessons of some pioneering protocols/projects, the architecture and mechanism of Merlion are
designed to meet the demands of an increasingly diversified multi/cross-chain world. Merlion's tokenomics is elaborately
formulated to not only incentivize early adopters, but also encourage long-term cooperation to grow bigger and stronger.

### Features

Here’s a glance at some key features of Merlion:

- System native stablecoin **MerUSD**, or **USM**
- Pegging through partial backing/collateral and partial algorithmic
- ve(3,3) staking of native **Lion** and proportional incentive
- Web3 and EVM compatibility
- High throughput and fast transaction finality via [Tendermint Core](https://github.com/tendermint/tendermint)
- Horizontal scalability via [IBC](https://cosmos.network/ibc) and [Gravity Bridge](https://www.gravitybridge.net)
- Full decentralized on-chain governance and flexible dynamical policy regulation

## Mer and Lion

The merlion system consists of two main native coins, Mer and Lion.

- **Mer**: Stablecoins that track the price of fiat currencies, and they are named for their fiat counterparts. In the
  early stage of the mainnet launch, it will mainly issue **MerUSD**, or **USM**, which tracks/pegs the price of $USD.
- **Lion**: Merlion blockchain's native staking coin that partially absorbs the price volatility of Mer. Users stake
  Lion to validators to add blocks of transactions to the blockchain, and earn various fees and rewards. Holders of Lion
  also can vote on proposals and participate in on-chain governance.

## The stablecoin protocol

### Stablecoin

Like [Terra](https://terra.money), Merlion is also a DeFi-specific blockchain with built-in algorithmic stablecoin
protocol. But the difference is that Merlion does not rely on any purely algorithmic design which is difficult to grow
and exhibits extreme periods of volatility. Merlion stablecoin will be minted in two ways:
fractional-backing-algorithmic and over-collateralized-catalytic.

The **fractional-backing-algorithmic** (inspired by [Frax](https://frax.finance)), or **FBA**, is with parts of its
backing assets and parts of the algorithmic supply. The ratio of backing and algorithmic depends on the market's pricing
of the Mer stablecoin. We named the ratio **CR (Collateral Ratio)**. If MerUSD is trading at above $1, the system
decreases the ratio. If Mer is trading at under $1, the system increases the ratio. At any point in time, CR is
determined. If users want to mint Mer, they must **spend** a certain amount of backing assets and a certain amount of
Lion coins, which will enter the unified swap pool. Conversely, when a user wants to acquire backing assets and Lion
coins in the swap pool, he must spend a certain amount of Mer coins in exchange, and can only get the proportional coins
that follow the CR ratio at the market price.

The **over-collateralized-catalytic** (partially inspired by [MIM](https://abracadabra.money)), or **OCC**, is over
collateralized for interest-bearing lending, and loan-to-value maximized by catalytic Lion. Each kind of supported
over-collateralized asset forms a separate pool. Users must pre-deposit over-collateralized assets and then lend Mer (
actually minted directly by the system) when needed. The maximum ratio that can be lent (called loan-to-value, or LTV)
depends on the parameters set by the system for this collateral pool, and the additional Lion-based assisted minting (
called catalytic) added by the user. If users want to redeem their collateral, they are obliged to repay the principal
and interest of the lent assets, and they must always pay attention to the price fluctuation of the collateral assets,
to avoid triggering the liquidation mechanism.

Why two mintage ways? Well, some users have a lot of real demand for stablecoins and don't want to hold volatile
assets (like Lion, ETH, etc.). Other users want to hold volatile assets for a long time in order to obtain their
appreciation, but have short-term liquidity needs, like to borrow stablecoins, and are willing to pay interest and bear
possible liquidation risks. Exactly, FBA is suitable for the former, and OCC is suitable for the latter.

### Oracle, maker and arbitrage

The built-in oracle module is responsible for providing fair and true on-chain prices for specified assets. It accepts
near real-time quotes from active validators, removes deviations, and takes the median as the final standard on-chain
price. Since the price is the most important indicator for the whole system to reflect the real market situation, and it
is also an important reliance on whether the stablecoin can remain pegged, the oracle module will periodically reward
validators who consistently quote correctly.

The most important Mer minting/burning activities are handled by the maker module. Whether it is FBA-style spending
mintage or OCC-style collateralized mintage, the maker module will accept the assets deposited by users and securely
host them. The maker module defines a series of transaction types to facilitate the user to mint any desired amount of
Mer coins according to the deterministic parameters of the current system and the token price data from the oracle
module. The maker module will charge a certain seigniorage-like fee for the oracle module to incentivize validators who
provide quotes. For OCC-style collateralized mintage, the maker module is also responsible for settling interest and
providing possible liquidation channels.

When possible market fluctuations or black swan events cause Mer prices to decouple, the arbitrage behavior of
arbitrageurs helps Mer return to their target prices. For FBA-style mintage, when MerUSD is lower than $1, the system
target CR value automatically increases accordingly. At this time, the system shows that there is a lack of backing
assets, and a surplus of Lion coins. Arbitrageurs can use a certain amount of backing assets to buy Lion coins from the
system at a discounted price, so that the actual CR value tends to the target CR value. An increase in the reserve of
backing assets will increase people's confidence and bring MerUSD back to the price of $1. For OCC-style mintage, when
MerUSD is below $1, users tend to buy low-priced MerUSD from the market and repay the system with MerUSD's nominal value
of $1 to unlock the collateralized assets. This will reduce the circulation of MerUSD in the market, and will also bring
MerUSD back to the price of $1. And for MerUSD above $1, it is obvious and easy to get it back to $1, so we won't go
into details here.

## The ve(3,3) mechanism and staking

## Smart contracting and virtual machine

## Inter-blockchain communication

## Proof-of-Stake and validators

## Governance

## Tokenomics
