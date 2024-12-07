package main

import (
	"fmt"
	"slices"
)

func og_starter() {
	fmt.Println("Christ is King!")
}

func generate_all_possible_hands(num_cards_in_deck int) [][]int {
	assert(num_cards_in_deck == 40, "deck has to have a length of 40 cards")

	hands := [][]int{}
	for c1 := 1; c1 <= num_cards_in_deck; c1++ {
		for c2 := c1 + 1; c2 <= num_cards_in_deck; c2++ {
			for c3 := c2 + 1; c3 <= num_cards_in_deck; c3++ {
				for c4 := c3 + 1; c4 <= num_cards_in_deck; c4++ {
					for c5 := c4 + 1; c5 <= num_cards_in_deck; c5++ {
						hands = append(hands, []int{c1, c2, c3, c4, c5})
					}
				}
			}
		}
	}
	fmt.Println("all hands generated...")
	fmt.Println()

	assert(len(hands) == 658008, "possible hands generated is incorrect")

	return hands
}

func find_successful_hands(hands [][]int) int {
	successes := 0
	for _, hand := range hands {
		if !slices.Contains(hand, 1) {
			continue
		}
		if !slices.Contains(hand, 2) {
			continue
		}
		successes++
	}
	return successes
}
