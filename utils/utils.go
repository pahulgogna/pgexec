package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pahulgogna/pgexec/config"
)

func GetTagFromLanguage(language string) string {
	var tag string = ""

	if language == "python" {
		tag = "python:alpine"
	}

	return tag
}

func GetInstallCommandFromLanguage(language string) string {

	var command string = ""

	if language == "python" {
		command = "pip install "
	}

	return command
}

func GetFileExtension(language string) string {
	var extension string = ""

	if language == "python" {
		extension = "py"
	}

	return extension
}

func GetRunCodeCommand(language string) string {
	var runCommand string = ""
	var filePath string = fmt.Sprintf("%s/%s.%s", config.RootDir, config.CodeFileName, GetFileExtension(language))

	if language == "python" {
		runCommand = fmt.Sprintf("python %s", filePath)
	}

	return runCommand
}

func WriteToLogFile(id string, data string) {
	if !config.GenerateLogFile {
		return
	}
	
	f, err := os.OpenFile(config.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("ERROR:", err)
	}
	if _, err := fmt.Fprintf(f, "%s, %s, %s\n", time.Now().String(), id, data); err != nil {
		log.Println("ERROR:", err)
	}
	if err := f.Close(); err != nil {
		log.Println("ERROR:", err)
	}
}
