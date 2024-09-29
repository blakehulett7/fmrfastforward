package main

import (
	"fmt"
	"os"
	"testing"
)

func TestInstaller(t *testing.T) {
	defer os.Remove(storageDirectory + "/database.db")
	fmt.Println(tableExists("probabilities"))
	fmt.Println(tableExists("cards"))
	fmt.Println(tableExists("fusions"))
	getFmrData()
}
