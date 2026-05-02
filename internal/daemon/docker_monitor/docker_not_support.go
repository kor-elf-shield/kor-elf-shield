package docker_monitor

import (
	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/daemon/docker_monitor/firewall"
)

type dockerNotSupport struct{}

func NewDockerNotSupport() Docker {
	return &dockerNotSupport{}
}

func (d *dockerNotSupport) NftReload(_ firewall.NFTDocker) error {
	return nil
}

func (d *dockerNotSupport) Run() {

}

func (d *dockerNotSupport) Close() error {
	return nil
}
