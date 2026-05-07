package entity

const (
	MetadataKeyVersion         = "Version"
	MetadataKeyFirewallFileNft = "firewall-file-nft" // checksum of the firewall file
)

type Metadata struct {
	Value string `json:"Value"`
}

func KeySetting(name string) string {
	return "setting-" + name
}
