package main

import (
	"os"
	"testing"
)

func TestInstaller(t *testing.T) {
	defer os.Remove(dbPath)
	getFmrData()
}
