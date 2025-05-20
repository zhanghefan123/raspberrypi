package network

import (
	"fmt"
	"github.com/vishvananda/netlink"
)

func AddRoute(DestinationNetworkSegment, Gateway string) error {
	// 进行路由的添加
	// 1. 进行网段的解析
	network, err := netlink.ParseIPNet(DestinationNetworkSegment)
	if err != nil {
		return fmt.Errorf("failed to parse network %s: %v", DestinationNetworkSegment, err)
	}
	// 2. 进行网关的解析
	gatewayIP, err := netlink.ParseAddr(Gateway)
	if err != nil {
		return fmt.Errorf("failed to parse gateway %s: %v", Gateway, err)
	}
	// 3. 创建路由对象
	route := &netlink.Route{
		Dst: network,
		Gw:  gatewayIP.IP,
	}
	err = netlink.RouteAdd(route)
	if err != nil {
		return fmt.Errorf("failed to add route %s via %s: %v", DestinationNetworkSegment, Gateway, err)
	}
	return nil
}
