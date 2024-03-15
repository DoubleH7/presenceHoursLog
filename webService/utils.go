package webService

import (
	"log"
	"os"
)

func getWeblogger() *log.Logger {
	file, err := openLogFile("./webService/weblog.log")
	if err != nil {
		log.Fatal(err)
	}

	logger := log.New(file, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)

	return logger
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
