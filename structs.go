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

	FusionInfo string
}

type Probability struct {
	Id     string
	Duel   string
	CardId string

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

type CharactersQuery struct {
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

type Deck []string

type WikiSection []string

type DuelText []string

type DuelTableEntry [3]string

type DuelTable []DuelTableEntry
