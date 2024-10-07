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
	data, err := outputSql(sqlQuery)
	if err != nil {
		panic("Something went wrong checking a table's existence")
	}
	if strings.ReplaceAll(string(data), "\n", "") == "0" {
		return false
	}
	return true
}

func tableIsEmpty(tableName string) bool {
	sqlQuery := fmt.Sprintf("SELECT count(*) FROM %v;", tableName)
	data, err := outputSql(sqlQuery)
	if err != nil {
		panic("Something went wrong checking if table is empty")
	}
	if strings.ReplaceAll(string(data), "\n", "") != "0" {
		return false
	}
	return true
}

func cardExists(cardName string) bool {
	sqlQuery := fmt.Sprintf("SELECT count(*) FROM cards WHERE name = '%v';", cardName)
	data, err := outputSql(sqlQuery)
	if err != nil {
		panic("Something went wrong checking a card's existence")
	}
	if strings.ReplaceAll(string(data), "\n", "") == "0" {
		return false
	}
	return true
}

func probabilityExists(duel, cardId string) bool {
	sqlQuery := fmt.Sprintf("SELECT count(*) FROM probabilities WHERE duel = '%v' AND card_id = '%v';", duel, cardId)
	data, err := outputSql(sqlQuery)
	if err != nil {
		panic("Something went wrong checking if a probability has been written")
	}
	if strings.ReplaceAll(string(data), "\n", "") == "0" {
		return false
	}
	return true
}

func getCardId(cardName string) string {
	assert(cardExists(cardName), "can't get an id for a card not yet in the db")
	sqlQuery := fmt.Sprintf("SELECT id FROM cards WHERE name = '%v';", cardName)
	data, err := outputSql(sqlQuery)
	if err != nil {
		message := fmt.Sprintf("Something went wrong getting %v's card id", cardName)
		panic(message)
	}
	return strings.ReplaceAll(string(data), "\n", "")
}

func getProbabilityId(duel, cardId string) string {
	assert(probabilityExists(duel, cardId), "can't get an id for a probability not yet in the db")
	sqlQuery := fmt.Sprintf("SELECT id FROM probabilities WHERE duel = '%v' AND card_id = '%v';", duel, cardId)
	data, err := outputSql(sqlQuery)
	if err != nil {
		panic("Something went wrong getting this probability id")
	}
	return strings.ReplaceAll(string(data), "\n", "")
}

func initializeDB() {
	assert(!fileExists(storageDirectory+"/database.db"), "db file already exists...")
	sqlQuery := `
CREATE TABLE probabilities (
    id TEXT PRIMARY KEY,
    duel TEXT,
    card TEXT,
    deck INTEGER,
    sapow INTEGER,
    satec INTEGER,
    bcd INTEGER
    );
`
	runSql(sqlQuery)
	assert(tableExists("probabilities"), "db file not initialized properly")
}

func initializeCardsDB() {
	assert(!tableExists("cards"), "cards table already exists, shouldn't be calling this")
	sqlQuery := `
CREATE TABLE cards (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE,
    atk INTEGER,
    def INTEGER
    );
    `
	runSql(sqlQuery)
	assert(tableExists("cards"), "cards table failed to initialize properly")
}

func initializeFusionsDB() {
	assert(!tableExists("fusions"), "fusions table already exists...")
	sqlQuery := `
CREATE TABLE fusions (
    id TEXT PRIMARY KEY,
    card_id TEXT,
    used_in TEXT,
    material_1 TEXT,
    material_2 TEXT,
    FOREIGN KEY(card_id, used_in, material_1, material_2) REFERENCES cards(id, id, id, id));`
	runSql(sqlQuery)
	assert(tableExists("fusions"), "fusions table failed to properly intialize")
}

func initialize_rate_table(table_name string) {
	assert(!tableExists(table_name), fmt.Sprintf("%v table already exists, should not call this function", table_name))
	sql_query := fmt.Sprintf("CREATE TABLE %v (id TEXT PRIMARY KEY, duel TEXT, card TEXT, rate INTEGER);", table_name)
	runSql(sql_query)
	assert(tableExists(table_name), fmt.Sprintf("failed to initialize the %v table", table_name))
}

func initializeCard(cardName string) {
	assert(!cardExists(cardName), "card already present in the db")
	id := uuid.NewString()
	sqlQuery := fmt.Sprintf("INSERT INTO cards(id, name) VALUES ('%v', '%v');", id, cardName)
	runSql(sqlQuery)
	assert(cardExists(cardName), "card was not saved to the db properly")
}

func initializeProbability(duel, cardId string) {
	assert(!probabilityExists(duel, cardId), "probability already present in the db")
	id := uuid.NewString()
	sqlQuery := fmt.Sprintf("INSERT INTO probabilities(id, duel, card_id) VALUES('%v', '%v', '%v');", id, duel, cardId)
	runSql(sqlQuery)
	assert(probabilityExists(duel, cardId), "probability was not saved to the db properly")
}

func WriteProbabilities(entries []Probability) {
	values_string := ""
	for _, entry := range entries {
		entry_string := fmt.Sprintf(", ('%v', '%v', '%v', %v)", entry.Id, entry.Duel, entry.Card, entry.Rate)
		values_string += entry_string
	}
	fmt.Println(values_string)
}
