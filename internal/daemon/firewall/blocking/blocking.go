package blocking

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/entity"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/db/repository"
	nftFirewall "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/nft/block"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/types"
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type API interface {
	NftReload(nft nftFirewall.NFT, blockListIP block.ListIP, blockListIPWithPort block.ListIPWithPort) error
	BlockIP(block BlockIP) (bool, error)
	BlockIPWithPorts(block BlockIPWithPorts) (bool, error)
	UnblockAllIPs() error
	UnblockIP(ip net.IP) error
	ClearDBData() error
}

type blocking struct {
	blockingRepository  repository.BlockingRepository
	blockListIP         block.ListIP
	blockListIPWithPort block.ListIPWithPort
	logger              log.Logger

	mu sync.Mutex
}

type BlockIP struct {
	IP          net.IP
	TimeSeconds uint32
	Reason      string
}

type BlockIPWithPorts struct {
	IP          net.IP
	TimeSeconds uint32
	Reason      string
	Ports       []types.L4Port
}

func New(blockingRepository repository.BlockingRepository, logger log.Logger) API {
	return &blocking{
		blockingRepository: blockingRepository,
		logger:             logger,
		mu:                 sync.Mutex{},
	}
}

func (b *blocking) NftReload(nft nftFirewall.NFT, blockListIP block.ListIP, blockListIPWithPort block.ListIPWithPort) error {
	b.mu.Lock()
	b.blockListIP = blockListIP
	b.blockListIPWithPort = blockListIPWithPort
	b.mu.Unlock()

	batchBuilder, err := nft.NewBuildBatch()
	if err != nil {
		return err
	}
	defer func() {
		if err := batchBuilder.Close(); err != nil {
			b.logger.Warn(err.Error())
		}
	}()

	isExpiredEntries := false
	nowUnix := time.Now().Unix()
	err = b.blockingRepository.List(func(e entity.Blocking) error {
		ip := net.ParseIP(e.IP)
		if ip == nil {
			b.logger.Error(fmt.Sprintf("Failed to parse IP address: %s", e.IP))
			return nil
		}

		blockSeconds := uint32(0)
		if e.ExpireAtUnix > 0 {
			if e.ExpireAtUnix < nowUnix {
				isExpiredEntries = true
				return nil
			}
			blockSeconds = uint32(e.ExpireAtUnix - nowUnix)
		}

		if e.IsPorts() {
			l4Ports, err := e.ToL4Ports()
			if err != nil {
				b.logger.Error(fmt.Sprintf("Failed to parse ports: %s", err))
				return nil
			}
			if err := b.blockListIPWithPort.AddBatchIP(batchBuilder, ip, l4Ports, blockSeconds); err != nil {
				b.logger.Error(fmt.Sprintf("Failed to add IP %s to block list: %s", ip.String(), err))
			}

			return nil
		}

		if err := b.blockListIP.AddBatchIP(batchBuilder, ip, blockSeconds); err != nil {
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

	if err != nil {
		return err
	}

	return nft.RunBatchAndMoveFile(batchBuilder, "/var/lib/kor-elf-shield/firewall/tmp/block.nft")
}

func (b *blocking) BlockIP(block BlockIP) (bool, error) {
	if block.IP.IsLoopback() {
		return false, fmt.Errorf("loopback IP address %s cannot be blocked", block.IP.String())
	}

	if err := b.blockListIP.AddIP(block.IP, block.TimeSeconds); err != nil {
		return false, err
	}

	expireAtUnix := int64(0)
	if block.TimeSeconds > 0 {
		expire := time.Now().Add(time.Duration(int64(block.TimeSeconds)) * time.Second)
		expireAtUnix = expire.Unix()
	}
	data := entity.Blocking{
		IP:           block.IP.String(),
		ExpireAtUnix: expireAtUnix,
		Reason:       block.Reason,
	}
	if err := b.blockingRepository.Add(data); err != nil {
		return true, fmt.Errorf("the IP is blocked, but not recorded in the database. Failed to add IP %s to database: %w", block.IP.String(), err)
	}

	return true, nil
}

func (b *blocking) BlockIPWithPorts(block BlockIPWithPorts) (bool, error) {
	if block.IP.IsLoopback() {
		return false, fmt.Errorf("loopback IP address %s cannot be blocked", block.IP.String())
	}

	if err := b.blockListIPWithPort.AddIP(block.IP, block.Ports, block.TimeSeconds); err != nil {
		return false, err
	}

	var l4Ports []entity.BlockingPort
	for _, port := range block.Ports {
		l4Ports = append(l4Ports, entity.BlockingPort{
			Number:   port.Number(),
			Protocol: port.ProtocolString(),
		})
	}

	expireAtUnix := int64(0)
	if block.TimeSeconds > 0 {
		expire := time.Now().Add(time.Duration(int64(block.TimeSeconds)) * time.Second)
		expireAtUnix = expire.Unix()
	}
	data := entity.Blocking{
		IP:           block.IP.String(),
		ExpireAtUnix: expireAtUnix,
		Reason:       block.Reason,
		Ports:        l4Ports,
	}
	if err := b.blockingRepository.Add(data); err != nil {
		return true, fmt.Errorf("the IP is blocked, but not recorded in the database. Failed to add IP %s to database: %w", block.IP.String(), err)
	}

	return true, nil
}

func (b *blocking) UnblockIP(ip net.IP) error {
	err := b.blockingRepository.DeleteByIP(ip, func(e entity.Blocking) error {
		if e.IsPorts() {
			l4Ports, err := e.ToL4Ports()
			if err != nil {
				return err
			}
			return b.removeIPWithPorts(ip, l4Ports)
		}

		if err := b.blockListIP.DeleteIP(ip); err != nil {
			if strings.Contains(err.Error(), "element does not exist") {
				return nil
			}
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *blocking) UnblockAllIPs() error {
	err := b.blockingRepository.List(func(e entity.Blocking) error {
		ip := net.ParseIP(e.IP)
		if ip == nil {
			return fmt.Errorf("failed to parse IP address: %s", e.IP)
		}

		if e.IsPorts() {
			l4Ports, err := e.ToL4Ports()
			if err != nil {
				return err
			}
			for _, port := range l4Ports {
				if err := b.blockListIPWithPort.DeleteIP(ip, port); err != nil {
					if strings.Contains(err.Error(), "element does not exist") ||
						strings.Contains(err.Error(), "Error: Could not process rule: No such file or directory") {
						continue
					}
					return err
				}
			}
		}

		if err := b.blockListIP.DeleteIP(ip); err != nil {
			if strings.Contains(err.Error(), "element does not exist") {
				return nil
			}
			return err
		}
		return nil
	})

	if err != nil {
		_ = b.blockingRepository.Clear()
		return err
	}

	return b.blockingRepository.Clear()
}

func (b *blocking) ClearDBData() error {
	return b.blockingRepository.Clear()
}

func (b *blocking) removeIPWithPorts(ip net.IP, l4Ports []types.L4Port) error {
	for _, port := range l4Ports {
		if err := b.blockListIPWithPort.DeleteIP(ip, port); err != nil {
			if strings.Contains(err.Error(), "element does not exist") ||
				strings.Contains(err.Error(), "Error: Could not process rule: No such file or directory") {
				continue
			}
			return err
		}
	}
	return nil
}
