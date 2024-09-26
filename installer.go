package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const storageDirectory = "fmrfastforward"

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
	//Check for the characters table in the db and create it if it is not there
}

func getFmrCharacters() {
	path := storageDirectory + "/characters.json"
	assert(!fileExists(path))
	fmrCharactersUrl := generateApiUrl("Portal:Yu-Gi-Oh!_Forbidden_Memories_characters")
	req, err := http.NewRequest("GET", fmrCharactersUrl, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Couldn't generate request to get character list, error:", err)
		bufio.NewScanner(os.Stdin).Scan()
		return
	}
	godotenv.Load()
	email := os.Getenv("EMAIL")
	userAgent := fmt.Sprint("speedrun bot, email: ", email)
	req.Header.Add("User-Agent", userAgent)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Couldn't get a response from the yugipedia api, error:", err)
		bufio.NewScanner(os.Stdin).Scan()
		return
	}
	defer res.Body.Close()
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't read json response from the character list page, error:", err)
		bufio.NewScanner(os.Stdin).Scan()
		return
	}
	os.WriteFile(path, resData, 0777)
	assert(fileExists(path))
}
