# OpenX

[![Build Status](https://travis-ci.com/YaleOpenLab/openx.svg?branch=master)](https://travis-ci.com/YaleOpenLab/openx)
![Github Build Status](https://github.com/YaleOpenLab/openx/workflows/.github/workflows/go.yml/badge.svg)
[![Codecov](https://codecov.io/gh/YaleOpenLab/openx/branch/master/graph/badge.svg)](https://codecov.io/gh/YaleOpenLab/openx)
[![Go Report Card](https://goreportcard.com/badge/github.com/YaleOpenLab/openx)](https://goreportcard.com/report/github.com/YaleOpenLab/openx)

Openx is a project hosted at the MIT Media Lab Digital Currency Initiative and the Yale Open Innonvation Lab. The openx model seeks to implement the paradigm of investing in and developing projects without hassles, and enabling smart ownership with the help of the blockchain. Openx can be thought of as a platform of platforms and houses multiple platforms within it.

The goal of openx is to have a common interface between all parties that relate to a project: Investors, Recipients, Developers, an dOriginators. A running pilot of openx is at [opensolar](https://github.com/YaleOpenLab/opensolar)

## Related Repositories

[Openx](https://github.com/YaleOpenLab/openx)  
[Opensolar](https://github.com/YaleOpenLab/opensolar)  
[Opensolar Frontend](https://github.com/YaleOpenLab/openx-frontend)  
[API Docs](https://github.com/YaleOpenLab/openx-apidocs)  
[Wiki](https://github.com/YaleOpenLab/openxdocs)  
[Create-openx-app](https://github.com/YaleOpenLab/create-openx-app)  
[Openx-CLI](https://github.com/Varunram/openx-cli)

## Related Websites

[Demo](www.openx.solar)  
[Openx API](api.openx.solar)  
[Opensolar API](api2.openx.solar)  
[API docs](apidocs.openx.solar)  
[Wiki](api.openx.solar)  
[Builds](builds.openx.solar)  
[MQTT broker](mqtt.openx.solar)  
[Pilot Dashboard](dashboard.openx.solar)

## Getting Started

### Download

Openx builds are available at [builds.openx.solar](builds.openx.solar)

Docker image available at [Docker Hub](https://hub.docker.com/repository/docker/varunramg/openx)

### Installing from PPA

Openx is available on PPA [here](https://launchpad.net/~varunram/+archive/ubuntu/openx). Warning: The PPA might not be up to date.

## Building from source

Requirements:

1. Go 1.11 and above
2. Standard build tools depending on architecture (build-essential for linux, brew and xcode dev tools for mac)

IMPORTANT: Inline with the Golang dev team, we will not be supporting versions of go that are more than two releases old. If you have a version of Go that is older than 1.11, please upgrade to the latest version of Go before continuing.

```
go get -v github.com/YaleOpenLab/openx
cd $GOPATH/src/github.com/YaleOpenLab/openx/
go mod download
go mod verify
go build
```

Make sure you have the necessary permissions to write to `$HOME`. Start openx using `./openx`

## Installing IPFS

ipfs is used by some parts of the program to store data that needs to be publicly verified. Download a release from https://github.com/ipfs/go-ipfs/releases and run install.sh. If you face conflicts between multiple ipfs versions, you might need to run [fs-repo-migrations](https://github.com/ipfs/fs-repo-migrations/blob/master/run.md) to migrate to the newer version. If you don't have anything worth storing, you can delete the ipfs home directory and run `ipfs init` again (this will delete the data stored in ipfs prior to deletion)

You need to keep your peer key (`ipfs.key` usually) in a safe place for future reference. Start ipfs using `ipfs daemon` and you can test it by creating a file `test.txt` and run `ipfs add test.txt`. The resultant hash can be decrypted using `curl "http://127.0.0.1:8080/ipfs/hash"` where `127:0.0.1:8080` is the endpoint of the ipfs server. If you need more help on setting up ipfs, you can also refer to [this helpful tutorial](https://michalzalecki.com/set-up-ipfs-node-on-the-server/).

## Running tests

Running tests is simple with `go test` but tests have flags since some require running daemons in the background (eg. `ipfs`). There are two kinds of flags - `travis` and `all`.

If you need coverage stats as well, you need to install `cover` as well. `go get golang.org/x/tools/cmd/cover` if you don't have the package. Running `go test --tags="travis" -coverprofile=test.txt ./...` should run all the tests and provide coverage data on each specific package.

## Contributing

Please feel free to open Pull Requests and Issues with your changes and suggestions. Before working on a major feature, please describe the same in an issue so everyone can understand what you're building.

## Security

<img src="security/discloseio.png" width="50">  

For security related issues, *DO NOT OPEN A GITHUB ISSUE!*. Please disclose the information responsibly by sending a (preferably PGP Encrypted) email to the lead developer `Varunram Ganesh` (`contact@varunram.com`). [His PGP Key](https://pgp.mit.edu/pks/lookup?op=vindex&fingerprint=on&search=0x708C606504A49970) fingerprint is `C98F 0014 9A99 36E4 E56D  2471 708C 6065 04A4 9970`

In addition to this, openx is fully fully compliant with the [disclose.io](https://disclose.io) core terms followed by [bugcrowd](https://www.bugcrowd.com/resource/what-is-responsible-disclosure/). For more info, please checkout [SECURITY.md](SECURITY.md)

## License
[GPL3](https://github.com/YaleOpenLab/openx/blob/master/LICENSE)
