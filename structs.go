package main

type Character struct {
	Id   string
	Name string
	Deck Deck
}

type Deck struct {
	Id        string
	Name      string
	PlayTable []Card
	DropTable []Card
}

type Card struct {
	Id   string
	Name string
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
