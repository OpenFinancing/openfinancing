# Stellar PoC

[![Build Status](https://travis-ci.org/OpenFinancing/openfinancing.svg?branch=master)](https://travis-ci.org/OpenFinancing/openfinancing)
[![Codecov](https://codecov.io/gh/OpenFinancing/openfinancing/branch/master/graph/badge.svg)](https://codecov.io/gh/OpenFinancing/openfinancing)

This sub-repo contains the WIP Stellar Proof of Concept. This has some design
tradeoffs compared to the Ethereum implementation, especially in the concept of
"state" in Ethereum.

There are various stages to a given solar contract, this can be expanded upon to furhter sub stages, but right now are defined to be  7 for convenience

 - Stage 0: A pre origin contract, which is given by an originator to a recipient
 - Stage 1: An originated contract, which is the evolution of a pre origin contract at the behest of the recipient, who may require a legal contract to be signed by the originator for doing so
 - Stage 2: A proposed contract, which is a contract given by a contractor who is willing to lend his services for the particular project at a fee that will be paid by the recipient and in turn, the investor as well
 - Stage 3: A finalized contract, which is open for investor funds. A contract can be promtoed from stage 2 to 3 only by a recipient, who chooses the best contract that fits in with his needs and requirements
 - Stage 4: A contract that has raised the money required for investment from investors and is ready to be installed. The control of the contract now passes over from the recipient to the developer, who is responsible to take this project to stage 5.
 - Stage 5: the stage where the developer should install all the required devices and the contractor must do his job of supervising.
 - Stage 6: Power Generation s tarts and the recipient starts to payback the investor in monthly instalments similar to an electricity bill
 - Stage 7: The recipient has finished paying back the investors and can now own the solar panels installed in his space.

Roughly, in the financing of a solar contract ,there are many entities involved who perform various roles in the entire system.
 - ISSUER - the issuer is the server on which the orders are advertised on
 - RECIPIENT - the recipient is the school which receives funding from investors through an order
 - ORIGINATOR - the person who proposes a pre-origin contract to a recipient and requires approval from the recipient to originate the contract
 - CONTRACTOR - takes in a originated order and proposes a contract which fits in with the requirements described in the originated contract
 - INVESTOR - the investor is a person / group of persons who invest in a particular finalized order
 - DEVELOPER - the person who is responsible for installing the hardware and making sure the solar panels work and produce electricity
 - GUARANTOR - the person who is liable to make the recipient pay or pays itself in the case of a delayed payment from the recipient's side

 The rough workflow for promotion of a contract is:
 Originator proposes contract -> recipient reviews all proposed contracts -> recipient chooses a particular contract to originate for the project -> contractors see all originated contracts -> contractors propose contracts to install the systems -> recipient reviews all proposed contracts -> recipient chooses a particular proposed contract to be final -> finalize contract -> investors see all finalised contracts -> choose one -> contract ready to be invested in -> solar panels installed -> power generation -> full ownership

These three entities interact use the three tokens detailed below for proofs of investment, proof of debt and proof of payback respectively:

1. INVTokens - INVTokens are sent from the ISSUER to the INVESTOR one time, as proof of investment
2. DEBTokens - DEBTokens are sent from the ISSUER to the RECIPIENT the time the order is created and then each month, the RECIPIENT pays back the ISSUER each month with a number of DEBTokens like an electricity bill. (This has to be capped at a minimum of the electricity bill in the future, along with an oracle that can attest to the price of electricity at the place)
3. PBTokens (optional)- PBTokens are sent from the ISSUER to the RECIPIENT at multiple intervals:
  - At the time of order confirmation, 1 each for each month the school opts in to pay the amount back. eg. if the school opts in for a period of 5 years, the ISSUER issues 60 PBTokens. We also calculate expected_paid_amount as `amount_invested /  payback_amount`
  - Each time the ISSUER receives DEBTokens, the ISSUER confirms the transaction and sends back PBTokens as `paid_amount / expected_paid_amount` eg you pay 210 towards an expected amount of 120. You get paid back 1.75 PBTokens.

Any disparity / failure on the ISSUER's part can be argued with, since the transactions are on chain. The PBToken simplifies this, since one doesn't need to go back in history and calculate how much a given school has paid a person. Percentage paid is simply `PBToken_balance / PBToken_total`.

# Installing IPFS

ipfs is used by some parts of the program to store legal contracts, files that the user might want to store permanently. Download a release from https://github.com/ipfs/go-ipfs/releases and run install.sh. In case you face an issue with migration between various ipfs versions, you might need to run [fs-repo-migrations](https://github.com/ipfs/fs-repo-migrations/blob/master/run.md) to migrate to a newer version. If you don't have anything valuable, you can delete the directory and run `ipfs init` again.

You need to keep your peer key (`ipfs.key` usually) in a safe place for future reference. Start ipfs using `ipfs daemon` and you can test if it worked by creating a test file `test.txt` and run `ipfs add test.txt` to see if it succeeds. The resultant hash can be decrypted using `curl "http://127.0.0.1:8080/ipfs/hash"` where 8080 is the endpoint of the ipfs server or by doing `cat /ipfs/hash` directly. You can also refer to [this helpful tutorial](https://michalzalecki.com/set-up-ipfs-node-on-the-server/) in order to get easily started with ipfs.

# Running the code in this repo

In order to be able to run this, you need to have the latest version of go installed. [Here](https://medium.com/@patdhlk/how-to-install-go-1-9-1-on-ubuntu-16-04-ee64c073cd79) is a quick tutorial on how to get go installed on a Linux / macOS machine.

Once you have go installed, you need to get the packages in this repo. This can be done using `go get -v ./...`. The `stellar/go` package might print out errors due to its problems with `go get`, in which case you need to get the package separately and then run `dep ensure -v` on the project, proceeded by `go get`ing the other packages as normal.

Then you need to build the stellar package `go build` and then run the executable like `./stellar`.

You will be faced with a CLI interface with which you can interact with

# Running tests

Running tests is mostly simple with `go test` but the tests have flags since some require running other software in the background (`ipfs`). There are two kinds of flags right now - `travis` and `all`. If you need the coverage stats as well, you need to
```
go get golang.org/x/tools/cmd/cover
```
if you already don't have the package. Then running `go test --tags="all" -coverprofile=test.txt ./...` should run all the tests and provide coverage data. Running with the tag `travis` will omit the tests in `ipfs/` which requires [a local `go-ipfs` node running](https://michalzalecki.com/set-up-ipfs-node-on-the-server/) as described above.
