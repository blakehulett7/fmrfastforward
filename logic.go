package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
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

	fusions := []Fusion{}
	for _, card := range cards {
		for _, target_card := range cards {
			for _, fusion := range card.m1_potential {
				if slices.Contains(target_card.m2_potential, fusion) {
					//fmt.Printf("Found fusion! m1: %v, m2: %v, result: %v\n", card.Name, target_card.Name, fusion)
					fusion_card := get_card(strings.Split(fusion, "|")[0]) //Huge performace hit here just fyi
					num, _ := strconv.Atoi(strings.Split(fusion, "|")[1])
					assert(num > 0, "bad conversion")
					fusion := Fusion{
						Name:          fusion_card.Name,
						Attack:        fusion_card.Attack,
						Defense:       fusion_card.Defense,
						fusion_number: num,
						m1:            card.Name,
						m2:            target_card.Name,
					}
					if !slices.Contains(fusions, fusion) {
						fusions = append(fusions, fusion)
					}
				}
			}
		}
	}

	fusions_by_atk := fusions
	slices.SortFunc(fusions_by_atk, func(a Fusion, b Fusion) int {
		return b.Attack - a.Attack
	})
	var chance_to_draw float64 = 2.0
	var percent_chance float64
	for _, card := range fusions {
		chance_to_not_draw := 40 - chance_to_draw
		first_hand_percent_change_to_draw := (1 - ((chance_to_not_draw / 40.0) * ((chance_to_not_draw - 1) / 39.0))) * 100
		percent_chance += first_hand_percent_change_to_draw
		fmt.Println(card.Name, card.Attack, card.fusion_number, card.m1, card.m2)
	}
	fmt.Println()

	deck_by_atk := starting_deck
	slices.SortFunc(deck_by_atk, func(a Card, b Card) int {
		return b.Attack - a.Attack
	})
	chance_to_draw = 0
	for _, card := range deck_by_atk {
		chance_to_draw++
		chance_to_not_draw := 40 - chance_to_draw
		first_hand_percent_change_to_draw := (1 - ((chance_to_not_draw / 40.0) * ((chance_to_not_draw - 1) / 39.0) * ((chance_to_not_draw - 2) / 38.0) * ((chance_to_not_draw - 3) / 37.0) * ((chance_to_not_draw - 4) / 36.0))) * 100
		fmt.Println(card.Name, card.Attack, fmt.Sprintf("%.2f%%", first_hand_percent_change_to_draw))
	}

	fmt.Println()
	sim := simulation{starting_seed: 1, current_seed: 1} //temporary
	sim.draw_cards(starting_deck, 5)
}
