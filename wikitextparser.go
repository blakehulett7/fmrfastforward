package main

import (
	"fmt"
	"strings"
)

func splitWikitext(wikitext string) (deckSlice, dropSlice []string) {
	//assert something about the wikitext
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

func splitWikiSlice(wikiSlice []string) [][]string {
	//asert something about the wikiSlice
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
	fmt.Println(indices)
	for idx := range indices {
		if idx == len(indices)-1 {
			fmt.Println(idx, indices[idx], "stop")
			fmt.Println(wikiSlice[indices[idx]:])
			break
		}
		fmt.Println(idx, indices[idx])
	}
	return wikiSlices
}
