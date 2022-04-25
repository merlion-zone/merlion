---
order: 2
---

# Accounts

This document describes the in-built account and public key system of Merlion. {synopsis}

## Account Definition

In Merlion, an _account_ designates a pair of _public key_ `PubKey` and _private key_ `PrivKey`. The `PubKey` can be
derived to generate various `Addresses`, which are used to identify users (among other parties) in the
application. `Addresses` are also associated with `message`s to identify the sender of the `message`. The `PrivKey` is
used to generate [digital signatures](#signatures) to prove that an `Address` associated with the `PrivKey` approved of
a given `message`.

## <a name="signatures"></a>Keys, accounts, addresses, and signatures

The principal way of authenticating a user is done
using [digital signatures](https://en.wikipedia.org/wiki/Digital_signature). Users sign transactions using their own
private key. Signature verification is done with the associated public key. For on-chain signature verification
purposes, we store the public key in an `Account` object (alongside other data required for a proper transaction
validation).

Merlion supports the following digital key schemes for creating digital signatures:

* `eth_secp256k1`, deriving keys and addresses which identify **accounts** and **validator operators**.
* `ed25519`, deriving keys and addresses which identify **consensus nodes** only for the consensus validation.

> Not like in the Cosmos SDK, Merlion do not support the original digital key schemes: `secp256k1` and `secp256r1`. This is because of compatibility with EVM.

## Addresses

`Addresses` and `PubKey`s are both public information that identifies actors in the application. `Account` is used to
store authentication information.

Each account is identified using `Address` which is a sequence of bytes derived from a public key. Like in the Cosmos
SDK, we define 3 types of addresses that specify a context where an account is used:

* `AccAddress` identifies users (the sender of a `message`).
* `ValAddress` identifies validator operators.
* `ConsAddress` identifies validator nodes that are participating in consensus. Validator nodes are derived using
  the **`ed25519`** curve.

For user interaction, addresses are formatted using [Bech32](https://en.bitcoin.it/wiki/Bech32) and implemented by
the `String` method. The Bech32 method is the only supported format to use when interacting with Merlion. The Bech32
human-readable part (Bech32 prefix) is used to denote an address type. Example:

|                    | Address Bech32 Prefix |
| ------------------ |-----------------------|
| Accounts           | mer                   |
| Validator Operator | mervaloper            |
| Consensus Nodes    | mervalcons            |

## Keyring

A `Keyring` is an object that stores and manages accounts. Specifically, Merlion has a default implementation
of `Keyring`, and users can manage their accounts/keys through the CLI `merliond keys`.

> Broadly speaking, the concept of `Keyring` is just an abstract of keys management. For DApp users, it may be presented as a wallet or a keystore, etc.

## Hierarchical Deterministic Wallets

The `Keyring` API in Merlion supports the HD key derivation, following the
standard [BIP32](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
and [BIP44](https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki). BIP32 describes a set of accounts derived
from an initial secret seed. A seed is usually created from a 12- or 24-word mnemonic. A single seed can derive any
number of `PrivKey`s using a one-way cryptographic function. Then, a `PubKey` can be derived from the `PrivKey`.
Naturally, the mnemonic is the most sensitive information, as private keys can always be re-generated if the mnemonic is
preserved. BIP44 defines a logical hierarchy/levels in BIP32 path:

```
m / purpose' / coin_type' / account' / change / address_index
```

Although HD key derivation is a general solution for controlling multiple accounts with a single key,
Merlion's `Keyring`
API primarily uses it to derive just one account, with a specific HD path, e.g., `m / 44' / 60' / 0' / 0 / 0`.
