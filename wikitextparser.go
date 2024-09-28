package main

import (
	"fmt"
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
	fmt.Println(strings.ReplaceAll(deckslice[0], "===", ""))
	for _, line := range deckslice {
		if strings.Contains(line, ";") {
			fmt.Println(line)
		}
	}
}
