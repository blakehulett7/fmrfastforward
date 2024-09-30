package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func parse_wikitext(wikitext string) {
	assert(wikitext != "")
	decksection, dropsection := splitWikitext(wikitext)
	assert(len(decksection) != 0)
	assert(len(dropsection) != 0)
	deckTextByDuel := splitByDuels(decksection)
	//_, drop_text_by_duel := splitByDuels(dropsection)
	assert(len(deckTextByDuel) != 0)
	//assert(len(drop_text_by_duel) != 0)
	entries := []Probability{}
	for _, duelText := range deckTextByDuel {
		duelTable := getDuelTable(duelText)
		entries = append(entries, parse_deck_table(duelTable)...)
	}
	//parse_drop_text(drop_text_by_duel)
}

func read_character_data() CharactersQuery {
	data, err := os.ReadFile(storageDirectory + "/characterdata.json")
	if err != nil {
		fmt.Println("Couldn't load the character json data from disk, error:", err)
		panic(err)
	}
	charactersQuery := CharactersQuery{}
	err = json.Unmarshal(data, &charactersQuery)
	if err != nil {
		fmt.Println("Couldn't decode the character json data on the disk, error:", err)
		panic(err)
	}
	return charactersQuery
}

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
		panic("Couldn't split by duels")
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

func split_by_table(wikiSection WikiSection) []WikiSection {
	assert(len(wikiSection) != 0)
	indices := []int{}
	for idx, line := range wikiSection {
		if !strings.HasPrefix(line, "|") {
			continue
		}
		if strings.HasPrefix(line, "| n") {
			continue
		}
		indices = append(indices, idx)
	}
	assert(len(indices) == 3) // These are the 3 possible drop table sections, has to be 3
	sections := []WikiSection{}
	for idx := range indices {
		if idx == len(indices)-1 {
			sections = append(sections, wikiSection[indices[idx]:])
			break
		}
		sections = append(sections, wikiSection[indices[idx]:indices[idx+1]])
	}
	assert(len(sections) == 3)
	return sections
}

func getDuelTable(deckslice []string) DuelTable {
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

func parseDuelTableEntry(duelTableEntry DuelTableEntry) (duel, cardName string, probability int) {
	for _, value := range duelTableEntry {
		assert(value != "")
	}
	duel = duelTableEntry[0]
	cardName = duelTableEntry[1]
	probability, err := strconv.Atoi(duelTableEntry[2])
	if err != nil {
		panic(err)
	}
	assert(duel != "")
	assert(cardName != "")
	assert(probability != 0)
	return duel, cardName, probability
}

func parse_deck_table(duelTable DuelTable) []Probability {
	assert(len(duelTable) != 0)
	probabilities := []Probability{}
	for _, entry := range duelTable {
		duel, cardName, deck := parseDuelTableEntry(entry)
		id := uuid.NewString()
		probabilities = append(probabilities, Probability{
			Id:   id,
			Duel: duel,
			Card: cardName,
			Deck: deck,
		})
	}
	return probabilities
}

func parse_drop_table(drop_table WikiSection) [][2]string {
	drop_table_entries := [][2]string{}
	for _, line := range drop_table[1:] {
		if !strings.Contains(line, ";") {
			continue
		}
		values := strings.Split(line, ";")
		assert(len(values) == 2)
		drop_table_entries = append(drop_table_entries, [2]string{strings.TrimSpace(values[0]), strings.TrimSpace(values[1])})
	}
	return drop_table_entries
}

func parse_drop_text(drop_text []WikiSection) []Probability {
	assert(len(drop_text) != 0)
	entries := []Probability{}
	for _, duel_text := range drop_text {
		fmt.Println(duel_text)
		duel := strings.TrimSpace(strings.ReplaceAll(duel_text[0], "===", ""))
		duel_text_by_table := split_by_table(duel_text)
		for _, drop_text := range duel_text_by_table {
			if strings.HasPrefix(drop_text[0], "| pow") {
				values := parse_drop_table(drop_text)
				fmt.Println(values)
				continue
			}
			if strings.HasPrefix(drop_text[0], "| tec") {
				fmt.Println("satec")
				continue
			}
			if strings.HasPrefix(drop_text[0], "| bcd") {
				fmt.Println("bcd")
				continue
			}
			panic("We should not get here, something is wrong parsing drop text")
		}
		fmt.Println(duel)
	}
	return entries
}

/*
func parse_sapow_table(drop_table DuelTable, duel string) {
	assert(len(drop_table) != 0)
	assert(duel != "")
	probabilities := []Probability{}
	for _, entry := range drop_table {

	}
}
*/
