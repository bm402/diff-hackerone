package main

import (
	"log"
)

func main() {
	log.Print("== diff-hackerone ==")

	connectToDatabase()
	directory := getDirectory()
	storedDirectoryCount := getStoredDirectoryCount()

	if storedDirectoryCount > 0 {
		updateDirectory(directory)
	} else {
		insertFullDirectory(directory)
	}

	log.Print("== end diff-hackerone ==")
	log.Print("")
}
