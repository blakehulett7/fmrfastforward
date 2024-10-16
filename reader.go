package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func read_character_data(path string) CharactersQuery {
	data, err := os.ReadFile(storageDirectory + path)
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

func read_cards_data(directory string) map[string]Page {
	cards_files, err := os.ReadDir(storageDirectory + directory)
	if err != nil {
		panic(err)
	}
	cards_wikimap := map[string]Page{}
	for _, card_file := range cards_files {
		path := "/cards/" + card_file.Name()
		file_map := read_character_data(path)
		for key, value := range file_map.Query.Pages {
			cards_wikimap[key] = value
		}
	}
	return cards_wikimap
}
