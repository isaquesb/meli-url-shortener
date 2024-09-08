package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const LOCAL = "local"

func Environment() string {
	appEnv, hasEnv := os.LookupEnv("ENVIRONMENT")
	if !hasEnv {
		return LOCAL
	}
	return appEnv
}

func IsLocal() bool {
	return Environment() == LOCAL
}

func LoadEnv() {
	if !IsLocal() {
		return
	}

	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func GetEnv(args ...string) string {
	key := args[0]
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	defArg := ""
	if len(args) > 1 {
		defArg = args[1]
	}
	if "" == defArg {
		panic("Environment variable " + key + " not found")
	}
	return defArg
}

func GetIntEnv(args ...string) int {
	env := GetEnv(args...)
	intEnv, err := strconv.Atoi(env)
	if err != nil {
		panic(err)
	}
	return intEnv
}
