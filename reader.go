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
