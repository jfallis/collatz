# Collatz Conjecture Solver

This project is a Go application that attempts to solve the Collatz Conjecture using two methods: a seed-based approach and a brute force approach.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go version 1.22 or higher

### Installing

1. Clone the repository to your local machine.
1. Navigate to the project directory.
1. Run `go mod download` to download the necessary dependencies.

## Running the Application

The application can be run in two modes: seed and bruteforce.

1. `make build`: Builds the project and outputs the binary in `out/bin/`.
1. Run the binary with the following commands:
   - Seed mode: `./collatz seed <number>`
   - Bruteforce mode: `./collatz bruteforce <number>`
   
   Replace `<number>` with the number you want to start the sequence from.

![sample screenshot](/readme/screenshot.png)

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
