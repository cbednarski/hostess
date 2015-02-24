package hostess

func ContainsHostname(hostnames []Hostname, hostname Hostname) bool {
	for _, v := range hostnames {
		if v.Ip == hostname.Ip && v.Domain == hostname.Domain {
			return true
		}
	}
	return false
}

func ContainsDomain(hostnames []Hostname, domain string) bool {
	for _, v := range hostnames {
		if v.Domain == domain {
			return true
		}
	}
	return false
}

func ContainsIp(hostnames []Hostname, ip string) bool {
	for _, v := range hostnames {
		if v.Ip == ip {
			return true
		}
	}
	return false
}
