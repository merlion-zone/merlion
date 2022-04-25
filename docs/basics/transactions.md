---
order: 3
---

# Transactions

End-users create `transactions` to trigger state changes in Merlion. {synopsis}

## Pre-requisite Readings

- [SDK Transactions](https://docs.cosmos.network/master/core/transactions.html)
- [SDK Transaction Lifecycle](https://docs.cosmos.network/master/basics/tx-lifecycle.html) {prereq}
- [Ethereum Transactions](https://ethereum.org/en/developers/docs/transactions/)

Transactions are cryptographically signed instructions from accounts. An account will initiate a transaction to update
the state of the Merlion network. For example, if Bob sends Alice 1 LION, Bob's account must be debited and Alice's must
be credited. This state-changing action takes place within a transaction.

Transactions, which change the state of the blockchain, need to be broadcast to the whole network. Any node can
broadcast a request for a transaction to be executed; after this happens, a validator will execute the transaction and
propagate the resulting state change to the rest of the network.

Merlion must handle two transaction types:

- Cosmos SDK Transaction
- EVM Transaction

## Cosmos SDK Transaction

Transactions are comprised of metadata held in `contexts`
and [`sdk.Msg`s](https://docs.cosmos.network/master/building-modules/messages-and-queries.html)
that trigger state changes within a module through the module's
gRPC [`Msg` service](https://docs.cosmos.network/master/building-modules/msg-services.html).

When users want to interact with Merlion and make state changes (e.g., sending coins), they create transactions.
Each `sdk.Msg` of a transaction must be signed using the private key associated with the appropriate account(s), before
the transaction is broadcasted to the network. A transaction must then be included in a block, validated, and approved
by the network through the consensus process.

## EVM Transaction

Due to relying on the underlying Tendermint consensus process, the Merlion network can only accept the Cosmos SDK
transactions. For receiving and handling EVM transactions, Merlion attempts to
mimic [geth](https://github.com/ethereum/go-ethereum) `Transaction` structure and treat it as unique `sdk.Msg` type. It
is implemented by the `evm` module, and the unique `sdk.Msg` type is `MsgEthereumTx`. All relevant Ethereum transaction
information is contained in this message. This includes the signature, gas, payload, amount, etc.

Being that Merlion implements the Tendermint ABCI application interface, as transactions are consumed, they are passed
through a series of handlers. Once such handler, the `AnteHandler`, is responsible for performing preliminary message
execution business logic such as fee payment, signature verification, etc. This is particular to Cosmos SDK
transactions. EVM transactions will bypass this as the EVM handles the same business logic.

### Broadcast EVM Transaction through JSON-RPC

Merlion exposes the standard Web3 [JSON-RPC APIs](../api/json-rpc/server.md) for connection by existing Web3 tooling.
How the full node accepts and processes an EVM transaction through the JSON-RPC endpoint? Well, when a Merlion node
receives such a transaction, it will wrap it into a unique `MsgEthereumTx` message with no gas fees and no signatures,
since the original EVM transaction already has its own gas setting and signature. And then, the node broadcast the
Cosmos SDK transaction, which containing the message, to the network. As mentioned above, the custom `AnteHandler` will
handle this kind of transaction specially, that is, it will be handed over to the EVM for processing.
