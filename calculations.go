package main

import (
	"fmt"
	"slices"

	"github.com/montanaflynn/stats"
)

func add_fusion_odds(fusion_1, fusion_2 Fusion) float64 {
	odds_1 := odds_of_drawing_fusion(fusion_1, 5)
	odds_2 := odds_of_drawing_fusion(fusion_2, 5)

	overlapping_m1 := []string{}
	for _, m1 := range fusion_1.m1_components {
		if slices.Contains(fusion_2.m1_components, m1) {
			overlapping_m1 = append(overlapping_m1, m1)
		}
	}
	fmt.Println(overlapping_m1)

	overlapping_m2 := []string{}
	for _, m2 := range fusion_1.m2_components {
		if slices.Contains(fusion_2.m2_components, m2) {
			overlapping_m2 = append(overlapping_m2, m2)
		}
	}
	fmt.Println(overlapping_m2)

	return odds_1 + odds_2
}

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
