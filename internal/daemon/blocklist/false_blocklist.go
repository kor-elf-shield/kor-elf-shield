package blocklist

type FalseBlocklist struct {
}

func NewFalseBlocklist() Blocklist {
	return &FalseBlocklist{}
}

func (b *FalseBlocklist) NftReload(_ newBlocklist) error {
	return nil
}

func (b *FalseBlocklist) Run() {}

func (b *FalseBlocklist) Close() error {
	return nil
}
