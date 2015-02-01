package hostess

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
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
}

type Hostfile struct {
	Path  string
	Hosts map[string]*Hostname
	data  string
}

func NewHostfile(path string) *Hostfile {
	return &Hostfile{path, make(map[string]*Hostname), ""}
}

func (h *Hostfile) Read() string {
	data, err := ioutil.ReadFile(h.Path)
	if err != nil {
		fmt.Println("Can't read ", h.Path)
		os.Exit(1)
	}
	h.data = string(data)
	return h.data
}

func Dump(hostnames []Hostname) string {
	return ""
}

func DumpToFile(path string) {

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

func Read(hostfile string) []Hostname {
	var hosts = make([]Hostname, 0)
	return hosts
}

func ReadFile(path string) {

}

func (h *Hostfile) Add(host Hostname) {
	h.Hosts[host.Domain] = &host
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
	hostfile.Read()
	hostfile.Add(Hostname{"localhost", "127.0.0.1", true})
	hostfile.Enable("localhost")
}
