package executor

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pahulgogna/pgexec/config"
	"github.com/pahulgogna/pgexec/customTypes"
	"github.com/pahulgogna/pgexec/utils"
)

func executeCommand(id string, name string, cmdString ...string) (string, error) {

	if name == "" {
		return "", fmt.Errorf("no command provided")
	}

	var cmd *exec.Cmd

	command := fmt.Sprintf("%s %s", name, strings.Join(cmdString, " "))

	utils.WriteToLogFile(id, fmt.Sprintf("executing command: %s", command))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.RequestTimeout)*time.Second)
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

func StartDockerContainer(id string) (string, error) {
	output, err := executeCommand(id, "docker", "run", "-itd", "--rm", id)
	if err != nil {
		return output, err
	}

	output = strings.Trim(output, "\n ")

	executeCommandInDocker(output, "mkdir", "-p", config.RootDir)

	return output, err
}

func executeCommandInDocker(id string, command ...string) (string, error) {
	cmdString := strings.Join(command, " ")

	args := []string{"exec", "-i", id, "sh", "-c", cmdString}
	output, err := executeCommand(id, "docker", args...)

	return output, err
}

func StopDockerContainer(id string) bool {
	_, err := executeCommand(id, "docker", "kill", id)
	if err != nil {
		utils.WriteToLogFile(id, "could not stop the container, check if it is running\n")
		return false
	}

	utils.WriteToLogFile(id, "stopped docker container")
	return true

}

func InstallDependencies(id string, snippet *customtypes.Snippet) string {

	if len(snippet.Dependencies) == 0 {
		return "None"
	}

	installCommand := utils.GetInstallCommandFromLanguage(snippet.Language)
	if installCommand == "" {
		return "unsupported language"
	}

	for _, dep := range snippet.Dependencies {
		executeCommandInDocker(id, installCommand, dep)
	}

	return "Done"
}

func RunCode(id string, snippet *customtypes.Snippet) (string, error) {
	extension := utils.GetFileExtension(snippet.Language)
	if extension == "" {
		return "", fmt.Errorf("unsupported language")
	}

	filePath := fmt.Sprintf("%s/%s.%s", config.RootDir, config.CodeFileName, extension)

	_, err := executeCommandInDocker(id, "echo", fmt.Sprintf("\"%s\"", strings.ReplaceAll(snippet.Code, "\"", "\\\"")), ">", filePath)
	if err != nil {
		return "Code file could not be created.", err
	}

	output, err := executeCommandInDocker(id, utils.GetRunCodeCommand(snippet.Language))

	return output, err
}
