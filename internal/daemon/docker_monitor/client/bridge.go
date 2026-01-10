package client

import (
	"encoding/json"
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

func (d *docker) BridgeInfo(bridgeID string) (DockerBridgeInspect, error) {
	args := []string{"network", "inspect", bridgeID}
	result, err := d.command(args...)
	if err != nil {
		return DockerBridgeInspect{}, fmt.Errorf("failed to get bridge name: %s", err.Error())
	}

	var info []DockerBridgeInspect
	if err := json.Unmarshal(result, &info); err != nil {
		return DockerBridgeInspect{}, err
	}

	return info[0], nil
}
