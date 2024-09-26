package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const storageDirectory = "fmrfastforward"
const dbPath = storageDirectory + "/sql.db"

func generateApiUrl(pagetoFetch string) string {
	return fmt.Sprintf("https://yugipedia.com/api.php?action=query&prop=revisions&titles=%v&rvprop=content&format=json", pagetoFetch)
}

func getFmrData() {
	if !directoryExists(storageDirectory) {
		err := os.Mkdir(storageDirectory, 0777)
		if err != nil {
			fmt.Println(err)
		}
		assert(directoryExists(storageDirectory))
	}
	if !fileExists(storageDirectory + "/characters.json") {
		getFmrCharacters()
		assert(fileExists(storageDirectory + "/characters.json"))
	}
	if !fileExists(dbPath) {
		parseCharacterList()
	}
	//Check for the characters table in the db and create it if it is not there
}

func getFmrCharacters() {
	path := storageDirectory + "/characters.json"
	assert(!fileExists(path))
	fmrCharactersUrl := generateApiUrl("Portal:Yu-Gi-Oh!_Forbidden_Memories_characters")
	req, err := http.NewRequest("GET", fmrCharactersUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Couldn't generate request to get character list, error:", err)
		return
	}
	godotenv.Load()
	email := os.Getenv("EMAIL")
	userAgent := fmt.Sprint("speedrun bot, email: ", email)
	req.Header.Add("User-Agent", userAgent)
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
	assert(fileExists(path))
}

func parseCharacterList() []string {
	assert(!fileExists(dbPath))
	assert(fileExists(storageDirectory + "/characters.json"))
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
	fmt.Println(characters)
	return characters
}
