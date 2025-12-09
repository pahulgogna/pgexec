package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var RootDir string
var CodeFileName string
var LogFile string
var GenerateLogFile bool
var RequestTimeout int

func Init() {

	_ = godotenv.Load()

	RootDir = getEnv("RootDir", "/home/code")
	CodeFileName = getEnv("CodeFileName", "Main")
	LogFile = getEnv("LogFile", "logs.txt")
	GenerateLogFile = getEnv("GenerateLogFile", "true") == "true"
	RequestTimeout, _ = strconv.Atoi(getEnv("RequestTimeout", "30"))
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	fmt.Printf("fallback to default %s\n", key)
	return fallback
}
