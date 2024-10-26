package main

import "github.com/montanaflynn/stats"

func odds_of_drawing_fusion(fusion Fusion, draws int) float64 {
	ways_to_draw_an_m1 := ways_to_draw_at_least_n_cards(len(fusion.m1_components), draws)
	ways_to_draw_an_m2 := ways_to_draw_at_least_n_cards(len(fusion.m2_components), draws)
	possible_hand_combinations := stats.Ncr(40, draws)
	return (float64(ways_to_draw_an_m1) / float64(possible_hand_combinations)) * (float64(ways_to_draw_an_m2) / float64(possible_hand_combinations))
}

func ways_to_draw_at_least_n_cards(n, draws int) int {
	ways_to_draw := 0
	for i := 1; i <= n; i++ {
		ways_to_draw += stats.Ncr(n, i) * stats.Ncr((40-n), (draws-i))
	}
	return ways_to_draw
}
