package main

type Character struct {
	Id   string
	Name string
}

type Card struct {
	Id        string
	Name      string
	Type      string
	Attack    int
	Defense   int
	StarChips int
}

type Target struct {
	Id     string
	Equip  string
	Target string
}

type Star struct {
	Id        string
	Name      string
	Strong_To string
	Weak_To   string
}

type Card_Star struct {
	Id   string
	Card string
	Star string
}

type Material struct {
	Id               string
	Resulting_Fusion string
	Fusion_Number    int
	Card             string
}

type Probability struct {
	Id   string
	Duel string
	Card string
	Rate int
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

type Card_Page_JSON struct {
	Query struct {
		Pages map[string]Page `json:"pages"`
	} `json:"query"`
}

type Deck []string

type WikiSection []string

type DuelText []string

type DuelTableEntry [3]string

type DuelTable []DuelTableEntry

type DeckText []string
