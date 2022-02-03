package socket

import (
	"fmt"
	"net"
	"syscall"
)

func tcpserver_syscall() {

	var (
		sockfd int
		addr   syscall.SockaddrInet4
		err    error
	)

	if sockfd, err = syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, syscall.IPPROTO_TCP); err != nil {
		fmt.Println("syscall.Socket() error: ", err.Error())
		return
	}
	defer syscall.Shutdown(sockfd, syscall.SHUT_RDWR)
	if err = syscall.SetNonblock(sockfd, true); err != nil {
		syscall.Close(sockfd)
		return
	}
	//addr := syscall.SockaddrInet4{Port: 10101}
	copy(addr.Addr[:], net.ParseIP("0.0.0.0").To4())
	addr.Port = 10101

	syscall.Bind(sockfd, &addr)
	syscall.Listen(sockfd, 10)
}
