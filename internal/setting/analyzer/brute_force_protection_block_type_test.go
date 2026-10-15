package analyzer

import (
	"fmt"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/testutil/fakers"
)

func fakeRandBruteForceProtectionBlockTypeAndPorts() (blockType string, ports []string) {
	blockType = fakeRandBruteForceProtectionBlockType()
	if blockType == "ip_port" {
		ports = fakeRandBruteForceProtectionPorts()
	}

	return blockType, ports
}

func fakeRandBruteForceProtectionBlockType() string {
	items := []string{"", "ip", "ip_port"}

	return fakers.RandItem(items)
}

func fakeRandBruteForceProtectionPorts() []string {
	var ports []string
	for i := fakers.RandInt(1, 5); i > 0; i-- {
		port := fakers.RandInt(0, 65535)
		if isProtocol := fakers.RandBool(); !isProtocol {
			ports = append(ports, fmt.Sprintf("%d", port))
			continue
		}

		protocol := fakers.RandItem([]string{"tcp", "udp"})

		ports = append(ports, fmt.Sprintf("%d/%s", port, protocol))
	}

	return ports
}
