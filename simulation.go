package main

import "fmt"

func generate_starting_deck() {
	sql_query := "SELECT * from starting_deck_rates WHERE pool = pool_sub_1100;"
	data, err := outputSql(sql_query)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
