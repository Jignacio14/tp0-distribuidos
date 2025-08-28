package common

import (
	"fmt"
	"strings"
)

// 8kb lumit for each batch
const Limit = 8000

// Batch entity that encapsulates a group of bets to be sent together
type Batch struct {
	bets     []string
	currSize int
	size     int
}

// Creates a new batch with a given size limit
// if limit is greater than the constant Limit, Limit is used instead
func NewBatch(size int) *Batch {
	return &Batch{
		bets:     make([]string, 0, size),
		currSize: 0,
		size:     size,
	}
}

// private method to check if a new bet can be appended to the batch
func (b *Batch) canAppend(serialize string) bool {
	nextSize := b.currSize + len(serialize)
	return nextSize < b.size && nextSize < Limit
}

// Adds a new bet if it does not exceed the batch size limit
// returns an error if the bet cannot be added
func (b *Batch) AddBet(bet Bet) error {
	serialize := bet.serialize()

	if !b.canAppend(serialize) {
		log.Info("action: add_bet | result: fail | error: batch_size_exceeded")
		return fmt.Errorf("batch size exceeded")
	}

	b.bets = append(b.bets, serialize)
	b.currSize += len(serialize)
	log.Infof("action: add_bet | result: success | current_batch_size: %v | bet: %v", len(b.bets), serialize)
	return nil
}

// Serializes the batch into a string, with bets separated by new lines
func (b Batch) Serialize() string {
	return strings.Join(b.bets, "\n")
}
