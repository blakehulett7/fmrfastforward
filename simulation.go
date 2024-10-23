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

func (sim simulation) generate_starting_deck() {
	pool_sub_1100 := get_starting_deck_rates("pool_sub_1100")
	pool_1100_1600 := get_starting_deck_rates("pool_1100_1600")

	fmt.Println(sim.get_cards_from_pool(pool_sub_1100, 16))
	fmt.Println(sim.get_cards_from_pool(pool_1100_1600, 16))
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
		cards = append(cards, sim.drop_card(pool))
	}
	return cards
}
