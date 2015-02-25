package hostess

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var ipv4_pattern = regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`)
var ipv6_pattern = regexp.MustCompile(`^[a-z0-9:]+$`)

func LooksLikeIpv4(ip string) bool {
	return ipv4_pattern.MatchString(ip)
}

func LooksLikeIpv6(ip string) bool {
	if !strings.Contains(ip, ":") {
		return false
	}
	return ipv6_pattern.MatchString(ip)
}

type Hostname struct {
	Domain  string
	Ip      net.IP
	Enabled bool
	Ipv6    bool
}

func NewHostname(domain, ip string, enabled bool) (hostname *Hostname) {
	IP := net.ParseIP(ip)
	hostname = &Hostname{domain, IP, enabled, LooksLikeIpv6(ip)}
	return
}

func (h *Hostname) Equal(n *Hostname) bool {
	return h.Ip.Equal(n.Ip) && h.Domain == n.Domain
}

func (h *Hostname) EqualIp(ip net.IP) bool {
	return h.Ip.Equal(ip)
}

func (h *Hostname) IsValid() bool {
	return h.Domain != "" && h.Ip != nil
}

func (h *Hostname) Format() string {
	r := fmt.Sprintf("%s %s", h.Ip.String(), h.Domain)
	if !h.Enabled {
		r = "# " + r
	}
	return r
}

func (a *Hostname) Equals(b Hostname) bool {
	if a.Domain == b.Domain && a.Ip.Equal(b.Ip) {
		return true
	}
	return false
}
