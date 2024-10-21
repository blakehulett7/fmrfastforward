package main

import (
	"fmt"
)

func main() {
	fmt.Println("Christ is King!")
	if !fileExists(dbPath) {
		fmt.Println("Installing game data...")
		install()
	}
}
