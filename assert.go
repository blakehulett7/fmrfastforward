package main

import (
	"errors"
	"io/fs"
	"os"
)

func assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}

func fileExists(path string) bool {
	_, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

func directoryExists(path string) bool {
	_, err := os.ReadDir(path)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}
