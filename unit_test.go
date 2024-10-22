package main

import (
	"os"
	"testing"
)

func TestInstaller(t *testing.T) {
	if fileExists(dbPath) {
		os.Rename(dbPath, dbHoldingPath)
		t.Cleanup(func() {
			os.Rename(dbHoldingPath, dbPath)
		})
	}
	defer os.Remove(dbPath)
	install(dbPath)
}
