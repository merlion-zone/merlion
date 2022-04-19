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

Merlion's vision is to become the most decentralized, permissionless, high-throughput, low-cost, easy-to-use cross-chain
assets settlement center and experimental DeFi innovation zone.

### Features

Here’s a glance at some key features of Merlion:

- System native stablecoin **MerUSD**, or **USM**
- Pegging through partial backing/collateral and partial algorithmic
- ve(3,3) time-locked of native **Lion** and proportional incentive
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
depends on the parameters set by the system for this collateral pool, and the additional Lion-boosting minting (
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

## Proof-of-Stake and validators

Inherited from [Tendermint Core](https://github.com/tendermint/tendermint)
and [Cosmos SDK](https://github.com/cosmos/cosmos-sdk), Merlion uses BFT (Byzantine Fault Tolerance)
consensus protocol to securely and consistently replicate states/blocks/transactions on many machines (or validators)
all over the world. Validators run Merlion programs called full nodes, and take turns proposing blocks of transactions
and voting on them. Blocks are committed in a chain, with one block at each height. A block may fail to be committed, in
which case the protocol moves to the next round, and a new validator gets to propose a block for that height. Two stages
of voting are required to successfully commit a block; Tendermint call them **pre-vote** and **pre-commit**. A block is
committed when more than 2/3 of validators pre-commit for the same block in the same round. Merlion can tolerate up to a
1/3 of validators failures, and those failures can include arbitrary behaviour, e.g., hacking and malicious attacks.

Not all validators will have the same "weight" in the consensus protocol. Every validator will have some voting power,
which may not be uniformly distributed across individual validators. Merlion denominates the weight or stake in our
native Lion coin, and hence the system is often referred to as **Proof-of-Stake**. Validators will be forced to "bond"
or stake their Lion holdings in the security deposit that can be slashed if they're found to misbehave in the consensus
protocol. This adds an economic element to the security of the protocol, allowing one to quantify the cost of violating
the assumption that less than one-third of voting power is Byzantine.

Merlion allows up to the top 120 validators to participate in consensus. A validator’s rank is determined by their stake
or the total amount of Lion bonded to them. Although validators can bond Lion to themselves, they mainly amass larger
stakes from delegators. Validators with larger stakes get chosen more often to propose new blocks and earn
proportionally more rewards.

Delegators are users who want to receive rewards from consensus without running a full node. Any user that stakes Lion
is a delegator. Delegators stake their Lion to a validator, adding to a validator's weight, or total stake. In return,
delegators receive a portion of system fees as staking rewards.

## The ve(3,3) mechanism and time-locked voting escrow

In addition to the normal Lion staking through the Proof-of-Stake consensus protocol, Merlion brings in another enhanced
time-locked voting escrow mechanism, called **ve(3,3)**. We use ve(3,3) to incentivize various innovative DeFi DApps in
the Merlion system as well as getting as many users/investors involved as possible in the governance of the network. We
define the NFT token **veLion** as the vote-escrowed Lion, which is simply Lion locked for a period of time, from 1 week
to 4 years. veLion token holders will receive a multiplied amount of voting power, compared to normal PoS staking. Along
with that, they will gain more staking rewards and voting power on governance proposals. In essence, what is more
important and further is that they will have certain voting rights to incentivize various innovative DeFi DApps in the
system. Of course, they will also be rewarded with extra Lion coins, which come from the reserve in the treasury
according to Merlion's tokenomics.

The rewarded Lion coins will be distributed weekly. The emission amount are adjusted as a percentage of circulating
supply (`circulating_supply = total_supply - (ve_locked + normal_staked)`). Meaning, assuming the maximum weekly
emission is 500,000, if 0% of the coin is staked or ve-locked, the weekly emission would be 500,000. If 50% of the coin
is staked or ve-locked, the weekly emission would be 250,000. If 100% of the coin is staked or ve-locked, the weekly
emission would be 0.

To ensure that veLion holders are never diluted, their holdings will be increased proportional to the weekly emission.
And since the lock position veLion is tokenized as NFT, it allows veLion to be traded on future secondary markets, as
well as to allow participants to borrow against their veLion in future lending marketplaces. This addresses the capital
inefficiency problem of ve assets, as well as addresses concerns over future liquidity (should it ever be required).

## Smart contracting and virtual machine

As a general-purpose DeFi-specific blockchain platform, Merlion must integrate certain smart contract virtual machines
to facilitate the deployment of innovative DeFi protocols or DApps by various third-party developers/teams. Fortunately,
nowadays we have many options for mature virtual machine protocols/modules. First and foremost, Merlion will integrate
the **evm** module from [Ethermint](https://github.com/tharsis/ethermint), to provide native Web3/EVM capabilities and
be compatible with the huge Ethereum ecosystem.

The growth of EVM-based chains (e.g., Ethereum), however, has uncovered several scalability challenges that are often
referred to as
the [Trilemma of decentralization, security, and scalability](https://vitalik.ca/general/2021/04/07/sharding.html).
Developers are frustrated by high gas fees, slow transaction speed & throughput, and chain-specific governance that can
only undergo slow change because of its wide range of deployed applications. A solution is required that eliminates
these concerns for developers, who build applications within a familiar EVM environment.

The evm module provides the EVM familiarity on a scalable, high-throughput Proof-of-Stake blockchain. It allows for the
deployment of smart contracts, interaction with the EVM state machine (state transitions), and the use of EVM tooling.
It alleviates the aforementioned concerns through high transaction throughput
via [Tendermint Core](https://github.com/tendermint/tendermint), fast transaction finality, and horizontal scalability.

## Inter-blockchain communication

## Governance

## Tokenomics
