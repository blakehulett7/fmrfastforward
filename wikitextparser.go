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
		"Sadin",
		"Servant",
		"Joey (FMR)",
		"Tea Gardner (FMR)",
		"Yugi (FMR)",
		"Prince (FMR)",
		"Bakura (FMR)",
		"",
	}
	for id, value := range wikimap {
		if strings.Contains(id, "-") {
			continue
		}
		if slices.Contains(characters_to_remove, value.Title) {
			continue
		}
		wikitext := value.Revisions[0].Body
		assert(wikitext != "", "no wikitext present to parse...")
		decksection, dropsection := splitWikitext(wikitext)
		if slices.Contains(decksection, "**") {
			fmt.Println(value.Title, "couldn't split into a deck and duel section")
			continue
		}
		assert(len(decksection) != 0, "couldn't find a deck section")
		assert(len(dropsection) != 0, "couldn't find a drops section")
		deckText := splitByDuels(decksection)
		assert(len(deckText) != 0, "didn't find any duels to parse...")
		new_deck_entries := parse_deck_text(deckText, value.Title)
		drop_text := splitByDuels(dropsection)
		assert(len(drop_text) == len(deckText), "there should be the same number of duels for the deck and drops sections")
		new_sapow_entries, new_satec_entries, new_bcd_entries := parse_drop_text(drop_text, value.Title)
		assert(len(new_sapow_entries) != 0, "didn't get sapow drop rates")
		assert(len(new_satec_entries) != 0, "didn't get satec drop rates")
		assert(len(new_bcd_entries) != 0, "didn't get bcd drop rates")
		deck_entries = append(deck_entries, new_deck_entries...)
		sapow_entries = append(sapow_entries, new_sapow_entries...)
		satec_entries = append(satec_entries, new_satec_entries...)
		bcd_entries = append(bcd_entries, new_bcd_entries...)
	}
	assert(len(deck_entries) != 0, "we are returning an empty slice of deck entries, should never do that")
	assert(len(sapow_entries) != 0, "we are returning an empty slice of sapow entries, should never do that")
	assert(len(satec_entries) != 0, "we are returning an empty slice of satec entries, should never do that")
	assert(len(bcd_entries) != 0, "we are returning an empty slice of bcd entries, should never do that")
	return deck_entries, sapow_entries, satec_entries, bcd_entries
}

func splitWikitext(wikitext string) (deckSlice, dropSlice WikiSection) {
	//assert something about the wikitext
	assert(wikitext != "", "no wikitext to parse, shouldn't call this function")
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
			assert(len(deckSlice) != 0, "didn't find a deck section")
			dropSlice := wikitextslice[dropIdx:dialogueIdx]
			assert(len(dropSlice) != 0, "didn't find a drop section")
			return deckSlice, dropSlice
		}
	}
	deckSlice = wikitextslice[deckIdx:dropIdx]
	dropSlice = wikitextslice[dropIdx:]
	dropSlice = dropSlice[:len(dropSlice)-2]
	assert(len(deckSlice) != 0, "didn't find a deck section")
	assert(len(dropSlice) != 0, "didn't find a drop section")
	return deckSlice, dropSlice
}

func splitByDuels(wikiSection WikiSection) []WikiSection {
	assert(len(wikiSection) != 0, "no section to parse")
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
	assert(len(wikiSlices) != 0, "failed to split by duels")
	return wikiSlices
}

func split_by_table(wikiSection WikiSection) (sapow_text, satec_text, bcd_text WikiSection) {
	assert(len(wikiSection) != 0, "no text to parse")
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
	assert(len(indices) == 3, "there should always be 3 drop tables") // These are the 3 possible drop table sections, has to be 3
	sections := []WikiSection{}
	for idx := range indices {
		if idx == len(indices)-1 {
			sections = append(sections, wikiSection[indices[idx]:])
			break
		}
		sections = append(sections, wikiSection[indices[idx]:indices[idx+1]])
	}
	assert(len(sections) == 3, "we got the indices right, but somehow did not end up with 3 tables")
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
	assert(len(sapow_text) != 0, "didn't successfully find the sapow text")
	assert(len(satec_text) != 0, "didn't successfully find the satec text")
	assert(len(bcd_text) != 0, "didn't successfully find the bcd text")
	return sapow_text, satec_text, bcd_text
}

func parse_entry_text(line string) (string, int) {
	assert(strings.Contains(line, ";"), "these entry texts are in an improper format, need a ; character")
	values := strings.Split(line, ";")
	rate, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		panic("Couldn't convert rate to an int type, something is wrong")
	}
	assert(values[0] != "", "failed to get the card name")
	assert(rate != 0, "failed to get the card rate")
	return strings.TrimSpace(values[0]), rate
}

func parse_deck_text(deck_text_by_duel []WikiSection, duel_title string) []Probability {
	assert(len(deck_text_by_duel) != 0, "no deck text to parse")
	entries := []Probability{}
	for _, duel_text := range deck_text_by_duel {
		duel := strings.TrimSpace(strings.ReplaceAll(duel_text[0], "===", ""))
		if strings.Contains(duel, "=") {
			duel = duel_title
		}
		assert(duel != "", "shouldn't find a blank duel section")
		for _, line := range duel_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			assert(card != "", "failed to parse the card name")
			assert(rate != 0, "failed to parse the card rate")
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
	if duel_title != "Duel Master K" {
		assert(len(entries) != 0, "no deck table data was parsed")
	}
	return entries
}

func parse_drop_text(drop_text []WikiSection, duel_title string) (sapow_entries, satec_entries, bcd_entries []Probability) {
	assert(len(drop_text) != 0, "no drop text to parse")
	sapow_entries = []Probability{}
	satec_entries = []Probability{}
	bcd_entries = []Probability{}
	for _, duel_text := range drop_text {
		duel := strings.TrimSpace(strings.ReplaceAll(duel_text[0], "===", ""))
		if strings.Contains(duel, "=") {
			duel = duel_title
		}
		assert(duel != "", "shouldn't find a blank duel section")
		sapow_text, satec_text, bcd_text := split_by_table(duel_text)
		for _, line := range sapow_text {
			if !strings.Contains(line, ";") {
				continue
			}
			card, rate := parse_entry_text(line)
			assert(card != "", "failed to parse the card name")
			assert(rate != 0, "failed to parse the card rate")
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
			assert(card != "", "failed to parse the card name")
			assert(rate != 0, "failed to parse the card rate")
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
			assert(card != "", "failed to parse the card name")
			assert(rate != 0, "failed to parse the card rate")
			entry := Probability{
				Id:   uuid.NewString(),
				Duel: duel,
				Card: card,
				Rate: rate,
			}
			bcd_entries = append(bcd_entries, entry)
		}
	}
	assert(len(sapow_entries) != 0, "failed to parse sapow table")
	assert(len(satec_entries) != 0, "failed to parse satec table")
	assert(len(bcd_entries) != 0, "failed to parse bcd table")
	return sapow_entries, satec_entries, bcd_entries
}

func parse_cards(wikimap map[string]Page) ([]Card, []Target, []Card_Star) {
	cards := []Card{}
	targets := []Target{}
	cards_stars := []Card_Star{}
	for _, card_text := range wikimap {
		if card_text.Title == "Ray" {
			continue
		}
		id := uuid.NewString()
		card := Card{Id: id, Name: card_text.Title}
		my_targets := []Target{}
		my_stars := []Card_Star{}
		lines := strings.Split(card_text.Revisions[0].Body, "\n")
		for _, line := range lines {
			if !strings.HasPrefix(line, "|") {
				if !strings.HasPrefix(line, "*") {
					continue
				}
			}
			if strings.HasPrefix(line, "| type") {
				card.Type = strings.TrimSpace(strings.Split(line, "=")[1])
			}
			if strings.HasPrefix(line, "| atk") {
				atk, err := strconv.Atoi(strings.TrimSpace(strings.Split(line, "=")[1]))
				assert(err == nil, "conversion of card's atk to an integer failed")
				card.Attack = atk
			}
			if strings.HasPrefix(line, "| def") {
				def, err := strconv.Atoi(strings.TrimSpace(strings.Split(line, "=")[1]))
				assert(err == nil, "conversion of card's def to an integer failed")
				card.Defense = def
			}
			if strings.HasPrefix(line, "| star_chips") {
				sc, err := strconv.Atoi(strings.TrimSpace(strings.Split(line, "=")[1]))
				assert(err == nil, "conversion of card's star chips to an integer failed")
				card.StarChips = sc
			}
			if strings.HasPrefix(line, "| guardian_stars") {
				stars := strings.Split(strings.ReplaceAll(strings.Split(line, "=")[1], " ", ""), "/")
				my_stars = append(my_stars, Card_Star{
					Id:   uuid.NewString(),
					Card: card_text.Title,
					Star: stars[0],
				})
				my_stars = append(my_stars, Card_Star{
					Id:   uuid.NewString(),
					Card: card_text.Title,
					Star: stars[1],
				})

			}
			if strings.HasPrefix(line, "*") {
				my_targets = append(my_targets, Target{
					Id:     uuid.NewString(),
					Equip:  card_text.Title,
					Target: strings.TrimSpace(line[1:]),
				})
			}
		}
		cards = append(cards, card)
		targets = append(targets, my_targets...)
		cards_stars = append(cards_stars, my_stars...)
	}
	return cards, targets, cards_stars
}

func parse_fusions(wikimap map[string]Page) {
	for key, value := range wikimap {
		if value.Title == "Glitch fusion" {
			continue
		}
		fmt.Println(key, value.Title)
	}
	fmt.Println(wikimap["1009580"])
	v := wikimap["1009580"].Revisions[0].Body
	v_slices := strings.Split(v, "\n")
	indices := []int{}
	for index, line := range v_slices {
		if strings.HasPrefix(line, "=") {
			indices = append(indices, index)
		}
	}
	v_slices_by_fusion := splitter(v_slices, indices)
	indices = []int{}
	for index, line := range v_slices_by_fusion[0] {
		if strings.HasPrefix(strings.ReplaceAll(line, " ", ""), "|f") {
			indices = append(indices, index)
		}
	}
	lines_by_fusion_number := splitter(v_slices_by_fusion[0], indices)
	for _, section := range lines_by_fusion_number {
		normalized := strings.ReplaceAll(section[0], " ", "")
		fusion_number_rune := normalized[2]
		fusion_number := int(fusion_number_rune - '0')
		material_number_rune := normalized[5]
		material_number := int(material_number_rune - '0')
		fmt.Println(fusion_number, material_number)
		fmt.Println()
	}
}

func splitter(string_slice []string, indices []int) [][]string {
	split_slice := [][]string{}
	for index := range indices {
		if index == len(indices)-1 {
			split_slice = append(split_slice, string_slice[indices[index]:])
			break
		}
		split_slice = append(split_slice, string_slice[indices[index]:indices[index+1]])
	}
	return split_slice
}
