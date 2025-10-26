package firewall

import "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"

type Config struct {
	SavesRules     bool
	SavesRulesPath string
	MetadataNaming ConfigMetadata
	Policy         ConfigPolicy
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
