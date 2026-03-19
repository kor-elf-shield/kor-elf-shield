package blocklist

import (
	"context"
	"fmt"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type newBlocklist func(name string) (block.Blocklist, error)

type Blocklist interface {
	NftReload(newBlocklist newBlocklist) error
	Run()
	Close() error
}

type updateSource struct {
	forcedly bool
	source   *SourceConfig
}

type blocklist struct {
	Sources             []*SourceConfig
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
	return &blocklist{
		Sources:             config.Sources,
		blocklistRepository: config.BlocklistRepository,
		logger:              logger,
		ctx:                 ctx,

		nftBlocklists: map[string]block.Blocklist{},
		mu:            sync.Mutex{},

		launchChannel: make(chan updateSource, 50),
	}, nil
}

func (b *blocklist) NftReload(newBlocklist newBlocklist) error {
	b.logger.Debug("Reload blocklist")
	for _, source := range b.Sources {
		b.logger.Debug(fmt.Sprintf("Reload blocklist from %s", source.Name))
		if source.Name == "" {
			continue
		}

		nftBlocklist, err := newBlocklist("blocklist_" + source.Name)
		if err != nil {
			b.logger.Error(fmt.Sprintf("Failed to create blocklist: %s", err))
			continue
		}

		b.mu.Lock()
		b.nftBlocklists[source.Name] = nftBlocklist
		b.mu.Unlock()

		if listEntity, err := b.blocklistRepository.Get(source.Name); err != nil {
			b.logger.Error(fmt.Sprintf("Failed to get blocklist %s: %s", source.Name, err))
		} else if listEntity.IsFresh(source.Interval) {
			if err := nftBlocklist.ReplaceElementsIPv4(listEntity.IPsV4); len(listEntity.IPsV4) > 0 && err != nil {
				b.logger.Error(fmt.Sprintf("Failed to replace elements (IPv4): %s", err))
			}

			if err := nftBlocklist.ReplaceElementsIPv6(listEntity.IPsV6); len(listEntity.IPsV6) > 0 && err != nil {
				b.logger.Error(fmt.Sprintf("Failed to replace elements (IPv6): %s", err))
			}
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

	for _, src := range b.Sources {
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
			} else if listEntity.IsFresh(updSource.source.Interval) {
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
			IPsV4:         nil,
			IPsV6:         nil,
		}

		if len(ipsV4) > 0 {
			if err := nftBlocklist.ReplaceElementsIPv4(ipsV4); err != nil {
				b.logger.Error(fmt.Sprintf("Failed to replace elements (IPv4): %s", err))
			} else {
				listEntity.IPsV4 = ipsV4
			}
		}

		if len(ipsV6) > 0 {
			if err := nftBlocklist.ReplaceElementsIPv6(ipsV6); err != nil {
				b.logger.Error(fmt.Sprintf("Failed to replace elements (IPv6): %s", err))
			} else {
				listEntity.IPsV6 = ipsV6
			}
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
