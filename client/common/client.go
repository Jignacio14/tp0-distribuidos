package common

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	MaxBatchSize  int
}

// Client Entity that encapsulates how
type Client struct {
	config    ClientConfig
	sigChan   chan os.Signal
	isRunning bool
	protocol  *Protocol
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {

	protocol, err := NewProtocol(config.ServerAddress)

	if err != nil {
		log.Criticalf("action: init | result: fail | client_id: %v | error: %v",
			config.ID,
			err,
		)
		return nil
	}

	client := &Client{
		sigChan:   make(chan os.Signal, 1),
		isRunning: true,
		protocol:  protocol,
		config:    config,
	}

	signal.Notify(client.sigChan, syscall.SIGTERM)
	return client
}

func (c *Client) Shutdown() {
	<-c.sigChan
	close(c.sigChan)
	c.isRunning = false
	c.protocol.Shutdown()
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	go c.Shutdown()
	defer c.protocol.Shutdown()

	filepath := fmt.Sprintf(".data/agency-%v.csv", c.config.ID)
	batchGenerator, err := NewBatchGenerator(filepath)

	if err != nil {
		log.Errorf("action: batch_generator_init | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	defer batchGenerator.Close()

	for batchGenerator.IsReading() {

		batch, err := batchGenerator.Read(c.config.MaxBatchSize)

		if err != nil {
			log.Errorf("action: read_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		batchStr := batch.Serialize()

		err = c.protocol.SendBatch(batchStr)

		if err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		bets_processed, err, status := c.protocol.ReceivedConStatus()

		if err != nil {
			log.Errorf("action: receive_confirmation | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		if !status {
			log.Errorf("action: apuesta_recibida | result: fail | cantidad: %v ", bets_processed)
			break
		}
	}

	err = c.protocol.EndSedingBets()

	if err != nil {
		log.Errorf("action: end_sending_batches | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	status, err := c.protocol.ReceivedEnd()

	if err != nil {
		log.Errorf("action: receive_confirmation | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	if !status {
		log.Errorf("action: receive_confirmation | result: fail | client_id: %v", c.config.ID)
		return
	}

	log.Infof("action: complete | result: success | client_id: %v", c.config.ID)
}
