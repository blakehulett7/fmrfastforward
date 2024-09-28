package main

import (
	"fmt"
	"testing"
)

func TestInstaller(t *testing.T) {
	fmt.Println(tableExists("probabilities"))
	fmt.Println(tableExists("cards"))
	fmt.Println(tableExists("fusions"))
	buildProbabilitiesTable()
}
