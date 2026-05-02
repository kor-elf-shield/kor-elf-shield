package data

type NftOutput struct {
	Nftables []NftElement `json:"nftables"`
}
type NftElement struct {
	Rule *Rule `json:"rule,omitempty"`
}

type Rule struct {
	Handle  uint64 `json:"handle"`
	Comment string `json:"comment"`
}
