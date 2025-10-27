package firewall

import "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"

type Config struct {
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
