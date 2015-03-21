package hostess_test

import (
	"github.com/cbednarski/hostess"
	"net"
	"strings"
	"testing"
)

const asserts = `
--- Expected ---
%s
---- Actual ----
%s`

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

func TestGetHostsPath(t *testing.T) {
	path := hostess.GetHostsPath()
	const expected = "/etc/hosts"
	if path != expected {
		t.Error("Hosts path should be " + expected)
	}
}

func TestHostfile(t *testing.T) {
	hostfile := hostess.NewHostfile("./hosts")
	hostname := hostess.NewHostname(domain, ip, true)
	hostfile.Hosts.Add(hostname)
	if !hostfile.Hosts[0].IP.Equal(net.IP(ip)) {
		t.Errorf("Hostsfile should have %s pointing to %s", domain, ip)
	}

	hostfile.Hosts.Disable(domain)
	if hostfile.Hosts[0].Enabled != false {
		t.Errorf("%s should be disabled", domain)
	}

	hostfile.Hosts.Enable(domain)
	if hostfile.Hosts[0].Enabled != true {
		t.Errorf("%s should be enabled", domain)
	}

	hostfile.Hosts.RemoveDomain(domain)
	if hostfile.Hosts[0] != nil {
		t.Errorf("Did not expect to find %s", domain)
	}
}

func TestHostFileDuplicates(t *testing.T) {
	hostfile := hostess.NewHostfile("./hosts")

	const exp_duplicate = "Duplicate hostname entry for localhost -> 127.0.0.1"
	hostfile.Hosts.Add(hostess.NewHostname(domain, ip, true))
	err := hostfile.Hosts.Add(hostess.NewHostname(domain, ip, true))
	if err.Error() != exp_duplicate {
		t.Errorf(asserts, exp_duplicate, err)
	}

	const exp_conflict = "Conflicting hostname entries for localhost -> 127.0.1.1 and -> 127.0.0.1"
	err2 := hostfile.Hosts.Add(hostess.NewHostname(domain, "127.0.1.1", true))
	if err2.Error() != exp_conflict {
		t.Errorf(asserts, exp_conflict, err2)
	}

	// @TODO Add an additional test case here: Adding a domain twice with one
	// enabled and one disabled should just add the domain once enabled.
}

func TestFormatHostfile(t *testing.T) {
	// The sort order here is a bit weird.
	// 1. We want localhost entries at the top
	// 2. The rest are sorted by IP as STRINGS, not numeric values, so 10
	//    precedes 8
	const expected = `127.0.0.1 localhost devsite
127.0.1.1 ip-10-37-12-18
10.37.12.18 devsite.com m.devsite.com
# 8.8.8.8 google.com`

	hostfile := hostess.NewHostfile("./hosts")
	hostfile.Hosts.Add(hostess.NewHostname("localhost", "127.0.0.1", true))
	hostfile.Hosts.Add(hostess.NewHostname("ip-10-37-12-18", "127.0.1.1", true))
	hostfile.Hosts.Add(hostess.NewHostname("devsite", "127.0.0.1", true))
	hostfile.Hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", false))
	hostfile.Hosts.Add(hostess.NewHostname("devsite.com", "10.37.12.18", true))
	hostfile.Hosts.Add(hostess.NewHostname("m.devsite.com", "10.37.12.18", true))
	f := hostfile.Format()
	if f != expected {
		t.Errorf(asserts, expected, f)
	}
}

func TestTrimWS(t *testing.T) {
	const expected = `  candy

	`
	got := hostess.TrimWS(expected)
	if got != "candy" {
		t.Errorf(asserts, expected, got)
	}
}

func TestListDomainsByIp(t *testing.T) {
	hostfile := hostess.NewHostfile("./hosts")
	hostfile.Hosts.Add(hostess.NewHostname("devsite.com", "10.37.12.18", true))
	hostfile.Hosts.Add(hostess.NewHostname("m.devsite.com", "10.37.12.18", true))
	hostfile.Hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", false))

	names := hostfile.Hosts.ListDomainsByIP(net.ParseIP("10.37.12.18"))
	if !(names[0].Domain == "devsite.com" && names[1].Domain == "m.devsite.com") {
		t.Errorf("Expected devsite.com and m.devsite.com. Got %s", names)
	}

	hostfile2 := hostess.NewHostfile("./hosts")
	hostfile2.Hosts.Add(hostess.NewHostname("localhost", "127.0.0.1", true))
	hostfile2.Hosts.Add(hostess.NewHostname("ip-10-37-12-18", "127.0.1.1", true))
	hostfile2.Hosts.Add(hostess.NewHostname("devsite", "127.0.0.1", true))

	names2 := hostfile2.Hosts.ListDomainsByIP(net.ParseIP("127.0.0.1"))
	if !(names2[0].Domain == "localhost" && names2[1].Domain == "devsite") {
		t.Errorf("Expected localhost and devsite. Got %s", names2)
	}
}

func TestParseLine(t *testing.T) {
	var hosts = hostess.NewHostlist()

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
	if !hosts.ContainsHostname(hostess.NewHostname("test.domain.com", "66.33.99.11", false)) ||
		len(hosts) != 1 {
		t.Error("Expected to find test.domain.com (disabled)")
	}

	hosts = hostess.ParseLine("#  66.33.99.11	test.domain.com	domain.com")
	if !hosts.ContainsHostname(hostess.NewHostname("test.domain.com", "66.33.99.11", false)) ||
		!hosts.ContainsHostname(hostess.NewHostname("domain.com", "66.33.99.11", false)) ||
		len(hosts) != 2 {
		t.Error("Expected to find domain.com and test.domain.com (disabled)")
		t.Errorf("Found %s", hosts)
	}

	// Not Commented stuff
	hosts = hostess.ParseLine("255.255.255.255 broadcasthost test.domain.com	domain.com")
	if !hosts.ContainsHostname(hostess.NewHostname("broadcasthost", "255.255.255.255", true)) ||
		!hosts.ContainsHostname(hostess.NewHostname("test.domain.com", "255.255.255.255", true)) ||
		!hosts.ContainsHostname(hostess.NewHostname("domain.com", "255.255.255.255", true)) ||
		len(hosts) != 3 {
		t.Error("Expected to find broadcasthost, domain.com, and test.domain.com (enabled)")
	}

	// Ipv6 stuff
	hosts = hostess.ParseLine("::1             localhost")
	if !hosts.ContainsHostname(hostess.NewHostname("localhost", "::1", true)) ||
		len(hosts) != 1 {
		t.Error("Expected to find localhost ipv6 (enabled)")
	}

	hosts = hostess.ParseLine("ff02::1 ip6-allnodes")
	if !hosts.ContainsHostname(hostess.NewHostname("ip6-allnodes", "ff02::1", true)) ||
		len(hosts) != 1 {
		t.Error("Expected to find ip6-allnodes ipv6 (enabled)")
	}
}

func TestLoadHostfile(t *testing.T) {
	hostfile := hostess.NewHostfile(hostess.GetHostsPath())
	data := hostfile.Load()
	if !strings.Contains(data, domain) {
		t.Errorf("Expected to find %s", domain)
	}

	hostfile.Parse()
	hostname := hostess.NewHostname(domain, ip, enabled)
	_, found := hostfile.Hosts[hostname.Domain]
	if !found {
		t.Errorf("Expected to find %s", hostname)
	}
}
