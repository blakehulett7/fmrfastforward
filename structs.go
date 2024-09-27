package main

type Character struct {
	Id   string
	Name string
}

type Card struct {
	Id   string
	Name string

	Attack  int
	Defense int
}

type CardTable struct {
	Id          string
	CharacterId string
	CardId      string

	Deck  int
	SaPow int
	SaTec int
	Bcd   int
}

type CharactersPageJSON struct {
	Query struct {
		Pages struct {
			Num369496 struct {
				Revisions []struct {
					Body string `json:"*"`
				} `json:"revisions"`
			} `json:"369496"`
		} `json:"pages"`
	} `json:"query"`
}

type CharacterJSON struct {
	Query struct {
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}

type Page struct {
	Title     string `json:"title"`
	Revisions []struct {
		Body string `json:"*"`
	} `json:"revisions"`
}
