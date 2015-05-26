package hostess_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cbednarski/hostess"
)

const ipv4Pass = `
127.0.0.1
127.0.1.1
10.200.30.50
99.99.99.99
999.999.999.999
0.1.1.0
`

const ipv4Fail = `
1234.1.1.1
123.5.6
12.12
76.76.67.67.45
`

const ipv6 = ``
const domain = "localhost"
const ip = "127.0.0.1"
const enabled = true

func Diff(expected, actual string) string {
	return fmt.Sprintf(`
---- Expected ----
%s
----- Actual -----
%s
`, expected, actual)
}

func TestGetHostsPath(t *testing.T) {
	path := hostess.GetHostsPath()
	const expected = "/etc/hosts"
	if path != expected {
		t.Error("Hosts path should be " + expected)
	}
}

func TestFormatHostfile(t *testing.T) {
	// The sort order here is a bit weird.
	// 1. We want localhost entries at the top
	// 2. The rest are sorted by IP as STRINGS, not numeric values, so 10
	//    precedes 8
	const expected = `127.0.0.1 localhost devsite
127.0.1.1 ip-10-37-12-18
# 8.8.8.8 google.com
10.37.12.18 devsite.com m.devsite.com
`

	hostfile := hostess.NewHostfile()
	hostfile.Path = "./hosts"
	hostfile.Hosts.Add(hostess.NewHostname("localhost", "127.0.0.1", true))
	hostfile.Hosts.Add(hostess.NewHostname("ip-10-37-12-18", "127.0.1.1", true))
	hostfile.Hosts.Add(hostess.NewHostname("devsite", "127.0.0.1", true))
	hostfile.Hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", false))
	hostfile.Hosts.Add(hostess.NewHostname("devsite.com", "10.37.12.18", true))
	hostfile.Hosts.Add(hostess.NewHostname("m.devsite.com", "10.37.12.18", true))
	f := string(hostfile.Format())
	if f != expected {
		t.Errorf("Hostfile output is not formatted correctly: %s", Diff(expected, f))
	}
}

func TestTrimWS(t *testing.T) {
	const expected = `  candy

	`
	actual := hostess.TrimWS(expected)
	if actual != "candy" {
		t.Errorf("Output was not trimmed correctly: %s", Diff(expected, actual))
	}
}

func TestParseLine(t *testing.T) {
	var hosts hostess.Hostlist

	// Blank line
	hosts = hostess.ParseLine("")
	if len(hosts) > 0 {
		t.Error("Expected to find zero hostnames")
	}

	// Comment
	hosts = hostess.ParseLine("# The following lines are desirable for IPv6 capable hosts")
	if len(hosts) > 0 {
		t.Error("Expected to find zero hostnames")
	}

	// Single word comment
	hosts = hostess.ParseLine("#blah")
	if len(hosts) > 0 {
		t.Error("Expected to find zero hostnames")
	}

	hosts = hostess.ParseLine("#66.33.99.11              test.domain.com")
	if !hosts.Contains(hostess.NewHostname("test.domain.com", "66.33.99.11", false)) ||
		len(hosts) != 1 {
		t.Error("Expected to find test.domain.com (disabled)")
	}

	hosts = hostess.ParseLine("#  66.33.99.11	test.domain.com	domain.com")
	if !hosts.Contains(hostess.NewHostname("test.domain.com", "66.33.99.11", false)) ||
		!hosts.Contains(hostess.NewHostname("domain.com", "66.33.99.11", false)) ||
		len(hosts) != 2 {
		t.Error("Expected to find domain.com and test.domain.com (disabled)")
		t.Errorf("Found %s", hosts)
	}

	// Not Commented stuff
	hosts = hostess.ParseLine("255.255.255.255 broadcasthost test.domain.com	domain.com")
	if !hosts.Contains(hostess.NewHostname("broadcasthost", "255.255.255.255", true)) ||
		!hosts.Contains(hostess.NewHostname("test.domain.com", "255.255.255.255", true)) ||
		!hosts.Contains(hostess.NewHostname("domain.com", "255.255.255.255", true)) ||
		len(hosts) != 3 {
		t.Error("Expected to find broadcasthost, domain.com, and test.domain.com (enabled)")
	}

	// Ipv6 stuff
	hosts = hostess.ParseLine("::1             localhost")
	if !hosts.Contains(hostess.NewHostname("localhost", "::1", true)) ||
		len(hosts) != 1 {
		t.Error("Expected to find localhost ipv6 (enabled)")
	}

	hosts = hostess.ParseLine("ff02::1 ip6-allnodes")
	if !hosts.Contains(hostess.NewHostname("ip6-allnodes", "ff02::1", true)) ||
		len(hosts) != 1 {
		t.Error("Expected to find ip6-allnodes ipv6 (enabled)")
	}
}

func TestLoadHostfile(t *testing.T) {
	hostfile := hostess.NewHostfile()
	hostfile.Read()
	if !strings.Contains(string(hostfile.GetData()), domain) {
		t.Errorf("Expected to find %s", domain)
	}

	hostfile.Parse()
	hostname := hostess.NewHostname(domain, ip, enabled)
	found := hostfile.Hosts.Contains(hostname)
	if !found {
		t.Errorf("Expected to find %s", hostname)
	}
}
