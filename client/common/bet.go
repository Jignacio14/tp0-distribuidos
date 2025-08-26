package common

import (
	"fmt"
	"os"
)

type Bet struct {
	agency    string
	name      string
	lastname  string
	document  string
	birthdate string
	number    string
}

func newBet() *Bet {

	agency := os.Getenv("CLI_ID")
	name := os.Getenv("CLI_NAME")
	lastname := os.Getenv("CLI_LASTNAME")
	document := os.Getenv("CLI_DOCUMENT")
	birthdate := os.Getenv("CLI_BIRTHDATE")
	number := os.Getenv("CLI_NUMBER")

	return &Bet{
		agency:    agency,
		name:      name,
		lastname:  lastname,
		document:  document,
		birthdate: birthdate,
		number:    number,
	}
}

func (bet Bet) serialize() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", bet.agency, bet.name, bet.lastname, bet.document, bet.birthdate, bet.number)
}
