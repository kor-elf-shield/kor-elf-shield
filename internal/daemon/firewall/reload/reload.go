package reload

import (
	"fmt"

	nftChain "git.kor-elf.net/kor-elf-shield/go-nftables-client/chain"
	"git.kor-elf.net/kor-elf-shield/go-nftables-client/family"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/config"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/chain"
	dataTable "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/table"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"

	nft "git.kor-elf.net/kor-elf-shield/go-nftables-client/contract"
)

type Reload interface {
	Run() (dataTable.Table, error)
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

func (r *reload) Run() (dataTable.Table, error) {
	var dockerChains firewall.NFTDockerChains
	if r.config.Options.DockerSupport {
		dockerChains = firewall.NewNFTChains(r.nft, r.table.family, r.table.name)
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

	if err := r.clear(batchBuilder); err != nil {
		return nil, err
	}
	packetFilter, err := r.packetFilter(batchBuilder)
	if err != nil {
		return nil, err
	}

	blockList, err := r.input(batchBuilder, packetFilter)
	if err != nil {
		return nil, err
	}

	if err := r.output(batchBuilder, packetFilter); err != nil {
		return nil, err
	}
	if err := r.forward(batchBuilder, dockerChains); err != nil {
		return nil, err
	}

	if err := r.nft.RunBatch(batchBuilder); err != nil {
		return nil, err
	}

	return dataTable.New(blockList, dockerChains), nil
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
		if err := builder.Table().Clear(r.table.family, r.table.name); err != nil {
			return fmt.Errorf("failed to clear table: %w", err)
		}
		break
	default:
		return fmt.Errorf("unknown clear mode: %d", r.config.Options.ClearMode)
	}

	return nil
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
