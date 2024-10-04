package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func parse_wikitext(wikimap map[string]Page) (deck_entries, sapow_entries, satec_entries, bcd_entries []Probability) {
	characters_to_remove := []string{
		"Kemo (FMR)",
		"Card shop owner (FMR)",
		"Duel Master K",
	}
	for _, value := range wikimap {
		if slices.Contains(characters_to_remove, value.Title) {
			continue
		}
		wikitext := value.Revisions[0].Body
		fmt.Println(wikitext)
		assert(wikitext != "")
		decksection, dropsection := splitWikitext(wikitext)
		assert(len(decksection) != 0)
		assert(len(dropsection) != 0)
		deckText := splitByDuels(decksection)
		assert(len(deckText) != 0)
		new_deck_entries := parse_deck_text(deckText)
		drop_text := splitByDuels(dropsection)
		assert(len(drop_text) == len(deckText))
		new_sapow_entries, new_satec_entries, new_bcd_entries := parse_drop_text(drop_text)
		assert(len(sapow_entries) != 0)
		assert(len(satec_entries) != 0)
		assert(len(bcd_entries) != 0)
		deck_entries = append(deck_entries, new_deck_entries...)
		sapow_entries = append(sapow_entries, new_sapow_entries...)
		satec_entries = append(satec_entries, new_satec_entries...)
		bcd_entries = append(bcd_entries, new_bcd_entries...)
	}
	return deck_entries, sapow_entries, satec_entries, bcd_entries
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
			assert(len(deckSlice) != 0)
			dropSlice := wikitextslice[dropIdx:dialogueIdx]
			assert(len(dropSlice) != 0)
			return deckSlice, dropSlice
		}
		deckSlice := wikitextslice[deckIdx:dropIdx] // This
		dropSlice := wikitextslice[dropIdx:]        // Doesn't
		return deckSlice, dropSlice                 // Look
	}
	fmt.Println("**")
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
		panic("We should never get here, something went wrong separating drop text by table")
	}
	assert(len(sapow_text) != 0)
	assert(len(satec_text) != 0)
	assert(len(bcd_text) != 0)
	return sapow_text, satec_text, bcd_text
}

func parse_entry_text(line string) (string, int) {
	assert(strings.Contains(line, ";"))
	values := strings.Split(line, ";")
	rate, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		panic("Couldn't convert rate to an int type, something is wrong")
	}
	assert(values[0] != "")
	assert(rate != 0)
	return strings.TrimSpace(values[0]), rate
}

func parse_deck_text(deck_text_by_duel []WikiSection) []Probability {
	assert(len(deck_text_by_duel) != 0)
	entries := []Probability{}
	for _, duel_text := range deck_text_by_duel {
		duel := strings.TrimSpace(strings.ReplaceAll(duel_text[0], "===", ""))
		assert(duel != "")
		assert(!strings.Contains(duel, "="))
		for _, line := range duel_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			assert(card != "")
			assert(rate != 0)
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
	assert(len(entries) != 0)
	return entries
}

func parse_drop_text(drop_text []WikiSection) (sapow_entries, satec_entries, bcd_entries []Probability) {
	assert(len(drop_text) != 0)
	sapow_entries = []Probability{}
	satec_entries = []Probability{}
	bcd_entries = []Probability{}
	for _, duel_text := range drop_text {
		duel := strings.TrimSpace(strings.ReplaceAll(duel_text[0], "===", ""))
		assert(duel != "")
		assert(!strings.Contains(duel, "="))
		sapow_text, satec_text, bcd_text := split_by_table(duel_text)
		for _, line := range sapow_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			assert(card != "")
			assert(rate != 0)
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
			assert(card != "")
			assert(rate != 0)
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
			assert(card != "")
			assert(rate != 0)
			entry := Probability{
				Id:   uuid.NewString(),
				Duel: duel,
				Card: card,
				Rate: rate,
			}
			bcd_entries = append(bcd_entries, entry)
		}
	}
	assert(len(sapow_entries) != 0)
	assert(len(satec_entries) != 0)
	assert(len(bcd_entries) != 0)
	return sapow_entries, satec_entries, bcd_entries
}
