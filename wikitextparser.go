package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func splitWikitext(wikitext string) (deckSlice, dropSlice WikiSection) {
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

func splitByDuels(wikiSection WikiSection) []WikiSection {
	assert(len(wikiSection) != 0)
	indices := []int{}
	for idx, line := range wikiSection {
		if !strings.HasPrefix(line, "===") {
			continue
		}
		indices = append(indices, idx)
	}
	if len(indices) == 0 {
		return []WikiSection{wikiSection}
	}
	wikiSlices := []WikiSection{}
	for idx := range indices {
		if idx == len(indices)-1 {
			wikiSlices = append(wikiSlices, wikiSection[indices[idx]:])
			break
		}
		wikiSlices = append(wikiSlices, wikiSection[indices[idx]:indices[idx+1]])
	}
	assert(len(wikiSlices) != 0)
	return wikiSlices
}

func getDuelTables(deckslice []string) DuelTable {
	assert(len(deckslice) != 0)
	duelTable := DuelTable{}
	for _, line := range deckslice {
		duel := strings.TrimSpace(strings.ReplaceAll(deckslice[0], "===", ""))
		if strings.Contains(line, ";") {
			entryValues := strings.Split(line, ";")
			cardName := strings.TrimSpace(entryValues[0])
			probability := strings.TrimSpace(entryValues[1])
			duelTable = append(duelTable, [3]string{duel, cardName, probability})
		}
	}
	return duelTable
}

func parseDuelTable(duelTableEntry DuelTableEntry) (duel, cardId string, probability int) {
	for _, value := range duelTableEntry {
		assert(value != "")
	}
	card := duelTableEntry[1]
	if !cardExists(card) {
		initializeCard(card)
	}
	duel = duelTableEntry[0]
	cardId = getCardId(card)
	probability, err := strconv.Atoi(duelTableEntry[2])
	if err != nil {
		panic(err)
	}
	assert(duel != "")
	assert(cardId != "")
	assert(probability != 0)
	return duel, cardId, probability
}
