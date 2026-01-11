package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"git.kor-elf.net/kor-elf-shield/kor-elf-shield/internal/log"
)

type Docker interface {
	FetchBridges() (Bridges, error)
	FetchBridge(bridgeID string) (Bridge, error)
	FetchContainers(bridgeID string) (Containers, error)
	FetchContainer(containerID string) (Container, error)

	Events() <-chan Event
	EventsClose() error
}

type docker struct {
	path   string
	ctx    context.Context
	logger log.Logger

	cmd *exec.Cmd
	mu  sync.Mutex
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
	list, err := d.bridges()
	if err != nil {
		return nil, err
	}

	for _, bridgeId := range list {
		bridge, err := d.FetchBridge(bridgeId)
		if err != nil {
			d.logger.Error(err.Error())
			continue
		}

		bridges = append(bridges, bridge)
	}

	return bridges, nil
}

func (d *docker) FetchBridge(bridgeID string) (Bridge, error) {

	bridgeInfo, err := d.bridgeInfo(bridgeID)
	if err != nil {
		return Bridge{}, err
	}

	var containers Containers
	containers, err = d.FetchContainers(bridgeID)
	if err != nil {
		d.logger.Error(err.Error())
	}

	bridgeName := fmt.Sprintf("br-%s", bridgeID)
	if bridgeInfo.Options.Name != "" {
		bridgeName = bridgeInfo.Options.Name
	}

	var bridgeSubnet []string
	if bridgeInfo.IPAM.Config != nil {
		for _, config := range bridgeInfo.IPAM.Config {
			bridgeSubnet = append(bridgeSubnet, config.Subnet)
		}
	}

	return Bridge{
		ID:         bridgeInfo.ID,
		Name:       bridgeName,
		Subnets:    bridgeSubnet,
		Containers: containers,
	}, nil
}

func (d *docker) FetchContainers(bridgeID string) (Containers, error) {
	containers := Containers{}

	list, err := d.containers(bridgeID)
	if err != nil {
		return nil, err
	}
	for _, containerID := range list {
		container, err := d.FetchContainer(containerID)
		if err != nil {
			d.logger.Error(err.Error())
			continue
		}

		containers = append(containers, container)
	}

	return containers, nil
}

func (d *docker) FetchContainer(containerID string) (Container, error) {
	info, err := d.containerNetworks(containerID)
	if err != nil {
		return Container{}, err
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

	return Container{
		ID:       containerID,
		Networks: networks,
	}, nil
}

func (d *docker) command(args ...string) ([]byte, error) {
	cmd := exec.CommandContext(d.ctx, d.path, args...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(string(result))
	}
	return result, nil
}

func (d *docker) Events() <-chan Event {
	eventsChan := make(chan Event)

	d.logger.Debug("Starting docker monitor")
	go func() {
		defer close(eventsChan)
		for {
			select {
			case <-d.ctx.Done():
				return
			default:
				if err := d.watch(eventsChan); err != nil {
					d.logger.Error(fmt.Sprintf("Docker monitor exited with error: %v", err))
				}

				// Pause before restarting to avoid CPU load during persistent errors
				select {
				case <-d.ctx.Done():
					return
				case <-time.After(15 * time.Second):
					d.logger.Warn("Docker connection lost. Restarting in 15s...")
					continue
				}
			}
		}
	}()

	return eventsChan
}

func (d *docker) watch(eventsChan chan Event) error {
	args := []string{
		"events",

		"--filter", "type=container",
		"--filter", "event=start",
		"--filter", "event=die",

		"--filter", "type=network",
		"--filter", "event=create",
		"--filter", "event=destroy",

		"--format",
		"{{json .}}",
	}
	cmd := exec.CommandContext(d.ctx, d.path, args...)
	d.mu.Lock()
	d.cmd = cmd
	d.mu.Unlock()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return err
		}

		if scanner.Text() == "" {
			return fmt.Errorf("empty line")
		}

		var dockerEvent DockerEvent
		if err := json.Unmarshal([]byte(scanner.Text()), &dockerEvent); err != nil {
			return fmt.Errorf("failed to unmarshal docker event: %v", err)
		}

		if dockerEvent.Type == "" || dockerEvent.Action == "" || dockerEvent.Actor.ID == "" {
			continue
		}

		eventsChan <- Event{
			Type:    dockerEvent.Type,
			Action:  dockerEvent.Action,
			ID:      dockerEvent.Actor.ID,
			Message: scanner.Text(),
		}
	}

	return scanner.Err()
}

func (d *docker) EventsClose() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cmd != nil && d.cmd.Process != nil {
		d.logger.Debug("Stopping docker monitor")

		// Force docker monitor to quit on shutdown
		return d.cmd.Process.Kill()
	}

	d.logger.Debug("Docker monitor stopped")

	return nil
}
