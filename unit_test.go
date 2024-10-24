package main

import (
	//"os"
	"testing"
	//"github.com/google/uuid"
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

// This is a scuffed solution to a problem I don't want to actually solve right now
func TestWriteCard(t *testing.T) {
}

func TestGenerateStarterDeck(t *testing.T) {
	sim := simulation{starting_seed: 1, current_seed: 1}
	sim.generate_starting_deck()
}

func TestEvaluateStarterDeck(t *testing.T) {
	sim := simulation{starting_seed: 1, current_seed: 1}
	evaluate_starting_deck(sim.generate_starting_deck())
}
