package common

import (
	"fmt"
	"os"
	"strings"
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

func betFromString(bet string) (*Bet, error) {
	parts := strings.Split(bet, ",")

	if len(parts) != 6 {
		return nil, fmt.Errorf("invalid bet format")
	}

	agency := os.Getenv("CLI_ID")

	return &Bet{
		agency:    agency,
		name:      parts[0],
		lastname:  parts[1],
		document:  parts[2],
		birthdate: parts[3],
		number:    parts[4],
	}, nil
}

func (bet Bet) serialize() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s", bet.agency, bet.name, bet.lastname, bet.document, bet.birthdate, bet.number)
}
