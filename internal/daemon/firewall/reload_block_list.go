package firewall

func (f *firewall) reloadBlockList() error {
	listBlockedIP, err := f.chains.NewBlockListIP("blocked_ip")
	if err != nil {
		return err
	}
	if err := listBlockedIP.AddRuleToChain(f.chains.BeforeLocalInput().AddRule, "drop"); err != nil {
		return err
	}

	listBlockedIPWithPort, err := f.chains.NewBlockListIPWithPort("blocked_ip_port")
	if err != nil {
		return err
	}
	if err := listBlockedIPWithPort.AddRuleToChain(f.chains.BeforeLocalInput().AddRule, "drop"); err != nil {
		return err
	}

	if err := f.blockingService.NftReload(listBlockedIP, listBlockedIPWithPort); err != nil {
		return err
	}

	return nil
}
