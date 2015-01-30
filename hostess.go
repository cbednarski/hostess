package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

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

func writeHosts(path string, contents string) {

}

func parseLine(line string) {
	// return (Hostname, err)
}

func parseHosts(hostfile string) []Hostname {
	var hosts = make([]Hostname, 0)
	return hosts
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

func getHostsPath() string {
	path := os.Getenv("HOSTESS_FILE")
	if path == "" {
		path = "/etc/hosts"
	}
	return path
}

func getCommand() string {
	return os.Args[1]
}

func getArgs() []string {
	return os.Args[2:]
}

func main() {
	hostfile := NewHostfile(getHostsPath())
	hostfile.Read()
	hostfile.Add(Hostname{"localhost", "127.0.0.1", true})
	hostfile.Enable("localhost")

	fmt.Println(getArgs())
}
