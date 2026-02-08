package analyzer

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/analyzer/config"
)

type Login struct {
	Enabled bool `mapstructure:"enabled"`
	Notify  bool `mapstructure:"notify"`

	SSHEnable bool `mapstructure:"ssh_enable"`
	SSHNotify bool `mapstructure:"ssh_notify"`

	LocalEnable bool `mapstructure:"local_enable"`
	LocalNotify bool `mapstructure:"local_notify"`

	SuEnable bool `mapstructure:"su_enable"`
	SuNotify bool `mapstructure:"su_notify"`

	SudoEnable bool `mapstructure:"sudo_enable"`
	SudoNotify bool `mapstructure:"sudo_notify"`
}

func defaultLogin() Login {
	return Login{
		Enabled: true,
		Notify:  true,

		SSHEnable: true,
		SSHNotify: true,

		LocalEnable: true,
		LocalNotify: true,

		SuEnable: true,
		SuNotify: true,

		SudoEnable: false,
		SudoNotify: true,
	}
}

func (l Login) Validate() error {
	return nil
}

func (l Login) ToSources() ([]*config.Source, error) {
	var sources []*config.Source

	if !l.Enabled {
		return sources, nil
	}

	if l.SSHEnable {
		sources = append(sources, config.NewLoginSSH(l.Notify && l.SSHNotify)...)
	}

	if l.LocalEnable {
		sources = append(sources, config.NewLoginLocal(l.Notify && l.LocalNotify)...)
	}

	if l.SuEnable {
		sources = append(sources, config.NewLoginSu(l.Notify && l.SuNotify)...)
	}

	if l.SudoEnable {
		sources = append(sources, config.NewLoginSudo(l.Notify && l.SudoNotify)...)
	}

	return sources, nil
}
