package fakers

func RandFirewallProtocol() string {
	protocols := []string{"tcp", "udp"}
	return RandItem(protocols)
}

func RandFirewallProtocols() []string {
	var protocols []string
	for i := RandInt(1, 2); i > 0; i-- {
		protocols = append(protocols, RandFirewallProtocol())
	}

	return protocols
}

func RandFirewallDirection() string {
	directions := []string{"in", "out"}
	return RandItem(directions)
}

func RandFirewallDirections() []string {
	var directions []string
	for i := RandInt(1, 2); i > 0; i-- {
		directions = append(directions, RandFirewallDirection())
	}

	return directions
}

func RandFirewallAction() string {
	actions := []string{"accept", "drop", "reject"}
	return RandItem(actions)
}

func RandFirewallLimitRate() string {
	rates := []string{
		"",
		"1/s",
		"10/minute",
		"100/hour",
		"1000/day",
		"5/week",
		"0.5/second",
		"over 1/s",
		"10 packets/second",
		"1 kbytes/min",
		"2 mbytes/h",
		"1/s burst 5",
		"1/s burst 10 packets",
		"2/min burst 1.5 mbytes",
	}
	return RandItem(rates)
}

func RandFirewallDrop() string {
	drops := []string{"drop", "reject"}
	return RandItem(drops)
}

func RandFirewallPorts() []int {
	var ports []int
	for i := RandInt(1, 5); i > 0; i-- {
		ports = append(ports, RandInt(0, 65535))
	}

	return ports
}
