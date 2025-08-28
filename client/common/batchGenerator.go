package common

import (
	"bufio"
	"os"
)

type BatchGenerator struct {
	file         *os.File
	scanner      *bufio.Scanner
	currLine     int
	isReading    bool
	lastLineRead string
}

func NewBatchGenerator(filePath string) (*BatchGenerator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &BatchGenerator{
		file:         file,
		scanner:      bufio.NewScanner(file),
		currLine:     1,
		isReading:    true,
		lastLineRead: "",
	}, nil
}

func (bg *BatchGenerator) IsReading() bool {
	return bg.isReading
}

func (bg *BatchGenerator) Read(batchSize int) (*Batch, error) {
	batch := NewBatch(batchSize)

	if bg.lastLineRead != "" {
		bet, err := betFromString(bg.lastLineRead)
		if err != nil {
			return nil, err
		}

		err = batch.AddBet(*bet)

		if err != nil {
			return batch, nil
		}

		bg.lastLineRead = ""
	}

	for bg.scanner.Scan() {

		betStr := bg.scanner.Text()
		bet, err := betFromString(betStr)
		if err != nil {
			return nil, err
		}

		err = batch.AddBet(*bet)

		if err != nil {
			bg.lastLineRead = betStr
			return batch, nil
		}

	}

	bg.isReading = false
	return batch, nil
}

func (bg *BatchGenerator) Close() error {
	return bg.file.Close()
}
