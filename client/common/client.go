package common

import (
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
	if c.sigChan != nil {
		close(c.sigChan)
	}
	c.isRunning = false
	c.protocol.Shutdown()
	log.Infof("action: shutdown | result: success | client_id: %v", c.config.ID)
}

func (c *Client) handleSignal() {
	<-c.sigChan
	c.Shutdown()
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	go c.handleSignal()

	if !c.isRunning {
		return
	}

	bet := newBet()
	err := c.protocol.SendBet(bet.serialize())

	if err != nil {
		log.Errorf("action: send_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}

	err = c.protocol.ReceiveConfirmation()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}

	log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", bet.document, bet.number)
	c.protocol.Shutdown()
}
