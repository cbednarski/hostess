package hostess

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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

func TrimWS(s string) string {
	return strings.Trim(s, " \n\t")
}

func parseLine(line string) []Hostname {
	hostnames := make([]Hostname, 0)

	// Parse leading comment for disabled lines
	enabled := true
	if line[0:1] == "#" {
		enabled = false
		line = line[1:]
	}

	// Parse other comment for actual comments
	line = strings.Split(line, "#")[0]

	// Replace tabs and multispaces with single spaces throughout
	line = strings.Replace(line, "\t", " ", -1)
	line = strings.Replace(line, "  ", " ", -1)

	// Break line into words
	words := strings.Split(line, " ")

	// Separate the first bit (the ip) from the other bits (the domains)
	ip := words[0]
	domains := words[1:]

	if LooksLikeIpv4(ip) || LooksLikeIpv6(ip) {
		for _, v := range domains {
			hostnames = append(hostnames, Hostname{v, ip, enabled})
		}
	}

	// 1. Split on # to discard comments.
	//    - What if it's a legit line but commented out?
	//    - Disabled lines should have # at the beginning of the line
	// 2. Split on first space to find the IP
	//    - Pass to LooksLikeIpv4 or LooksLikeIpv6 to check
	//    - If it doesn't look like either of these then it's probably a random
	//      comment.
	// 3. Split remainder of line on whitespace to find
	//    domain names
	//    - Watch out for comments.

	return hostnames
}

var ipv4_pattern = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)

func LooksLikeIpv4(ip string) bool {
	return false
}

var ipv6_pattern = regexp.MustCompile(``)

func LooksLikeIpv6(ip string) bool {
	return false
}

func ContainsHostname(hostnames []Hostname, hostname Hostname) bool {
	for _, v := range hostnames {
		if v.Ip == hostname.Ip && v.Domain == hostname.Domain {
			return true
		}
	}
	return false
}

func ContainsDomain(hostnames []Hostname, domain string) bool {
	for _, v := range hostnames {
		if v.Domain == domain {
			return true
		}
	}
	return false
}

func ContainsIp(hostnames []Hostname, ip string) bool {
	for _, v := range hostnames {
		if v.Ip == ip {
			return true
		}
	}
	return false
}

type Hostname struct {
	Domain  string
	Ip      string
	Enabled bool
	// Ipv6    bool
}

func (h *Hostname) Format() string {
	r := fmt.Sprintf("%s %s", h.Ip, h.Domain)
	if !h.Enabled {
		r = "# " + r
	}
	return r
}

// Hostfile represents /etc/hosts (or a similar file, depending on OS), and
// includes a list of Hostnames. Hostfile includes
type Hostfile struct {
	Path  string
	Hosts map[string]*Hostname
	data  string
}

// NewHostFile creates a new Hostfile object from the specified file.
func NewHostfile(path string) *Hostfile {
	return &Hostfile{path, make(map[string]*Hostname), ""}
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

func (h *Hostfile) Parse() {
	for _, v := range strings.Split(h.data, "\n") {
		fmt.Println(v)
	}
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

// ListDomainsByIp will look through Hostfile to find domains that match the
// specified Ip and return them in a sorted slice.
func (h *Hostfile) ListDomainsByIp(ip string) []string {
	names := make([]string, 0)
	for _, v := range h.Hosts {
		if v.Ip == ip {
			names = append(names, v.Domain)
		}
	}
	sort.Strings(names)

	// Magic for localhost only, to make sure it's the first domain on its line
	if ip == "127.0.0.1" {
		names = MoveToFront(names, "localhost")
	}

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
		if hostname.Ip[0:4] == "127." {
			localhosts[hostname.Ip] = append(localhosts[hostname.Ip], hostname.Domain)
		} else {
			ips[hostname.Ip] = append(ips[hostname.Ip], hostname.Domain)
		}
	}

	localhosts_keys := getSortedMapKeys(localhosts)
	ips_keys := getSortedMapKeys(ips)
	out := make([]string, 0)

	for _, ip := range localhosts_keys {
		enabled := ip
		enabled_b := false
		disabled := "# " + ip
		disabled_b := false
		for _, domain := range h.ListDomainsByIp(ip) {
			hostname := *h.Hosts[domain]
			if hostname.Ip == ip {
				if hostname.Enabled {
					enabled += " " + hostname.Domain
					enabled_b = true
				} else {
					disabled += " " + hostname.Domain
					disabled_b = true
				}
			}
		}
		if enabled_b {
			out = append(out, enabled)
		}
		if disabled_b {
			out = append(out, disabled)
		}
	}

	for _, ip := range ips_keys {
		enabled := ip
		enabled_b := false
		disabled := "# " + ip
		disabled_b := false
		for _, domain := range h.ListDomainsByIp(ip) {
			hostname := *h.Hosts[domain]
			if hostname.Ip == ip {
				if hostname.Enabled {
					enabled += " " + hostname.Domain
					enabled_b = true
				} else {
					disabled += " " + hostname.Domain
					disabled_b = true
				}
			}
		}
		if enabled_b {
			out = append(out, enabled)
		}
		if disabled_b {
			out = append(out, disabled)
		}
	}

	return strings.Join(out, "\n")
}

func (h *Hostfile) Save() error {
	// h.Format(h.Path)
	return nil
}

func (h *Hostfile) Add(host Hostname) error {
	host_f, found := h.Hosts[host.Domain]
	if found {
		if host_f.Ip == host.Ip {
			return errors.New(fmt.Sprintf("Duplicate hostname entry for %s -> %s",
				host.Domain, host.Ip))
		} else {
			return errors.New(fmt.Sprintf("Conflicting hostname entries for %s -> %s and -> %s",
				host.Domain, host.Ip, host_f.Ip))
		}
	} else {
		h.Hosts[host.Domain] = &host
	}
	return nil
}

func (h *Hostfile) Delete(domain string) {
	delete(h.Hosts, domain)
}

func (h *Hostfile) Enable(domain string) {
	_, ok := h.Hosts[domain]
	if ok {
		h.Hosts[domain].Enabled = true
	}
}

func (h *Hostfile) Disable(domain string) {
	_, ok := h.Hosts[domain]
	if ok {
		h.Hosts[domain].Enabled = false
	}
}

func GetHostsPath() string {
	path := os.Getenv("HOSTESS_FILE")
	if path == "" {
		path = "/etc/hosts"
	}
	return path
}
