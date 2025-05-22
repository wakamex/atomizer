# Atomizer Project

This repository contains the Atomizer project and its related examples.
It utilizes the `ryskV12-cli` as a submodule located in the `sdk` directory.

## Progress

sdk actions:
- [x] approve - leave in the cli
- [x] balances - debugging
- [x] connect - working w/ default channel id
- [ ] positions
- [x] quote - maker_quote_response.go
- [ ] transfer
combo:
- [x] maker_quote_response.go (connect and quote)


## Prerequisites

- Git
- Go (version as specified in `examples/go.mod` and `sdk/go.mod`)

## Setup

1. Clone the repository:
   ```bash
   git clone <repository_url> atomizer
   cd atomizer
   ```

2. Initialize and update the submodule:
   ```bash
   git submodule update --init --recursive
   ```

## Maker Quote Responder Application

The primary application in this repository is the `maker_quote_responder`, located in the `cmd/maker_quote_responder/` directory. This application connects to the Rysk Finance API, listens for RFQs, and responds with quotes.

For detailed instructions on how to build, configure, and run the `maker_quote_responder`, please refer to its dedicated README file:

[**Maker Quote Responder Instructions (`cmd/maker_quote_responder/README.md`)**](./cmd/maker_quote_responder/README.md)

This guide covers:
- Prerequisites
- Building the executable
- Setting up the `.env` file
- Running the application using the `run.sh` script

## Submodule

The `sdk` directory is a submodule pointing to the [ryskV12-cli](https://github.com/wakamex/ryskV12-cli) repository.

To update the submodule to the latest commit from its remote:
```bash
cd sdk
git pull
cd ..
git add sdk
git commit -m "Update sdk submodule"
```
