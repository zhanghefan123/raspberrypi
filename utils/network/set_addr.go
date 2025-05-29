//go:build linux

package network

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

//// 设置 ipv6 地址
//ipv6Addr := normalNode.IfNameToInterfaceMap[ifName].SourceIpv6Addr
//ipv6, _ := netlink.ParseAddr(ipv6Addr)
//if err = netlink.AddrAdd(veth, ipv6); err != nil {
//fmt.Printf("netlink.AddrAdd(%s) failed: %v", ipv6, err)
//return fmt.Errorf("netlink.AddrAdd(%s) failed: %w", ipv6, err)
//}

func SetAddr(interfaceName string, interfaceCidrAddr string, addrType string) error {
	networkIntf, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return fmt.Errorf("failed to get interface %s: %v", interfaceName, err)
	}

	if addrType == "ipv4" {
		err = SetIpv4Addr(&networkIntf, interfaceName, interfaceCidrAddr)
		if err != nil {
			return fmt.Errorf("failed to set ipv4 addr: %v", err)
		}
	} else if addrType == "ipv6" {
		err = SetIPv6Addr(&networkIntf, interfaceName, interfaceCidrAddr)
		if err != nil {
			return fmt.Errorf("failed to set ipv6 addr: %v", err)
		}
	} else {
		return fmt.Errorf("unsupported address type: %s", addrType)
	}

	return nil
}

func SetIpv4Addr(networkIntf *netlink.Link, interfaceName string, ipv4Addr string) error {
	var addrs []netlink.Addr
	var err error
	addrs, err = netlink.AddrList(*networkIntf, netlink.FAMILY_V4)
	for _, addr := range addrs {
		if addr.IP.String() == ipv4Addr[:len(ipv4Addr)-3] {
			fmt.Println("address already exists")
			return nil
		}
	}
	ipv4, err := netlink.ParseAddr(ipv4Addr)

	if err != nil {
		return fmt.Errorf("failed to parse address %s: %v", ipv4Addr, err)
	}
	if err = netlink.AddrAdd(*networkIntf, ipv4); err != nil {
		return fmt.Errorf("failed to add address %s to interface %s: %v", ipv4Addr, interfaceName, err)
	}
	return nil
}

func SetIPv6Addr(networkIntf *netlink.Link, interfaceName string, ipv6Addr string) error {
	var addrs []netlink.Addr
	var err error
	addrs, err = netlink.AddrList(*networkIntf, netlink.FAMILY_V6)
	for _, addr := range addrs {
		if addr.IP.String() == ipv6Addr {
			fmt.Println("address already exists")
			return nil
		}
	}
	ipv6, err := netlink.ParseAddr(ipv6Addr)

	if err != nil {
		return fmt.Errorf("failed to parse address %s: %v", ipv6Addr, err)
	}
	if err = netlink.AddrAdd(*networkIntf, ipv6); err != nil {
		return fmt.Errorf("failed to add address %s to interface %s: %v", ipv6Addr, interfaceName, err)
	}
	return nil
}
