# Collatz Conjecture Solver

This project is a Go application that attempts to solve the Collatz Conjecture using two methods: a seed-based approach and a brute force approach.

The calculation utilises `big.Int` and can represent integers as large as can fit in your computer's memory. There is no theoretical limit to the size of a `big.Int` number, the limit is only practical and depends on the amount of memory your computer has.

## Contributing

Please see [`CONTRIBUTING.md`](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Usage

1. Download the [collatz.zip](https://raw.githubusercontent.com/jfallis/collatz/main/collatz.zip) file from the repository.
1. Extract the `collatz.zip` file to get the `collatz` binary.
1. Open a terminal and navigate to the directory containing the `collatz` binary.
1. Run the binary with the desired mode (seed or bruteforce) and the necessary arguments.

### Calculate the Collatz Conjecture

![seed screenshot](/readme/screenshot_seed.png)

```bash
./collatz seed -number=<number>
```
- Replace `<number>` with the number you want to calculate the Collatz Conjecture.

Example:
```bash
./collatz seed -number=27
```

### Bruteforce Calculate the Collatz Conjecture

![bruteforce screenshot](/readme/screenshot_bruteforce.png)


```bash
./collatz bruteforce begin=<start_number> end=<end_number>
```

- Replace `<start_number>` with the number you want to start the sequence from.
- Replace `<end_number>` with the number you want to end the sequence at.

Example:
```bash
./collatz bruteforce begin=0 end=1234567890
```

```bash
./collatz bruteforce -end=1000000 -print-all | sed -E 's/level=INFO msg="([^"]+)"/\1/' | sed -E 's/number: ([0-9]+), steps: ([0-9]+), max: ([0-9]+)/\1,\2,\3/' > example.txt
```

### Help
More options can be found by running the following command:

```bash
./collatz --help
```

### Prerequisites

- Go version 1.23 or higher

## Installation

1. Ensure you have Go version 1.23 or higher installed on your machine.
1. Clone the repository to your local machine.
1. Navigate to the project directory.
1. Run `go mod download` to download the necessary dependencies.

## Building the Application

The application can be run in two modes: seed and bruteforce.

`make build`: Builds the project and outputs the binary in `out/bin/`.

## Using the Makefile

The Makefile includes several commands that help with building, testing, and linting the application.

Builds the project and outputs the binary in `out/bin/`.
```bash
make build
```

Removes build related files.
```bash
make clean
```

Runs the tests.
```bash
make test
```

Runs the tests and exports the coverage.
```bash
make coverage
```

Runs golang linting.
```bash
make lint
```

Shows help information for the Makefile commands.
```bash
make help
```
