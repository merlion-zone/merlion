---
order: 0 title: "ERC20 Overview"
parent:
title: "erc20"
---

# `erc20`

## Abstract

This document specifies the erc20 module of the Merlion blockchain.

The erc20 module enables Merlion to support an automatic on-chain bidirectional, instant mapping or synchronization of
tokens between the EVM and Cosmos runtimes, specifically the `x/evm` and `x/bank` modules. This allows token holders on
Merlion to instantaneously use their native Cosmos-style `sdk.Coins` (in this document referred to as "Coin(s)") as
ERC-20 tokens (aka "Token(s)"), and vice versa.

Unlike `Evmos`'s erc20 module, Merlion do not use transaction-triggered conversion of coins/tokens. This is because we
believe that users should always feel that various asset operations are performed on one blockchain, rather than two
"logical" chains. When a user purchases a certain ERC-20 token in the DApp based on the EVM smart contracts, he can
immediately see the balance of this token in any Cosmos SDK based wallet. When a user transfers a certain coin from
other Cosmos SDK based blockchains to Merlion using the IBC protocol, he can immediately see the ERC-20 mapping token of
this coin in the corresponding EVM smart contract and perform arbitrary transaction operations.

How the erc20 module implements the mapping of coins/tokens? Since EVM and Cosmos are two runtimes that are not
compatible by default, we need to insert hooks in the right places to handle this mapping in a synchronous way.
Fortunately, the same wallet private key can derive the same address, which has only different representations in the
two runtimes. We have different mapping methods for native `sdk.Coin` and native ERC-20 tokens.

For native `sdk.Coin`, the necessary prerequisite is that the certain `sdk.Coin` must have a registered `DenomMetaData`.
For coins that have not registered `DenomMetaData` in advance, the erc20 module requires that the
missing `DenomMetaData`
must be registered through a governance proposal, otherwise the coin transaction operation using the bank module will
also report an error. After registering `DenomMetaData`, when the `MintCoins` transaction RPC method of the bank module
is called, the erc20 module will automatically create and deploy an ERC-20 smart contract mapped with this `sdk.Coin`,
and synchronize the account balances. And then, other transactions of this native coin through the bank module are
always synchronized to the ERC-20 contract call conducted on the EVM state machine. Surely the EVM-based ERC-20
transactions through JSON-RPC API are also synchronized to the native operations in the bank module. It is not only
applied to native staking/gov coins, but to all IBC vouchers.

For native ERC-20 token, there are no governance precondition or other preprocessing that is needed to map the token
to `sdk.Coin`. Actually, when user sends transactions through EVM JSON-RPC, there isn't any operation conducted in the
bank module. But when user sends bank transactions through Cosmos gRPC API, the operation is also conducted on the EVM
state machine. However, for any query gRPC call, the bank module always proxy the query to the EVM state store. This
approach maximizes savings in gas costs and actual runtime overhead.

With the `x/erc20` users on Merlion can

- use existing native Cosmos assets (like $OSMO or $ATOM) on EVM-based chains, e.g., for Trading IBC tokens on DeFi
  protocols, buying NFT, etc.
- transfer existing tokens on Ethereum and other EVM-based chains to Merlion to take advantage of application-specific
  chains in the Cosmos ecosystem.
- build new applications that are based on ERC-20 smart contracts and have access to the Cosmos ecosystem.