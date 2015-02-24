package hostess_test

import (
	"github.com/cbednarski/hostess"
	"testing"
)

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
