package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
)

const storageDirectory = "fmrfastforward"
const dbPath = storageDirectory + "/database.db"
const apiHeader = "speedrun bot, email: blake.hulett7@gmail.com"
const known_character_length = 42
const known_deck_table_length = 3649
const known_sapow_table_length = 3066
const known_satec_table_length = 2917
const known_bcd_table_length = 2683
const max_permissions = 0777

func getFmrData() { //TODO: function is too long, need to break this up
	if !directoryExists(storageDirectory) {
		os.Mkdir(storageDirectory, max_permissions)
		assert(directoryExists(storageDirectory), "storage directory not properly initialized")
	}
	if !fileExists(storageDirectory + "/characters.json") {
		getFmrCharacters()
		assert(fileExists(storageDirectory+"/characters.json"), "character list not retrieved")
	}
	if !fileExists(storageDirectory + "/characterdata.json") {
		charactersToFetch := parseCharacterList()
		assert(len(charactersToFetch) == known_character_length, "there should be 42 characters here")
		getCharacterData(charactersToFetch)
		assert(fileExists(storageDirectory+"/characterdata.json"), "failed to retrieve data for each character")
	}
	if !fileExists(dbPath) {
		charactersQuery := read_character_data("/characterdata.json")
		wikimap := charactersQuery.Query.Pages
		deck_entries, sapow_entries, satec_entries, bcd_entries := parse_wikitext(wikimap)

		assert(len(deck_entries) == known_deck_table_length, "incorrect number of deck entries, most likely missing cards...")
		assert(len(sapow_entries) == known_sapow_table_length, "incorrect number of sapow entries, most likely missing cards...")
		assert(len(satec_entries) == known_satec_table_length, "incorrect number of satec entries, most likely missing cards...")
		assert(len(bcd_entries) == known_bcd_table_length, "incorrect number of bcd entries, most likely missing cards...")

		initialize_rate_table("decks")
		WriteProbabilities(deck_entries, "decks")
		assert(table_has_length("decks", known_deck_table_length), "deck table incorrectly written, we are missing cards most likely...")

		initialize_rate_table("sapow")
		WriteProbabilities(sapow_entries, "sapow")
		assert(table_has_length("sapow", known_sapow_table_length), "sapow table incorrectly written, we are missing cards most likely...")

		initialize_rate_table("satec")
		WriteProbabilities(satec_entries, "satec")
		assert(table_has_length("satec", known_satec_table_length), "satec table incorrectly written, we are missing cards most likely...")

		initialize_rate_table("bcd")
		WriteProbabilities(bcd_entries, "bcd")
		assert(table_has_length("bcd", known_bcd_table_length), "bcd table incorrectly written, we are missing cards most likely...")

		if !directoryExists(storageDirectory + "/cards") {
			cards_to_fetch := generate_cards_fetch_list([][]Probability{deck_entries, sapow_entries, satec_entries, bcd_entries})
			get_cards_data(cards_to_fetch)
		}

		assert(!tableExists("cards"), "cards table should not exist yet")
	}

	if !tableExists("cards") { //Will likely change this to an assert
		initializeCardsDB()
		assert(tableExists("cards"), "failed to initialize cards table")
	}
	if !tableExists("fusions") { //Will likely change this to an assert
		initializeFusionsDB()
		assert(tableExists("fusions"), "failed to initialize fusions table")
	}
}

func getFmrCharacters() {
	path := storageDirectory + "/characters.json"
	assert(!fileExists(path), "should not call this function with character list data already present")
	fmrCharactersUrl := generateApiUrl("Portal:Yu-Gi-Oh!_Forbidden_Memories_characters")
	req, err := http.NewRequest("GET", fmrCharactersUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Couldn't generate request to get character list, error:", err)
		return
	}
	req.Header.Add("User-Agent", apiHeader)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Couldn't get a response from the yugipedia api, error:", err)
		return
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't read json response from the character list page, error:", err)
		return
	}
	os.WriteFile(path, resData, max_permissions)
	assert(fileExists(path), "character list data was not written properly")
}

func parseCharacterList() []string {
	assert(fileExists(storageDirectory+"/characters.json"), "need to fetch the character list data first")
	data, err := os.ReadFile(storageDirectory + "/characters.json")
	if err != nil {
		fmt.Println("Couldn't read characters json stored on disk, error:", err)
		return []string{}
	}
	dataStruct := CharactersPageJSON{}
	err = json.Unmarshal(data, &dataStruct)
	if err != nil {
		fmt.Println("Couldn't decode characters json data stored on disk, error:", err)
		return []string{}
	}
	charactersTableXwikiString := dataStruct.Query.Pages.Num369496.Revisions[0].Body
	charactersTableRows := strings.Split(charactersTableXwikiString, "\n")[1 : len(strings.Split(charactersTableXwikiString, "\n"))-3]
	characters := []string{}
	for _, row := range charactersTableRows {
		row := strings.TrimSpace(
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.Split(row, "|")[1], "[[", "",
				), "]]", "",
			),
		)
		characters = append(characters, row)
	}
	assert(len(characters) == 42, "Expecting 42 characters")
	return characters
}

func getCharacterData(fetchList []string) {
	assert(len(fetchList) == 42, "Expecting 42 characters to fetch")
	titles := ""
	for _, character := range fetchList {
		character = strings.ReplaceAll(character, " ", "_")
		titles = titles + "|" + character
	}
	url := generateApiUrl(titles)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Couldn't generate get char data request, error:", err)
		return
	}
	req.Header.Add("User-Agent", apiHeader)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Couldn't execute get character data request, error:", err)
		return
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't write the get character data response json, error:", err)
	}
	os.WriteFile(storageDirectory+"/characterdata.json", resData, 0777)
	assert(fileExists(storageDirectory+"/characterdata.json"), "character data was not written properly")
}

func generate_cards_fetch_list(entries_array [][]Probability) [][]string {
	fetch_slice := []string{}
	for _, entries := range entries_array {
		for _, entry := range entries {
			if !slices.Contains(fetch_slice, entry.Card) {
				fetch_slice = append(fetch_slice, entry.Card)
			}
		}
	}
	// len(fetch_slice) = 716
	batch_1 := fetch_slice[:49]
	batch_2 := fetch_slice[50:99]
	batch_3 := fetch_slice[100:149]
	batch_4 := fetch_slice[150:199]
	batch_5 := fetch_slice[200:249]
	batch_6 := fetch_slice[250:299]
	batch_7 := fetch_slice[300:349]
	batch_8 := fetch_slice[350:399]
	batch_9 := fetch_slice[400:449]
	batch_10 := fetch_slice[450:499]
	batch_11 := fetch_slice[500:549]
	batch_12 := fetch_slice[550:599]
	batch_13 := fetch_slice[600:649]
	batch_14 := fetch_slice[650:699]
	batch_15 := fetch_slice[700:]
	return [][]string{batch_1, batch_2, batch_3, batch_4, batch_5, batch_6, batch_7, batch_8, batch_9, batch_10, batch_11, batch_12, batch_13, batch_14, batch_15}
}

func get_cards_data(cards_to_fetch [][]string) {
	assert(!directoryExists(storageDirectory+"/cards"), "cards directory should not exist yet")
	os.Mkdir(storageDirectory+"/cards", 0777)
	for idx, fetch_list := range cards_to_fetch {
		cards_string := ""
		for _, card := range fetch_list {
			cards_string = cards_string + "|" + strings.ReplaceAll(card, " ", "_")
		}
		cards_string = cards_string[1:]
		path := fmt.Sprintf("/cards/cards%v.json", idx+1)
		fetch_data(cards_string, path)
		assert(fileExists(storageDirectory+path), "cards data was not written by the fetch_data function")
	}
}
