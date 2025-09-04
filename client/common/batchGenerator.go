package common

import (
	"bufio"
	"os"
)

// BatchGenerator Reads bets from a file and generates batches
type BatchGenerator struct {
	file         *os.File
	scanner      *bufio.Scanner
	isReading    bool
	lastLineRead string
}

// / Creates a new BatchGenerator for the given file path
func NewBatchGenerator(filePath string) (*BatchGenerator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &BatchGenerator{
		file:         file,
		scanner:      bufio.NewScanner(file),
		isReading:    true,
		lastLineRead: "",
	}, nil
}

// Indicates if there are more bets to read
func (bg *BatchGenerator) IsReading() bool {
	return bg.isReading
}

// / Reads the next batch of bets from the file
func (bg *BatchGenerator) Read(batchSize int) (*Batch, error) {

	batch := NewBatch(batchSize)

	err := bg.processLastLine(batch)

	if err != nil {
		return batch, err
	}

	for bg.scanner.Scan() {

		betStr := bg.scanner.Text()

		if betStr == "" {
			continue
		}

		bet, err := betFromString(betStr)

		if err != nil {
			continue
		}

		err = batch.AddBet(bet)

		if err != nil {
			bg.lastLineRead = betStr
			return batch, nil
		}

	}

	bg.isReading = false
	return batch, nil
}

// / Processes the last line read that couldn't be added to the previous batch
func (bg *BatchGenerator) processLastLine(batch *Batch) error {
	if bg.lastLineRead == "" {
		return nil
	}

	bet, err := betFromString(bg.lastLineRead)

	if err != nil {
		bg.lastLineRead = ""
		return nil
	}

	err = batch.AddBet(bet)

	if err != nil {
		return err
	}

	bg.lastLineRead = ""

	return nil
}

// Closes the underlying file
func (bg *BatchGenerator) Close() {
	bg.file.Close()
}
