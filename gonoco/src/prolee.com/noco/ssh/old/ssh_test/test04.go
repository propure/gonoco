package main

import (
	"bytes"
	"fmt"
	"io"
  "io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	host     string = "10.1.6.199"
	port     int    = 22
	username string = "cisco"
	password string = "cisco"
)

func main() {

	config := &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers:      []string{"3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc"},
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
      ssh.PublicKeysCallback(
        func()([]ssh.Signer, error){
            key, err := ioutil.ReadFile(prikeyFile) //客户端保存私钥，服务端保存公钥,私钥就是id_rsa，公钥就是id_rsa.pub
            if err != nil {
                err = fmt.Errorf("unable to read private key: %v", err)
                return nil, err
            }
            
            signer, err := ssh.ParsePrivateKey(key)
            if err != nil {
                err = fmt.Errorf("unable to parse private key: %v", err)
                return nil, err
            }

            return []ssh.Signer{signer}, nil
        },
      ),
		},
		//HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    HostKeyCallback: ssh.FixedHostKey()
		BannerCallback: func(message string) error {
			return nil
		},
		ClientVersion:     "",
		HostKeyAlgorithms: []string{"ssh-rsa"},
		Timeout:           time.Second * 5,
	}
	fmt.Println("Connecting to ", host+":"+strconv.Itoa(port))
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer session.Close()

	//设置terminalmodes的方式
	modes := ssh.TerminalModes{
		ssh.ECHO:          0, // 0：不回显
		ssh.TTY_OP_ISPEED: 28800,
		ssh.TTY_OP_OSPEED: 28800,
	}
	//建立伪终端
	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal("创建requestpty出错", err)
		os.Exit(1)
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer stdin.Close()

	// 此步必须在StdinPipe之后，否则产生错误：StdinPipe after process started
	if err := session.Shell(); err != nil {
		log.Fatal("shell Error: ", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	var e bytes.Buffer
	session.Stdout = &buf
	session.Stderr = &e

	time.Sleep(time.Millisecond * 500)
	fmt.Printf("buf len: %d\n", buf.Len())
	base_prompt, _ := buf.ReadString(0)
	fmt.Printf("base_prompt: %s\n", strings.Trim(base_prompt, "\n\r \t"))

	if _, err := io.WriteString(stdin, "terminal length 0\n"); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 500)
	buf.ReadString(0) //读取buf，直到字符0为止，返回为string, err，err有值，一般是io.EOF

	commands := []string{"show version", "show running-config"}
	rc := make(chan string)
	go func() {
		for _, cmd := range commands {
			if _, err := io.WriteString(stdin, cmd+"\n"); err != nil {
				log.Fatal(err)
			}
			for {
				time.Sleep(time.Second * 1)
				rc <- buf.String()
				_, err := buf.ReadString(0)
				if err != nil {
					if err == io.EOF {
						break
					} else {
						return
					}
				}
			}
		}
	}()

	result := ""
	tmp := ""
	for {
		select {
		case tmp = <-rc:
			result += tmp
		case <-time.After(time.Millisecond * 2000):
			goto exit
		}
	}

exit:
	fmt.Println(result)

}
