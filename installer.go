package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const storageDirectory = "fmrfastforward"
const dbPath = storageDirectory + "/sql.db"
const apiHeader = "speedrun bot, email: blake.hulett7@gmail.com"

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
	assert(tableIsEmpty("probabilities"), "there is old data in the probablities table")
	buildProbabilitiesTable()
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

func buildProbabilitiesTable() {
	assert(fileExists(storageDirectory+"/characterdata.json"), "need to fetch character data first...")
	assert(tableIsEmpty("probabilities"), "probablities table should be empty before this buildProbabilitiesTable function is called")
	charactersQuery := read_character_data()
	wikimap := charactersQuery.Query.Pages
	deck_entries, sapow_entries, satec_entries, bcd_entries := parse_wikitext(wikimap)
	fmt.Println(deck_entries, sapow_entries, satec_entries, bcd_entries)
	/*
		for key, value := range themap {
			if len(value.Revisions) == 0 {
				fmt.Println(strings.ToUpper(value.Title))
				continue
			}
			fmt.Println(strings.TrimSpace(
				strings.ReplaceAll(
					value.Title, "(FMR)", "",
				),
			),
				value.Revisions[0].Body,
			)
			fmt.Println(key)
		}
	*/
}
