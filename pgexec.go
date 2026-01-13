package pgexec

import (
	"fmt"
	"strings"

	"github.com/pahulgogna/pgexec/docker"
)

type Snippet struct {

	language     string
	code         string
	dependencies []string
}

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


