package main

import (
	"fmt"
	"os"
	"os/exec"
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

func tableExists(tableName string) bool {
	sqlQuery := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='%v';", tableName)
	data, _ := outputSql(sqlQuery)
	if strings.ReplaceAll(string(data), "\n", "") == "0" {
		return false
	}
	return true
}

func initializeDB() {
	assert(!fileExists(storageDirectory + "/database.db"))
	sqlQuery := `
CREATE TABLE cardTables (
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
	assert(tableExists("cardTables"))
}
