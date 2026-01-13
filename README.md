# pgexec

`pgexec` is a lightweight Go package designed for secure and isolated code execution using Docker. It provides a simple API to run arbitrary code snippets in language-specific containers, ensuring environment isolation and dependency management.

## Features
- **Isolated Execution**: Runs code inside ephemeral Docker containers (currently using `python:alpine`).
- **Dependency Management**: Automatically installs required packages (e.g., via `pip`) before execution.
- **Configurable**: Easily adjust timeouts, root directories, and logging via environment variables.
- **Simple API**: Focus on code execution with minimal boilerplate.

## Prerequisites

- **Go**: Version 1.22 or higher.
- **Docker**: Must be installed and running on the host machine.

## Installation

```bash
go get github.com/pahulgogna/pgexec
```

## Supported Languages

| Language | Docker Image | Package Manager |
| :--- | :--- | :--- |
| **Python** | `python:alpine` | `pip` |

## Usage

Here is a quick example of how to use `pgexec` to run a Python snippet with dependencies:

```go
package main

import (
	"fmt"
	"log"

	"github.com/pahulgogna/pgexec"
)

func main() {
	// Define a snippet with language, code, and dependencies
	code := `import requests; print(requests.get("https://google.com").status_code)`
	snippet, err := pgexec.NewSnippet("python", code, "requests")
	if err != nil {
		log.Fatal(err)
	}

	// Execute the code
	output, err := snippet.Execute()
	if err != nil {
		fmt.Printf("Execution Error: %v\n", err)
	}

	fmt.Println("Output:")
	fmt.Println(output)
}
```

## Configuration

`pgexec` can be configured using environment variables (or a `.env` file):

| Variable | Description | Default |
| :--- | :--- | :--- |
| `RequestTimeout` | Maximum time (in seconds) for a command to run. | `30` |
| `WriteLog` | Whether to log detailed execution info to stdout. | `false` |
| `RootDir` | The directory inside the container where code is stored. | `/home/code` |
| `CodeFileName` | The filename used for the snippet inside the container. | `Main` |

## License

This project is licensed under the MIT License.
