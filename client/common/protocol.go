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
	sendBetCode      = 0x01
	confirmationCode = 0x02
)

func NewProtocol(serverAdr string) (*Protocol, error) {
	conn, err := net.Dial("tcp", serverAdr)
	if err != nil {
		return nil, err
	}

	protocol := &Protocol{conn: conn}
	return protocol, nil
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

func (p *Protocol) htonsUint32(val uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, val)
	return bytes
}

func (p *Protocol) SendBet(bet string) error {
	lenData := p.htonsUint32(uint32(len(bet)))
	opCode := []byte{sendBetCode}

	err := p.sendAll(opCode)

	if err != nil {
		return err
	}

	err = p.sendAll(lenData)

	if err != nil {
		return err
	}

	err = p.sendAll([]byte(bet))

	if err != nil {
		return err
	}

	return nil
}

func (p *Protocol) ReceiveConfirmation() error {
	opCode := make([]byte, 1)
	err := p.receiveAll(opCode)

	if err != nil {
		return err
	}

	if opCode[0] != confirmationCode {
		return fmt.Errorf("invalid op code received: %v", opCode[0])
	}

	return nil
}

func (p *Protocol) Shutdown() {
	if p.conn == nil {
		return
	}
	p.conn.Close()
}
