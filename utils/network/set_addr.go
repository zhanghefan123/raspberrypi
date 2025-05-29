//go:build linux

package network

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func SetAddr(interfaceName string, interfaceCidrAddr string) error {
	networkIntf, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return fmt.Errorf("failed to get interface %s: %v", interfaceName, err)
	}
	addrs, err := netlink.AddrList(networkIntf, netlink.FAMILY_V4)
	for _, addr := range addrs {
		if addr.IP.String() == interfaceCidrAddr[:len(interfaceCidrAddr)-3] {
			fmt.Println("address already exists")
			return nil
		}
	}
	ipv4, err := netlink.ParseAddr(interfaceCidrAddr)

	if err != nil {
		return fmt.Errorf("failed to parse address %s: %v", interfaceCidrAddr, err)
	}
	if err = netlink.AddrAdd(networkIntf, ipv4); err != nil {
		return fmt.Errorf("failed to add address %s to interface %s: %v", interfaceCidrAddr, interfaceName, err)
	}
	return nil
}
