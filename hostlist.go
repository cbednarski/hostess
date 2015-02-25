package hostess

import (
	"net"
)

// TODO:
// Add
// Remove
// Sort
// Other things that maybe implemented in hostfile

func ContainsHostname(hostnames []*Hostname, b *Hostname) bool {
	for _, a := range hostnames {
		if a.Equal(b) {
			return true
		}
	}
	return false
}

func ContainsDomain(hostnames []*Hostname, domain string) bool {
	for _, hostname := range hostnames {
		if hostname.Domain == domain {
			return true
		}
	}
	return false
}

func ContainsIp(hostnames []*Hostname, ip net.IP) bool {
	for _, hostname := range hostnames {
		if hostname.EqualIp(ip) {
			return true
		}
	}
	return false
}
