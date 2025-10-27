package pkg

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type resolvData struct {
	addresses []string
}

var Resolv = resolvData{}

func (r *resolvData) Addresses() ([]string, error) {
	if len(r.addresses) > 0 {
		return r.addresses, nil
	}

	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return []string{}, fmt.Errorf("Failed to open /etc/resolv.conf: %s", err)
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	addresses := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				addresses = append(addresses, fields[1])
			}
		}
	}
	r.addresses = addresses
	return r.addresses, nil
}
