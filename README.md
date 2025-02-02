# Proof of Work Server

## Overview
This project implements a Proof of Work (PoW) server to protect against DDoS attacks. The client must perform a computationally expensive operation to receive a response from the server, making large-scale automated requests more costly.

For more details, see the task description: [task.pdf](task.pdf).

## Architecture
The project follows the **Hexagonal Architecture** (Ports and Adapters), ensuring modularity and testability.

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

