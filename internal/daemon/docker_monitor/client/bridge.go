package client

import (
	"fmt"
	"strings"
)

func (d *docker) Bridges() ([]string, error) {
	args := []string{"network", "ls", "-q", "--filter", "Driver=bridge"}
	result, err := d.command(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker bridge names: %s", err.Error())
	}

	output := strings.TrimSpace(string(result))
	if output == "" {
		return []string{}, nil
	}

	lines := strings.Split(output, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}

	return lines, nil
}

func (d *docker) BridgeNames() ([]string, error) {
	bridges, err := d.Bridges()
	if err != nil {
		return nil, err
	}

	var names []string

	for _, bridge := range bridges {
		bridgeName, err := d.BridgeName(bridge)
		if err != nil {
			d.logger.Error(err.Error())
			continue
		}
		names = append(names, bridgeName)
	}

	return names, nil
}

func (d *docker) BridgeName(bridgeID string) (string, error) {
	format := fmt.Sprintf(`{{"br-%s" | or (index .Options "com.docker.network.bridge.name")}}`, bridgeID)
	args := []string{"network", "inspect", "-f", format, bridgeID}
	result, err := d.command(args...)
	if err != nil {
		return "", fmt.Errorf("failed to get bridge name: %s", err.Error())
	}
	return strings.TrimSpace(string(result)), nil
}

func (d *docker) BridgeSubnet(bridgeID string) (string, error) {
	format := fmt.Sprintf(`{{range .IPAM.Config}}{{.Subnet}}{{end}}`)
	args := []string{"network", "inspect", "-f", format, bridgeID}
	result, err := d.command(args...)
	if err != nil {
		return "", fmt.Errorf("failed to get bridge subnet: %s", err.Error())
	}
	return strings.TrimSpace(string(result)), nil
}
