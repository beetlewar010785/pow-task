# Proof of Work Server

## Overview
This project implements a Proof of Work (PoW) server to protect against DDoS attacks. The client must perform a computationally expensive operation to receive a response from the server, making large-scale automated requests more costly.

For more details, see the task description: [task.pdf](task.pdf).

## Architecture
The project follows the **Hexagonal Architecture** (Ports and Adapters), ensuring modularity and testability.

### Application Layer
The core business logic is located in the **application layer**:

- **POWVerifier** ([./internal/application/pow_verifier.go](./internal/application/pow_verifier.go))
    - Generates PoW.
    - Waits for a nonce from the server.
    - Verifies the nonce and either returns a **grant** (a quote from Word Of Wisdom) to the client or an error.

- **POWSolver** ([./internal/application/pow_solver.go](./internal/application/pow_solver.go))
    - Reads the PoW challenge.
    - Searches for the correct nonce and submits the solution (client-side logic).
    - Receives the **grant** as a result and returns it upstream.

## Installation
Ensure you have **Go** and **Docker** installed before proceeding.

## Commands
The project provides a `Makefile` to simplify development tasks. Below are the available commands:

### Linting
```sh
make lint
```
Runs `golangci-lint` to check code quality.

### Running Tests
```sh
make test
```
Executes the test suite.

### Building Docker Images
```sh
make build-server
```
Builds the Docker image for the server.

```sh
make build-client
```
Builds the Docker image for the client.

### Running Containers
```sh
make run-server
```
Starts the server container.

```sh
make run-client
```
Starts the client container.

## Running from Source
You can also run the server and client directly from the source code:

```sh
go run ./cmd/server/main.go
```
Runs the server.

```sh
go run ./cmd/client/main.go
```
Runs the client.

## License
This project is licensed under the MIT License.

