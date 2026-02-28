package firewall

func (f *firewall) reloadBlockList() error {
	listBan, err := f.chains.NewBlockListIP("ban")
	if err != nil {
		return err
	}
	if err := listBan.AddRuleToChain(f.chains.BeforeLocalInput().AddRule, "drop"); err != nil {
		return err
	}

	if err := f.blockingService.NftReload(listBan); err != nil {
		return err
	}

	return nil
}
