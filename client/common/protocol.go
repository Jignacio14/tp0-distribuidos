package common

import (
	"net"
)

type Protocol struct {
	conn net.Conn
}

func NewProtocol(serverAdr string) (*Protocol, error) {
	conn, err := net.Dial("tcp", serverAdr)
	if err != nil {
		return nil, err
	}

	protocol := &Protocol{conn: conn}
	return protocol, nil
}

func (p *Protocol) SendClientInfo(clientInfo string) error {
	clientInfo += "\n"
	data := []byte(clientInfo)
	return p.sendAll(data)
}

func (p *Protocol) SendBatch(batchStr string) error {
	batchStr += "\t"
	data := []byte(batchStr)
	return p.sendAll(data)
}

func (p *Protocol) SendEndOfBatch() error {
	data := []byte("\000")
	return p.sendAll(data)
}

func (p *Protocol) ReceiveConfirmation() bool {
	response := make([]byte, 2)

	err := p.receiveAll(response)

	if err != nil {
		return false
	}

	return string(response) == "OK"
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
