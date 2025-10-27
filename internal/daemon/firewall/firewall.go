package firewall

import (
	"fmt"
	"os"

	nftableCchain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"

	nftables "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	nftablesFamily "git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
)

type API interface {
	// Reload Clear all rules and set new rules.
	Reload() error

	// SavesRules Save rules to file.
	SavesRules()

	// ClearRules Clear all rules.
	ClearRules()
}

type firewall struct {
	nft    nftables.NFT
	logger log.Logger
	config *Config
}

func New(pathNFT string, logger log.Logger, config Config) (API, error) {
	nft, err := nftables.NewWithPath(pathNFT)
	if err != nil {
		return nil, fmt.Errorf("failed to create nft client: %w %s", err, pathNFT)
	}

	return &firewall{
		nft:    nft,
		logger: logger,
		config: &config,
	}, nil
}

func (f *firewall) Reload() error {
	f.logger.Debug("Reload nftables rules")
	if err := f.nft.Clear(); err != nil {
		return err
	}
	if err := f.nft.Table().Add(nftablesFamily.INET, f.config.MetadataNaming.TableName); err != nil {
		return err
	}
	if err := f.packetFilter(); err != nil {
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

	f.logger.Debug("Reload nftables rules done")
	return nil
}

func (f *firewall) ClearRules() {
	f.logger.Debug("Clear nftables rules")
	if err := f.nft.Clear(); err != nil {
		f.logger.Error(fmt.Sprintf("Failed to clear rules: %s", err))
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

// packetFilter Drop out of order packets and packets in an INVALID state in nftables connection tracking.
func (f *firewall) packetFilter() error {
	if f.config.Options.PacketFilter == false {
		return nil
	}

	family := nftablesFamily.INET
	tableName := f.config.MetadataNaming.TableName
	chainName := "INVDROP"

	if err := f.nft.Chain().Add(family, tableName, chainName, nftableCchain.TypeNone); err != nil {
		return err
	}
	if err := f.nft.Rule().Add(family, tableName, chainName, "counter drop"); err != nil {
		return err
	}

	chainName = "INVALID"
	if err := f.nft.Chain().Add(family, tableName, chainName, nftableCchain.TypeNone); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "ct state invalid counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags ! fin,syn,rst,psh,ack,urg counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (fin | syn | rst | psh | ack | urg) == fin | syn | rst | psh | ack | urg counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (fin | syn) == fin | syn counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (syn | rst) == syn | rst counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (fin | rst) == fin | rst counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (fin | ack) == fin counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (psh | ack) == psh counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (ack | urg) == urg counter jump INVDROP"); err != nil {
		return err
	}

	if err := f.nft.Rule().Add(family, tableName, chainName, "tcp flags & (fin | syn | rst | ack) != syn ct state new counter jump INVDROP"); err != nil {
		return err
	}

	return nil
}
