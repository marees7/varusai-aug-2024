package loggers

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	InfoLog  *log.Logger
	FatalLog *log.Logger
	ErrorLog *log.Logger
	WarnLog  *log.Logger
)

func recoverLoadEnv() {
	if res := recover(); res != nil {
		fmt.Println("recovered from ", res)
	}
}

func ForLogs() {
	defer recoverLoadEnv()

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(filepath.Join(filepath.Dir(workingDir), os.Getenv("LOGGERS_PATH")), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		panic(err)
	}

	InfoLog = log.New(file, "INFO:", log.LstdFlags|log.Llongfile)
	FatalLog = log.New(file, "FATAL:", log.LstdFlags|log.Llongfile)
	ErrorLog = log.New(file, "ERROR:", log.LstdFlags|log.Llongfile)
	WarnLog = log.New(file, "WARNING:", log.LstdFlags|log.Llongfile)
}
