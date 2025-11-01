package ip

import (
	"errors"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall"
)

func ToDirection(direction string) (firewall.Direction, error) {
	switch strings.ToLower(direction) {
	case "in":
		return firewall.DirectionIn, nil
	case "out":
		return firewall.DirectionOut, nil
	default:
		return firewall.DirectionIn, errors.New("invalid direction. Must be in or out")
	}
}

func ToProtocol(protocol string) (firewall.Protocol, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return firewall.ProtocolTCP, nil
	case "udp":
		return firewall.ProtocolUDP, nil
	default:
		return firewall.ProtocolTCP, errors.New("invalid protocol. Must be tcp or udp")
	}
}

func ToAction(action string) (firewall.Action, error) {
	switch strings.ToLower(action) {
	case "accept":
		return firewall.ActionAccept, nil
	case "drop":
		return firewall.ActionDrop, nil
	case "reject":
		return firewall.ActionReject, nil
	default:
		return firewall.ActionAccept, errors.New("invalid action. Must be accept, drop or reject")
	}
}
