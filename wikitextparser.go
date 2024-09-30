package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func parse_wikitext(wikitext string) (deck_entries []Probability) {
	assert(wikitext != "")
	decksection, dropsection := splitWikitext(wikitext)
	assert(len(decksection) != 0)
	assert(len(dropsection) != 0)
	deckText := splitByDuels(decksection)
	assert(len(deckText) != 0)
	deck_entries = parse_deck_text(deckText)
	drop_text := splitByDuels(dropsection)
	assert(len(drop_text) == len(deckText))
	fmt.Println(parse_drop_text(drop_text))
	return deck_entries
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
		deckSlice := wikitextslice[deckIdx:dropIdx] // This
		dropSlice := wikitextslice[dropIdx:]        // Doesn't
		return deckSlice, dropSlice                 // Look
	}
	panic("Should never get here, something went wrong parsing the wikitext!") // Right?
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

func split_by_table(wikiSection WikiSection) (sapow_text, satec_text, bcd_text WikiSection) {
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
	for _, section := range sections {
		if strings.HasPrefix(section[0], "| pow") {
			sapow_text = section
			continue
		}
		if strings.HasPrefix(section[0], "| tec") {
			satec_text = section
			continue
		}
		if strings.HasPrefix(section[0], "| bcd") {
			bcd_text = section
			continue
		}
		panic("We should never get here, something went wrong seperating drop text by table")
	}
	return sapow_text, satec_text, bcd_text
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

func parse_entry_text(line string) (string, int) {
	assert(strings.Contains(line, ";"))
	values := strings.Split(line, ";")
	rate, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		panic("Couldn't convert rate to an int type, something is wrong")
	}
	return strings.TrimSpace(values[0]), rate
}

func parse_deck_text(deck_text_by_duel []WikiSection) []Probability {
	entries := []Probability{}
	for _, duel_text := range deck_text_by_duel {
		duel := strings.TrimSpace(strings.ReplaceAll(duel_text[0], "===", ""))
		for _, line := range duel_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			id := uuid.NewString()
			entry := Probability{
				Id:   id,
				Duel: duel,
				Card: card,
				Rate: rate,
			}
			entries = append(entries, entry)
		}
	}
	return entries
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

func parse_drop_text(drop_text []WikiSection) (sapow_entries, satec_entries, bcd_entries []Probability) {
	assert(len(drop_text) != 0)
	sapow_entries = []Probability{}
	satec_entries = []Probability{}
	bcd_entries = []Probability{}
	for _, duel_text := range drop_text {
		duel := strings.TrimSpace(strings.ReplaceAll(duel_text[0], "===", ""))
		sapow_text, satec_text, bcd_text := split_by_table(duel_text)
		for _, line := range sapow_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			entry := Probability{
				Id:   uuid.NewString(),
				Duel: duel,
				Card: card,
				Rate: rate,
			}
			sapow_entries = append(sapow_entries, entry)
		}
		for _, line := range satec_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			entry := Probability{
				Id:   uuid.NewString(),
				Duel: duel,
				Card: card,
				Rate: rate,
			}
			satec_entries = append(satec_entries, entry)
		}
		for _, line := range bcd_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			entry := Probability{
				Id:   uuid.NewString(),
				Duel: duel,
				Card: card,
				Rate: rate,
			}
			bcd_entries = append(bcd_entries, entry)
		}
	}
	return sapow_entries, satec_entries, bcd_entries
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
