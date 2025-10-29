package firewall

type ip6 struct {
	Enable     bool `mapstructure:"enable"`
	IcmpStrict bool `mapstructure:"icmp_strict"`
}

func defaultIp6() ip6 {
	return ip6{
		Enable:     false,
		IcmpStrict: false,
	}
}

func (i ip6) Validate() error {
	return nil
}
