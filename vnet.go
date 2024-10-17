package vnet

import (
	"fmt"
	"net"
	"sudp"
	"utun"

	"golang.org/x/net/ipv4"
)

type VnetSwitch struct {
	self  net.IP
	iface *utun.Utun
	sock  *sudp.ServerConn
	route *RouteTable
}
type VnetClient struct {
	self  net.IP
	iface *utun.Utun
	sock  *sudp.ClientConn
}

func NewVnetSwitch(cird string, peer string, serverConfig string) (*VnetSwitch, error) {
	laddr, raddrs, e := sudp.ParseConfig(serverConfig)
	if e != nil {
		return nil, e
	}

	server, e := sudp.Listen(laddr, raddrs)
	if e != nil {
		return nil, e
	}
	iface, e := createIf(cird, peer)
	if e != nil {
		return nil, e
	}

	self, _, e := net.ParseCIDR(cird)

	vnets := VnetSwitch{
		self:  self,
		iface: iface,
		sock:  server,
		route: &RouteTable{},
	}
	go vnets.tx()
	go vnets.rx()
	return &vnets, nil
}

func NewVnetClient(cird string, peer string, laddr *sudp.LocalAddr, raddr *sudp.RemoteAddr) (*VnetClient, error) {
	client, e := sudp.Connect(laddr, raddr)
	if e != nil {
		return nil, e
	}
	iface, e := createIf(cird, peer)
	if e != nil {
		return nil, e
	}

	self, _, e := net.ParseCIDR(cird)
	vnetc := VnetClient{
		self:  self,
		iface: iface,
		sock:  client,
	}
	go vnetc.rx()
	go vnetc.tx()
	return &vnetc, nil
}

func (v *VnetClient) tx() {
	buff := make([]byte, v.iface.MTU)
	for {
		n, e := v.iface.Read(buff)
		if e != nil {
			panic(e)
		}
		fmt.Println(ipv4.ParseHeader(buff[4:n]))
		if e := v.sock.Send(buff[4:n]); e != nil {
			// Drop
			fmt.Println(e)
			continue
		}
	}
}

func (v *VnetClient) rx() {
	for {
		buff, e := v.sock.Recv()
		if e != nil {
			panic(e)
		}
		fmt.Println("RX", buff)
		if _, e := v.iface.Write(buff); e != nil {
			panic(e)
		}
	}
}

func (v *VnetSwitch) tx() {
	buff := make([]byte, v.iface.MTU)
	for {
		n, e := v.iface.Read(buff)
		if e != nil {
			panic(e)
		}
		ip, e := ipv4.ParseHeader(buff)
		fmt.Println("TX:", ip)
		if e != nil {
			// Drop
			//continue
		}
		vaddr, ok := v.route.GetDestination(ip.Dst)
		if !ok {
			// Drop
			//continue
		}
		vaddr = 1001
		if e := v.sock.SendTo(buff[0:n], vaddr); e != nil {
			// Drop
			fmt.Println(e)
			continue
		}
	}
}

func (v *VnetSwitch) rx() {
	for {
		buff, vaddr, e := v.sock.RecvFrom()
		if e != nil {
			panic(e)
		}
		ip, e := ipv4.ParseHeader(buff)
		if e != nil {
			continue
		}
		fmt.Println(ip)
		v.route.SetDestination(ip.Src, vaddr)
		if vaddr, ok := v.route.GetDestination(ip.Dst); !ok || ip.Dst.Equal(v.self) {
			fmt.Println("Escribiendo en la interfaz")
			if _, e := v.iface.Write(buff); e != nil {
				panic(e)
			}
		} else {
			if e := v.sock.SendTo(buff, vaddr); e != nil {
				fmt.Println(e)
				continue
			}
		}
	}
}
