package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

func runSql(sqlQuery string) {
	os.WriteFile("query.sql", []byte(sqlQuery), 0777)
	defer os.Remove("query.sql")
	command := "cat query.sql | sqlite3 fmrfastforward/database.db"
	exec.Command("bash", "-c", command).Run()
}

func outputSql(sqlQuery string) ([]byte, error) {
	os.WriteFile("query.sql", []byte(sqlQuery), 0777)
	defer os.Remove("query.sql")
	command := "cat query.sql | sqlite3 fmrfastforward/database.db"
	return exec.Command("bash", "-c", command).Output()
}

func tableExists(tableName string) bool {
	sqlQuery := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='%v';", tableName)
	data, _ := outputSql(sqlQuery)
	if strings.ReplaceAll(string(data), "\n", "") == "0" {
		return false
	}
	return true
}

func tableIsEmpty(tableName string) bool {
	sqlQuery := fmt.Sprintf("SELECT count(*) FROM %v;", tableName)
	data, _ := outputSql(sqlQuery)
	if strings.ReplaceAll(string(data), "\n", "") != "0" {
		return false
	}
	return true
}

func cardExists(cardName string) bool {
	sqlQuery := fmt.Sprintf("SELECT count(*) FROM cards WHERE name = '%v';", cardName)
	data, _ := outputSql(sqlQuery)
	if strings.ReplaceAll(string(data), "\n", "") == "0" {
		return false
	}
	return true
}

func initializeDB() {
	assert(!fileExists(storageDirectory + "/database.db"))
	sqlQuery := `
CREATE TABLE probabilities (
    id TEXT PRIMARY KEY,
    duel TEXT,
    card_id TEXT,
    deck INTEGER,
    sapow INTEGER,
    satec INTEGER,
    bcd INTEGER
    );
`
	runSql(sqlQuery)
	assert(tableExists("probabilities"))
}

func initializeCardsDB() {
	assert(!tableExists("cards"))
	sqlQuery := `
CREATE TABLE cards (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE,
    atk INTEGER,
    def INTEGER,
    );
    `
	runSql(sqlQuery)
	assert(tableExists("cards"))
}

func initializeFusionsDB() {
	assert(!tableExists("fusions"))
	sqlQuery := `
CREATE TABLE fusions (
    id TEXT PRIMARY KEY,
    card_id TEXT,
    used_in TEXT,
    material_1 TEXT,
    material_2 TEXT,
    FOREIGN KEY(card_id, used_in, material_1, material_2) REFERENCES cards(id, id, id, id));`
	runSql(sqlQuery)
	assert(tableExists("fusions"))
}

func initializeCard(cardName string) {
	assert(!cardExists(cardName))
	id := uuid.NewString()
	sqlQuery := fmt.Sprintf("INSERT INTO cards(id, name) VALUES ('%v', '%v');", id, cardName)
	runSql(sqlQuery)
	assert(cardExists(cardName))
}

func writeDuelTables(duelTables []DuelTable) {
	for _, duelTable := range duelTables {
		card := duelTable[1]
		if !cardExists(card) {
			initializeCard(card)
		}
		fmt.Println()
	}
}
