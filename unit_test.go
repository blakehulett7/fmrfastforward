package main

import (
	"testing"
)

func TestInstaller(t *testing.T) {
	chars := parseCharacterList()
	getCharacterData(chars)
}
