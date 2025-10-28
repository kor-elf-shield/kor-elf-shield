package firewall

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
}
