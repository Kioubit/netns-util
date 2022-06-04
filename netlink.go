package main

import (
	"github.com/vishvananda/netlink"
)

func setLinkToNetns(link string, netns string) error {
	iface, err := netlink.LinkByName(link)
	if err != nil {
		return err
	}

	fd, err := getNsByName(netns)
	if err != nil {
		return err
	}
	netlink.LinkSetNsFd(iface, int(fd))
	if err != nil {
		return err
	}
	return nil
}
