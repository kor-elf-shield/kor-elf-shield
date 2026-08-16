package firewall

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/setting/validate"
	"github.com/spf13/viper"
)

type Setting struct {
	Ports          []Port
	IPS            []IP
	IP4            ip4
	IP6            ip6
	Options        options
	MetadataNaming metadataNaming
	Policy         policy
	PortKnocking   []portKnocking
	RulesGuard     RulesGuard
}

func InitSetting(path string) (Setting, error) {
	if err := validate.IsTomlFile(path, "otherSettingsPath.firewall"); err != nil {
		return Setting{}, err
	}

	setting := settingDefault()

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return Setting{}, err
	}
	if err := v.Unmarshal(&setting); err != nil {
		return Setting{}, err
	}
	if err := setting.Validate(); err != nil {
		return Setting{}, err
	}

	return setting, nil
}

func settingDefault() Setting {
	return Setting{
		IPS:            defaultIPs(),
		Ports:          defaultPorts(),
		IP4:            defaultIp4(),
		IP6:            defaultIp6(),
		Options:        defaultOptions(),
		MetadataNaming: defaultMetadataNaming(),
		Policy:         defaultPolicy(),
		PortKnocking:   defaultPortKnocking(),
		RulesGuard:     defaultRulesGuard(),
	}
}

func (s Setting) Validate() error {
	if err := s.IP4.Validate(); err != nil {
		return err
	}
	if err := s.IP6.Validate(); err != nil {
		return err
	}
	if err := s.MetadataNaming.Validate(); err != nil {
		return err
	}
	if err := s.Policy.Validate(); err != nil {
		return err
	}
	if err := s.Options.Validate(); err != nil {
		return err
	}
	if err := s.RulesGuard.Validate(); err != nil {
		return err
	}
	return nil
}

func (s Setting) ToPorts() (InPorts []config.ConfigPort, OutPorts []config.ConfigPort, error error) {
	for _, port := range s.Ports {
		addInPorts, addOutPorts, err := port.ToPorts()
		if err != nil {
			error = err
			return
		}
		InPorts = append(InPorts, addInPorts...)
		OutPorts = append(OutPorts, addOutPorts...)
	}

	return
}

func (s Setting) ToIPs() (IPs IPs, error error) {
	for _, ips := range s.IPS {
		addIPs, err := ips.ToIPs()
		if err != nil {
			error = err
			return
		}
		IPs.InIP4 = append(IPs.InIP4, addIPs.InIP4...)
		IPs.OutIP4 = append(IPs.OutIP4, addIPs.OutIP4...)

		IPs.InIP6 = append(IPs.InIP6, addIPs.InIP6...)
		IPs.OutIP6 = append(IPs.OutIP6, addIPs.OutIP6...)
	}

	return
}

func (s Setting) ToConfigPortKnocking() ([]config.ConfigPortKnocking, error) {
	var configPortKnocking []config.ConfigPortKnocking

	portKnockingNames := make(map[string]string)

	for _, portKnocking := range s.PortKnocking {
		if _, ok := portKnockingNames[portKnocking.Name]; ok {
			return nil, fmt.Errorf("port knocking name %s is duplicated", portKnocking.Name)
		}
		portKnockingNames[portKnocking.Name] = portKnocking.Name

		addPortKnocking, err := portKnocking.ToPortKnocking()
		if err != nil {
			return nil, err
		}
		configPortKnocking = append(configPortKnocking, addPortKnocking)
	}
	return configPortKnocking, nil
}
