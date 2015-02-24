package hostess

import (
	"fmt"
	"regexp"
	"strings"
)

func TrimWS(s string) string {
	return strings.Trim(s, " \n\t")
}

var ipv4_pattern = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)

func LooksLikeIpv4(ip string) bool {
	return ipv4_pattern.MatchString(ip)
}

var ipv6_pattern = regexp.MustCompile(`^[a-z0-9:]+$`)

func LooksLikeIpv6(ip string) bool {
	if !strings.Contains(ip, ":") {
		return false
	}
	return ipv6_pattern.MatchString(ip)
}

type Hostname struct {
	Domain  string
	Ip      string
	Enabled bool
	// Ipv6    bool
}

func (h *Hostname) Format() string {
	r := fmt.Sprintf("%s %s", h.Ip, h.Domain)
	if !h.Enabled {
		r = "# " + r
	}
	return r
}

func (a *Hostname) Equals(b Hostname) bool {
	if a.Domain == b.Domain && a.Ip == b.Ip {
		return true
	}
	return false
}
