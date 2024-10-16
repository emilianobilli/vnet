package vnet

import (
	"sudp"
	"utun"

	"golang.org/x/net/ipv4"
)

const (
	udplen   = 8
	overhead = ipv4.HeaderLen + 8 + sudp.HeaderLen + sudp.DataHeaderLen
)

func mtu(m int) int {
	return m - overhead
}

func createIf(cird string, peer string) (*utun.Utun, error) {
	iface, err := utun.OpenUtun()
	if err != nil {
		return nil, err
	}

	if err = iface.SetMTU(mtu(1500)); err != nil {
		return nil, err
	}
	if err = iface.SetIP(cird, peer); err != nil {
		return nil, err
	}
	return iface, nil
}
