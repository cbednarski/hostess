package hostess_test

import (
	"github.com/cbednarski/hostess"
	"net"
	"testing"
)

func TestContainsDomainIp(t *testing.T) {
	hosts := []*hostess.Hostname{
		&hostess.Hostname{domain, ip, false, false},
		&hostess.Hostname{"google.com", net.ParseIP("8.8.8.8"), true, false},
	}

	if !hostess.ContainsDomain(hosts, domain) {
		t.Errorf("Expected to find %s", domain)
	}

	const extra_domain = "yahoo.com"
	if hostess.ContainsDomain(hosts, extra_domain) {
		t.Errorf("Did not expect to find %s", extra_domain)
	}

	if !hostess.ContainsIp(hosts, ip) {
		t.Errorf("Expected to find %s", ip)
	}

	var extra_ip = net.ParseIP("1.2.3.4")
	if hostess.ContainsIp(hosts, extra_ip) {
		t.Errorf("Did not expect to find %s", extra_ip)
	}

	hostname := &hostess.Hostname{domain, ip, true, false}
	if !hostess.ContainsHostname(hosts, hostname) {
		t.Errorf("Expected to find %s", hostname)
	}

	extra_hostname := &hostess.Hostname{"yahoo.com", net.ParseIP("4.3.2.1"), false, false}
	if hostess.ContainsHostname(hosts, extra_hostname) {
		t.Errorf("Did not expect to find %s", extra_hostname)
	}
}
