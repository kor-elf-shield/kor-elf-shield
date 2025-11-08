package setting

type binaryLocations struct {
	Nftables string `mapstructure:"nftables"`
}

func binaryLocationsDefault() *binaryLocations {
	return &binaryLocations{
		Nftables: "/usr/sbin/nft",
	}
}
