package hostess_test

import (
	"fmt"
	"github.com/cbednarski/hostess"
	"net"
	"testing"
)

func TestMakeSurrogateIP(t *testing.T) {
	const orig = "127.0.0.1"
	const expected1 = "0.0.0.1"
	IP1 := string(hostess.MakeSurrogateIP(net.IP(orig)))
	if IP1 != expected1 {
		t.Errorf("Expected %s to convert to %s; got %s", orig, expected1, IP1)
	}

	const expected2 = "10.20.30.40"
	IP2 := string(hostess.MakeSurrogateIP(net.IP(expected2)))
	if IP2 != expected2 {
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

	hosts.Remove(1)
	if len(*hosts) > 1 {
		t.Errorf("Expected hostlist to have 1 item, found %d", len(*hosts))
	}

	if hosts.ContainsDomain("google.com") {
		t.Errorf("Expected not to find google.com")
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

	hosts.Sort()
	if (*hosts)[0].Domain != "localhost" {
		t.Error("Expected localhost to be first")
		t.Error(hosts.Format())
	}
	if (*hosts)[1].Domain != "google2.com" {
		t.Error("Expected google2 to be second")
		t.Error(hosts.Format())
	}
	if (*hosts)[2].Domain != "google.com" {
		t.Error("Expected google3 to be third")
		t.Error(hosts.Format())
	}
	if (*hosts)[3].Domain != "blah" {
		t.Error("Expected blah to be fourth")
		t.Error(hosts.Format())
	}
	if (*hosts)[4].Domain != "blah2" {
		t.Error("Expected blah2 to be fifth")
		t.Error(hosts.Format())
	}
	if (*hosts)[5].Domain != "blah3" {
		t.Error("Expected blah3 to be sixth")
		t.Error(hosts.Format())
	}
	if (*hosts)[6].Domain != "blah33" {
		t.Error("Expected blah33 to be seventh")
		t.Error(hosts.Format())
	}
	// IPv6 Domains
	if (*hosts)[7].Domain != "google3.com" {
		t.Error("Expected google3 to be eigth")
		t.Error(hosts.Format())
	}
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
