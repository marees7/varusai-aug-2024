package internals

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// recover from panic
func recoverLoadEnv() {
	if res := recover(); res != nil {
		fmt.Println("recovered from ", res)
	}
}

// load the env file
func LoadEnvFile() {
	defer recoverLoadEnv()

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	err = godotenv.Load(filepath.Join(filepath.Dir(workingDir), ".env"))
	if err != nil {
		panic(err)
	}
}
