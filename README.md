# 🚀 Proof of Work Server 🔒

![Coverage](https://codecov.io/gh/beetlewar010785/pow-task/branch/main/graph/badge.svg)

## 📌 Overview
This project implements a **Proof of Work (PoW) server** to protect against **DDoS attacks**. The client must perform a computationally expensive operation to receive a response from the server, making large-scale automated requests more costly.

📄 **For more details, see the task description:** [task.pdf](task.pdf).

---

## 🏗️ Architecture
The project follows the **Hexagonal Architecture** (Ports and Adapters), ensuring modularity and testability.

### 🧩 Application Layer
The core business logic is located in the **application layer**:

- 🔍 **POWVerifier** ([./internal/application/pow_verifier.go](./internal/application/pow_verifier.go))
  - Generates PoW.
  - Waits for a nonce from the server.
  - Verifies the nonce and either returns a **grant** (a quote from Word Of Wisdom) to the client or an error.

- 🛠️ **POWSolver** ([./internal/application/pow_solver.go](./internal/application/pow_solver.go))
  - Reads the PoW challenge.
  - Searches for the correct nonce and submits the solution (client-side logic).
  - Receives the **grant** as a result and returns it upstream.

- ### 🌐 Server (POWServer)
- 📂 **Implementation:** [./internal/adapter/pow_server.go](./internal/adapter/pow_server.go)
- **Description:**
  - `POWServer` is a **TCP server** that enforces PoW protection.
  - When a client connects, it **sends a challenge** (a computational puzzle).
  - The client must solve the challenge and send back a **valid nonce**.
  - If the nonce is correct, the server **grants access** by sending a **quote from "Word of Wisdom"**.
  - If verification fails, the server returns a **failure message**.

- ### 🖥️ Client (POWClient)
- 📂 **Implementation:** [./internal/adapter/pow_client.go](./internal/adapter/pow_client.go)
- **Description:**
  - The client connects to the `POWServer` over TCP.
  - It receives a **PoW challenge** and computes the correct **nonce**.
  - Once solved, the client sends the nonce back to the server.
  - If successful, the server responds with a **grant** (a motivational quote).
  - The client **prints the quote to the console** and **exits**.

---

## ⚙️ Installation
Ensure you have **Go** and **Docker** installed before proceeding.

---

## 🎯 Commands
The project provides a `Makefile` to simplify development tasks. Below are the available commands:

### ✅ Linting
```sh
make lint
```
🔍 Runs `golangci-lint` to check code quality.

### 🧪 Running Tests
```sh
make test
```
🛠️ Executes the test suite.

### 🧪 Running Integration Tests
```sh
./integration-test.sh
```
🛠️ Executes the integration test suite (server and client docker images must be built).

### 📦 Building Docker Images
```sh
make build-server
```
🐳 Builds the Docker image for the server.

```sh
make build-client
```
🐳 Builds the Docker image for the client.

### 🚀 Running Containers
```sh
make run-server
```
▶️ Starts the server container.

```sh
make run-client
```
▶️ Starts the client container.

---

## 🏃 Running from Source
You can also run the server and client directly from the source code:

```sh
go run ./cmd/server/main.go
```
🖥️ Runs the server.

```sh
go run ./cmd/client/main.go
```
🖥️ Runs the client.

---

## 📌 Areas for Improvement
⚠️ **Potential improvements and known issues:**
- 📝 **StringReadWriter serializer** ([./internal/adapter/string_read_writer.go](./internal/adapter/string_read_writer.go)) is **not optimized** and may require performance improvements.
- ❌ **Large portions of negative test scenarios** are **not covered**, including timeout handling for PoW verification.

---

## 📜 License
This project is licensed under the **MIT License**.

🛠️ Happy Coding! 🚀

