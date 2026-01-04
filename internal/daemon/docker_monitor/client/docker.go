package client

import (
	"context"
	"fmt"
	"os/exec"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Docker interface {
	FetchBridges() (Bridges, error)
	FetchContainers(bridgeID string) (Containers, error)

	Bridges() ([]string, error)
	BridgeNames() ([]string, error)
	BridgeName(bridgeID string) (string, error)
	BridgeSubnet(bridgeID string) (string, error)

	Containers(bridgeID string) ([]string, error)
	ContainerNetworks(containerID string) (DockerContainerInspect, error)
}

type docker struct {
	path   string
	ctx    context.Context
	logger log.Logger
}

func NewDocker(path string, ctx context.Context, logger log.Logger) Docker {
	return &docker{
		path:   path,
		ctx:    ctx,
		logger: logger,
	}
}

func (d *docker) FetchBridges() (Bridges, error) {
	bridges := Bridges{}
	list, err := d.Bridges()
	if err != nil {
		return nil, err
	}

	for _, bridgeId := range list {
		bridgeName, err := d.BridgeName(bridgeId)
		if err != nil {
			d.logger.Error(err.Error())
			continue
		}

		bridgeSubnet, err := d.BridgeSubnet(bridgeId)
		if err != nil {
			d.logger.Error(err.Error())
			continue
		}

		var containers Containers
		containers, err = d.FetchContainers(bridgeId)
		if err != nil {
			d.logger.Error(err.Error())
		}

		bridges = append(bridges, Bridge{
			ID:         bridgeId,
			Name:       bridgeName,
			Subnet:     bridgeSubnet,
			Containers: containers,
		})
	}

	return bridges, nil
}

func (d *docker) FetchContainers(bridgeID string) (Containers, error) {
	containers := Containers{}

	list, err := d.Containers(bridgeID)
	if err != nil {
		return nil, err
	}
	for _, containerID := range list {
		info, err := d.ContainerNetworks(containerID)
		if err != nil {
			d.logger.Error(err.Error())
			continue
		}

		networks := ContainerNetworks{
			IPAddresses: []IPInfo{},
			Ports:       d.parsePorts(info),
		}
		for _, networkData := range info.NetworkSettings.Networks {
			if networkData.IPAddress != "" {
				ipVesion, err := ipVersion(networkData.IPAddress)
				if err != nil {
					d.logger.Error(err.Error())
					continue
				}
				networks.IPAddresses = append(networks.IPAddresses, IPInfo{Address: networkData.IPAddress, Version: ipVesion})
			}
		}

		containers = append(containers, Container{
			ID:       containerID,
			Networks: networks,
		})
	}

	return containers, nil
}

func (d *docker) command(args ...string) ([]byte, error) {
	cmd := exec.CommandContext(d.ctx, d.path, args...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(string(result))
	}
	return result, nil
}
