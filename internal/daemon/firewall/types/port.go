package types

import (
	"errors"
	"strconv"
)

type L4Port interface {
	Number() uint16
	NumberString() string
	ProtocolString() string
}

type l4Port struct {
	number   uint16
	protocol string
}

func NewL4Port(number uint16, protocol Protocol) (L4Port, error) {
	if protocol != ProtocolTCP && protocol != ProtocolUDP {
		return nil, errors.New("invalid protocol")
	}

	return &l4Port{number: number, protocol: protocol.String()}, nil
}

func (p *l4Port) Number() uint16 {
	return p.number
}

func (p *l4Port) NumberString() string {
	port := p.Number()
	return strconv.Itoa(int(port))
}

func (p *l4Port) ProtocolString() string {
	return p.protocol
}
