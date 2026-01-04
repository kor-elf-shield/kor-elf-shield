package client

import (
	"errors"
	"net"
)

type Bridges []Bridge

type Bridge struct {
	ID         string
	Name       string
	Subnet     string
	Containers Containers
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
	Address string
	Version int // "4" or "6"
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
