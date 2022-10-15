package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

var TickenEnv *Env = nil

const envFilename = ".env"

const (
	DevEnv   = "dev"
	StageEnv = "stage"
	ProdEnv  = "prod"
	TestEnv  = "test"
)

const (
	ExecEnvKey       = "ENV"
	ConnStringEnvKey = "CONN_STRING"
)

type Env struct {
	Env        string
	ConnString string
}

func Load() (*Env, error) {
	if fileExists(envFilename) {
		err := godotenv.Load(envFilename)
		if err != nil {
			return nil, err
		}
	}

	execEnv, err := getEnvOrError(ExecEnvKey)
	if err != nil {
		return nil, err
	}

	connString, err := getEnvOrError(ConnStringEnvKey)
	if err != nil {
		return nil, err
	}

	env := &Env{
		Env:        execEnv,
		ConnString: connString,
	}

	TickenEnv = env
	return env, nil
}

func getEnvOrError(key string) (string, error) {
	envVal := os.Getenv(key)
	if len(envVal) == 0 {
		return "", fmt.Errorf("env var %s is mandatory", key)
	}
	return envVal, nil
}

func (env *Env) IsDev() bool {
	return env.Env == DevEnv
}

func (env *Env) IsProd() bool {
	return env.Env == ProdEnv
}

func (env *Env) IsTest() bool {
	return env.Env == TestEnv
}

func (env *Env) IsStage() bool {
	return env.Env == StageEnv
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
