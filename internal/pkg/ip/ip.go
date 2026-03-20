package ip

import (
	"errors"
	"net"
)

type Version int8

const (
	IPv4 Version = iota + 1
	IPv6
)

func (v Version) ToNft() string {
	switch v {
	case IPv4:
		return "ip"
	case IPv6:
		return "ip6"
	default:
		return "unknown"
	}
}

func (v Version) ToNftForSet() string {
	switch v {
	case IPv4:
		return "ipv4_addr"
	case IPv6:
		return "ipv6_addr"
	default:
		return "unknown"
	}
}

func DetermineIPVersion(ip string) (ipNet string, version Version, err error) {
	ipNet, version, err = parseCIDR(ip)
	if err != nil {
		return parseIP(ip)
	}

	return
}

func IPVersion(ip string) (Version, error) {
	_, version, err := parseIP(ip)
	return version, err
}

func parseCIDR(parseIP string) (ipNet string, version Version, err error) {
	_, parseIPNet, err := net.ParseCIDR(parseIP)
	if err != nil {
		return
	}

	if parseIPNet.IP.To4() == nil && parseIPNet.IP.To16() == nil {
		err = errors.New("invalid ip address")
		return
	}

	ipNet = parseIPNet.String()
	if parseIPNet.IP.To4() != nil {
		version = IPv4
		return ipNet, version, nil
	}

	version = IPv6
	return ipNet, version, nil
}

func parseIP(parseIP string) (ip string, version Version, err error) {
	ipParse := net.ParseIP(parseIP)
	if ipParse == nil || (ipParse.To4() == nil && ipParse.To16() == nil) {
		err = errors.New("invalid ip address")
		return
	}

	if ipParse.To4() != nil {
		version = IPv4
		return ipParse.To4().String(), version, nil
	}

	version = IPv6
	return ipParse.To16().String(), version, nil
}
