package firewall

type ip4 struct {
	IcmpIn            bool   `mapstructure:"icmp_in"`
	IcmpInRate        string `mapstructure:"icmp_in_rate"`
	IcmpOut           bool   `mapstructure:"icmp_out"`
	IcmpOutRate       string `mapstructure:"icmp_out_rate"`
	IcmpTimestampDrop bool   `mapstructure:"icmp_timestamp_drop"`
}

func defaultIp4() ip4 {
	return ip4{
		IcmpIn:            true,
		IcmpInRate:        "1/s",
		IcmpOut:           true,
		IcmpOutRate:       "0",
		IcmpTimestampDrop: false,
	}
}
