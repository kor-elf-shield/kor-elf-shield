package table

import "git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"

type Table interface {
	DockerChains() firewall.NFTDockerChains
}

type table struct {
	dockerChains firewall.NFTDockerChains
}

func New(dockerChains firewall.NFTDockerChains) Table {
	return &table{
		dockerChains: dockerChains,
	}
}

func (t *table) DockerChains() firewall.NFTDockerChains {
	return t.dockerChains
}
