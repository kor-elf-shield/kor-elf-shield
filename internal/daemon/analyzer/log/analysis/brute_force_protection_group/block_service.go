package brute_force_protection_group

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/firewall/blocking"

type BlockService interface {
	BlockIP(blockIP blocking.BlockIP) (bool, error)
	BlockIPWithPorts(blockIP blocking.BlockIPWithPorts) (bool, error)
}

type BlockIPFunc func(blockIP blocking.BlockIP) (bool, error)
type BlockIPWithPortsFunc func(blockIP blocking.BlockIPWithPorts) (bool, error)

type blockService struct {
	blockIPFunc          BlockIPFunc
	blockIPWithPortsFunc BlockIPWithPortsFunc
}

func NewBlockService(blockIPFunc BlockIPFunc, blockIPWithPortsFunc BlockIPWithPortsFunc) BlockService {
	return &blockService{
		blockIPFunc:          blockIPFunc,
		blockIPWithPortsFunc: blockIPWithPortsFunc,
	}
}

func (b *blockService) BlockIP(blockIP blocking.BlockIP) (bool, error) {
	return b.blockIPFunc(blockIP)
}

func (b *blockService) BlockIPWithPorts(blockIP blocking.BlockIPWithPorts) (bool, error) {
	return b.blockIPWithPortsFunc(blockIP)
}
