package utils

import (
	"fmt"
	"os"
)

func GetEnvOrDefault(key string, defaultValue string) string {
	value, defined := os.LookupEnv(key)
	if !defined {
		return defaultValue
	}

	return value
}

func GetEnvOrError(key string) (string, error) {
	value, defined := os.LookupEnv(key)
	if !defined {
		return "", fmt.Errorf("env var %s not defined", key)
	}

	return value, nil
}

func GetEnvOrPanic(key string) string {
	value, err := GetEnvOrError(key)
	if err != nil {
		panic(err)
	}
	return value
}

const TickenEnvName = "TICKEN_ENV"

const (
	DevEnv  = "Dev"
	ProdEnv = "Prod"
	TestEnv = "Test"
)

func IsDev() bool {
	return GetEnvOrPanic(TickenEnvName) == DevEnv
}

func IsProd() bool {
	return GetEnvOrPanic(TickenEnvName) == ProdEnv
}

func Istest() bool {
	return GetEnvOrPanic(TickenEnvName) == TestEnv
}
