package main

import (
	//"os"
	"fmt"
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
	sim := simulation{starting_seed: 1, current_seed: 1}
	sim.generate_starting_deck()
}

func TestEvaluateStarterDeck(t *testing.T) {
	state_machine := state_machine{}
	sim := simulation{starting_seed: 1, current_seed: 1}
	state_machine.deck = sim.generate_starting_deck()
	fmt.Println(len(state_machine.deck))
	evaluate_starting_deck(state_machine.deck)
}
