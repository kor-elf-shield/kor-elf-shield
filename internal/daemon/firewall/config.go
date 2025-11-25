package firewall

import "fmt"

type Config struct {
	InPorts        []ConfigPort
	OutPorts       []ConfigPort
	IP4            ConfigIP4
	IP6            ConfigIP6
	Options        ConfigOptions
	MetadataNaming ConfigMetadata
	Policy         ConfigPolicy
}

type ConfigOptions struct {
	ClearMode      ClearMode
	SavesRules     bool
	SavesRulesPath string
	DnsStrict      bool
	DnsStrictNs    bool
	PacketFilter   bool
}

type ConfigMetadata struct {
	TableName        string
	ChainInputName   string
	ChainOutputName  string
	ChainForwardName string
}

type ConfigPolicy struct {
	DefaultAllowInput   bool
	DefaultAllowOutput  bool
	DefaultAllowForward bool
	InputDrop           PolicyDrop
	OutputDrop          PolicyDrop
	ForwardDrop         PolicyDrop
}

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

type ConfigIP4 struct {
	IcmpIn            bool
	IcmpInRate        string
	IcmpOut           bool
	IcmpOutRate       string
	IcmpTimestampDrop bool
	InIPs             []ConfigIP
	OutIPs            []ConfigIP
}

type ConfigIP6 struct {
	Enable     bool
	IcmpStrict bool
	InIPs      []ConfigIP
	OutIPs     []ConfigIP
}

type ConfigPort struct {
	Number    uint16
	Protocol  Protocol
	Action    Action
	LimitRate string
}

type ConfigIP struct {
	IP        string
	OnlyIP    bool // Port is not taken into account
	Port      uint16
	Action    Action
	Protocol  Protocol
	LimitRate string
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

type ClearMode int8

const (
	ClearModeGlobal ClearMode = iota + 1
	ClearModeOwn
)
