package hostess_test

import (
	"github.com/cbednarski/hostess"
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
	hostfile.Add(hostess.Hostname{domain, ip, true})
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

func TestHostFileDuplicates(t *testing.T) {
	hostfile := hostess.NewHostfile("./hosts")

	const exp_duplicate = "Duplicate hostname entry for localhost -> 127.0.0.1"
	hostfile.Add(hostess.Hostname{domain, ip, true})
	err := hostfile.Add(hostess.Hostname{domain, ip, true})
	if err.Error() != exp_duplicate {
		t.Errorf(asserts, exp_duplicate, err)
	}

	const exp_conflict = "Conflicting hostname entries for localhost -> 127.0.1.1 and -> 127.0.0.1"
	err2 := hostfile.Add(hostess.Hostname{domain, "127.0.1.1", true})
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
	hostfile.Add(hostess.Hostname{"localhost", "127.0.0.1", true})
	hostfile.Add(hostess.Hostname{"ip-10-37-12-18", "127.0.1.1", true})
	hostfile.Add(hostess.Hostname{"devsite", "127.0.0.1", true})
	hostfile.Add(hostess.Hostname{"google.com", "8.8.8.8", false})
	hostfile.Add(hostess.Hostname{"devsite.com", "10.37.12.18", true})
	hostfile.Add(hostess.Hostname{"m.devsite.com", "10.37.12.18", true})
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
	hostfile.Add(hostess.Hostname{"devsite.com", "10.37.12.18", true})
	hostfile.Add(hostess.Hostname{"m.devsite.com", "10.37.12.18", true})
	hostfile.Add(hostess.Hostname{"google.com", "8.8.8.8", false})

	names := hostfile.ListDomainsByIp("10.37.12.18")
	if !(names[0] == "devsite.com" && names[1] == "m.devsite.com") {
		t.Errorf("Expected devsite.com and m.devsite.com. Got %s", names)
	}

	hostfile2 := hostess.NewHostfile("./hosts")
	hostfile2.Add(hostess.Hostname{"localhost", "127.0.0.1", true})
	hostfile2.Add(hostess.Hostname{"ip-10-37-12-18", "127.0.1.1", true})
	hostfile2.Add(hostess.Hostname{"devsite", "127.0.0.1", true})

	names2 := hostfile2.ListDomainsByIp("127.0.0.1")
	if !(names2[0] == "localhost" && names2[1] == "devsite") {
		t.Errorf("Expected localhost and devsite. Got %s", names2)
	}
}

func TestParseLine(t *testing.T) {
	var hosts = []hostess.Hostname{}

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
	if !hostess.ContainsHostname(hosts, hostess.Hostname{"test.domain.com", "66.33.99.11", false}) ||
		len(hosts) != 1 {
		t.Error("Expected to find test.domain.com (disabled)")
	}

	hosts = hostess.ParseLine("#  66.33.99.11	test.domain.com	domain.com")
	if !hostess.ContainsHostname(hosts, hostess.Hostname{"test.domain.com", "66.33.99.11", false}) ||
		!hostess.ContainsHostname(hosts, hostess.Hostname{"domain.com", "66.33.99.11", false}) ||
		len(hosts) != 2 {
		t.Error("Expected to find domain.com and test.domain.com (disabled)")
		t.Errorf("Found %s", hosts)
	}

	// Not Commented stuff
	hosts = hostess.ParseLine("255.255.255.255 broadcasthost test.domain.com	domain.com")
	if !hostess.ContainsHostname(hosts, hostess.Hostname{"broadcasthost", "255.255.255.255", true}) ||
		!hostess.ContainsHostname(hosts, hostess.Hostname{"test.domain.com", "255.255.255.255", true}) ||
		!hostess.ContainsHostname(hosts, hostess.Hostname{"domain.com", "255.255.255.255", true}) ||
		len(hosts) != 3 {
		t.Error("Expected to find broadcasthost, domain.com, and test.domain.com (enabled)")
	}

	// Ipv6 stuff
	hosts = hostess.ParseLine("::1             localhost")
	if !hostess.ContainsHostname(hosts, hostess.Hostname{"localhost", "::1", true}) ||
		len(hosts) != 1 {
		t.Error("Expected to find localhost ipv6 (enabled)")
	}

	hosts = hostess.ParseLine("ff02::1 ip6-allnodes")
	if !hostess.ContainsHostname(hosts, hostess.Hostname{"ip6-allnodes", "ff02::1", true}) ||
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
	hostname := hostess.Hostname{domain, ip, enabled}
	_, found := hostfile.Hosts[hostname.Domain]
	if !found {
		t.Errorf("Expected to find %s", hostname)
	}
}
