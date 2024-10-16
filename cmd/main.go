package main

import (
	"fmt"
	"net"
	"sudp"
	"vnet"
)

func main() {
	pub, e := sudp.PublicKeyFromPemFile("sdtl_public.pem")
	if e != nil {
		fmt.Println(e)
		return
	}
	pri, e := sudp.PrivateFromPemFile("private.pem")
	if e != nil {
		fmt.Println(e)
		return
	}

	a, _ := net.ResolveUDPAddr("udp4", "18.212.245.20:7000")
	s, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:")
	laddr := sudp.LocalAddr{
		VirtualAddress: 1001,
		NetworkAddress: s,
		PrivateKey:     pri,
	}
	raddr := sudp.RemoteAddr{
		VirtualAddress: 0,
		NetworkAddress: a,
		PublicKey:      pub,
	}
	_, err := vnet.NewVnetClient("10.0.0.2/24", "10.0.0.1", &laddr, &raddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	end := make(chan bool)
	<-end
}
