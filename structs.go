package main

type Character struct {
	Id    string
	Name  string
	Deck  CardTable
	Drops DropTable
}

type Deck struct {
	Id        string
	Name      string
	PlayTable []Card
	DropTable []Card
}

type Card struct {
	Id      string
	Name    string
	Attack  int
	Defense int
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
