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

func (h *Hostfile) ReadFile(path string) string {
	data, err := ioutil.ReadFile(h.Path)
	if err != nil {
		fmt.Println("Can't read ", h.Path)
		os.Exit(1)
	}
	h.data = string(data)
	return h.data
}

var line_parser = regexp.MustCompile(``)

func parseLine(line string) (Hostname, error) {
	// 1. Split on # to discard comments.
	// 2. Split on first space to find the IP
	// 3. Split remainder of line on whitespace to find
	//    domain names
	// 4. Validate the IP (maybe -- could be ipv4 or ipv6)
	hostname := Hostname{}
	if false {
		return hostname, errors.New("Can't parse hostname")
	}
	return hostname, nil
}

func (h *Hostfile) Read(hostfile string) []Hostname {
	var hosts = make([]Hostname, 0)
	return hosts
}

// Dump takes the current list of Hostnames in this Hostfile and turns it into
// a string suitable for use as an /etc/hosts file.
// Sorting uses the following logic:
// 1. List is sorted by IP address
// 2. Commented items are left in place
// 3. 127.* appears at the top of the list (so boot resolvers don't break)
// 4. When present, localhost will always appear first in the domain list
func (h *Hostfile) Format() string {
	localhost := "127.0.0.1 localhost"

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

	out := make([]string, 0)
	for _, hostname := range h.Hosts {
		if hostname.Enabled && hostname.Ip == "127.0.0.1" {
			if hostname.Domain != "localhost" {
				localhost += " " + hostname.Domain
			}
		} else {
			out = append(out, hostname.Format())
		}

	}
	sort.Sort(sort.StringSlice(out))
	return localhost + "\n" + strings.Join(out, "\n")
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

func Hostess() {
	hostfile := NewHostfile(GetHostsPath())
	hostfile.ReadFile(hostfile.Path)
}
