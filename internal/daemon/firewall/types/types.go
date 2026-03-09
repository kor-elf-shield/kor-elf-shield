package types

import "fmt"

type PolicyDrop int8

const (
	Drop PolicyDrop = iota + 1
	Reject
)

func (p PolicyDrop) String() string {
	switch p {
	case Drop:
		return "drop"
	case Reject:
		return "reject"
	default:
		return "drop"
	}
}

type Action int8

const (
	ActionAccept Action = iota + 1
	ActionReject
	ActionDrop
)

func (a Action) String() string {
	switch a {
	case ActionAccept:
		return "accept"
	case ActionReject:
		return "reject"
	case ActionDrop:
		return "drop"
	default:
		return "drop"
	}
}

type Protocol int8

const (
	ProtocolTCP Protocol = iota + 1
	ProtocolUDP
)

func (p Protocol) String() string {
	switch p {
	case ProtocolTCP:
		return "tcp"
	case ProtocolUDP:
		return "udp"
	default:
		return fmt.Sprintf("Protocol(%d)", p)
	}
}

type Direction int8

const (
	DirectionIn Direction = iota + 1
	DirectionOut
)

func (d Direction) String() string {
	switch d {
	case DirectionIn:
		return "in"
	case DirectionOut:
		return "out"
	default:
		return fmt.Sprintf("Direction(%d)", d)
	}
}
