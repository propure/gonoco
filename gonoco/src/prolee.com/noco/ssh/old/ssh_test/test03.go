package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	host     = "10.1.6.199"
	port     = 22
	username = "cisco"
	password = "cisco"
)

type Device struct {
	Config  *ssh.ClientConfig
	Client  *ssh.Client
	Session *ssh.Session
	Stdin   io.WriteCloser
	Stdout  io.Reader
	Stderr  io.Reader
}

func (d *Device) Connect() error {
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), d.Config)
	if err != nil {
		return err
	}
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	sshIn, err := session.StdinPipe()
	if err != nil {
		return err
	}
	sshOut, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	sshErr, err := session.StderrPipe()
	if err != nil {
		return err
	}
	d.Client = client
	d.Session = session
	d.Stdin = sshIn
	d.Stdout = sshOut
	d.Stderr = sshErr
	return nil
}

func (d *Device) SendCommand(cmd string) error {
	if _, err := io.WriteString(d.Stdin, cmd+"\r\n"); err != nil {
		return err
	}
	return nil
}

func (d *Device) SendConfigSet(cmds []string) error {
	for _, cmd := range cmds {
		if _, err := io.WriteString(d.Stdin, cmd+"\r\n"); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (d *Device) PrintOutput() {
	r := bufio.NewReader(d.Stdout)
	for {
		text, err := r.ReadString('\n')
		fmt.Printf("%s", text)
		if err == io.EOF {
			break
		}
	}
}

func (d *Device) PrintErr() {
	r := bufio.NewReader(d.Stderr)
	for {
		text, err := r.ReadString('\n')
		fmt.Printf("%s", text)
		if err == io.EOF {
			break
		}
	}
}

func main() {

	config := &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers:      []string{"3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc"},
			KeyExchanges: []string{"diffie-hellman-group1-sha1"},
		},
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				// Just send the password back for all questions
				answers := make([]string, len(questions))
				for i, _ := range answers {
					answers[i] = password // replace this
				} //这种模式cisco_ios不支持

				return answers, nil
			}),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 5,
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

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer stdin.Close()

	var buf bytes.Buffer
	var e bytes.Buffer
	session.Stdout = &buf
	session.Stderr = &e

	//用以下方式读stdout导致阻塞
	// stdout, err := session.StdoutPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)
	// }
	// if _, err := io.Copy(&buf, stdout); err != nil {
	// 	log.Fatalf("reading filed: %s", err)
	// }

	// stderr, err := session.StderrPipe()
	// if err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)
	// }
	// if _, err := io.Copy(&e, stderr); err != nil {
	// 	log.Fatalf("reading filed: %s", err)
	// }

	//设置terminalmodes的方式
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 28800,
		ssh.TTY_OP_OSPEED: 28800,
	}
	//建立伪终端
	if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal("创建requestpty出错", err)
		os.Exit(1)
	}
	if err := session.Shell(); err != nil {
		log.Fatal("shell Error: ", err)
		os.Exit(1)
	}

	// if err := session.Run("terminal length 0"); err != nil {
	// 	log.Fatal(err)
	// 	os.Exit(1)

	// }

	// if err := session.Run("show version"); err != nil { //和另session already started
	// 	log.Fatal(err)
	// 	os.Exit(1)

	// }

	// for _, cmd := range commands {
	// 	if _, err := io.WriteString(stdin, cmd+"\r\n"); err != nil {
	// 		log.Fatal(err)
	// 		os.Exit(1)
	// 	}
	// 	time.Sleep(time.Second)
	// }
	//fmt.Println("--------------------------------------------------------------------")
	// if _, err := io.WriteString(stdin, "tmsh\r\n"); err != nil {
	// 	log.Fatal(err)
	// }

	//time.Sleep(time.Millisecond * 5000) //100ms以上才能显示完list ltm virtual就有问题了。
	//fmt.Println(buf.String())
	//fmt.Println("--------------------------------------------------------------------")
	// if _, err := io.WriteString(stdin, "modify cli preference pager disabled\r\n"); err != nil {
	// 	log.Fatal(err)
	// }
	if _, err := io.WriteString(stdin, "terminal length 0\n"); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Millisecond * 2000) //100ms以上才能显示完list ltm virtual就有问题了。
	//Stdnull, _ := os.OpenFile("/dev/null", os.O_WRONLY, os.ModePerm)
	//fmt.Fprintf(Stdnull, buf.String()) //无法清空
	//buf.Reset() //无法清空
	//buf.Truncate(0)//无法清空
	// b, _ := rsyncStd(bufio.NewReader(stdout))
	// fmt.Println(b)
	buf.ReadString(0)

	commands := []string{"show version", "show running-config"}
	//commands := []string{"list ltm virtual", "list ltm pool"}

	rc := make(chan string)
	go func() {
		for _, cmd := range commands {
			if _, err := io.WriteString(stdin, cmd+"\n"); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond * 800) //命令输出越多，这个值也必须越大
			// b, _ := rsyncStd(bufio.NewReader(stdout))
			// rc <- string(b)
			rc <- buf.String()
			buf.Reset()
		}
	}()

	// fmt.Print("%s", <-rc) //输出不完整

	result := ""
	tmp := ""
	for {
		select {
		case tmp = <-rc:
			//fmt.Println(result)
			result += tmp
			//fmt.Println("-----------------------------------------------------------------------------------------")
		case <-time.After(time.Millisecond * 2000):
			//break
			goto exit
			// default:
			// 	runtime.Gosched()
		}
	}

exit:
	fmt.Println(result)

	// session.Wait()

	// r := bufio.NewReader(stdout)
	// for {
	// 	text, err := r.ReadString('\n')
	// 	fmt.Printf("%s", text)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// }

	// r = bufio.NewReader(stderr)
	// for {
	// 	text, err := r.ReadString('\n')
	// 	fmt.Printf("%s", text)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// }
}
