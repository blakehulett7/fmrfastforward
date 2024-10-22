package main

import (
	"os"
	"testing"
)

func TestInstaller(t *testing.T) {
	if fileExists(dbPath) {
		os.Rename(dbPath, testdbPath)
		t.Cleanup(func() {
			os.Rename(testdbPath, dbPath)
		})
	}
	defer os.Remove(dbPath)
	install(dbPath)
}
