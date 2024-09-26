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
