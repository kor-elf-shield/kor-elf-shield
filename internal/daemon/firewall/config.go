package firewall

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/ip"
)

type Config struct {
	InPorts        []ConfigPort
	OutPorts       []ConfigPort
	IP4            ConfigIP4
	IP6            ConfigIP6
	Options        ConfigOptions
	MetadataNaming ConfigMetadata
	Policy         ConfigPolicy
	PortKnocking   []ConfigPortKnocking
}

type ConfigOptions struct {
	ClearMode      ClearMode
	SavesRules     bool
	SavesRulesPath string
	DnsStrict      bool
	DnsStrictNs    bool
	PacketFilter   bool
	DockerSupport  bool
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
	InputDrop           types.PolicyDrop
	InputPriority       int
	OutputDrop          types.PolicyDrop
	OutputPriority      int
	ForwardDrop         types.PolicyDrop
	ForwardPriority     int
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
	Port      types.L4Port
	Action    types.Action
	LimitRate string
}

type ConfigIP struct {
	IP        string
	OnlyIP    bool // Port is not taken into account
	Port      types.L4Port
	Action    types.Action
	LimitRate string
}

type ConfigPortKnocking struct {
	Name      string
	Port      types.L4Port
	IPVersion ip.Version
	Knocks    []*ConfigKnock
}

type ConfigKnock struct {
	Port    types.L4Port
	Action  types.KnockAction
	Timeout uint32
}

type ClearMode int8

const (
	ClearModeGlobal ClearMode = iota + 1
	ClearModeOwn
)
