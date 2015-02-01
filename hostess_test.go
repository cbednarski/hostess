package hostess

import (
	"testing"
)

const ipv4_pass = `
127.0.0.1
127.0.1.1
10.200.30.50
99.99.99.99
999.999.999.999
0.1.1.0
`

const ipv4_fail = `
1234.1.1.1
123.5.6
12.12
76.76.67.67.45
`

const ipv6 = ``

const domain = "localhost"
const ip = "127.0.0.1"
const enabled = true

func TestHostname(t *testing.T) {

	h := Hostname{}
	h.Domain = domain
	h.Ip = ip
	h.Enabled = enabled

	if h.Domain != domain {
		t.Errorf("Domain should be %s", domain)
	}
	if h.Ip != ip {
		t.Errorf("Domain should be %s", ip)
	}
	if h.Enabled != enabled {
		t.Errorf("Enabled should be %s", enabled)
	}
}

func TestGetHostsPath(t *testing.T) {
	path := GetHostsPath()
	const expected = "/etc/hosts"
	if path != expected {
		t.Error("Hosts path should be " + expected)
	}
}

func TestHostfile(t *testing.T) {
	hostfile := NewHostfile(GetHostsPath())
	hostfile.Add(Hostname{domain, ip, true})
	if hostfile.Hosts[domain].Ip != ip {
		t.Errorf("Hostsfile should have %s pointing to %s", domain, ip)
	}

	hostfile.Disable(domain)
	if hostfile.Hosts[domain].Enabled != false {
		t.Errorf("%s should be disabled", domain)
	}

	hostfile.Enable(domain)
	if hostfile.Hosts[domain].Enabled != true {
		t.Errorf("%s should be enabled", domain)
	}

	hostfile.Delete(domain)
	if hostfile.Hosts[domain] != nil {
		t.Errorf("Did not expect to find %s", domain)
	}

}
