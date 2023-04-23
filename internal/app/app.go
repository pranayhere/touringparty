package app

import (
	"os"
)

// GetEnv returns the current environment, prod, dev etc.
func GetEnv() string {
	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "dev"
	}

	return environment
}

// IsTestMode return true if current execution environment is test mode
func IsTestMode() bool {
	env := GetEnv()
	if env == "test" || env == "dev" || env == "dev_docker" {
		return true
	}

	return false;
}