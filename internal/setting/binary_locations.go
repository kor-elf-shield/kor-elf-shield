package setting

type binaryLocations struct {
	Nftables   string `mapstructure:"nftables"`
	Journalctl string `mapstructure:"journalctl"`
}

func binaryLocationsDefault() *binaryLocations {
	return &binaryLocations{
		Nftables:   "/usr/sbin/nft",
		Journalctl: "/bin/journalctl",
	}
}
