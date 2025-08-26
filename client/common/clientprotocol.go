package common

import "net"

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

func (p *Protocol) ReceiveConfirmation() bool {
	response := make([]byte, 2)
	lenght, err := p.receiveAll(2, response)

	if err == nil || lenght != 2 {
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

func (p *Protocol) receiveAll(len uint, array []byte) (int, error) {
	received := 0

	for received < int(len) {
		n, err := p.conn.Read(array[received:])
		if err != nil {
			return 0, err
		}
		received += n
	}

	return received, nil
}

func (p *Protocol) Shutdown() {
	if p.conn == nil {
		return
	}
	p.conn.Close()
}
