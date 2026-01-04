package firewall

func (f *firewall) reloadForward() error {
	f.logger.Debug("Reloading forward chain")
	err := f.chains.NewForward(f.config.MetadataNaming.ChainForwardName, f.config.Policy.DefaultAllowForward, f.config.Policy.ForwardPriority)
	if err != nil {
		return err
	}
	chain := f.chains.Forward()

	if f.config.Options.DockerSupport {
		if err := f.docker.NftChains().ForwardFilterJump(chain.AddRule); err != nil {
			return err
		}
	}

	if f.config.Policy.DefaultAllowForward == false {
		drop := f.config.Policy.ForwardDrop.String()
		if err := chain.AddRule(drop); err != nil {
			return err
		}
	}

	return nil
}
