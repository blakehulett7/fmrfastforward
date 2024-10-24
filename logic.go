package main

import (
	"fmt"
)

func evaluate_starting_deck(starting_deck []Card) {
	assert(len(starting_deck) == 40, "Decks must have exactly 40 cards")
	for _, card := range starting_deck {
		card.m1_potential = get_potential_fusions(card.Name, "m1")
		card.m2_potential = get_potential_fusions(card.Name, "m2")
		fmt.Printf("%v:\n%v\n%v\n\n", card.Name, card.m1_potential, card.m2_potential)
	}
}
