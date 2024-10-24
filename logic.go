package main

import (
	"fmt"
	"slices"
)

func evaluate_starting_deck(starting_deck []Card) {
	assert(len(starting_deck) == 40, "Decks must have exactly 40 cards")

	cards := []Card{}
	for _, card := range starting_deck {
		card.m1_potential = get_potential_fusions(card.Name, "m1")
		card.m2_potential = get_potential_fusions(card.Name, "m2")
		//fmt.Printf("%v:\n%v\n%v\n\n", card.Name, card.m1_potential, card.m2_potential)
		cards = append(cards, card)
	}

	for _, card := range cards {
		for _, target_card := range cards {
			for _, fusion := range card.m1_potential {
				if slices.Contains(target_card.m2_potential, fusion) {
					fmt.Printf("Found fusion! m1: %v, m2: %v, result: %v\n", card.Name, target_card.Name, fusion)
				}
			}
		}
	}

	deck_by_atk := starting_deck
	slices.SortFunc(deck_by_atk, func(a Card, b Card) int {
		if a.Attack < b.Attack {
			return 1
		}
		if a.Attack == b.Attack {
			return 0
		}
		return -1
	})
	for _, card := range deck_by_atk {
		fmt.Println(card.Name, card.Attack)
	}
}
