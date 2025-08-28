package common

import (
	"fmt"
	"strings"
)

// 8kb lumit for each batch
const Limit = 8000

// Batch entity that encapsulates a group of bets to be sent together
type Batch struct {
	bets []string
	size int
	max  int
}

// Creates a new batch with a given size limit
// if limit is greater than the constant Limit, Limit is used instead
func NewBatch(max int) *Batch {
	return &Batch{
		bets: make([]string, 0, Limit),
		max:  max,
		size: Limit,
	}
}

// private method to check if a new bet can be appended to the batch
func (b *Batch) canAppend(serialize string) bool {
	return len(b.bets)+1 <= b.max && len(serialize) < b.size
}

// Adds a new bet if it does not exceed the batch size limit
// returns an error if the bet cannot be added
func (b *Batch) AddBet(bet Bet) error {
	serialize := bet.serialize()
	serialize = strings.Trim(serialize, "\n")
	if !b.canAppend(serialize) {
		return fmt.Errorf("batch size exceeded")
	}

	b.bets = append(b.bets, serialize)
	b.size -= len(serialize)
	log.Infof("action: add_bet | result: success | current_batch_size: %v | bet: %v", len(b.bets), serialize)
	return nil
}

// Serializes the batch into a string, with bets separated by new lines
func (b *Batch) Serialize() string {
	return strings.Join(append(b.bets, ""), "\n")
}
