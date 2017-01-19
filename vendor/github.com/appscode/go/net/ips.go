package net

import (
	"bytes"
	"errors"
	"net"
	"regexp"
)

type IPRange struct {
	from net.IP
	to   net.IP
}

func NewIPRange(from string, to string) IPRange {
	return IPRange{net.ParseIP(from).To4(), net.ParseIP(to).To4()}
}

var privateIPRanges []IPRange = []IPRange{
	NewIPRange("10.0.0.0", "10.255.255.255"),
	NewIPRange("172.16.0.0", "172.31.255.255"),
	NewIPRange("192.168.0.0", "192.168.255.255"),
}

func IsPrivateIP(ip net.IP) bool {
	for _, r := range privateIPRanges {
		if bytes.Compare(ip, r.from) >= 0 && bytes.Compare(ip, r.to) <= 0 {
			return true
		}
	}
	return false
}

var (
	knownLocalBridges = regexp.MustCompile(`^(docker|cbr|cni)[0-9]+$`)
)

func routableIPs() []net.IP {
	result := make([]net.IP, 0, 2)
	ifaces, err := net.Interfaces()
	if err != nil {
		return result
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		if knownLocalBridges.MatchString(iface.Name) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			result = append(result, ip)
		}
	}
	return result
}

var (
	InterfaceDownErr       = errors.New("Interface down")
	LoopbackInterfaceErr   = errors.New("Loopback interface")
	KnownLocalInterfaceErr = errors.New("Known local interface")
	NotFoundErr            = errors.New("No IPV4 address found!")
)

/*
NodeIP returns a IPv4 address for a given set of interface names. It always prefers a private IP over a public IP.
If no interface name is given, all interfaces are checked.
*/
func NodeIP(interfaceName ...string) (string, net.IP, error) {
	var err error
	var ifaces []net.Interface

	if len(interfaceName) == 0 {
		ifaces, err = net.Interfaces()
		if err != nil {
			return "", nil, err
		}
	} else {
		ifaces = make([]net.Interface, len(interfaceName))
		for i, name := range interfaceName {
			d, err := net.InterfaceByName(name)
			if err != nil {
				return name, nil, err
			}
			ifaces[i] = *d
		}
	}

	type ipData struct {
		ip    net.IP
		iface string
	}
	internalIPs := make([]ipData, 0)
	externalIPs := make([]ipData, 0)
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			if len(ifaces) == 1 {
				return iface.Name, nil, InterfaceDownErr
			} else {
				continue
			}
		}
		if iface.Flags&net.FlagLoopback != 0 {
			if len(ifaces) == 1 {
				return iface.Name, nil, LoopbackInterfaceErr
			} else {
				continue
			}
		}
		if knownLocalBridges.MatchString(iface.Name) {
			if len(ifaces) == 1 {
				return iface.Name, nil, KnownLocalInterfaceErr
			} else {
				continue
			}
		}
		addrs, err := iface.Addrs()
		if err != nil {
			if len(ifaces) == 1 {
				return iface.Name, nil, err
			} else {
				continue
			}
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // Not an ipv4 address
			}
			if IsPrivateIP(ip) {
				internalIPs = append(internalIPs, ipData{ip: ip, iface: iface.Name})
			} else {
				externalIPs = append(externalIPs, ipData{ip: ip, iface: iface.Name})
			}
		}
	}
	if len(internalIPs) > 0 {
		return internalIPs[0].iface, internalIPs[0].ip, nil
	} else if len(externalIPs) > 0 {
		return externalIPs[0].iface, externalIPs[0].ip, nil
	} else {
		return "", nil, NotFoundErr
	}
}

func GetInternalIP() string {
	ips := routableIPs()
	for _, ip := range ips {
		if IsPrivateIP(ip) {
			return ip.String()
		}
	}
	return ""
}

func GetExternalIPs() []string {
	ret := make([]string, 0)
	ips := routableIPs()
	for _, ip := range ips {
		if !IsPrivateIP(ip) {
			ret = append(ret, ip.String())
		}
	}
	return ret
}
