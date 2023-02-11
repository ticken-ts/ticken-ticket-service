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
	ExecEnvKey         = "ENV"
	DbConnStringEnvKey = "DB_CONN_STRING"
	BusConnStringKey   = "BUS_CONN_STRING"
	ConfigFilePath     = "CONFIG_FILE_PATH"
	ConfigFileName     = "CONFIG_FILE_NAME"
	HSMEncryptionKey   = "HSM_ENCRYPTION_KEY"
	TickenWalletKey    = "TICKEN_WALLET_KEY"
)

type Env struct {
	Env              string
	DbConnString     string
	BusConnString    string
	ConfigFilePath   string
	ConfigFileName   string
	HSMEncryptionKey string
	TickenWalletKey  string
}

func Load() (*Env, error) {
	if fileExists(envFilename) {
		err := godotenv.Load(envFilename)
		if err != nil {
			return nil, err
		}
	}

	env := &Env{
		Env:              getEnvOrPanic(ExecEnvKey),
		DbConnString:     getEnvOrPanic(DbConnStringEnvKey),
		BusConnString:    getEnvOrPanic(BusConnStringKey),
		ConfigFilePath:   getEnvOrPanic(ConfigFilePath),
		ConfigFileName:   getEnvOrPanic(ConfigFileName),
		HSMEncryptionKey: getEnvOrPanic(HSMEncryptionKey),
		TickenWalletKey:  getEnvOrPanic(TickenWalletKey),
	}

	TickenEnv = env
	return env, nil
}

func getEnvOrPanic(key string) string {
	envVal := os.Getenv(key)
	if len(envVal) == 0 {
		panic(fmt.Errorf("env var %s is mandatory", key))
	}
	return envVal
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
