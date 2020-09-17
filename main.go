package main

import (
	"log"
	"os"
)

var flog *log.Logger

func main() {
	f, err := os.OpenFile("diff.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger(err)
	}
	defer f.Close()
	flog = log.New(f, "", log.LstdFlags)

	logger("== diff-hackerone ==")

	connectToDatabase()
	directory := getDirectory()
	storedDirectoryCount := getStoredDirectoryCount()

	if storedDirectoryCount > 0 {
		updateDirectory(directory)
	} else {
		insertFullDirectory(directory)
	}

	logger("== end diff-hackerone ==")
	logger("")
}

func logger(message interface{}) {
	switch message.(type) {
	case string, error:
		log.Print(message)
		flog.Print(message)

	default:
		logger("Unknown log type")
	}
}
