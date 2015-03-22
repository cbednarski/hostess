package hostess_test

import (
	"fmt"
	"github.com/cbednarski/hostess"
	"net"
	"testing"
)

func TestAddDuplicate(t *testing.T) {
	list := hostess.NewHostlist()

	hostname := hostess.NewHostname("mysite", "1.2.3.4", false)
	err := list.Add(hostname)
	if err != nil {
		t.Error("Expected no errors when adding a hostname for the first time")
	}

	hostname.Enabled = true
	err = list.Add(hostname)
	if err == nil {
		t.Error("Expected error when adding a duplicate")
	}
	if !(*list)[0].Enabled {
		t.Error("Expected hostname to be in enabled state")
	}
}

func TestAddConflict(t *testing.T) {
	hostnameA := hostess.NewHostname("mysite", "1.2.3.4", true)
	hostnameB := hostess.NewHostname("mysite", "5.2.3.4", false)

	list := hostess.NewHostlist()
	list.Add(hostnameA)
	err := list.Add(hostnameB)
	if err == nil {
		t.Error("Expected conflict error")
	}
	if !(*list)[0].Equal(hostnameB) {
		t.Error("Expected second hostname to overwrite")
	}
	if (*list)[0].Enabled {
		t.Error("Expected second hostname to be disabled")
	}
}

func TestMakeSurrogateIP(t *testing.T) {
	original := net.ParseIP("127.0.0.1")
	expected1 := net.ParseIP("0.0.0.1")
	IP1 := hostess.MakeSurrogateIP(original)
	if !IP1.Equal(expected1) {
		t.Errorf("Expected %s to convert to %s; got %s", original, expected1, IP1)
	}

	expected2 := net.ParseIP("10.20.30.40")
	IP2 := hostess.MakeSurrogateIP(expected2)
	if !IP2.Equal(expected2) {
		t.Errorf("Expected %s to remain unchanged; got %s", expected2, IP2)
	}
}

func TestContainsDomainIp(t *testing.T) {
	hosts := hostess.NewHostlist()
	hosts.Add(hostess.NewHostname(domain, ip, false))
	hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", true))

	if !hosts.ContainsDomain(domain) {
		t.Errorf("Expected to find %s", domain)
	}

	const extra_domain = "yahoo.com"
	if hosts.ContainsDomain(extra_domain) {
		t.Errorf("Did not expect to find %s", extra_domain)
	}

	var first_ip = net.ParseIP(ip)
	if !hosts.ContainsIP(first_ip) {
		t.Errorf("Expected to find %s", ip)
	}

	var extra_ip = net.ParseIP("1.2.3.4")
	if hosts.ContainsIP(extra_ip) {
		t.Errorf("Did not expect to find %s", extra_ip)
	}

	hostname := hostess.NewHostname(domain, ip, true)
	if !hosts.Contains(hostname) {
		t.Errorf("Expected to find %s", hostname)
	}

	extra_hostname := hostess.NewHostname("yahoo.com", "4.3.2.1", false)
	if hosts.Contains(extra_hostname) {
		t.Errorf("Did not expect to find %s", extra_hostname)
	}
}

func TestFormat(t *testing.T) {
	hosts := hostess.NewHostlist()
	hosts.Add(hostess.NewHostname(domain, ip, false))
	hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", true))

	expected := `# 127.0.0.1 localhost
8.8.8.8 google.com
`
	if string(hosts.Format()) != expected {
		t.Error("Formatted hosts list is not formatted correctly")
	}
}

func TestRemove(t *testing.T) {
	hosts := hostess.NewHostlist()
	hosts.Add(hostess.NewHostname(domain, ip, false))
	hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", true))

	removed := hosts.Remove(1)
	if removed != 1 {
		t.Error("Expected to remove 1 item")
	}
	if len(*hosts) > 1 {
		t.Errorf("Expected hostlist to have 1 item, found %d", len(*hosts))
	}
	if hosts.ContainsDomain("google.com") {
		t.Errorf("Expected not to find google.com")
	}

	hosts.Add(hostess.NewHostname(domain, "::1", enabled))
	removed = hosts.RemoveDomain(domain)
	if removed != 2 {
		t.Error("Expected to remove 2 items")
	}
}

func TestRemoveDomain(t *testing.T) {
	hosts := hostess.NewHostlist()
	h1 := hostess.NewHostname("google.com", "127.0.0.1", false)
	h2 := hostess.NewHostname("google.com", "::1", true)
	hosts.Add(h1)
	hosts.Add(h2)

	hosts.RemoveDomainV("google.com", 4)
	if hosts.Contains(h1) {
		t.Error("Should not contain ipv4 hostname")
	}
	if !hosts.Contains(h2) {
		t.Error("Should still contain ipv6 hostname")
	}

	hosts.RemoveDomainV("google.com", 6)
	if len(*hosts) != 0 {
		t.Error("Should no longer contain any hostnames")
	}
}

func CheckIndexDomain(t *testing.T, index int, domain string, hosts *hostess.Hostlist) {
	if (*hosts)[index].Domain != domain {
		t.Errorf("Expected %s to be in position %d. Found: %s", domain, index, (*hosts)[index].FormatHuman())
	}
}

func TestSort(t *testing.T) {
	// Getting 100% coverage on this is kinda tricky. It's pretty close and
	// this is already too long.

	hosts := hostess.NewHostlist()
	hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", true))
	hosts.Add(hostess.NewHostname("google3.com", "::1", true))
	hosts.Add(hostess.NewHostname(domain, ip, false))
	hosts.Add(hostess.NewHostname("google2.com", "8.8.4.4", true))
	hosts.Add(hostess.NewHostname("blah2", "10.20.1.1", true))
	hosts.Add(hostess.NewHostname("blah3", "10.20.1.1", true))
	hosts.Add(hostess.NewHostname("blah33", "10.20.1.1", true))
	hosts.Add(hostess.NewHostname("blah", "10.20.1.1", true))
	hosts.Add(hostess.NewHostname("hostname", "127.0.1.1", true))
	hosts.Add(hostess.NewHostname("devsite", "127.0.0.1", true))

	hosts.Sort()

	CheckIndexDomain(t, 0, "localhost", hosts)
	CheckIndexDomain(t, 1, "devsite", hosts)
	CheckIndexDomain(t, 2, "hostname", hosts)
	CheckIndexDomain(t, 3, "google2.com", hosts)
	CheckIndexDomain(t, 4, "google.com", hosts)
	CheckIndexDomain(t, 5, "blah", hosts)
	CheckIndexDomain(t, 6, "blah2", hosts)
	CheckIndexDomain(t, 7, "blah3", hosts)
	CheckIndexDomain(t, 8, "blah33", hosts)
	CheckIndexDomain(t, 9, "google3.com", hosts)
}

func ExampleHostlist_1() {
	hosts := hostess.NewHostlist()
	hosts.Add(hostess.NewHostname("google.com", "127.0.0.1", false))
	hosts.Add(hostess.NewHostname("google.com", "::1", true))

	fmt.Printf("%s\n", hosts.Format())
	// Output:
	// # 127.0.0.1 google.com
	// ::1 google.com
}
