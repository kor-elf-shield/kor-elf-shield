package blocking

import (
	"fmt"
	"net"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/chain/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type API interface {
	NftReload(blockListIP block.ListIP) error
	BlockIP(blockIP BlockIP) (bool, error)
	ClearDBData() error
}

type blocking struct {
	blockingRepository repository.BlockingRepository
	blockListIP        block.ListIP
	logger             log.Logger

	mu sync.Mutex
}

type BlockIP struct {
	IP          net.IP
	TimeSeconds uint32
	Reason      string
}

func New(blockingRepository repository.BlockingRepository, logger log.Logger) API {
	return &blocking{
		blockingRepository: blockingRepository,
		logger:             logger,
		mu:                 sync.Mutex{},
	}
}

func (b *blocking) NftReload(blockListIP block.ListIP) error {
	b.mu.Lock()
	b.blockListIP = blockListIP
	b.mu.Unlock()

	isExpiredEntries := false
	nowUnix := time.Now().Unix()
	err := b.blockingRepository.List(func(e entity.Blocking) error {
		ip := net.ParseIP(e.IP)
		if ip == nil {
			b.logger.Error(fmt.Sprintf("Failed to parse IP address: %s", e.IP))
			return nil
		}

		banSeconds := uint32(0)
		if e.ExpireAtUnix > 0 {
			if e.ExpireAtUnix < nowUnix {
				isExpiredEntries = true
				return nil
			}
			banSeconds = uint32(e.ExpireAtUnix - nowUnix)
			fmt.Printf("Now %s ExpireAtUnix %d banSeconds %d\n", time.Now().Format(time.RFC3339), e.ExpireAtUnix, banSeconds)
		}

		if err := b.blockListIP.AddIP(ip, banSeconds); err != nil {
			b.logger.Error(fmt.Sprintf("Failed to add IP %s to block list: %s", ip.String(), err))
			return nil
		}

		return nil
	})

	if isExpiredEntries {
		go func() {
			deleteCount, err := b.blockingRepository.DeleteExpired(100)
			if err != nil {
				b.logger.Error(fmt.Sprintf("Failed to delete expired entries from database: %s", err))
			}
			b.logger.Debug(fmt.Sprintf("Deleted %d expired entries from database", deleteCount))
		}()
	}

	return err
}

func (b *blocking) BlockIP(blockIP BlockIP) (bool, error) {
	if blockIP.IP.IsLoopback() {
		return false, fmt.Errorf("loopback IP address %s cannot be blocked", blockIP.IP.String())
	}

	if err := b.blockListIP.AddIP(blockIP.IP, blockIP.TimeSeconds); err != nil {
		return false, err
	}

	expireAtUnix := int64(0)
	if blockIP.TimeSeconds > 0 {
		expire := time.Now().Add(time.Duration(int64(blockIP.TimeSeconds)) * time.Second)
		expireAtUnix = expire.Unix()
	}
	data := entity.Blocking{
		IP:           blockIP.IP.String(),
		ExpireAtUnix: expireAtUnix,
		Reason:       blockIP.Reason,
	}
	if err := b.blockingRepository.Add(data); err != nil {
		return true, fmt.Errorf("the IP is blocked, but not recorded in the database. Failed to add IP %s to database: %w", blockIP.IP.String(), err)
	}

	return true, nil
}

func (b *blocking) ClearDBData() error {
	return b.blockingRepository.Clear()
}
