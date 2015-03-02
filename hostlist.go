package hostess

import (
	"errors"
	"fmt"
	"net"
	"sort"
)

var InvalidVersionArgumentError = errors.New("Version argument must be 4 or 6")

type Hostlist []*Hostname

// NewHostList initializes a new Hostlist
func NewHostlist() *Hostlist {
	return &Hostlist{}
}

// Len returns the number of hostnames in the list, part of sort.Interface
func (h Hostlist) Len() int {
	return len(h)
}

// Less determines the sort order of two hostnames, part of sort.Interface
func (h Hostlist) Less(i, j int) bool {
	// Sort 127.0.0.1, 127.0.1.1 and "localhost" at the top
	if h[i].Domain == "localhost" {
		return true
	}
	if h[j].Domain == "localhost" {
		return false
	}

	// Sort ipv4 before ipv6
	if h[i].Ipv6 && !h[j].Ipv6 {
		return false
	}
	if !h[i].Ipv6 && h[j].Ipv6 {
		return true
	}

	// Compare the the ip addresses (byte array)
	for c, _ := range h[i].Ip {
		if h[i].Ip[c] < h[j].Ip[c] {
			return true
		} else if h[i].Ip[c] > h[j].Ip[c] {
			return false
		}
	}

	// Prep for domain sorting
	ilen := len(h[i].Domain)
	jlen := len(h[j].Domain)
	max := ilen
	if jlen > max {
		max = jlen
	}

	// Sort domains alphabetically
	// Note: This works best if domains are lowercased. However, we do not
	// enforce lowercase because of UTF-8 domain names, which may be broken by
	// case folding. There is a way to do this correctly but it's completed so
	// I'm not going to do it right now.
	for c := 0; c < max; c++ {
		if c > ilen {
			return true
		}
		if c > jlen {
			return false
		}
		if h[i].Domain[c] < h[j].Domain[c] {
			return true
		}
		if h[i].Domain[c] > h[j].Domain[c] {
			return false
		}
	}

	// Seems like everything was the same, so it can't be Less
	return false
}

// Swap changes the position of two hostnames, part of sort.Interface
func (h Hostlist) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Sort this list of Hostnames, according to Hostlist sorting rules:
// 1. localhost comes first
// 2. ipv4 comes first
// 3. ips are sorted in numerical order
// 4. domains are sorted in alphabetical
func (h *Hostlist) Sort() {
	sort.Sort(*h)
}

// Contains returns true if this Hostlist has the specified Hostname
func (h *Hostlist) Contains(b *Hostname) bool {
	for _, a := range *h {
		if a.Equal(b) {
			return true
		}
	}
	return false
}

// ContainsDomain returns true if a Hostname in this Hostlist matches domain
func (h *Hostlist) ContainsDomain(domain string) bool {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			return true
		}
	}
	return false
}

// ContainsDomain returns true if a Hostname in this Hostlist matches ip
func (h *Hostlist) ContainsIp(ip net.IP) bool {
	for _, hostname := range *h {
		if hostname.EqualIp(ip) {
			return true
		}
	}
	return false
}

// Add a new Hostname to this hostlist. If a Hostname with the same domain name
// and Ip version is found, it will be replaced and an error will be returned.
// If you try to add an identical Hostname, an error will be returned.
// Note that in normal operation, you will sometimes expect an error, and the
// error data is mainly to alert you that you mis-entered information, not that
// the application has a problem.
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

// IndexOf will indicate the index of a Hostname in Hostlist, or -1 if it is
// not found.
func (h *Hostlist) IndexOf(host *Hostname) int {
	for index, found := range *h {
		if found.Equal(host) {
			return index
		}
	}
	return -1
}

// IndexOfDomainV will indicate the index of a Hostname in Hostlist that has
// the same domain and IP version, or -1 if it is not found.
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) IndexOfDomainV(domain string, version int) int {
	if version != 4 && version != 6 {
		panic(InvalidVersionArgumentError)
	}
	for index, hostname := range *h {
		if hostname.Domain == domain && hostname.Ipv6 == (version == 6) {
			return index
		}
	}
	return -1
}

// Remove will delete the Hostname at the specified index. If index is out of
// bounds (i.e. -1), Remove silently no-ops.
func (h *Hostlist) Remove(index int) {
	if index > -1 && index < len(*h) {
		*h = append((*h)[:index], (*h)[index+1:]...)
	}
}

// RemoveDomain removes both Ipv4 and Ipv6 Hostname entries matching domain.
func (h *Hostlist) RemoveDomain(domain string) {
	h.Remove(h.IndexOfDomainV(domain, 4))
	h.Remove(h.IndexOfDomainV(domain, 6))
}

// RemoveDomainV removes a Hostname entry matching the domain and IP version.
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) RemoveDomainV(domain string, version int) {
	if version != 4 && version != 6 {
		panic(InvalidVersionArgumentError)
	}
	h.Remove(h.IndexOfDomainV(domain, version))
}

// Enable will change any Hostnames matching domain to be enabled.
func (h *Hostlist) Enable(domain string) {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			hostname.Enabled = true
		}
	}
}

// EnableV will change a Hostname matching domain and IP version to be enabled.
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) EnableV(domain string, version int) {
	if version != 4 && version != 6 {
		panic(InvalidVersionArgumentError)
	}
	for _, hostname := range *h {
		if hostname.Domain == domain && hostname.Ipv6 == (version == 6) {
			hostname.Enabled = true
		}
	}
}

// Enable will change any Hostnames matching domain to be disabled.
func (h *Hostlist) Disable(domain string) {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			hostname.Enabled = false
		}
	}
}

// Enable will change any Hostnames matching domain and IP version to be disabled.
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) DisableV(domain string, version int) {
	if version != 4 && version != 6 {
		panic(InvalidVersionArgumentError)
	}
	for _, hostname := range *h {
		if hostname.Domain == domain && hostname.Ipv6 == (version == 6) {
			hostname.Enabled = false
		}
	}
}

// Format takes the current list of Hostnames in this Hostfile and turns it
// into a string suitable for use as an /etc/hosts file.
// Sorting uses the following logic:
// 1. List is sorted by IP address
// 2. Commented items are left in place
// 3. 127.* appears at the top of the list (so boot resolvers don't break)
// 4. When present, localhost will always appear first in the domain list
func (h *Hostlist) Format() string {
	h.Sort()
	out := ""
	for _, hostname := range *h {
		out += hostname.Format() + "\n"
	}
	return out
}
