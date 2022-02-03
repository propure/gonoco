package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	DeviceType string
	Username   string
	Password   string
	Host       string
	Port       int
	BasePrompt string
	Client     *ssh.Client
	Session    *ssh.Session
	In         io.WriteCloser
	Out        *bytes.Buffer //不能是非指针
}

var device_type_list = []string{
	"cisco_ios",
	"cisco_nxos",
	"cisco_asa",
	"huawei_vrpv8",
}

var device_diable_pager = map[string]string{
	"cisco_ios":  "terminal length 0",
	"cisco_nxos": "terminal length 0",
	"cisco_asa":  "terminal pager 0",
	"huawei":     "screen-length 0 temporary",
	"h3c":        "screen-length disable",
}

var device_set_config = map[string]string{
	"cisco":  "end",
	"huawei": "return",
	"h3c":    "return",
}

func NewSSHClient(host string, port int, username string, password string, device_type string) (*SSHClient, error) {
	client := new(SSHClient)
	client.DeviceType = device_type
	err := client.SSHConnect(host, port, username, password)
	if err != nil {

	}
	return client, nil
}

func (this *SSHClient) SSHConnect(host string, port int, username string, password string) error {
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
	/**
		 * Dial Error: ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain
	     * ssh: handshake failed: ssh: unable to authenticate, attempted methods [none password], no supported methods remain
	*/
	client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), &config)
	if err != nil {
		fmt.Printf("Dial Error: %s\n", err)
		return err
	}
	//defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		fmt.Printf("NewSession Error: %s\n", err)
		return err
	}
	//defer session.Close()

	in, err := session.StdinPipe()
	if err != nil {
		fmt.Printf("StdinPipe Error: %s\n", err)
		return err
	}
	//defer in.Close()

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

	//处理banner
	var s string
	b := make([]byte, 1024)
	for i := 0; i < 100; i++ { //100将来替换为banner_timeout
		time.Sleep(time.Second)
		n, err := out.Read(b)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF")
				break
			}
			fmt.Print(err)
		}
		s += string(b[:n])

	}
	// time.Sleep(time.Second * 5)
	// s, err := out.ReadString(0)
	// if err != nil && err != io.EOF {
	// 	fmt.Printf("err1 : %s\n", err)
	// }
	s = strings.Trim(s, " \n\r\t")
	if len(s) < 1 {
		fmt.Printf("err2 : %s\n", errors.New("no banner"))
	} else {
		//fmt.Printf("banner: %s\n\n", s)
	}

	this.Client = client
	this.Session = session
	this.In = in
	this.Out = &out

	var ch = make(chan int, 1)
	go this.getPrompt(ch)

	<-ch

	return nil

}

func RegexESCAP(s string) string {
	r := regexp.QuoteMeta(strconv.Quote(s))
	i := strings.Replace(r, `\\`, `\`, -1)
	i = strings.Replace(i, `"`, ``, -1)
	return i
}

func (this *SSHClient) SendCommand(command string) (string, error) {

	if _, err := io.WriteString(this.In, command+"\n"); err != nil {
		fmt.Printf("WriteString error: %s\n", err)
		return "", err
	}

	// time.Sleep(time.Millisecond * 1000)
	// s, err := this.Out.ReadString(0) //只有这种会清理buf
	// if err != nil && err != io.EOF {
	// 	fmt.Printf("ReadString error: %s\n", err)
	// 	return "", err
	// }

	// r, _, err := this.Out.ReadRune()
	// if err != nil && err != io.EOF {
	// 	fmt.Printf("ReadString error: %s\n", err)
	// 	return "", err
	// }
	// s := string(r)

	s := ""
	b := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		n, err := this.Out.Read(b)
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("sendcommand(EOF)")
				break
			} else {
				return "", err
			}
		}
		if n > 0 {
			s += string(b[:n])
			//this.Out.Reset()
			//this.Out.Truncate(0)
		}
		//fmt.Printf("Out len %d\n", this.Out.Len())
	}

	return s, nil
}

func (this *SSHClient) SendCommandExpect(command string, pattern string) (string, error) {

	if _, err := io.WriteString(this.In, command+"\n"); err != nil {
		//fmt.Printf("WriteString error: %s\n", err)
		return "", err
	}
	// s := ""
	// b := make([]byte, 1024)
	// for i := 0; i < 100; i++ {
	// 	time.Sleep(time.Second)
	// 	n, err := this.Out.Read(b)
	// 	if err != nil {
	// 		if err == io.EOF {
	// 			fmt.Printf("sendcommand(EOF)")
	// 			break
	// 		} else {
	// 			return "", err
	// 		}
	// 	}
	// 	s += string(b[:n])
	// 	m, err := regexp.Match(pattern, b[:n])
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	if m {
	// 		break
	// 	}
	// }
	time.Sleep(time.Millisecond * 1000)
	s, err := this.Out.ReadString(0) //只有这种会清理buf
	if err != nil && err != io.EOF {
		fmt.Printf("ReadString error: %s\n", err)
		return "", err
	}
	m, err := regexp.MatchString(pattern, s)
	if err != nil {
		return "", err
	}
	if m {
		return s, nil
	}

	return "", nil
}

func (this *SSHClient) SendCommandAsync(command string, ch chan string, timeout time.Duration) error {

	if _, err := io.WriteString(this.In, command+"\n"); err != nil {
		log.Fatal(err)
		return err
	}

	time.Sleep(time.Millisecond * 5000)
	s, err := this.Out.ReadString(0)
	if err != nil {
		if err == io.EOF {

		} else {
			return err
		}
	}
	ch <- s

	return nil
}

func (this *SSHClient) SendMultiCommand(commands []string) ([]string, error) {
	return []string{""}, nil
}

func (this *SSHClient) SetConfig(commands []string) ([]string, error) {
	return []string{""}, nil
}

func (this *SSHClient) Close() {
	if this.In != nil {
		this.In.Close()
	}

	if this.Session != nil {
		this.Session.Close()
	}

	if this.Client != nil {
		this.Client.Close()
	}
}

func EnableMode(in io.WriteCloser, out bytes.Buffer) error {
	return nil
}

func (this *SSHClient) bannerWait(banner_timeout time.Duration, ch chan int) {

	// time.Sleep(banner_timeout)
	// s, err := this.Out.ReadString(0)
	// if err != nil && err != io.EOF {
	// 	fmt.Printf("err1 : %s\n", err)
	// 	return
	// }
	// s = strings.Trim(s, " \n\r\t")
	// if len(s) < 1 {
	// 	fmt.Printf("err2 : %s\n", errors.New("no banner"))
	// 	return
	// }

	s := ""
	b := make([]byte, 8)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		n, err := this.Out.Read(b)
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("Banner(EOF)")
				break
			} else {
				return
			}
		}
		if n > 0 {
			s += string(b[:n])
		}
		//fmt.Printf("Out len %d\n", this.Out.Len())
	}

	ch <- 1
}

func (this *SSHClient) getPrompt(ch chan int) {
	if _, err := io.WriteString(this.In, "\n"); err != nil {
		log.Fatal(err)
		return
	}
	// time.Sleep(time.Millisecond * 1000)
	// s, err := this.Out.ReadString(0)
	// if err != nil && err != io.EOF {
	// 	fmt.Printf("err1 : %s\n", err)
	// 	return
	// }
	s := ""
	b := make([]byte, 1024)
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		n, err := this.Out.Read(b)
		if err != nil {
			if err == io.EOF {
				//fmt.Printf("sendcommand(EOF)")
				break
			} else {
				return
			}
		}
		if n > 0 {
			s += string(b[:n])
		}
		//fmt.Printf("Out len %d\n", this.Out.Len())
	}
	s = strings.Trim(s, " \n\r\t")
	if len(s) < 1 {
		fmt.Printf("err2 : %s\n", errors.New("no Prompt"))
		return
	}

	this.BasePrompt = s
	ch <- 1
}

func DisablePager(in io.WriteCloser, out bytes.Buffer) error {
	return nil
}
