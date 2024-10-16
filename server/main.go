package main

import (
	"fmt"
	"utun"
	"vnet"
)

func main() {
	_, e := vnet.NewVnetSwitch("10.0.0.1/24", utun.NOPEER, "config.json")
	if e != nil {
		fmt.Println(e)
		return
	}
	end := make(chan bool)
	<-end

}
