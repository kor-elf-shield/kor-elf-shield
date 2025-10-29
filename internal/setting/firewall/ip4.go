package firewall

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
)

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

func (i ip4) Validate() error {
	if err := validateIcmpRate(i.IcmpInRate, "icmp_in_rate"); err != nil {
		return err
	}
	if err := validateIcmpRate(i.IcmpOutRate, "icmp_out_rate"); err != nil {
		return err
	}
	return nil
}

func validateIcmpRate(rate string, parameterName string) error {
	if rate == "0" {
		return nil
	}

	return validate.NftLimitRate(rate, parameterName)
}
