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
	config   ClientConfig
	sigChan  chan os.Signal
	protocol *Protocol
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
		sigChan:  make(chan os.Signal, 1),
		protocol: protocol,
		config:   config,
	}

	signal.Notify(client.sigChan, syscall.SIGTERM)
	return client
}

// signal handler
func (c *Client) handleSignal() {
	<-c.sigChan
	c.Shutdown()
}

// Shutdown Gracefully shuts down the client
func (c *Client) Shutdown() {
	if c.sigChan != nil {
		signal.Stop(c.sigChan)
	}
	if c.protocol != nil {
		c.protocol.Shutdown()
	}
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	go c.handleSignal()
	defer c.Shutdown()

	filepath := fmt.Sprintf(".data/agency-%v.csv", c.config.ID)
	batchGenerator, err := NewBatchGenerator(filepath)

	if err != nil {
		log.Errorf("action: batch_generator_init | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	defer batchGenerator.Close()

	err = c.protocol.SendLoteryId(c.config.ID)

	if err != nil {
		log.Errorf("action: send_lotery_id | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	err, should_wait_end := c.loop(batchGenerator)

	if err != nil && !should_wait_end {
		return
	}

	if err != nil && should_wait_end {
		return
	}

	_ = c.finishCommunication()

}

// loop sends batches to the server until there are no more bets to send or the server fails to process a bet
func (c *Client) loop(batchGenerator *BatchGenerator) (error, bool) {

	for batchGenerator.IsReading() {

		batch, err := batchGenerator.Read(c.config.MaxBatchSize)

		if err != nil {
			log.Errorf("action: read_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return err, false
		}

		err = c.protocol.SendBatch(batch)

		if err != nil {
			log.Errorf("action: send_batch | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return err, false
		}

		bets_processed, err, status := c.protocol.ReceivedConStatus()

		if err != nil {
			log.Errorf("action: receive_confirmation | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return err, false
		}

		if !status {
			log.Errorf("action: apuesta_recibida | result: fail | cantidad: %v ", bets_processed)
			return fmt.Errorf("apuestas con error"), true
		}
	}

	return nil, true

}

// In case all batches were sent successfully, finish the communication with the server
func (c *Client) finishCommunication() error {

	err := c.protocol.EndSedingBets()

	if err != nil {
		log.Errorf("action: end_sending_batches | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}

	winners, err := c.protocol.ReceiveWinners()

	if err != nil {
		log.Errorf("action: wait_for_ending | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}

	c.LogWinners(winners)

	log.Infof("action: complete | result: success | client_id: %v", c.config.ID)

	return nil
}

// / Receives the operation end from server
func (c *Client) waitForEnding() error {

	status, err := c.protocol.ReceivedEnd()

	if err != nil {
		log.Errorf("action: receive_confirmation | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}

	if !status {
		log.Errorf("action: receive_confirmation | result: fail | client_id: %v", c.config.ID)
		return err
	}

	return nil

}

func (c *Client) LogWinners(winners []string) {
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", len(winners))
}
