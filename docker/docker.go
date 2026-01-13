package docker

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pahulgogna/pgexec/config"
)

type Environment struct {
	tag      string
	language string
}

func NewEnvironment(language string, dependencies ...string) (*Environment, error) {

	var tag string

	writeToLog(fmt.Sprintf("NEW: %s", language), "INFO: starting docker container")
	tag, err := executeCommandHost("docker", "run", "-itd", "--rm", getImageNameFromLanguage(language))
	if err != nil {
		return nil, err
	}
	tag = strings.Trim(tag, "\n ")
	writeToLog(tag, "INFO: started container")

	env := Environment{
		tag:      tag,
		language: language,
	}

	if err := env.setupEnvironment(dependencies...); err != nil {
		writeToLog(tag, fmt.Sprintf("ERROR: error setting up the environment: %s", err))
		return nil, err
	}

	return &env, nil
}

func executeCommandHost(name string, cmdString ...string) (string, error) {

	if name == "" {
		return "", fmt.Errorf("no command provided")
	}

	var cmd *exec.Cmd

	ctx, cancel := context.WithTimeout(context.Background(), config.Env.RequestTimeout)
	defer cancel()

	cmd = exec.CommandContext(ctx, name, cmdString...)

	var out bytes.Buffer

	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("command timed out")
	}

	if err != nil {
		return out.String(), err
	}

	return out.String(), nil
}

func (e *Environment) StopEnvironment() bool {
	_, err := executeCommandHost("docker", "kill", e.tag)
	if err != nil {
		return false
	}

	return true
}

func (e *Environment) executeCommand(command ...string) (string, error) {
	cmdString := strings.Join(command, " ")

	args := []string{"exec", "-i", e.tag, "sh", "-c", cmdString}
	output, err := executeCommandHost("docker", args...)

	return output, err
}

func (e *Environment) setupEnvironment(dependencies ...string) error {

	// create root dir
	if out, err := e.executeCommand("mkdir", config.Env.Docker.RootDir); err != nil {
		return fmt.Errorf("ERROR: could not create root dir, err: %s || output: %s", err, out)
	}

	// dependencies
	if len(dependencies) == 0 {
		return nil
	}

	installCommand := getInstallCommandFromLanguage(e.language)
	if installCommand == "" {
		defer e.StopEnvironment()
		return fmt.Errorf("ERROR: unsupported language\n")
	}

	for _, dep := range dependencies {
		writeToLog(e.tag, fmt.Sprintf("INFO: installing dependency: %s", dep))
		_, err := e.executeCommand(installCommand, dep)
		if err != nil {
			writeToLog(e.tag, fmt.Sprintf("ERROR: error while installing dependency: %s, error: %s", dep, err))
			continue
		}

		writeToLog(e.tag, fmt.Sprintf("INFO: installed dependency: %s", dep))
	}

	return nil
}

func (e *Environment) run(code string) (string, error) {

	extension := getFileExtension(e.language)
	if extension == "" {
		return "", fmt.Errorf("unsupported language")
	}

	filePath := fmt.Sprintf("%s/%s.%s", config.Env.Docker.RootDir, config.Env.Docker.CodeFileName, extension)

	_, err := e.executeCommand("echo", fmt.Sprintf("\"%s\"", strings.ReplaceAll(code, "\"", "\\\"")), ">", filePath)
	if err != nil {
		return "ERROR: Code file could not be created inside the docker container.", err
	}

	output, err := e.executeCommand(getRunCodeCommand(e.language))

	return output, err
}

func (e *Environment) Run(code string) (string, error) {
	defer e.StopEnvironment()
	return e.run(code)
}

func (e *Environment) RunAndKeep(code string) (string, error) {
	return e.run(code)
}
