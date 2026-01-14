package analyzer

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
