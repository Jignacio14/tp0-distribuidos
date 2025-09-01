package common

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Protocol struct {
	conn net.Conn
}

const (
	sendBatchCode          byte = 0x01
	receivedBatchOKCode    byte = 0x02
	receivedBatchNotOKCode byte = 0x03
	endOfBatch             byte = 0x04
)

func NewProtocol(serverAdr string) (*Protocol, error) {
	conn, err := net.Dial("tcp", serverAdr)
	if err != nil {
		return nil, err
	}

	protocol := &Protocol{conn: conn}
	return protocol, nil
}

func (p *Protocol) SendBatch(batchStr string) error {

	data := []byte(batchStr)

	opCode := []byte{sendBatchCode}
	err := p.sendAll(opCode)

	if err != nil {
		return err
	}

	length := p.htonsUint32(uint32(len(data)))

	err = p.sendAll(length)

	if err != nil {
		return err
	}

	return p.sendAll(data)
}

func (p *Protocol) htonsUint32(val uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, val)
	return bytes
}

func (p *Protocol) ntohsUint32(data []byte) int {
	return int(binary.BigEndian.Uint32(data))
}

func (p *Protocol) ReceiveConfirmation() bool {
	response := make([]byte, 1)
	err := p.receiveAll(response)

	if err != nil {
		return false
	}

	return response[0] == receivedBatchOKCode
}

func (p *Protocol) ReceivedConStatus() (int, error, bool) {

	responseCode := make([]byte, 1)
	if err := p.receiveAll(responseCode); err != nil {
		return 0, fmt.Errorf("failed to receive response code: %w", err), false
	}

	amountBytes := make([]byte, 4)
	if err := p.receiveAll(amountBytes); err != nil {
		return 0, fmt.Errorf("failed to receive amount: %w", err), false
	}

	total := p.ntohsUint32(amountBytes)

	if responseCode[0] == receivedBatchOKCode {
		return total, nil, true
	}

	if responseCode[0] == receivedBatchNotOKCode {
		return total, nil, false
	}

	return 0, fmt.Errorf("unknown response code: %v", responseCode[0]), false
}

func (p *Protocol) ReceivedEnd() (bool, error) {
	response := make([]byte, 1)
	err := p.receiveAll(response)

	if err != nil {
		return false, err
	}

	return response[0] == endOfBatch, nil
}

func (p *Protocol) EndSedingBets() error {
	opCode := []byte{endOfBatch}
	return p.sendAll(opCode)
}

func (p *Protocol) sendAll(data []byte) error {

	len := len(data)

	for sent := 0; sent < len; {
		n, err := p.conn.Write(data[sent:])
		if err != nil {
			return err
		}
		sent += n
	}

	return nil
}

func (p *Protocol) receiveAll(array []byte) error {
	len := len(array)
	received := 0
	for received < int(len) {
		n, err := p.conn.Read(array[received:])
		if err != nil {
			return err
		}
		received += n
	}

	return nil
}

func (p *Protocol) Shutdown() {
	if p.conn == nil {
		return
	}
	p.conn.Close()
}
