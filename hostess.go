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
	Hosts map[string]Hostname
}

func NewHostfile(path string) *Hostfile {
	return &Hostfile{path, make(map[string]Hostname)}
}

func readHosts(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Can't read ", path)
		os.Exit(1)
	}
	return string(data)
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

func listHosts(hosts []Hostname, host Hostname) {

}

func (h *Hostfile) Add(host Hostname) {
	if h.Hosts == nil {
		h.Hosts = make(map[string]Hostname)
	}
	h.Hosts[host.Domain] = host
}

func (h *Hostfile) Delete(host string) {
	delete(h.Hosts, host)
}

func (h *Hostfile) Enable(host string) {

}

func (h *Hostfile) Disable(host string) {

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
	// count := len(os.Args[2:])
	// args := make([]string, count)
	// if count == 0 {
	// 	return args
	// }

	// for i := 0; i < count; i++ {
	// 	args[i] = append(args, os.Args[i+2])
	// }
	// return args
}

func main() {
	hostfile := new(Hostfile)
	hostfile.Path = "/etc/hosts"
	h := Hostname{"localhost", "127.0.0.1", true}
	hostfile.Add(h)
	fmt.Println(hostfile)
	hostfile.Delete(h.Domain)
	fmt.Println(hostfile)

	// fmt.Println(getCommand())
	// fmt.Println(getArgs())
	// fmt.Println(readHosts(getHostsPath()))
}
