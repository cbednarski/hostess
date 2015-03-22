package hostess_test

import (
	"github.com/cbednarski/hostess"
	"net"
	"testing"
)

func TestHostname(t *testing.T) {
	h := hostess.NewHostname(domain, ip, enabled)

	if h.Domain != domain {
		t.Errorf("Domain should be %s", domain)
	}
	if !h.IP.Equal(net.ParseIP(ip)) {
		t.Errorf("IP should be %s", ip)
	}
	if h.Enabled != enabled {
		t.Errorf("Enabled should be %t", enabled)
	}
}

func TestEqual(t *testing.T) {
	a := hostess.NewHostname("localhost", "127.0.0.1", true)
	b := hostess.NewHostname("localhost", "127.0.0.1", false)
	c := hostess.NewHostname("localhost", "127.0.1.1", false)

	if !a.Equal(b) {
		t.Errorf("%s and %s should be equal", a, b)
	}
	if a.Equal(c) {
		t.Errorf("%s and %s should not be equal", a, c)
	}
}

func TestEqualIP(t *testing.T) {
	a := hostess.NewHostname("localhost", "127.0.0.1", true)
	c := hostess.NewHostname("localhost", "127.0.1.1", false)
	ip := net.ParseIP("127.0.0.1")

	if !a.EqualIP(ip) {
		t.Errorf("%s and %s should be equal", a.IP, ip)
	}
	if a.EqualIP(c.IP) {
		t.Errorf("%s and %s should not be equal", a.IP, c.IP)
	}
}

func TestIsValid(t *testing.T) {
	a := hostess.NewHostname("localhost", "127.0.0.1", true)
	d := hostess.NewHostname("", "127.0.0.1", true)
	e := hostess.NewHostname("localhost", "localhost", true)

	if !a.IsValid() {
		t.Errorf("%s should be a valid hostname", a)
	}
	if d.IsValid() {
		t.Errorf("%s should be invalid because the name is blank", d)
	}
	if e.IsValid() {
		t.Errorf("%s should be invalid because the ip is malformed", e)
	}
}

func TestFormatHostname(t *testing.T) {
	hostname := hostess.NewHostname(domain, ip, enabled)

	const exp_enabled = "127.0.0.1 localhost"
	if hostname.Format() != exp_enabled {
		t.Errorf("Hostname format doesn't match desired output: %s", Diff(hostname.Format(), exp_enabled))
	}

	hostname.Enabled = false
	const exp_disabled = "# 127.0.0.1 localhost"
	if hostname.Format() != exp_disabled {
		t.Errorf("Hostname format doesn't match desired output: %s", Diff(hostname.Format(), exp_disabled))
	}
}

func TestFormatEnabled(t *testing.T) {
	hostname := hostess.NewHostname(domain, ip, enabled)
	const expectedOn = "(On)"
	if hostname.FormatEnabled() != expectedOn {
		t.Errorf("Expected hostname to be turned %s", expectedOn)
	}
	const expectedHumanOn = "localhost -> 127.0.0.1 (On)"
	if hostname.FormatHuman() != expectedHumanOn {
		t.Errorf("Unexpected output%s", Diff(expectedHumanOn, hostname.FormatHuman()))
	}

	hostname.Enabled = false
	if hostname.FormatEnabled() != "(Off)" {
		t.Error("Expected hostname to be turned (Off)")
	}
}
