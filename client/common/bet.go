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
	name := os.Getenv("CLI_NAME")
	lastname := os.Getenv("CLI_LASTNAME")
	document := os.Getenv("CLI_DOCUMENT")
	birthdate := os.Getenv("CLI_BIRTHDATE")
	number := os.Getenv("CLI_NUMBER")

	return &Bet{
		name:      name,
		lastname:  lastname,
		document:  document,
		birthdate: birthdate,
		number:    number,
	}
}

func (bet Bet) serialize() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s\n", bet.name, bet.lastname, bet.document, bet.birthdate, bet.number)
}
