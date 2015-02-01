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

func TestHostname(t *testing.T) {
	const domain = "localhost"
	const ip = "127.0.0.1"
	const enabled = true

	h := Hostname{}
	h.Domain = domain
	h.Ip = ip
	h.Enabled = enabled

	if h.Domain != domain {
		t.Error("Domain should match " + domain)
	}
	if h.Ip != ip {
		t.Error("Domain should match " + ip)
	}
	if h.Enabled != enabled {
		t.Error("Enabled should be " + ip)
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
	hostfile.Add(Hostname{"localhost", "127.0.0.1", true})
}
