# pgexec

`pgexec` is a lightweight, secure remote code execution (RCE) engine written in Go. It provides a REST API to execute code snippets in isolated environments using Docker containers.

## Features

- **Isolated Execution**: Runs code inside ephemeral Docker containers to ensure security and isolation.
- **Dependency Management**: Supports installing language-specific dependencies before execution.
- **REST API**: Simple HTTP interface built with [Gin](https://github.com/gin-gonic/gin).
- **Logging**: Tracks execution requests and errors in a log file.

## Supported Languages

Currently, the following languages are supported:

- **Python** (Image: `python:alpine`)

## Prerequisites

- **Go**: Version 1.23 or higher.
- **Docker**: Must be installed and running on the host machine.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/pahulgogna/pgexec.git
   cd pgexec
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

## Usage

1. Start the server:
   ```bash
   go run main.go
   ```
   The server will start on `0.0.0.0:8080`.

2. Make a request to execute code.

### API Endpoints

#### `GET /ping`

Health check endpoint.

- **Response**: `200 OK`
  ```json
  {
    "message": "pong"
  }
  ```

#### `POST /run`

Executes a code snippet.

- **Request Body**:
  ```json
  {
    "language": "python",
    "code": "print('Hello from pgexec!')",
    "dependencies": []
  }
  ```

- **Example with Dependencies**:
  ```json
  {
    "language": "python",
    "code": "import requests\nprint(requests.get('https://google.com').status_code)",
    "dependencies": ["requests"]
  }
  ```

- **Response**:
  ```json
  {
    "data": "Hello from pgexec!\n",
    "error": null
  }
  ```

## Configuration

Configuration constants are located in [`config/config.go`](config/config.go).

- `RootDir`: Directory inside the container where code is stored (default: `/home/code`).
- `LogFile`: Path to the log file stored on the host machine (default: `logs.txt`).
- `GenerateLogFile`: Boolean to enable/disable file logging.
- `RequestTimeout` : Timeout for commands in seconds.

## Project Structure

- [`main.go`](main.go): Entry point, sets up the Gin router.
- [`executor/`](executor/): Contains logic for Docker container management and command execution.
- [`utils`](utils ): Helper functions for language mapping and logging.
- [`customTypes`](customTypes ): Defines data structures like `Snippet`.
- [`config`](config ): Project configuration.
