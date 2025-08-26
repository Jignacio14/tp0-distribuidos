package common

import (
	"fmt"
	"os"
)

type Bet struct {
	name      string
	lastname  string
	document  string
	birthdate string
	number    string
}

func newBet() *Bet {
	name := os.Getenv("NAME")
	lastname := os.Getenv("LASTNAME")
	document := os.Getenv("DOCUMENT")
	birthdate := os.Getenv("BIRTHDATE")
	number := os.Getenv("NUMBER")

	return &Bet{
		name:      name,
		lastname:  lastname,
		document:  document,
		birthdate: birthdate,
		number:    number,
	}
}

func (bet Bet) serialize() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s", bet.name, bet.lastname, bet.document, bet.birthdate, bet.number)
}
