package main

import (
	"fmt"
	"testing"
)

func TestInstaller(t *testing.T) {
	fmt.Println(tableExists("cardTables"))
}
