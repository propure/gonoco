package main

import (
	"fmt"
	"net"
)

func main() {
	ipaddr := "192.168.100.254"
	//protocol := "icmp"
	protocol := "tcp"
	netaddr, err := net.ResolveIPAddr("ip4", ipaddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	//如果第二個string參數netaddr設置成零值，那么就監控所有源的網絡包
	conn, err := net.ListenIP("ip4:"+protocol, netaddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	buf := make([]byte, 1024)
	for {
		numRead, recvAddr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		if recvAddr != nil {
			fmt.Printf("raddr: %v\n", recvAddr)
		}
		fmt.Printf("% X\n", buf[:numRead])
	}
}
