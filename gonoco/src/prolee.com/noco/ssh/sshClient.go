package ssh

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "strconv"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
)

var PrivatekeyPath = "/root/.ssh/id_rsa" //这个应该是针对特定服务器的PrivateKey，例如是id_rsa私钥文件，不是公钥文件

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

func NewSSHClient(host string, port int, username string, password string, device_type string) (*SSHClient, error) {
    client := new(SSHClient)
    client.DeviceType = device_type
    err := client.SSHConnect(host, port, username, password)
    if err != nil {
        log.Fatal("SSHConnect")
        return nil, err
    }
    return client, nil
}

func (this *SSHClient) SSHConnect(host string, port int, username string, password string) error {
    //连接ssh server的配置，包括加密算法、密钥交换算法、用户名、登录方式（Password、KeyboardInteractive）及密码、连接超时
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
            ssh.PublicKeys(func(privateKeyPath string, passpharase []byte) ssh.Signer {
                keyData, err := ioutil.ReadFile(keyPath)
                if err != nil {
                    return nil
                }
                //sugner, err := ssh.ParsePrivateKey(privateKey)
                signer, err := ssh.ParsePrivateKeyWithPassphrase(keyData, passpharase) //ParsePrivateKeyxxx有两个返回值，但ssh.PublicKeys只接受一个参数，所有要做成内联函数
                if err != nil {
                    log.Fatal("ParsePrivateKeyWithPassphrase")
                    return nil
                }
                return signer
            }(PrivatekeyPath, []byte(password))),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), //func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
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
    for i := 0; i < 100; i++ { //100将来替换为banner_timeout，此处类试ssh2在非阻塞时读内容
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

    //去掉banner首尾所有的回车、换行、制表符和空格
    s = strings.Trim(s, "\n\r\t ") //这里的字符串是一个个来用，相当于[\r\n\t ]
    t := strings.Split(s, "\n")    //而这里的字符串是一起使用，如果多个就是"\r\n"
    s = t[len(t)-1]
    if len(s) < 1 {
        fmt.Printf("err2 : %s\n", errors.New("no banner"))
    } else {
        //显示prompt
        fmt.Printf("Prompt: `%s`\n\n", s)
    }

    this.Client = client
    this.Session = session
    this.In = in
    this.Out = &out

    return nil

}
