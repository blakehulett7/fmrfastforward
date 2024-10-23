package main

import (
	//"os"
	"testing"
)

/*
func TestInstaller(t *testing.T) {
	if fileExists(dbPath) {
		os.Rename(dbPath, dbHoldingPath)
		defer os.Rename(dbHoldingPath, dbPath)
	}
	defer os.Remove(dbPath)
	install(dbPath)
}
*/

func TestGenerateStarterDeck(t *testing.T) {
	sim := simulation{seed: 1, current_seed: 1}
	sim.generate_starting_deck()
}
