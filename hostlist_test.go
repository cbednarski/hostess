package hostess_test

import (
	"github.com/cbednarski/hostess"
	"net"
	"testing"
)

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
	if !hosts.ContainsIp(first_ip) {
		t.Errorf("Expected to find %s", ip)
	}

	var extra_ip = net.ParseIP("1.2.3.4")
	if hosts.ContainsIp(extra_ip) {
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
	if hosts.Format() != expected {
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
	hosts := hostess.NewHostlist()
	hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", true))
	hosts.Add(hostess.NewHostname("google3.com", "::1", true))
	hosts.Add(hostess.NewHostname(domain, ip, false))
	hosts.Add(hostess.NewHostname("google2.com", "8.8.4.4", true))

	hosts.Sort()
	if (*hosts)[0].Domain != "localhost" {
		t.Error("Expected localhost to be first")
	}
	if (*hosts)[1].Domain != "google2.com" {
		t.Error("Expected google2 to be second")
	}
	if (*hosts)[2].Domain != "google.com" {
		t.Error("Expected google3 to be third")
	}
	if (*hosts)[3].Domain != "google3.com" {
		t.Error("Expected google3 to be last")
	}
	hosts.Format()
}
