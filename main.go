package main

import (
	"fmt"
)

func main() {
	fmt.Println("Christ is King!")
	if !fileExists(dbPath) {
		fmt.Println("Installing game data...")
		install(dbPath)
		assert(fileExists(dbPath), "install failed, shutting down...")
	}

}
