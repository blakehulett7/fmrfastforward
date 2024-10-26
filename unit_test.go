package main

import (
	//"os"
	"fmt"
	"testing"

	"github.com/montanaflynn/stats"
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
/*
func TestWriteCard(t *testing.T) {
	write_cards_to_db([]Card{{
		Id:        uuid.NewString(),
		Name:      "Dragoness the Wicked Knight (FMR)",
		Type:      "Warrior",
		Attack:    1200,
		Defense:   900,
		StarChips: 60,
	}}, "cards")
}
*/

/*
func TestGenerateStarterDeck(t *testing.T) {
	sim := simulation{starting_seed: 1, current_seed: 1}
	sim.generate_starting_deck()
}

func TestEvaluateStarterDeck(t *testing.T) {
	sim := simulation{starting_seed: 1, current_seed: 1}
	evaluate_starting_deck(sim.generate_starting_deck())
}
*/

func TestMath(t *testing.T) {
	fmt.Println(stats.Ncr(40, 5))
}
