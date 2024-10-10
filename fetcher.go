package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

func generateApiUrl(pagetoFetch string) string {
	return fmt.Sprintf("https://yugipedia.com/api.php?action=query&prop=revisions&titles=%v&rvprop=content&format=json", pagetoFetch)
}

func fetch_data(fetch_string, output_path string) {
	path := storageDirectory + output_path
	assert(!fileExists(path), "should not call this function, data is already written locally")
	url := generateApiUrl(fetch_string)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("")))
	if err != nil {
		fmt.Println("Couldn't generate request to fetch data, error:", err)
		return
	}
	req.Header.Add("User-Agent", apiHeader)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Couldn't get a response from the yugipedia api, error:", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode > 499 {
		fmt.Println("Recieved an error code from the server", res.Status)
		return
	}
	resData, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Couldn't read json response from yugipedia api, error:", err)
		return
	}
	os.WriteFile(path, resData, 0777)
	assert(fileExists(path), "data was not written properly")
}
