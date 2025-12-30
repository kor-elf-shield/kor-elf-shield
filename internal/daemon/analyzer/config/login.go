package config

type Login struct {
	Enabled bool
	Notify  bool
	SSH     LoginSSH
}

type LoginSSH struct {
	Enabled bool
	Notify  bool
}
