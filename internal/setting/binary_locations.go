package setting

type binaryLocations struct {
	Nftables   string `mapstructure:"nftables"`
	Journalctl string `mapstructure:"journalctl"`
	Docker     string `mapstructure:"docker"`
}

func binaryLocationsDefault() *binaryLocations {
	return &binaryLocations{
		Nftables:   "/usr/sbin/nft",
		Journalctl: "/bin/journalctl",
		Docker:     "/usr/bin/docker",
	}
}
