package main

import (
	"log"

	"golang.org/x/crypto/ssh"
)

func main() {
	//配置ssh.ClientConfig
	// config := ssh.ClientConfig{
	//     User:             "root",
	//     Auth: []ssh.AuthMethod{ssh.Password("admin123")},
	// 	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	// }

	//创建连接
	sshClient, err := ssh.Dial("tcp", "localhost:22", &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("admin123")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()

	//执行远程命令
	combo1, err := session.CombinedOutput(`uptime; whoami`)
	if err != nil {
		log.Fatal("远程执行cmd 失败", err)
	}
	log.Println("命令输出:", string(combo1))
}
