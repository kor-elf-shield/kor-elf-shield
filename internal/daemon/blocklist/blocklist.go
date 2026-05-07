package blocklist

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/pkg/filesystem"
)

type Blocklist interface {
	Names() []string
	NftReload(blocks map[string]block.Blocklist) error
	Run()
	Close() error
}

type updateSource struct {
	forcedly bool
	source   *SourceConfig
}

type blocklist struct {
	pathDir             string
	sources             []*SourceConfig
	blocklistRepository repository.BlocklistRepository
	logger              log.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	nftBlocklists map[string]block.Blocklist
	mu            sync.Mutex

	launchChannel chan updateSource
}

func New(config Config, ctx context.Context, logger log.Logger) (Blocklist, error) {
	if config.PathDir == "" {
		return nil, fmt.Errorf("pathDir is empty")
	}

	if err := filesystem.EnsureDir(config.PathDir); err != nil {
		return nil, err
	}

	return &blocklist{
		pathDir:             config.PathDir,
		sources:             config.Sources,
		blocklistRepository: config.BlocklistRepository,
		logger:              logger,
		ctx:                 ctx,

		nftBlocklists: map[string]block.Blocklist{},
		mu:            sync.Mutex{},

		launchChannel: make(chan updateSource, 50),
	}, nil
}

func (b *blocklist) Names() []string {
	var names []string
	for _, source := range b.sources {
		if source.Name != "" {
			names = append(names, source.Name)
		}
	}
	return names
}

func (b *blocklist) NftReload(blocks map[string]block.Blocklist) error {
	b.logger.Debug("Reload blocklist")

	b.mu.Lock()
	b.nftBlocklists = blocks
	b.mu.Unlock()

	for _, source := range b.sources {
		if nftBlocklist, ok := b.nftBlocklists[source.Name]; ok {
			if listEntity, err := b.blocklistRepository.Get(source.Name); err != nil {
				b.logger.Error(fmt.Sprintf("Failed to get blocklist %s: %s", source.Name, err))
			} else if b.isFresh(source, listEntity) {
				file, err := b.pathFile(source)
				if err != nil {
					b.logger.Error(fmt.Sprintf("Failed to get blocklist file path: %s", err))
					continue
				}

				if err := nftBlocklist.ReplaceElementsWithFile(file); err != nil {
					b.logger.Error(fmt.Sprintf("Failed to replace elements with file %s: %s", file, err))
					continue
				}
			}
		} else {
			b.logger.Error(fmt.Sprintf("NFTables sets blocklist %s not found", source.Name))
		}
	}

	return nil
}

func (b *blocklist) Run() {
	b.logger.Debug("Starting blocklist")
	if b.cancel != nil {
		// already started
		b.logger.Warn("Blocklist already started")
		return
	}
	b.ctx, b.cancel = context.WithCancel(b.ctx)
	go b.processUpdateData(b.ctx)

	for _, src := range b.sources {
		if src == nil || src.Name == "" {
			continue
		}

		interval := src.Interval
		if interval <= 0 {
			interval = 5 * time.Minute // дефолт
		}

		b.wg.Add(1)
		go b.runSourceWorker(src, interval)
	}
}

func (b *blocklist) runSourceWorker(sourceConfig *SourceConfig, interval time.Duration) {
	defer b.wg.Done()

	b.launchChannel <- updateSource{
		forcedly: false,
		source:   sourceConfig,
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-b.ctx.Done():
			b.logger.Debug(fmt.Sprintf("source %s stopped", sourceConfig.Name))
			return
		case <-ticker.C:
			b.logger.Debug(fmt.Sprintf("source %s tick", sourceConfig.Name))
			b.launchChannel <- updateSource{
				forcedly: true,
				source:   sourceConfig,
			}
		}
	}
}

func (b *blocklist) processUpdateData(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case updSource, ok := <-b.launchChannel:
			if !ok {
				// Channel closed
				return
			}

			if updSource.forcedly {
				b.refreshSource(updSource.source)
				continue
			}

			if listEntity, err := b.blocklistRepository.Get(updSource.source.Name); err != nil {
				b.logger.Error(fmt.Sprintf("Failed to get blocklist %s: %s", updSource.source.Name, err))
				continue
			} else if b.isFresh(updSource.source, listEntity) {
				b.logger.Debug(fmt.Sprintf("blocklist %s is fresh", updSource.source.Name))
				continue
			}

			b.refreshSource(updSource.source)
		}
	}
}

func (b *blocklist) refreshSource(sourceConfig *SourceConfig) {
	ipsV4, ipsV6, err := sourceConfig.Source.Get()
	if err != nil {
		b.logger.Error(fmt.Sprintf("Failed to get IPs from source %s: %s", sourceConfig.Name, err))
		return
	}

	if nftBlocklist, ok := b.nftBlocklists[sourceConfig.Name]; ok {
		listEntity := &entity.Blocklist{
			UpdatedAtUnix: time.Now().Unix(),
		}

		if err := filesystem.EnsureDir(b.pathDir); err != nil {
			b.logger.Error(fmt.Sprintf("Failed to ensure dir: %s", err))
		}
		file, err := b.pathFile(sourceConfig)
		if err != nil {
			b.logger.Error(fmt.Sprintf("Failed to get blocklist file path: %s", err))
			return
		}

		if err := nftBlocklist.ReplaceElements(ipsV4, ipsV6, file); err != nil {
			b.logger.Error(fmt.Sprintf("Failed to replace elements: %s", err))
		}
		listEntity.Checksum, err = filesystem.FileChecksum(file)
		if err != nil {
			b.logger.Error(fmt.Sprintf("Failed to calculate checksum for %s: %s", file, err))
			return
		}

		if err := b.blocklistRepository.Update(sourceConfig.Name, listEntity); err != nil {
			b.logger.Error(fmt.Sprintf("Failed to update blocklist %s: %s", sourceConfig.Name, err))
		}

	} else {
		b.logger.Error(fmt.Sprintf("NFTables sets blocklist %s not found", sourceConfig.Name))
		return
	}

	b.logger.Debug(fmt.Sprintf("refresh blocklist from %s", sourceConfig.Name))
}

func (b *blocklist) Close() error {
	b.logger.Debug("Stopping blocklist")
	if b.cancel != nil {
		b.cancel()
		b.wg.Wait()
		b.cancel = nil
	}
	close(b.launchChannel)
	return nil
}

func (b *blocklist) pathFile(sourceConfig *SourceConfig) (string, error) {
	if sourceConfig == nil {
		return "", fmt.Errorf("sourceConfig is nil")
	}

	if sourceConfig.Name == "" {
		return "", fmt.Errorf("sourceConfig.Name is empty")
	}

	return strings.TrimRight(b.pathDir, "/") + "/" + sourceConfig.Name + ".nft", nil
}

func (b *blocklist) isFresh(sourceConfig *SourceConfig, listEntity *entity.Blocklist) bool {
	if !listEntity.IsFresh(sourceConfig.Interval) {
		return false
	}

	file, err := b.pathFile(sourceConfig)
	if err != nil {
		b.logger.Error(fmt.Sprintf("Failed to get blocklist file path: %s", err))
		return false
	}
	if !filesystem.FileExists(file) {
		b.logger.Warn(fmt.Sprintf("Blocklist file %s not found", file))
		return false
	}

	fileChecksum, err := filesystem.FileChecksum(file)
	if err != nil {
		b.logger.Error(fmt.Sprintf("Failed to calculate checksum for %s: %s", file, err))
		return false
	}
	if listEntity.Checksum != fileChecksum {
		b.logger.Error(fmt.Sprintf("Blocklist file %s checksum is not equal to database checksum", file))
		return false
	}

	return true
}
