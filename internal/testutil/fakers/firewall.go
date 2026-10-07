package fakers

func RandFirewallProtocol() string {
	protocols := []string{"tcp", "udp"}
	return RandItem(protocols)
}

func RandFirewallDirection() string {
	directions := []string{"in", "out"}
	return RandItem(directions)
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
