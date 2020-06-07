package utils

import (
	"log"
	"os"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetCurrentWorkingDirectory gets the working directory and returns its path as a string
func GetCurrentWorkingDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return dir
}
