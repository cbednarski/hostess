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
	if !hosts.ContainsHostname(hostname) {
		t.Errorf("Expected to find %s", hostname)
	}

	extra_hostname := hostess.NewHostname("yahoo.com", "4.3.2.1", false)
	if hosts.ContainsHostname(extra_hostname) {
		t.Errorf("Did not expect to find %s", extra_hostname)
	}
}

func TestRemove(t *testing.T) {
	hosts := hostess.NewHostlist()
	hosts.Add(hostess.NewHostname(domain, ip, false))
	hosts.Add(hostess.NewHostname("google.com", "8.8.8.8", true))

	hosts.Remove(hosts.IndexOfDomainIpv4("google.com"))

	// if len(hosts) > 1 {
	// 	t.Errorf("Expected hostlist to have one item, found %s", len(hosts))
	// }

	if hosts.ContainsDomain("google.com") {
		t.Errorf("Expected not to find google.com")
	}
}
