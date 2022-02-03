package main

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	DeviceType string
	Host       string
	Port       int
	Username   string
	Password   string
	Client     *ssh.Client
	Session    *ssh.Session
	In         io.WriteCloser
	Out        *bytes.Buffer //不能是非指针
}

func sshClientConfig(username string, password string) ssh.ClientConfig {
	config := ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "arcfour128", "arcfour256", "arcfour",
				"3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc"},
			KeyExchanges: []string{"diffie-hellman-group1-sha1"},
		},
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(
				func(user, instruction string, questions []string, echos []bool) ([]string, error) {
					answers := make([]string, len(questions))
					for i, _ := range answers {
						answers[i] = password
					}
					return answers, nil
				},
			),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		BannerCallback: func(message string) error {
			return nil
		},
		ClientVersion:     "",
		HostKeyAlgorithms: []string{"ssh-rsa"},
		Timeout:           time.Second * 5,
	}

	return config

}

func (this *SSHClient) sshConnect(host string, port int, username string, password string) {

	config := sshClientConfig(username, password)
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), &config)
	if err != nil {
		fmt.Printf("Dial Error: %s\n", err)
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("NewSession Error: %s\n", err)
		return err
	}
	defer session.Close()

	in, err := session.StdinPipe()
	if err != nil {
		fmt.Printf("StdinPipe Error: %s\n", err)
		return err
	}
	defer in.Close()

	var out bytes.Buffer
	session.Stdout = &out

	//设置terminalmodes的方式
	modes := ssh.TerminalModes{
		ssh.ECHO:          0, // 0：不回显
		ssh.TTY_OP_ISPEED: 28800,
		ssh.TTY_OP_OSPEED: 28800,
	}
	//建立伪终端
	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Printf("创建requestpty出错: %s\n", err)
		return err
	}
	//登录到shell，此步必须在StdinPipe之后，否则产生错误：StdinPipe after process started
	if err := session.Shell(); err != nil {
		fmt.Printf("shell Error: %s\n", err)
		return err
	}

	this.Host = host
	this.Port = port
	this.Username = username
	this.Password = password
	this.Client = client
	this.Session = session
	this.In = in
	this.Out = &out

}
