package client

import (
	"errors"
	"net"
)

type Event struct {
	Type    string
	Action  string
	ID      string // Full 64-char ID (Actor.ID)
	Message string // debug
}

type DockerEvent struct {
	Type   string `json:"Type"`   // container, network
	Action string `json:"Action"` // start, die, create, destroy
	Actor  struct {
		ID string `json:"ID"`
	} `json:"Actor"`
}

type Bridges []Bridge

type Bridge struct {
	ID         string
	Name       string
	Subnets    []string
	Containers Containers
}

type DockerBridgeInspect struct {
	ID      string `json:"Id"`
	Options struct {
		Name string `json:"com.docker.network.bridge.name"`
	} `json:"Options"`
	IPAM struct {
		Config []struct {
			Subnet string `json:"Subnet"`
		} `json:"Config"`
	} `json:"IPAM"`
}

type Containers []Container

type Container struct {
	ID       string
	Networks ContainerNetworks
}

type ContainerNetworks struct {
	IPAddresses []IPInfo
	Ports       []ContainerPort
}

type IPInfo struct {
	Address   string
	Version   int // "4" or "6"
	NetworkID string
}

func (i IPInfo) NftPrefix() string {
	if i.Version == 6 {
		return "ip6"
	}
	return "ip"
}

type ContainerPort struct {
	Port     string
	Protocol string
	HostPort []HostPort
}

type HostPort struct {
	Port string
	IP   IPInfo
}

type DockerContainerInspect struct {
	NetworkSettings struct {
		Ports map[string][]struct {
			HostIp   string `json:"HostIp"`
			HostPort string `json:"HostPort"`
		} `json:"Ports"`
		Networks map[string]struct {
			IPAddress string `json:"IPAddress"`
			NetworkID string `json:"NetworkID"`
		} `json:"Networks"`
	} `json:"NetworkSettings"`
}

func ipVersion(ip string) (int, error) {
	ipParse := net.ParseIP(ip)
	if ipParse == nil || (ipParse.To4() == nil && ipParse.To16() == nil) {
		return 0, errors.New("invalid ip address")
	}

	if ipParse.To4() != nil {
		return 4, nil
	}

	return 6, nil
}
