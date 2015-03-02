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
	*h = append(*h, host)
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
	*h = append((*h)[:index], (*h)[index+1:]...)
}

func (h *Hostlist) RemoveIpv4Domain(domain string) {
	h.Remove(h.IndexOfDomainIpv4(domain))
}

func (h *Hostlist) RemoveIpv6Domain(domain string) {
	h.Remove(h.IndexOfDomainIpv6(domain))
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

func (h *Hostlist) Sort() {

}

// Format takes the current list of Hostnames in this Hostfile and turns it
// into a string suitable for use as an /etc/hosts file.
// Sorting uses the following logic:
// 1. List is sorted by IP address
// 2. Commented items are left in place
// 3. 127.* appears at the top of the list (so boot resolvers don't break)
// 4. When present, localhost will always appear first in the domain list
func (h *Hostlist) Format() string {
	out := ""
	for _, hostname := range *h {
		out += hostname.Format() + "\n"
	}
	return out
}
