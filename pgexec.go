// Package pgexec provides a simple API to execute code snippets in various
// programming languages using Docker containers for isolation.
package pgexec

import (
	"fmt"
	"strings"

	"github.com/pahulgogna/pgexec/docker"
)

/*
	Snippet represents an executable code snippet in a specific language with its dependencies.
*/
type Snippet struct {
	language     string
	code         string
	dependencies []string
}

/*
	NewSnippet creates a new executable code snippet. It returns an error if the
	language is not supported or if the code is empty.
*/
func NewSnippet(language string, code string, dependencies ...string) (*Snippet, error) {
	if !docker.IsSupportedLanguage(language) {
		return nil, fmt.Errorf("unsupported language")
	}

	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("the code snippet provided is empty")
	}

	return &Snippet{
		language:     language,
		code:         code,
		dependencies: dependencies,
	}, nil
}

/*
	Execute runs the code snippet within an isolated Docker container and returns
	the combined stdout and stderr. It handles the environment initialization and
	cleanup automatically.
*/
func (s *Snippet) Execute() (string, error) {
	environment, err := docker.NewEnvironment(s.language, s.dependencies...)
	if err != nil {
		return "", err
	}

	output, err := environment.Run(s.code)
	if err != nil {
		return output, err
	}

	return output, nil
}
