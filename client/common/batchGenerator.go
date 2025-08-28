package common

import (
	"bufio"
	"os"
)

type BatchGenerator struct {
	file      *os.File
	scanner   *bufio.Scanner
	currLine  int
	isReading bool
}

func NewBatchGenerator(filePath string) (*BatchGenerator, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &BatchGenerator{
		file:      file,
		scanner:   bufio.NewScanner(file),
		currLine:  1,
		isReading: true,
	}, nil
}

func (bg *BatchGenerator) IsReading() bool {
	return bg.isReading
}

func (bg *BatchGenerator) Read(batchSize int) (*Batch, error) {
	batch := NewBatch(batchSize)
	currLine := 1

	for bg.scanner.Scan() {

		if currLine < bg.currLine {
			currLine++
			continue
		}

		betStr := bg.scanner.Text()
		bet, err := betFromString(betStr)
		if err != nil {
			return nil, err
		}

		err = batch.AddBet(*bet)

		if err != nil {
			return batch, nil
		}

		bg.currLine++
		currLine++
	}

	bg.isReading = false
	return batch, nil
}

func (bg *BatchGenerator) Close() error {
	return bg.file.Close()
}
