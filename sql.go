package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
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

func list_sql(sql_query string) []string {
	data, err := outputSql(sql_query)
	if err != nil {
		panic(fmt.Sprintf("bug running the following sql query: %v", sql_query))
	}
	data_slice := strings.Split(string(data), "\n")
	return data_slice[:len(data_slice)-1]
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

func table_has_length(table_name string, length int) bool {
	sql_query := fmt.Sprintf("SELECT count(*) FROM %v", table_name)
	data, err := outputSql(sql_query)
	if err != nil {
		panic(fmt.Sprintf("Something went wrong reading the data to check the %v table length", table_name))
	}
	got_length, err := strconv.Atoi(strings.ReplaceAll(string(data), "\n", ""))
	if err != nil {
		panic(fmt.Sprintf("Couldn't convert the string to an int to check the %v table length", table_name))
	}
	if got_length != length {
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

func get_potential_fusions(card_name string, m_table string) []string {
	assert(cardExists(card_name), fmt.Sprintf("%v is not in the cards database, shutting down...", card_name))

	sql_query := fmt.Sprintf("SELECT resulting_fusion, fusion_number FROM %v WHERE card = \"%v\";", m_table, card_name)
	return list_sql(sql_query)
}

func get_starting_deck_rates(pool_name string) []Probability {
	sql_query := fmt.Sprintf("SELECT * FROM starting_deck_rates WHERE pool = '%v' ORDER BY card;", pool_name)
	data, err := outputSql(sql_query)
	if err != nil {
		panic(err)
	}
	data_list := strings.Split(string(data), "\n")
	data_list = data_list[:len(data_list)-1]
	pool_entries := []Probability{}
	for _, entry := range data_list {
		entry_array := strings.Split(entry, "|")
		rate, err := strconv.Atoi(strings.TrimSpace(entry_array[3]))
		if err != nil {
			panic(err)
		}
		pool_entries = append(pool_entries, Probability{
			Id:   entry_array[0],
			Duel: entry_array[1],
			Card: entry_array[2],
			Rate: rate,
		})
	}
	return pool_entries
}

func get_card(card_to_get string) Card {
	sql_query := fmt.Sprintf("SELECT * FROM cards WHERE name = \"%v\";", card_to_get)
	data, err := outputSql(sql_query)
	if err != nil {
		fmt.Println("err")
		return Card{}
	}

	// Temporary for manually adding missing cards, will fix
	if string(data) == "" {
		panic(fmt.Sprintf("Manually add %v", card_to_get))
	}

	stringified_data := strings.ReplaceAll(string(data), "\n", "")
	data_slice := strings.Split(stringified_data, "|")
	int_slice := []int{}
	for _, num_string := range data_slice[3:] {
		num, err := strconv.Atoi(num_string)
		if err != nil {
			int_slice = append(int_slice, 0)
			continue
		}
		int_slice = append(int_slice, num)
	}
	card := Card{
		Id:        data_slice[0],
		Name:      data_slice[1],
		Type:      data_slice[2],
		Attack:    int_slice[0],
		Defense:   int_slice[1],
		StarChips: int_slice[2],
	}
	return card
}

func get_cards(cards_to_get []string) []Card {
	cards := []Card{}
	in_string := ""
	for _, card := range cards_to_get {
		in_string += fmt.Sprintf(", \"%v\"", card)
	}
	in_string = in_string[2:]
	sql_query := fmt.Sprintf("SELECT * FROM cards WHERE name IN (%v);", in_string)
	data, err := outputSql(sql_query)
	if err != nil {
		panic(err)
	}
	data_list := strings.Split(string(data), "\n")
	data_list = data_list[:len(data_list)-1]
	for _, entry := range data_list {
		entry_array := strings.Split(entry, "|")
		int_array := []int{}
		for _, number_string := range entry_array[3:] {
			num, err := strconv.Atoi(strings.TrimSpace(number_string))
			if err != nil {
				int_array = append(int_array, 0)
				continue
			}
			int_array = append(int_array, num)
		}
		card := Card{
			Id:        entry_array[0],
			Name:      entry_array[1],
			Type:      entry_array[2],
			Attack:    int_array[0],
			Defense:   int_array[1],
			StarChips: int_array[2],
		}
		cards = append(cards, card)
	}
	return cards
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
    type TEXT,
    atk INTEGER,
    def INTEGER,
    star_chips INTEGER
    );
    `
	runSql(sqlQuery)
	assert(tableExists("cards"), "cards table failed to initialize properly")
}

func initializeFusionsDB() {
	assert(!tableExists("m1"), "m1 table already exists...")
	sqlQuery := `
CREATE TABLE m1 (
    id TEXT PRIMARY KEY,
    resulting_fusion TEXT,
    fusion_number INTEGER,
    card TEXT
    );`
	runSql(sqlQuery)
	assert(tableExists("m1"), "m1 table failed to properly intialize")

	assert(!tableExists("m2"), "m2 table already exists...")
	sqlQuery = `
CREATE TABLE m2 (
    id TEXT PRIMARY KEY,
    resulting_fusion TEXT,
    fusion_number INTEGER,
    card TEXT
    );`
	runSql(sqlQuery)
	assert(tableExists("m2"), "m2 table failed to properly intialize")
}

func initialize_rate_table(table_name string) {
	assert(!tableExists(table_name), fmt.Sprintf("%v table already exists, should not call this function", table_name))
	sql_query := fmt.Sprintf("CREATE TABLE %v (id TEXT PRIMARY KEY, duel TEXT, card TEXT, rate INTEGER, UNIQUE(duel, card));", table_name)
	runSql(sql_query)
	assert(tableExists(table_name), fmt.Sprintf("failed to initialize the %v table", table_name))
}

func initialize_targets_table() {
	assert(!tableExists("targets"), "targets table already exists")
	sql_query := `
CREATE TABLE targets (
    id TEXT PRIMARY KEY,
    equip TEXT,
    target TEXT
    );`
	runSql(sql_query)
	assert(tableExists("targets"), "targets table failed to properly initialize")
}

func initialize_cards_stars_table() {
	assert(!tableExists("cards_stars"), "cards_stars table already exists")
	sql_query := `
CREATE TABLE cards_stars (
    id TEXT PRIMARY KEY,
    card TEXT,
    star TEXT
    );`
	runSql(sql_query)
	assert(tableExists("cards_stars"), "cards_stars table failed to properly initialize")
}

func initialize_starting_deck_rates_table() {
	assert(!tableExists("starting_deck_rates"), "starting_deck_rates table already exists")
	sql_query := `
CREATE TABLE starting_deck_rates (
    id TEXT PRIMARY KEY,
    pool TEXT,
    card TEXT,
    rate INTEGER
    );`
	runSql(sql_query)
	assert(tableExists("starting_deck_rates"), "starting_deck_rates table failed to properly initialize")
}

func WriteProbabilities(entries []Probability, table_name string) {
	values_string := ""
	for _, entry := range entries {
		entry_string := fmt.Sprintf(", (\"%v\", \"%v\", \"%v\", %v)", entry.Id, entry.Duel, entry.Card, entry.Rate)
		values_string += entry_string
	}
	values_string = values_string[2:]
	sql_query := fmt.Sprintf("INSERT INTO %v VALUES %v;", table_name, values_string)
	runSql(sql_query)
	assert(!tableIsEmpty(table_name), fmt.Sprintf("Something went wrong writing to the %v table, no data was written", table_name))
}

func write_cards_to_db(entries []Card, table_name string) {
	values_string := ""
	for _, entry := range entries {
		entry_string := fmt.Sprintf(", (\"%v\", \"%v\", \"%v\", %v, %v, %v)", entry.Id, entry.Name, entry.Type, entry.Attack, entry.Defense, entry.StarChips)
		values_string += entry_string
	}
	values_string = values_string[2:]
	sql_query := fmt.Sprintf("INSERT INTO %v VALUES %v;", table_name, values_string)
	runSql(sql_query)
	assert(!tableIsEmpty("cards"), "Something went wrong writing to the cards tables, no data was written")
}

func write_targets_to_db(entries []Target, table_name string) {
	values_string := ""
	for _, entry := range entries {
		entry_string := fmt.Sprintf(",\n(\"%v\", \"%v\", \"%v\")", entry.Id, entry.Equip, entry.Target)
		values_string += entry_string
	}
	values_string = values_string[2:]
	sql_query := fmt.Sprintf("INSERT INTO %v VALUES %v;", table_name, values_string)
	runSql(sql_query)
	assert(!tableIsEmpty("targets"), "Something went wrong writing to the targets tables, no data was written")
}

func write_cards_stars_to_db(entries []Card_Star, table_name string) {
	values_string := ""
	for _, entry := range entries {
		entry_string := fmt.Sprintf(", (\"%v\", \"%v\", \"%v\")", entry.Id, entry.Card, entry.Star)
		values_string += entry_string
	}
	values_string = values_string[2:]
	sql_query := fmt.Sprintf("INSERT INTO %v VALUES %v;", table_name, values_string)
	runSql(sql_query)
	assert(!tableIsEmpty("cards_stars"), "Something went wrong writing to the cards_stars tables, no data was written")
}

func write_fusions_to_db(entries []Material, table_name string) {
	values_string := ""
	for _, entry := range entries {
		entry_string := fmt.Sprintf(", (\"%v\", \"%v\", %v, \"%v\")", entry.Id, entry.Resulting_Fusion, entry.Fusion_Number, entry.Card)
		values_string += entry_string
	}
	values_string = values_string[2:]
	sql_query := fmt.Sprintf("INSERT INTO %v VALUES %v;", table_name, values_string)
	runSql(sql_query)
	assert(!tableIsEmpty(table_name), fmt.Sprintf("Something went wrong writing the %v table", table_name))
}
