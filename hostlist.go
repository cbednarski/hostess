package hostess

import (
	"errors"
	"fmt"
	"net"
	"sort"
)

// ErrInvalidVersionArg is raised when a function expects IPv 4 or 6 but is
// passed a value not 4 or 6.
var ErrInvalidVersionArg = errors.New("Version argument must be 4 or 6")

// Hostlist is an ordered set of Hostnames. When in a Hostlist, Hostnames must
// follow some rules:
//
// 	- Hostlist may contain IPv4 AND IPv6 ("IP version" or "IPv") Hostnames.
// 	- Names are only allowed to overlap if IP version is different.
// 	- Adding a Hostname for an existing name will replace the old one.
//
// See docs for the Sort and Add for more details.
type Hostlist []*Hostname

// NewHostlist initializes a new Hostlist
func NewHostlist() *Hostlist {
	return &Hostlist{}
}

// Len returns the number of Hostnames in the list, part of sort.Interface
func (h Hostlist) Len() int {
	return len(h)
}

// MakeSurrogateIP takes an IP like 127.0.0.1 and munges it to 0.0.0.1 so we can
// sort it more easily.
func MakeSurrogateIP(IP net.IP) net.IP {
	if string(IP[0:3]) == "127" {
		return net.IP("0" + string(IP[3:]))
	}
	return IP
}

// Less determines the sort order of two Hostnames, part of sort.Interface
func (h Hostlist) Less(A, B int) bool {
	// Sort "localhost" at the top
	if h[A].Domain == "localhost" {
		return true
	}
	if h[B].Domain == "localhost" {
		return false
	}

	// Sort IPv4 before IPv6
	// A is IPv4 and B is IPv6. A wins!
	if !h[A].IPv6 && h[B].IPv6 {
		return true
	}
	// A is IPv6 but B is IPv4. A loses!
	if h[A].IPv6 && !h[B].IPv6 {
		return false
	}

	// Compare the the IP addresses (byte array)
	// We want to push 127. to the top so we're going to mark it zero.
	surrogateA := MakeSurrogateIP(h[A].IP)
	surrogateB := MakeSurrogateIP(h[B].IP)
	if !surrogateA.Equal(surrogateB) {
		for charIndex := range surrogateA {
			// A and B's IPs differ at this index, and A is less. A wins!
			if surrogateA[charIndex] < surrogateB[charIndex] {
				return true
			}
			// A and B's IPs differ at this index, and B is less. A loses!
			if surrogateA[charIndex] > surrogateB[charIndex] {
				return false
			}
		}
		// If we got here then the IPs are the same and we want to continue on
		// to the domain sorting section.
	}

	// Prep for sorting by domain name
	aLength := len(h[A].Domain)
	bLength := len(h[B].Domain)
	max := aLength
	if bLength > max {
		max = bLength
	}

	// Sort domains alphabetically
	// TODO: This works best if domains are lowercased. However, we do not
	// enforce lowercase because of UTF-8 domain names, which may be broken by
	// case folding. There is a way to do this correctly but it's complicated
	// so I'm not going to do it right now.
	for charIndex := 0; charIndex < max; charIndex++ {
		// This index is longer than A, so A is shorter. A wins!
		if charIndex >= aLength {
			return true
		}
		// This index is longer than B, so B is shorter. A loses!
		if charIndex >= bLength {
			return false
		}
		// A and B differ at this index and A is less. A wins!
		if h[A].Domain[charIndex] < h[B].Domain[charIndex] {
			return true
		}
		// A and B differ at this index and B is less. A loses!
		if h[A].Domain[charIndex] > h[B].Domain[charIndex] {
			return false
		}
	}

	// If we got here then A and B are the same -- by definition A is not Less
	// than B so we return false. Technically we shouldn't get here since Add
	// should not allow duplicates, but we'll guard anyway.
	return false
}

// Swap changes the position of two Hostnames, part of sort.Interface
func (h Hostlist) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// Sort this list of Hostnames, according to Hostlist sorting rules:
//
// 	1. localhost comes before other domains
// 	2. IPv4 comes before IPv6
// 	3. IPs are sorted in numerical order
// 	4. domains are sorted in alphabetical
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

// ContainsIP returns true if a Hostname in this Hostlist matches IP
func (h *Hostlist) ContainsIP(IP net.IP) bool {
	for _, hostname := range *h {
		if hostname.EqualIP(IP) {
			return true
		}
	}
	return false
}

// Add a new Hostname to this hostlist. If a Hostname with the same domain name
// and IP version is found, it will be replaced and an error will be returned.
// If you try to add an identical Hostname, an error will be returned.
// Note that in normal operation, you will sometimes expect an error, and the
// error data is mainly to alert you that you mis-entered information, not that
// the application has a problem.
func (h *Hostlist) Add(host *Hostname) error {
	for _, found := range *h {
		if found.Equal(host) {
			return fmt.Errorf("Duplicate hostname entry for %s -> %s",
				host.Domain, host.IP)
		} else if found.Domain == host.Domain && found.IPv6 == host.IPv6 {
			return fmt.Errorf("Conflicting hostname entries for %s -> %s and -> %s",
				host.Domain, host.IP, found.IP)
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
//
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) IndexOfDomainV(domain string, version int) int {
	if version != 4 && version != 6 {
		panic(ErrInvalidVersionArg)
	}
	for index, hostname := range *h {
		if hostname.Domain == domain && hostname.IPv6 == (version == 6) {
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

// RemoveDomain removes both IPv4 and IPv6 Hostname entries matching domain.
func (h *Hostlist) RemoveDomain(domain string) {
	h.Remove(h.IndexOfDomainV(domain, 4))
	h.Remove(h.IndexOfDomainV(domain, 6))
}

// RemoveDomainV removes a Hostname entry matching the domain and IP version.
//
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) RemoveDomainV(domain string, version int) {
	if version != 4 && version != 6 {
		panic(ErrInvalidVersionArg)
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
//
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) EnableV(domain string, version int) {
	if version != 4 && version != 6 {
		panic(ErrInvalidVersionArg)
	}
	for _, hostname := range *h {
		if hostname.Domain == domain && hostname.IPv6 == (version == 6) {
			hostname.Enabled = true
		}
	}
}

// Disable will change any Hostnames matching domain to be disabled.
func (h *Hostlist) Disable(domain string) {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			hostname.Enabled = false
		}
	}
}

// DisableV will change any Hostnames matching domain and IP version to be disabled.
//
// This function will panic if IP version is not 4 or 6.
func (h *Hostlist) DisableV(domain string, version int) {
	if version != 4 && version != 6 {
		panic(ErrInvalidVersionArg)
	}
	for _, hostname := range *h {
		if hostname.Domain == domain && hostname.IPv6 == (version == 6) {
			hostname.Enabled = false
		}
	}
}

// FilterByIP filters the list of hostnames by IP address.
func (h *Hostlist) FilterByIP(IP net.IP) (hostnames []*Hostname) {
	for _, hostname := range *h {
		if hostname.IP.Equal(IP) {
			hostnames = append(hostnames, hostname)
		}
	}
	return
}

// FilterByDomain filters the list of hostnames by Domain.
func (h *Hostlist) FilterByDomain(domain string) (hostnames []*Hostname) {
	for _, hostname := range *h {
		if hostname.Domain == domain {
			hostnames = append(hostnames, hostname)
		}
	}
	return
}

// FilterByDomainV filters the list of hostnames by domain and IPv4 or IPv6.
// This should never contain more than one item, but returns a list for
// consistency with other filter functions.
func (h *Hostlist) FilterByDomainV(domain string, version int) (hostnames []*Hostname) {
	for _, hostname := range *h {
		if hostname.Domain == domain && hostname.IPv6 == (version == 6) {
			hostnames = append(hostnames, hostname)
		}
	}
	return
}

// Format takes the current list of Hostnames in this Hostfile and turns it
// into a string suitable for use as an /etc/hosts file.
// Sorting uses the following logic:
//
// 1. List is sorted by IP address
// 2. Commented items are sorted displayed
// 3. 127.* appears at the top of the list (so boot resolvers don't break)
// 4. When present, "localhost" will always appear first in the domain list
func (h *Hostlist) Format() []byte {
	h.Sort()
	var out []byte
	for _, hostname := range *h {
		out = append(out, []byte(hostname.Format())...)
		out = append(out, []byte("\n")...)
	}
	return out
}
