package hostess

import (
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
	hostfile := NewHostfile("./hosts")
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

func TestHostFileDuplicates(t *testing.T) {
	hostfile := NewHostfile("./hosts")

	const exp_duplicate = "Duplicate hostname entry for localhost -> 127.0.0.1"
	hostfile.Add(Hostname{domain, ip, true})
	err := hostfile.Add(Hostname{domain, ip, true})
	if err.Error() != exp_duplicate {
		t.Errorf(asserts, exp_duplicate, err)
	}

	const exp_conflict = "Conflicting hostname entries for localhost -> 127.0.1.1 and -> 127.0.0.1"
	err2 := hostfile.Add(Hostname{domain, "127.0.1.1", true})
	if err2.Error() != exp_conflict {
		t.Errorf(asserts, exp_conflict, err2)
	}

}

func TestFormatHostname(t *testing.T) {
	hostname := Hostname{domain, ip, enabled}

	const exp_enabled = "127.0.0.1 localhost"
	if hostname.Format() != exp_enabled {
		t.Errorf(asserts, hostname.Format(), exp_enabled)
	}

	hostname.Enabled = false
	const exp_disabled = "# 127.0.0.1 localhost"
	if hostname.Format() != exp_disabled {
		t.Errorf(asserts, hostname.Format(), exp_disabled)
	}
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

	hostfile := NewHostfile("./hosts")
	hostfile.Add(Hostname{"localhost", "127.0.0.1", true})
	hostfile.Add(Hostname{"ip-10-37-12-18", "127.0.1.1", true})
	hostfile.Add(Hostname{"devsite", "127.0.0.1", true})
	hostfile.Add(Hostname{"google.com", "8.8.8.8", false})
	hostfile.Add(Hostname{"devsite.com", "10.37.12.18", true})
	hostfile.Add(Hostname{"m.devsite.com", "10.37.12.18", true})
	f := hostfile.Format()
	if f != expected {
		t.Errorf(asserts, expected, f)
	}
}

func TestTrimWS(t *testing.T) {
	const expected = `  candy

	`
	got := TrimWS(expected)
	if got != "candy" {
		t.Errorf(asserts, expected, got)
	}
}

func TestListDomainsByIp(t *testing.T) {
	hostfile := NewHostfile("./hosts")
	hostfile.Add(Hostname{"devsite.com", "10.37.12.18", true})
	hostfile.Add(Hostname{"m.devsite.com", "10.37.12.18", true})
	hostfile.Add(Hostname{"google.com", "8.8.8.8", false})

	names := hostfile.ListDomainsByIp("10.37.12.18")
	if !(names[0] == "devsite.com" && names[1] == "m.devsite.com") {
		t.Errorf("Expected devsite.com and m.devsite.com. Got %s", names)
	}

	hostfile2 := NewHostfile("./hosts")
	hostfile2.Add(Hostname{"localhost", "127.0.0.1", true})
	hostfile2.Add(Hostname{"ip-10-37-12-18", "127.0.1.1", true})
	hostfile2.Add(Hostname{"devsite", "127.0.0.1", true})

	names2 := hostfile2.ListDomainsByIp("127.0.0.1")
	if !(names2[0] == "localhost" && names2[1] == "devsite") {
		t.Errorf("Expected localhost and devsite. Got %s", names2)
	}
}
