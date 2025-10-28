package firewall

import "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"

type Config struct {
	IP4            ConfigIP4
	Options        ConfigOptions
	MetadataNaming ConfigMetadata
	Policy         ConfigPolicy
}

type ConfigOptions struct {
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
	Input   Policy
	Output  Policy
	Forward Policy
}

type ConfigIP4 struct {
	IcmpIn            bool
	IcmpInRate        string
	IcmpOut           bool
	IcmpOutRate       string
	IcmpTimestampDrop bool
}

type Policy int8

const (
	PolicyAccept Policy = iota + 1
	PolicyDrop
	PolicyReject
)

func (p Policy) ChainDefaultPolicy() chain.Policy {
	switch p {
	case PolicyAccept:
		return chain.PolicyAccept
	case PolicyDrop:
		return chain.PolicyDrop
	case PolicyReject:
		return chain.PolicyDrop
	default:
		return chain.PolicyDrop
	}
}

func (p Policy) String() string {
	switch p {
	case PolicyAccept:
		return "accept"
	case PolicyDrop:
		return "drop"
	case PolicyReject:
		return "reject"
	default:
		return "drop"
	}
}
