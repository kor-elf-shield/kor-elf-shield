package firewall

import (
	"fmt"
	"os"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"

	nftables "git.kor-elf.net/kor-elf-shield/go-nftables-client"
)

type API interface {
	// Reload Clear all rules and set new rules.
	Reload() error

	// SavesRules Save rules to file.
	SavesRules()

	// ClearRules Clear all rules.
	ClearRules()

	DockerSupport() bool
}

type firewall struct {
	nft    nftables.NFT
	logger log.Logger
	config *Config
	chains chain.Chains
	docker docker_monitor.Docker
}

func New(pathNFT string, logger log.Logger, config Config, docker docker_monitor.Docker) (API, error) {
	nft, err := nftables.NewWithPath(pathNFT)
	if err != nil {
		return nil, fmt.Errorf("failed to create nft client: %w %s", err, pathNFT)
	}

	return &firewall{
		nft:    nft,
		logger: logger,
		config: &config,
		docker: docker,
	}, nil
}

func (f *firewall) Reload() error {
	f.logger.Debug("Reload nftables rules")
	if f.config.Options.ClearMode == ClearModeGlobal {
		if err := f.nft.Clear(); err != nil {
			return err
		}
	}

	chains, err := chain.NewChains(f.nft, f.config.MetadataNaming.TableName)
	if err != nil {
		return err
	}
	f.chains = chains

	if err := f.docker.NftReload(f.chains.NewNoneChain); err != nil {
		return err
	}

	if err := f.chains.NewPacketFilter(f.config.Options.PacketFilter); err != nil {
		return err
	}
	if err := f.reloadInput(); err != nil {
		return err
	}
	if err := f.reloadOutput(); err != nil {
		return err
	}
	if err := f.reloadForward(); err != nil {
		return err
	}
	if f.config.Options.DockerSupport {
		if err := f.reloadDocker(); err != nil {
			return err
		}
	}

	f.logger.Debug("Reload nftables rules done")
	return nil
}

func (f *firewall) ClearRules() {
	f.logger.Debug("Clear nftables rules")

	switch f.config.Options.ClearMode {
	case ClearModeGlobal:
		if err := f.nft.Clear(); err != nil {
			f.logger.Error(fmt.Sprintf("Failed to clear rules: %s", err))
		}
		break
	case ClearModeOwn:
		if err := f.chains.ClearRules(); err != nil {
			f.logger.Error(fmt.Sprintf("Failed to clear rules: %s", err))
		}
		break
	}

	f.logger.Debug("Clear nftables rules done")
}

func (f *firewall) SavesRules() {
	if !f.config.Options.SavesRules {
		f.logger.Debug("SavesRules is false, skip")
		return
	}

	if f.config.Options.SavesRulesPath == "" {
		f.logger.Warn("SavesRulesPath is empty, skip")
		return
	}

	args := []string{"list", "ruleset"}
	output, err := f.nft.Command().RunWithOutput(args...)
	if err != nil {
		f.logger.Warn(fmt.Sprintf("Failed to save rules: %s", err))
		return
	}

	data := []byte("#!/usr/sbin/nft -f\n\nflush ruleset\n" + output)
	err = os.WriteFile(f.config.Options.SavesRulesPath, data, 0755)
	if err != nil {
		f.logger.Warn(fmt.Sprintf("Failed to save rules: %s", err))
		return
	}

	f.logger.Info("Save nftables rules")
}

func (f *firewall) DockerSupport() bool {
	return f.config.Options.DockerSupport
}
