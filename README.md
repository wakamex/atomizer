# Atomizer Project

This repository contains the Atomizer project and its related examples.
It utilizes the `ryskV12-cli` as a submodule located in the `sdk` directory.

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

## Running the Go Example

To run the example program located in the `examples` directory:

1. Navigate to the examples directory:
   ```bash
   cd examples
   ```

2. Run the main Go program:
   ```bash
   go run main.go
   ```
   This will compile and execute the example. Dependencies will be automatically downloaded if needed based on `examples/go.mod`.

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
