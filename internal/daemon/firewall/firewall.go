package firewall

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	nftables "git.kor-elf.net/kor-elf-shield/go-nftables-client"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/blocklist"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor"
	dockerFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/blocking"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/table"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/reload"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/filesystem"
)

type API interface {
	// Reload Clear all rules and set new rules.
	Reload(appVersion string, isCache bool) error

	// SavesRules Save rules to file.
	SavesRules()

	// ClearRules Clear all rules.
	ClearRules()

	// BlockIP Block IP address.
	BlockIP(blockIP blocking.BlockIP) (bool, error)

	// BlockIPWithPorts Block IP address with ports.
	BlockIPWithPorts(blockIP blocking.BlockIPWithPorts) (bool, error)

	// UnblockAllIPs Unblock all IP addresses.
	UnblockAllIPs() error

	// UnblockIP Unblock IP address.
	UnblockIP(ip net.IP) error

	// ClearDBData Clear all data from DB
	ClearDBData() error

	// DockerSupport Return true if docker support
	DockerSupport() bool
}

type firewall struct {
	nft             nftFirewall.NFT
	table           table.Table
	logger          log.Logger
	config          *config.Config
	blockingService blocking.API
	docker          docker_monitor.Docker
	blocklist       blocklist.Blocklist
	dataDir         string
	metadataService Metadata

	mu sync.Mutex
}

func New(
	pathNFT string,
	blockingService blocking.API,
	logger log.Logger,
	config config.Config,
	docker docker_monitor.Docker,
	blocklist blocklist.Blocklist,
	dataDir string,
	metadataService Metadata,
) (API, error) {
	nftClient, err := nftables.NewWithPath(pathNFT)
	if err != nil {
		return nil, fmt.Errorf("failed to create nft client: %w %s", err, pathNFT)
	}

	return &firewall{
		nft:             nftFirewall.New(nftClient, strings.TrimRight(dataDir, "/")+"/tmp"),
		logger:          logger,
		config:          &config,
		blockingService: blockingService,
		docker:          docker,
		blocklist:       blocklist,
		dataDir:         dataDir,
		metadataService: metadataService,

		mu: sync.Mutex{},
	}, nil
}

func (f *firewall) Reload(appVersion string, isCache bool) error {
	f.logger.Debug("Reload nftables rules")

	nftReload := reload.New(f.nft, f.logger, f.config)
	blocklistNames := f.blocklist.Names()

	var nftTable table.Table
	var err error
	if isCache {
		file := f.pathFileCacheNFT()
		nftTable, err = nftReload.RunWithCache(
			file,
			f.isValidCacheFile(appVersion),
			blocklistNames,
		)
		if err != nil {
			return err
		}

		checksum, err := filesystem.FileChecksum(file)
		if err != nil {
			f.logger.Error(fmt.Sprintf("Failed to calculate checksum for %s: %s", file, err))
		} else if err := f.metadataService.UpdateMetadata(appVersion, checksum); err != nil {
			f.logger.Error(fmt.Sprintf("Failed to update metadata: %s", err))
		}
	} else {
		nftTable, err = nftReload.Run(blocklistNames)
		if err != nil {
			return err
		}
	}

	f.mu.Lock()
	f.table = nftTable
	f.mu.Unlock()

	if f.config.Options.DockerSupport && nftTable.DockerChains() != nil {
		nftDocker := dockerFirewall.NewNFT(f.nft, nftTable.DockerChains())
		if err := f.docker.NftReload(nftDocker); err != nil {
			return err
		}
	}

	if err := f.blockingService.NftReload(f.nft, nftTable.BlockList().ListIP(), nftTable.BlockList().ListIPWithPort()); err != nil {
		return err
	}

	if err := f.blocklist.NftReload(nftTable.BlockList().Blocks()); err != nil {
		f.logger.Error(fmt.Sprintf("Failed to reload blocklist: %s", err))
	}

	f.logger.Debug("Reload nftables rules done")
	return nil
}

func (f *firewall) ClearRules() {
	f.logger.Debug("Clear nftables rules")

	switch f.config.Options.ClearMode {
	case config.ClearModeGlobal:
		if err := f.nft.NFT().Clear(); err != nil {
			f.logger.Error(fmt.Sprintf("Failed to clear rules: %s", err))
		}
		break
	case config.ClearModeOwn:
		if f.table == nil {
			f.logger.Error("table is nil")
			return
		}
		if err := f.table.Clear(); err != nil {
			f.logger.Error(fmt.Sprintf("Failed to clear rules: %s", err))
		}
		break
	}

	f.logger.Debug("Clear nftables rules done")
}

func (f *firewall) UnblockAllIPs() error {
	return f.blockingService.UnblockAllIPs()
}

func (f *firewall) UnblockIP(ip net.IP) error {
	return f.blockingService.UnblockIP(ip)
}

func (f *firewall) ClearDBData() error {
	return f.blockingService.ClearDBData()
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
	output, err := f.nft.NFT().Command().RunWithOutput(args...)
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

func (f *firewall) BlockIP(blockIP blocking.BlockIP) (bool, error) {
	isBanned, err := f.blockingService.BlockIP(blockIP)

	if err != nil {
		f.logger.Warn(fmt.Sprintf("Failed to block ip %s: %s", blockIP.IP.String(), err))
	}
	return isBanned, err
}

func (f *firewall) BlockIPWithPorts(blockIP blocking.BlockIPWithPorts) (bool, error) {
	isBanned, err := f.blockingService.BlockIPWithPorts(blockIP)

	if err != nil {
		f.logger.Warn(fmt.Sprintf("Failed to block ip %s: %s", blockIP.IP.String(), err))
	}
	return isBanned, err
}

func (f *firewall) DockerSupport() bool {
	return f.config.Options.DockerSupport
}

func (f *firewall) pathFileCacheNFT() string {
	return strings.TrimRight(f.dataDir, "/") + "/nftables.nft"
}

func (f *firewall) isValidCacheFile(appVersion string) bool {
	fileNFT := f.pathFileCacheNFT()
	if !filesystem.FileExists(fileNFT) {
		return false
	}

	metadata, err := f.metadataService.Metadata()
	if err != nil {
		f.logger.Error(fmt.Sprintf("Failed to get metadata: %s", err))
		return false
	}

	if metadata.Version != appVersion {
		f.logger.Warn(fmt.Sprintf("App version %s is not equal to metadata version %s", appVersion, metadata.Version))
		return false
	}

	if metadata.Checksum == "" {
		return false
	}

	checksum, err := filesystem.FileChecksum(fileNFT)
	if err != nil {
		f.logger.Error(fmt.Sprintf("Failed to calculate checksum for %s: %s", fileNFT, err))
		return false
	}

	if checksum != metadata.Checksum {
		f.logger.Warn(fmt.Sprintf("Checksum of %s is not equal to metadata checksum", fileNFT))
		return false
	}

	return true
}
