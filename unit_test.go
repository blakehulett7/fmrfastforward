package main

import (
	//"os"
	//"fmt"
	"testing"
	//"github.com/montanaflynn/stats"
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

func TestGenerateStarterDeck(t *testing.T) {
	sim := simulation{starting_seed: 1, current_seed: 1}
	sim.generate_starting_deck()
}

/*
func TestEvaluateStarterDeck(t *testing.T) {
	sim := simulation{starting_seed: 1, current_seed: 1}
	deck := sim.generate_starting_deck()
	evaluate_starting_deck(deck)
}

func TestMath(t *testing.T) {
	assert(stats.Ncr(40, 5) == 658008, "Ncr algo is broken")
}

func TestBasicFusionOdds(t *testing.T) {
	sim := simulation{starting_seed: 1, current_seed: 1}
	deck := sim.generate_starting_deck()
	fmt.Println()
	//for i := 0; i < 100; i++ {
	starting_hand := sim.draw_cards(deck[:], 5)
	get_best_fusion(starting_hand)
	//}
}
*/
