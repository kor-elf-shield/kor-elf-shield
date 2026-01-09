package client

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"sync"
	"time"

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
		eventsChan <- Event{
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
