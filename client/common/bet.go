package common

import (
	"fmt"
	"os"
	"strings"
)

// Bet Represents a bet made by a user
type Bet struct {
	agency    string
	name      string
	lastname  string
	document  string
	birthdate string
	number    string
}

// / Creates a bet from a string representation
func betFromString(bet string) (*Bet, error) {
	parts := strings.Split(bet, ",")

	if len(parts) != 5 {
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
