package hostess_test

import (
	"github.com/cbednarski/hostess"
	"testing"
)

func TestHostname(t *testing.T) {

	h := hostess.Hostname{}
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

func TestFormatHostname(t *testing.T) {
	hostname := hostess.Hostname{domain, ip, enabled}

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
