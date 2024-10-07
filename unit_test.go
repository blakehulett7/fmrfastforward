package main

import (
	"fmt"
	"os"
	"testing"
)

func TestInstaller(t *testing.T) {
	defer os.Remove(storageDirectory + "/database.db")
	getFmrData()
	good := []string{"Decks", "Bandit Keith", "No duel master k?"}
	fmt.Println(good)
}
