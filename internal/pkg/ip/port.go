package ip

import (
	"errors"
	"strings"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
)

func ToDirection(direction string) (types.Direction, error) {
	switch strings.ToLower(direction) {
	case "in":
		return types.DirectionIn, nil
	case "out":
		return types.DirectionOut, nil
	default:
		return types.DirectionIn, errors.New("invalid direction. Must be in or out")
	}
}

func ToProtocol(protocol string) (types.Protocol, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return types.ProtocolTCP, nil
	case "udp":
		return types.ProtocolUDP, nil
	default:
		return types.ProtocolTCP, errors.New("invalid protocol. Must be tcp or udp")
	}
}

func ToAction(action string) (types.Action, error) {
	switch strings.ToLower(action) {
	case "accept":
		return types.ActionAccept, nil
	case "drop":
		return types.ActionDrop, nil
	case "reject":
		return types.ActionReject, nil
	default:
		return types.ActionAccept, errors.New("invalid action. Must be accept, drop or reject")
	}
}

func ToKnockAction(action string) (types.KnockAction, error) {
	switch strings.ToLower(action) {
	case "accept":
		return types.KnockActionAccept, nil
	case "drop":
		return types.KnockActionDrop, nil
	case "reject":
		return types.KnockActionReject, nil
	case "return":
		return types.KnockActionReturn, nil
	default:
		return types.KnockActionDrop, errors.New("invalid action. Must be accept, return, drop or reject")
	}
}
