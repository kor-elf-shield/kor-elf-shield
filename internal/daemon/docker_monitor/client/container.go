package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (d *docker) containers(bridgeID string) ([]string, error) {
	args := []string{"ps", "-q", "--no-trunc", "--filter", fmt.Sprintf("network=%s", bridgeID)}
	result, err := d.command(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get docker containers: %s", err.Error())
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

func (d *docker) containerNetworks(containerID string) (DockerContainerInspect, error) {
	result, err := d.command("inspect", containerID)

	if err != nil {
		return DockerContainerInspect{}, err
	}

	var info []DockerContainerInspect
	if err := json.Unmarshal(result, &info); err != nil {
		return DockerContainerInspect{}, err
	}

	if len(info) == 0 {
		return DockerContainerInspect{}, fmt.Errorf("container %s not found", containerID)
	}

	return info[0], nil
}

func (d *docker) parsePorts(info DockerContainerInspect) []ContainerPort {
	var ports []ContainerPort

	for containerPortFull, hostConfigs := range info.NetworkSettings.Ports {

		parts := strings.Split(containerPortFull, "/")
		portNum := parts[0]
		protocol := "tcp" // default
		if len(parts) > 1 {
			protocol = parts[1]
		}

		cp := ContainerPort{
			Port:     portNum,
			Protocol: protocol,
		}

		for _, h := range hostConfigs {
			host := HostPort{
				Port: h.HostPort,
			}

			ipVersion, err := ipVersion(h.HostIp)
			if err != nil {
				d.logger.Error(err.Error())
				continue
			}
			host.IP = IPInfo{Address: h.HostIp, Version: ipVersion}

			cp.HostPort = append(cp.HostPort, host)
		}

		ports = append(ports, cp)
	}

	return ports
}
