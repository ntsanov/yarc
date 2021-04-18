# yarc (Yet Another Rosetta CLI)

## General

This cli was written for Harmony One implementation of rosetta-sdk. Although most of the code could be reused at this moment is only tested/working with Harmony

Reference to the API could be found [here](https://api.hmny.io/#9a6b2616-11bb-4f28-a851-ac554456c571)

## Build

There are a few dependencies inherited from harmony libraries and the easiest way to buld it is to use harmony/main docker image. To do this run the provided build script
```
    bash build_static.sh
```
The build binary is located in cli/build and is statically build

For development/debugging you would need to build and install (make & make install) [mcl](https://github.com/harmony-one/mcl)
 and [bls](https://github.com/harmony-one/bls)
    
## Signing
For signing transaction sender's wallet must locally exist

[How to create wallet](https://docs.harmony.one/home/network/wallets/harmony-cli/create-import-wallet)

## Usage

All API endpoints are mapped to a subcommand

### Subcommand list

|endpoint|subcommand |
|---|---|
|/network/list| data network list|
|/network/options| data network options|
|/network/status| data network status|
|/account/balance | data balance |
|/block | data block |
|/block/transaction | data block |
|/mempool | data mempool |
|/mempool/transaction| data mempool |
|/construction/derive| con derive |
|/construction/preprocess| con preprocess |
|/construction/metadata| con metadata |
|/construction/payloads| con payloads |
|/construction/parse | con parse |
|/construction/combine| con combine |
|/construction/hash| con hash |
|/construction/submit| con submit |

### Response

All responses are json encoded except command parameters validation. If error occurs return code is non 0

### Examples 

Preprocess transaction
```
yarc con preprocess \
    --node=https://rosetta.s0.b.hmny.io \
    --from-shard 0 \
    --to-shard 0 \
    --from one1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
    --to one1yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy \
    --amount 1
```

In directory scripts is a bash script that does all the steps to create and submit the result. Before running .env-example should be renamed to .env and filled
