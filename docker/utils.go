package docker 

import (
	"fmt"
	"log"
	"time"

	"github.com/pahulgogna/pgexec/config"
)

func IsSupportedLanguage(language string) bool {

	if language == "python" {
		return true
	}

	return false
}

func getImageNameFromLanguage(language string) string {

	if language == "python" {
		return "python:alpine"
	}

	return ""
}

func getInstallCommandFromLanguage(language string) string {

	if language == "python" {
		return"pip install "
	}

	return ""
}

func getFileExtension(language string) string {

	if language == "python" {
		return "py"
	}

	return ""
}

func getRunCodeCommand(language string) string {

	var filePath string = fmt.Sprintf("%s/%s.%s", config.Env.Docker.RootDir, config.Env.Docker.CodeFileName, getFileExtension(language))

	if language == "python" {
		return  fmt.Sprintf("python %s", filePath)
	}

	return ""
}

func writeToLog(id string, data string) {
	if !config.Env.WriteToLog {
		return
	}

	ts := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("|| %s || %s || %s\n", ts, id, data)
}
