package main

import (
	"fmt"
)

func evaluate_starting_deck(starting_deck []Card) {
	fmt.Println(len(starting_deck))
	assert(len(starting_deck) == 40, "Decks must have exactly 40 cards")
	for _, card := range starting_deck {
		potential_fusions := get_potential_fusions(card.Name)
		fmt.Printf("%v: %v\n", card.Name, potential_fusions)
	}
}
