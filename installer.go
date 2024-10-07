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
const dbPath = storageDirectory + "/sql.db"
const apiHeader = "speedrun bot, email: blake.hulett7@gmail.com"
const known_deck_table_length = 3649
const known_sapow_table_length = 3066
const known_satec_table_length = 2917
const known_bcd_table_length = 2683

func generateApiUrl(pagetoFetch string) string {
	return fmt.Sprintf("https://yugipedia.com/api.php?action=query&prop=revisions&titles=%v&rvprop=content&format=json", pagetoFetch)
}

func getFmrData() {
	if !directoryExists(storageDirectory) {
		os.Mkdir(storageDirectory, 0777)
		assert(directoryExists(storageDirectory), "storage directory not properly initialized")
	}
	if !fileExists(storageDirectory + "/characters.json") {
		getFmrCharacters()
		assert(fileExists(storageDirectory+"/characters.json"), "character list not retrieved")
	}
	if !fileExists(storageDirectory + "/characterdata.json") {
		charactersToFetch := parseCharacterList()
		assert(len(charactersToFetch) == 42, "there should be 42 characters here")
		getCharacterData(charactersToFetch)
		assert(fileExists(storageDirectory+"/characterdata.json"), "failed to retrieve data for each character")
	}
	if !fileExists(storageDirectory + "/database.db") {
		initializeDB()
		assert(fileExists(storageDirectory+"/database.db"), "database file not present")
		assert(tableExists("probabilities"), "failed to initialize probabilities table")
	}
	if !tableExists("cards") {
		initializeCardsDB()
		assert(tableExists("cards"), "failed to initialize cards table")
	}
	if !tableExists("fusions") {
		initializeFusionsDB()
		assert(tableExists("fusions"), "failed to initialize fusions table")
	}
	if !tableExists("decks") {
		initialize_rate_table("decks")
		assert(tableExists("decks"), "failed to initialize decks table")
	}
	if !tableExists("sapow") {
		initialize_rate_table("sapow")
		assert(tableExists("sapow"), "failed to initialize sapow table")
	}
	if !tableExists("satec") {
		initialize_rate_table("satec")
		assert(tableExists("satec"), "failed to initialize satec table")
	}
	if !tableExists("bcd") {
		initialize_rate_table("bcd")
		assert(tableExists("bcd"), "failed to initialize bcd table")
	}
	assert(tableIsEmpty("probabilities"), "there is old data in the probablities table")
	charactersQuery := read_character_data()
	wikimap := charactersQuery.Query.Pages
	deck_entries, sapow_entries, satec_entries, bcd_entries := parse_wikitext(wikimap)
	assert(len(deck_entries) == known_deck_table_length, "incorrect number of deck entries, most likely missing cards...")
	assert(len(sapow_entries) == known_sapow_table_length, "incorrect number of sapow entries, most likely missing cards...")
	assert(len(satec_entries) == known_satec_table_length, "incorrect number of satec entries, most likely missing cards...")
	assert(len(bcd_entries) == known_bcd_table_length, "incorrect number of bcd entries, most likely missing cards...")
	assert(tableIsEmpty("decks"), "there is old data in the decks table")
	WriteProbabilities(deck_entries, "decks")
	assert(table_has_length("decks", known_deck_table_length), "deck table incorrectly written, we are missing cards most likely...")
	WriteProbabilities(sapow_entries, "sapow")
	assert(table_has_length("sapow", known_sapow_table_length), "sapow table incorrectly written, we are missing cards most likely...")
	assert(tableIsEmpty("satec"), "there is old data in the satec table")
	WriteProbabilities(satec_entries, "satec")
	assert(table_has_length("satec", known_satec_table_length), "satec table incorrectly written, we are missing cards most likely...")
	assert(tableIsEmpty("bcd"), "there is old data in the bcd table")
	WriteProbabilities(bcd_entries, "bcd")
	assert(table_has_length("bcd", known_bcd_table_length), "bcd table incorrectly written, we are missing cards most likely...")
	cards_to_fetch := []string{}
	for _, entry := range deck_entries {
		if !slices.Contains(cards_to_fetch, entry.Card) {
			fmt.Println(entry.Card)
			cards_to_fetch = append(cards_to_fetch, entry.Card)
		}
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
	os.WriteFile(path, resData, 0777)
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
