package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RequestTimeout time.Duration
	WriteToLog     bool
	Docker         *DockerConfig
}

type DockerConfig struct {
	RootDir      string
	CodeFileName string
}

var Env *Config = initConfig()

func initConfig() *Config {

	_ = godotenv.Load()

	requestTimeout, _ := strconv.Atoi(getEnv("RequestTimeout", "30"))

	return &Config{
		RequestTimeout: time.Duration(requestTimeout) * time.Second,
		WriteToLog:     getEnv("WriteLog", "false") == "true",

		Docker: &DockerConfig{
			RootDir:      getEnv("RootDir", "/home/code"),
			CodeFileName: getEnv("CodeFileName", "Main"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
