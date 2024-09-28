package main

import (
	"fmt"
	"strconv"
	"strings"
)

func splitWikitext(wikitext string) (deckSlice, dropSlice []string) {
	//assert something about the wikitext
	assert(wikitext != "")
	wikitextslice := strings.Split(wikitext, "\n")
	deckIdx := 0
	dropIdx := 0
	dialogueIdx := 0
	for idx, line := range wikitextslice {
		normalized := strings.ReplaceAll(line, " ", "")
		if !strings.HasPrefix(normalized, "==D") {
			continue
		}
		if strings.HasPrefix(normalized, "==Deck") {
			deckIdx = idx
			continue
		}
		if strings.HasPrefix(normalized, "==Drop") {
			dropIdx = idx
			continue
		}
		if strings.HasPrefix(normalized, "==Dialogue") {
			dialogueIdx = idx
			deckSlice := wikitextslice[deckIdx:dropIdx]
			dropSlice := wikitextslice[dropIdx:dialogueIdx]
			return deckSlice, dropSlice
		}
		deckSlice := wikitextslice[deckIdx:dropIdx]
		dropSlice := wikitextslice[dropIdx:]
		return deckSlice, dropSlice
	}
	panic("Should never get here, something went wrong parsing the wikitext!")
}

func splitDuels(wikiSlice []string) [][]string {
	assert(len(wikiSlice) != 0)
	indices := []int{}
	for idx, line := range wikiSlice {
		if !strings.HasPrefix(line, "===") {
			continue
		}
		indices = append(indices, idx)
	}
	if len(indices) == 0 {
		return [][]string{wikiSlice}
	}
	wikiSlices := [][]string{}
	for idx := range indices {
		if idx == len(indices)-1 {
			wikiSlices = append(wikiSlices, wikiSlice[indices[idx]:])
			break
		}
		wikiSlices = append(wikiSlices, wikiSlice[indices[idx]:indices[idx+1]])
		fmt.Println()
	}
	assert(len(wikiSlices) != 0)
	return wikiSlices
}

func getDecks(deckslice []string) {
	for _, line := range deckslice {
		duel := strings.TrimSpace(strings.ReplaceAll(deckslice[0], "===", ""))
		if strings.Contains(line, ";") {
			fmt.Println(line, duel)
			entries := strings.Split(line, ";")
			cardName := strings.TrimSpace(entries[0])
			probability, err := strconv.Atoi(strings.TrimSpace(entries[1]))
			if err != nil {
				panic(err)
			}
			fmt.Println(cardName, probability)
		}
	}
}
