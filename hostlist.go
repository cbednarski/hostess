package hostess

import (
	"errors"
	"fmt"
	"net"
)

// TODO:
// Add
// Remove
// Sort
// Other things that maybe implemented in hostfile

type Hostlist []*Hostname

func NewHostlist() *Hostlist {
	return &Hostlist{}
}

func (h *Hostlist) ContainsHostname(b *Hostname) bool {
	for _, a := range *h {
		if a.Equal(b) {
			return true
		}
	}
	return false
}

func (h *Hostlist) ContainsDomain(domain string) bool {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			return true
		}
	}
	return false
}

func (h *Hostlist) ContainsIp(ip net.IP) bool {
	for _, hostname := range *h {
		if hostname.EqualIp(ip) {
			return true
		}
	}
	return false
}

func (h *Hostlist) Add(host *Hostname) error {
	for _, found := range *h {
		if found.Equal(host) {
			return errors.New(fmt.Sprintf("Duplicate hostname entry for %s -> %s",
				host.Domain, host.Ip))
		} else if found.Domain == host.Domain && found.Ipv6 == host.Ipv6 {
			return errors.New(fmt.Sprintf("Conflicting hostname entries for %s -> %s and -> %s",
				host.Domain, host.Ip, found.Ip))
		}
	}
	return nil
}

func (h *Hostlist) Get(host *Hostname) *Hostname {
	for _, found := range *h {
		if found.Equal(host) {
			return found
		}
	}
	return nil
}

func (h *Hostlist) IndexOf(host *Hostname) int {
	for index, found := range *h {
		if found.Equal(host) {
			return index
		}
	}
	return -1
}

func (h *Hostlist) IndexOfDomainIpv4(domain string) int {
	for index, found := range *h {
		if found.Domain == domain && found.Ipv6 == false {
			return index
		}
	}
	return -1
}

func (h *Hostlist) IndexOfDomainIpv6(domain string) int {
	for index, found := range *h {
		if found.Domain == domain && found.Ipv6 == true {
			return index
		}
	}
	return -1
}

func (h *Hostlist) Remove(index int) {
	// var a *Hostlist
	// copy(a, h[0:index])
	// a = append(a, *h[index:])
	// *h[index] = nil
	// // return a
}

func (h *Hostlist) RemoveIpv4(domain string) {

}

func (h *Hostlist) RemoveIpv6() {

}

func (h *Hostlist) Enable(domain string) {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			hostname.Enabled = true
		}
	}
}

func (h *Hostlist) Disable(domain string) {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			hostname.Enabled = false
		}
	}
}

func (h *Hostlist) Copy() *Hostlist {
	var n *Hostlist
	copy(*h, *n)
	return n
}
