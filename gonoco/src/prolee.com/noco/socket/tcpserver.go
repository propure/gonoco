package socket //尽量和目录名一样，否则引用可能是import "A"而使用B.Xyz()来调用

import (
	"fmt"
	"net"
	"os"
	"time"
)

func Tcpserver() { //如果要在其它地方import这个函数，首字母必须大写

	/**
	 * type TCPAddr struct {
	 *     IP   IP
	 *	   Port int
	 *	   Zone string // IPv6 scoped addressing zone
	 * }
	 **/
	tcpaddr, err := net.ResolveTCPAddr("tcp4", ":10101") //tcp、tcp4、tcp6
	if err != nil {
		fmt.Println("net.ResolveTCPAddr: ", err.Error())
		os.Exit(1) //只允许0-255

	}

	listen, err := net.ListenTCP("tcp", tcpaddr) //""、"tcp"、"tcp4"、"tcp6"
	/**
	 * 或
	 * listen, err := net.ListenTCP("tcp", &net.TCPAddr{
	 * 	IP:   net.IPv4(0, 0, 0, 0),
	 * 	Port: 10101,
	 * })
	 * 或
	 * listen, err := net.ListenTCP("tcp", "0.0.0.0:10101")
	 * 或
	 * listen, err := net.Listen("tcp", ":10101")
	 */
	if err != nil {
		fmt.Println("net.ListenTCP: ", err.Error())
		os.Exit(1) //Exit参数只允许0-255

	}

	defer listen.Close() //最后关闭listen

	for {
		connect, err := listen.Accept()
		//或connect, err := listen.AcceptTCP()
		if err != nil {
			continue
		}

		go handleClient(connect) //创建一个goroutinue处理
	}

}

func handleClient(connect net.Conn) {
	defer connect.Close()

	clientAddr := connect.RemoteAddr() //获取对端地址
	serverAddr := connect.LocalAddr()
	fmt.Printf("client: %s connect server: %s successed.\n", clientAddr.String(), serverAddr.String())
	connect.SetReadDeadline(time.Now().Add(2 * time.Minute)) //设置两分钟超时
	recv := make([]byte, 1024)

	for {
		len, err := connect.Read(recv)

		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("EOF")
				break
			}
			fmt.Print("socket recivied error: ", err.Error())
		}
		if len < 1 {
			break //说明客户端主动断开连接
		}
		fmt.Printf("%s\n", recv) //recv是[]byte类型，必须格式化为%s打印

		connect.Write([]byte("Hello\n"))

	}

}
