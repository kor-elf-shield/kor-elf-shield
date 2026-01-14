package config

type Login struct {
	Enabled bool
	Notify  bool
	SSH     LoginSSH
	Local   LoginLocal
}

type LoginSSH struct {
	Enabled bool
	Notify  bool
}

type LoginLocal struct {
	Enabled bool
	Notify  bool
}
