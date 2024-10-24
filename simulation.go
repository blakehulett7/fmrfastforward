package main

import (
	"fmt"
	"math/rand/v2"
)

type simulation struct {
	starting_seed uint64
	current_seed  uint64
}

func (sim *simulation) increment_seed() {
	sim.current_seed++
}

func (sim simulation) generate_starting_deck() []Card {
	cards_to_get := []string{}

	pool_sub_1100 := get_starting_deck_rates("pool_sub_1100")
	pool_1100_1600 := get_starting_deck_rates("pool_1100_1600")
	pool_1600_2100 := get_starting_deck_rates("pool_1600_2100")
	pool_over_2100 := get_starting_deck_rates("pool_over_2100")
	pool_pure_magic := get_starting_deck_rates("pool_pure_magic")
	pool_field_magic := get_starting_deck_rates("pool_field_magic")
	pool_equip_magic := get_starting_deck_rates("pool_equip_magic")

	cards_to_get = append(cards_to_get, sim.get_cards_from_pool(pool_sub_1100, 16)...)
	cards_to_get = append(cards_to_get, sim.get_cards_from_pool(pool_1100_1600, 16)...)

	cards_to_get = append(cards_to_get, sim.get_cards_from_pool(pool_1600_2100, 4)...)

	cards_to_get = append(cards_to_get, sim.get_cards_from_pool(pool_over_2100, 1)...)
	cards_to_get = append(cards_to_get, sim.get_cards_from_pool(pool_pure_magic, 1)...)
	cards_to_get = append(cards_to_get, sim.get_cards_from_pool(pool_field_magic, 1)...)
	cards_to_get = append(cards_to_get, sim.get_cards_from_pool(pool_equip_magic, 1)...)

	assert(len(cards_to_get) == 40, "bug generating the starting deck, did not get exactly 40 cards to grab from the db")
	deck := []Card{}
	for _, card := range cards_to_get {
		deck = append(deck, get_card(card))
	}
	return deck
}

func (sim *simulation) drop_card(drop_table []Probability) (card_name string) {
	source := rand.NewPCG(sim.current_seed, 0)
	rng := rand.New(source)
	table_selector := rng.IntN(2048)
	sim.increment_seed()
	for _, entry := range drop_table {
		table_selector = table_selector - entry.Rate
		if table_selector < 0 {
			return entry.Card
		}
	}
	panic("Something went wrong, no card was chosen from the drop table")
}

func (sim simulation) get_cards_from_pool(pool []Probability, num_cards int) []string {
	cards := []string{}
	for i := 0; i < num_cards; i++ {
		var card string
		for {
			card = sim.drop_card(pool)
			if count_instances_in_slice(cards, card) < 3 {
				break
			}
		}
		cards = append(cards, card)
	}

	for _, card := range cards {
		assert(count_instances_in_slice(cards, card) < 4, "Can't have more than 3 of the same card in a deck")
	}
	return cards
}

func (sim simulation) draw_cards(deck []Card, number int) {
	for i := 0; i < number; i++ {
		card_index := rand.IntN(len(deck))
		sim.increment_seed()
		fmt.Print(deck[card_index].Name, ", ")
		deck = append(deck[:card_index], deck[card_index+1:]...)
	}
}
