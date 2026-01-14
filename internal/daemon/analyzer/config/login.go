package config

type Login struct {
	Enabled bool
	Notify  bool
	SSH     LoginSSH
	Local   LoginLocal
	Su      LoginSu
}

type LoginSSH struct {
	Enabled bool
	Notify  bool
}

type LoginLocal struct {
	Enabled bool
	Notify  bool
}

type LoginSu struct {
	Enabled bool
	Notify  bool
}
