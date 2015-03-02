package hostess

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sort"
	"strings"
)

const default_osx = `
##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##

127.0.0.1       localhost
255.255.255.255 broadcasthost
::1             localhost
fe80::1%lo0     localhost
`

const default_linux = `
127.0.0.1   localhost
127.0.1.1   HOSTNAME

# The following lines are desirable for IPv6 capable hosts
::1     localhost ip6-localhost ip6-loopback
fe00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
ff02::3 ip6-allhosts
`

// Hostfile represents /etc/hosts (or a similar file, depending on OS), and
// includes a list of Hostnames. Hostfile includes
type Hostfile struct {
	Path  string
	Hosts Hostlist
	data  string
}

// NewHostFile creates a new Hostfile object from the specified file.
func NewHostfile(path string) *Hostfile {
	return &Hostfile{path, Hostlist{}, ""}
}

func (h *Hostfile) Load() string {
	data, err := ioutil.ReadFile(h.Path)
	if err != nil {
		fmt.Println("Can't read ", h.Path)
		os.Exit(1)
	}
	h.data = string(data)
	return h.data
}

func (h *Hostfile) Parse() []error {
	var errs []error
	for _, v := range strings.Split(h.data, "\n") {
		for _, hostname := range ParseLine(v) {
			err := h.Add(hostname)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func LoadHostFile() (*Hostfile, []error) {
	hostfile := NewHostfile(GetHostsPath())
	hostfile.Load()
	errs := hostfile.Parse()
	return hostfile, errs
}

func TrimWS(s string) string {
	return strings.Trim(s, " \n\t")
}

func ParseLine(line string) Hostlist {
	var hostnames Hostlist

	if len(line) == 0 {
		return hostnames
	}

	// Parse leading # for disabled lines
	enabled := true
	if line[0:1] == "#" {
		enabled = false
		line = TrimWS(line[1:])
	}

	// Parse other #s for actual comments
	line = strings.Split(line, "#")[0]

	// Replace tabs and multispaces with single spaces throughout
	line = strings.Replace(line, "\t", " ", -1)
	for strings.Contains(line, "  ") {
		line = strings.Replace(line, "  ", " ", -1)
	}

	// Break line into words
	words := strings.Split(line, " ")

	// Separate the first bit (the ip) from the other bits (the domains)
	ip := words[0]
	domains := words[1:]

	if LooksLikeIPv4(ip) || LooksLikeIPv6(ip) {
		for _, v := range domains {
			hostname := NewHostname(v, ip, enabled)
			hostnames = append(hostnames, hostname)
		}
	}

	return hostnames
}

func getSortedMapKeys(m map[string][]string) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i += 1
	}
	sort.Strings(keys)
	return keys
}

// MoveToFront looks for string in a slice of strings and if it finds it, moves
// it to the front of the slice.
// Note: this could probably be made faster using pointers to switch the values
// instead of copying a bunch of crap, but it works and speed is not a problem.
func MoveToFront(list []string, search string) []string {
	for k, v := range list {
		if v == search {
			list = append(list[:k], list[k+1:]...)
		}
	}
	return append([]string{search}, list...)
}

// ListDomainsByIP will look through Hostfile to find domains that match the
// specified IP and return them in a sorted slice.
func (h *Hostfile) ListDomainsByIP(ip net.IP) []string {
	var names []string
	for _, v := range h.Hosts {
		if v.IP.Equal(ip) {
			names = append(names, v.Domain)
		}
	}
	sort.Strings(names)

	// Magic for localhost only, to make sure it's the first domain on its line
	if ip.Equal(net.ParseIP("127.0.0.1")) {
		names = MoveToFront(names, "localhost")
	}

	return names
}

// ListDomains will return a list of domains in alphabetical order.
func (h *Hostfile) ListDomains() []string {
	var names []string
	for _, v := range h.Hosts {
		names = append(names, v.Domain)
	}
	sort.Strings(names)
	return names
}

// Format takes the current list of Hostnames in this Hostfile and turns it
// into a string suitable for use as an /etc/hosts file.
// Sorting uses the following logic:
// 1. List is sorted by IP address
// 2. Commented items are left in place
// 3. 127.* appears at the top of the list (so boot resolvers don't break)
// 4. When present, localhost will always appear first in the domain list
func (h *Hostfile) Format() string {
	// localhost := "127.0.0.1 localhost"

	localhosts := make(map[string][]string)
	ips := make(map[string][]string)

	// Map domains and IPs into slices of domains keyd by IP
	// 127.0.0.1 = [localhost, blah, blah2]
	// 2.2.2.3 = [domain1, domain2]
	for _, hostname := range h.Hosts {
		if hostname.IP.String()[0:4] == "127." {
			localhosts[hostname.IP.String()] = append(localhosts[hostname.IP.String()], hostname.Domain)
		} else {
			ips[hostname.IP.String()] = append(ips[hostname.IP.String()], hostname.Domain)
		}
	}

	// localhosts_keys := getSortedMapKeys(localhosts)
	// ips_keys := getSortedMapKeys(ips)
	var out []string

	// for _, ip := range localhosts_keys {
	// 	enabled := ip
	// 	enabled_b := false
	// 	disabled := "# " + ip
	// 	disabled_b := false
	// 	IP := net.ParseIP(ip)
	// 	for _, domain := range h.ListDomainsByIP(IP) {
	// 		hostname := *h.Hosts[domain]
	// 		if hostname.IP.Equal(IP) {
	// 			if hostname.Enabled {
	// 				enabled += " " + hostname.Domain
	// 				enabled_b = true
	// 			} else {
	// 				disabled += " " + hostname.Domain
	// 				disabled_b = true
	// 			}
	// 		}
	// 	}
	// 	if enabled_b {
	// 		out = append(out, enabled)
	// 	}
	// 	if disabled_b {
	// 		out = append(out, disabled)
	// 	}
	// }

	// for _, ip := range ips_keys {
	// 	enabled := ip
	// 	enabled_b := false
	// 	disabled := "# " + ip
	// 	disabled_b := false
	// 	IP := net.ParseIP(ip)
	// 	for _, domain := range h.ListDomainsByIP(IP) {
	// 		hostname := *h.Hosts[domain]
	// 		if hostname.IP.Equal(IP) {
	// 			if hostname.Enabled {
	// 				enabled += " " + hostname.Domain
	// 				enabled_b = true
	// 			} else {
	// 				disabled += " " + hostname.Domain
	// 				disabled_b = true
	// 			}
	// 		}
	// 	}
	// 	if enabled_b {
	// 		out = append(out, enabled)
	// 	}
	// 	if disabled_b {
	// 		out = append(out, disabled)
	// 	}
	// }

	return strings.Join(out, "\n")
}

func (h *Hostfile) Save() error {
	// h.Format(h.Path)
	return nil
}

func (h *Hostfile) Contains(b *Hostname) bool {
	for _, a := range h.Hosts {
		if a.Equal(b) {
			return true
		}
	}
	return false
}

func (h *Hostfile) ContainsDomain(search string) bool {
	for _, hostname := range h.Hosts {
		if hostname.Domain == search {
			return true
		}
	}
	return false
}

func (h *Hostfile) Add(host *Hostname) error {
	// host_f, found := h.Hosts[host.Domain]
	// if found {
	// 	if host_f.IP.Equal(host.IP) {
	// 		return errors.New(fmt.Sprintf("Duplicate hostname entry for %s -> %s",
	// 			host.Domain, host.IP))
	// 	} else {
	// 		return errors.New(fmt.Sprintf("Conflicting hostname entries for %s -> %s and -> %s",
	// 			host.Domain, host.IP, host_f.IP))
	// 	}
	// } else {
	// 	h.Hosts[host.Domain] = host
	// }
	return nil
}

func (h *Hostfile) Delete(domain string) {
	// delete(h.Hosts, domain)
}

func (h *Hostfile) Enable(domain string) {
	// _, ok := h.Hosts[domain]
	// if ok {
	// 	h.Hosts[domain].Enabled = true
	// }
}

func (h *Hostfile) Disable(domain string) {
	// _, ok := h.Hosts[domain]
	// if ok {
	// 	h.Hosts[domain].Enabled = false
	// }
}

func GetHostsPath() string {
	path := os.Getenv("HOSTESS_FILE")
	if path == "" {
		path = "/etc/hosts"
	}
	return path
}
