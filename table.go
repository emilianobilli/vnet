package vnet

import (
	"encoding/binary"
	"net"
)

type RouteTable struct {
	table map[uint32]uint16
	defgw uint16
}

func (r *RouteTable) SetDestination(ip net.IP, vaddr uint16) {
	if r.table == nil {
		r.table = make(map[uint32]uint16)
	}
	ipb := binary.BigEndian.Uint32(ip.To4())
	r.table[ipb] = vaddr
}
func (r *RouteTable) GetDestination(ip net.IP) (uint16, bool) {
	if r.table == nil {
		r.table = make(map[uint32]uint16)
	}
	vaddr, ok := r.table[binary.BigEndian.Uint32(ip.To4())]
	return vaddr, ok
}
