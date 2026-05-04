package reload

import (
	"fmt"

	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	dataTable "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/table"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
)

const blockedIP = "block_ip"
const blockedIPWithPort = "block_ip_with_port"

type Reload interface {
	RunWithCache(file string, isValidCacheFile bool, blockListNames []string) (dataTable.Table, error)
	Run(blockListNames []string) (dataTable.Table, error)
}

type table struct {
	name   string
	family family.Type
}

type reload struct {
	nft    nftFirewall.NFT
	logger log.Logger
	config *config.Config
	table  *table
}

func New(nft nftFirewall.NFT, logger log.Logger, config *config.Config) Reload {
	return &reload{
		nft:    nft,
		logger: logger,
		config: config,
		table: &table{
			name:   config.MetadataNaming.TableName,
			family: family.INET,
		},
	}
}

func (r *reload) RunWithCache(file string, isValidCacheFile bool, blockListNames []string) (dataTable.Table, error) {
	if file == "" {
		r.logger.Warn("file is empty, using default reload")
		return r.Run(blockListNames)
	}
	if isValidCacheFile {
		r.logger.Debug("use cache file")
		table, err := r.loadCache(file, blockListNames)
		if err == nil {
			return table, nil
		}
		r.logger.Warn(fmt.Sprintf("failed to load cache file: %s", err))
	}

	batchBuilder, err := r.nft.NewBuildBatch()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			r.logger.Warn(err.Error())
		}
	}()

	table, err := r.reload(batchBuilder, blockListNames)
	if err != nil {
		return nil, err
	}

	if err := r.nft.RunBatchAndMoveFile(batchBuilder, file); err != nil {
		return nil, err
	}

	return table, nil
}

func (r *reload) Run(blockListNames []string) (dataTable.Table, error) {
	batchBuilder, err := r.nft.NewBuildBatch()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			r.logger.Warn(err.Error())
		}
	}()

	table, err := r.reload(batchBuilder, blockListNames)
	if err != nil {
		return nil, err
	}

	if err := r.nft.RunBatch(batchBuilder); err != nil {
		return nil, err
	}

	return table, nil
}

func (r *reload) reload(batchBuilder nft.BatchBuilder, blockListNames []string) (dataTable.Table, error) {
	var dockerChains firewall.NFTDockerChains
	if r.config.Options.DockerSupport {
		dockerChains = firewall.NewNFTChains(r.nft, r.table.family, r.table.name)
	}

	if err := r.clear(batchBuilder); err != nil {
		return nil, err
	}
	packetFilter, err := r.packetFilter(batchBuilder)
	if err != nil {
		return nil, err
	}

	blockList, err := r.input(batchBuilder, packetFilter, blockListNames)
	if err != nil {
		return nil, err
	}

	if err := r.output(batchBuilder, packetFilter); err != nil {
		return nil, err
	}
	if err := r.forward(batchBuilder, dockerChains); err != nil {
		return nil, err
	}

	return dataTable.New(
		r.nft, r.table.family, r.table.name,
		blockList, dockerChains,
	), nil
}

func (r *reload) clear(builder nft.BatchBuilder) error {
	switch r.config.Options.ClearMode {
	case config.ClearModeGlobal:
		if err := builder.Clear(); err != nil {
			return fmt.Errorf("failed to clear global rules: %w", err)
		}
		if err := builder.Table().Add(r.table.family, r.table.name); err != nil {
			return fmt.Errorf("failed to add table: %w", err)
		}
		break
	case config.ClearModeOwn:
		if err := builder.Table().Add(r.table.family, r.table.name); err != nil {
			return fmt.Errorf("failed to add table: %w", err)
		}
		// clear does not clean completely
		if err := builder.Table().Delete(r.table.family, r.table.name); err != nil {
			return fmt.Errorf("failed to clear table: %w", err)
		}
		if err := builder.Table().Add(r.table.family, r.table.name); err != nil {
			return fmt.Errorf("failed to add table: %w", err)
		}
		break
	default:
		return fmt.Errorf("unknown clear mode: %d", r.config.Options.ClearMode)
	}

	return nil
}

func (r *reload) loadCache(file string, blockListNames []string) (dataTable.Table, error) {
	args := []string{"-f", file}
	if err := r.nft.NFT().Command().Run(args...); err != nil {
		return nil, err
	}

	var dockerChains firewall.NFTDockerChains
	if r.config.Options.DockerSupport {
		dockerChains = firewall.NewNFTChains(r.nft, r.table.family, r.table.name)
	}

	blocks := make(map[string]block.Blocklist)
	for _, blockListName := range blockListNames {
		if blockListName == "" {
			continue
		}
		blockList := block.NewBlocklistWithoutCommand(r.nft, r.table.family, r.table.name, getBlocklistName(blockListName))
		blocks[blockListName] = blockList
	}

	listBlockedIP := block.NewListIPWithoutCommand(r.nft, r.table.family, r.table.name, blockedIP)
	listBlockedIPWithPort := block.NewListIPWithPortWithoutCommand(r.nft, r.table.family, r.table.name, blockedIPWithPort)

	tableBlocklist := dataTable.NewBlockList(listBlockedIP, listBlockedIPWithPort, blocks)
	return dataTable.New(
		r.nft, r.table.family, r.table.name,
		tableBlocklist, dockerChains,
	), nil
}

func (r *reload) addChain(builder nft.BatchBuilder, chainName string, baseChain nftChain.ChainOptions) error {
	return builder.Chain().Add(r.table.family, r.table.name, chainName, baseChain)
}

func (r *reload) addRule(builder nft.BatchBuilder, chainName string, rule string) error {
	return builder.Rule().Add(r.table.family, r.table.name, chainName, rule)
}

func (r *reload) addChainWithReturn(builder nft.BatchBuilder, chainName string, baseChain nftChain.ChainOptions) (chain.Chain, error) {
	return chain.NewBatchChainWithOptions(builder, r.table.family, r.table.name, chainName, baseChain)
}

func getBlocklistName(blockListName string) string {
	return "blocklist_" + blockListName
}
