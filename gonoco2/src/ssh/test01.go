package ssh

import (
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "strconv"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
)

type sshConnect struct {
    // config    *ssh.ClientConfig
    client  *ssh.Client
    session *ssh.Session
    in      io.WriteCloser
    out     strings.Builder
    // username  string
    // password  string
    // idKeyPath string
}

type netDriver struct {
}

func (this *netDriver) NewDriver(host string, port int, username string, password string, idkeyfilepath string, device_type string) (*netDriver, error) {

    return this, nil
}

func NewSSHConnect(host string, port int, username string, password string, idkeyfilepath string) (*sshConnect, error) {

    config := ssh.ClientConfig{
        Config: ssh.Config{
            Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "arcfour128", "arcfour256", "arcfour",
                "3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc"},

            KeyExchanges: []string{"diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1", "ecdh-sha2-nistp256",
                "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "curve25519-sha256@libssh.org", "diffie-hellman-group-exchange-sha1",
                "diffie-hellman-group-exchange-sha256"},
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
            ssh.PublicKeys(func(keyPath string) ssh.Signer {
                pemCodePrivateKey, err := ioutil.ReadFile(keyPath)
                if err != nil {
                    return nil
                }
                signer, err := ssh.ParsePrivateKey(pemCodePrivateKey)
                if err != nil {
                    log.Fatal("ParsePrivateKeyWithPassphrase")
                    return nil
                }
                return signer
            }(idkeyfilepath)),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        BannerCallback: func(message string) error {
            return nil
        },
        ClientVersion: "",
        HostKeyAlgorithms: []string{"ssh-rsa", "rsa-sha2-256", "rsa-sha2-512", "ssh-ed25519", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384",
            "ecdsa-sha2-nistp521"},
        Timeout: time.Second * 5,
    }

    client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), &config)
    defer client.Close()
    if err != nil {
        log.Fatal("ssh Dial error")
        return nil, err
    }

    session, err := client.NewSession() //NewSession返回的是指针
    defer session.Close()
    if err != nil {
        log.Fatal("ssh NewSession error")
        return nil, err
    }

    in, err := session.StdinPipe() //in io.WriteCloser: 用来接收输入
    defer in.Close()
    if err != nil {
        log.Fatal("ssh session StdinPipe error")
        return nil, err
    }

    var out strings.Builder //bytes.Buffer或string.Builder类型变量out用于存储输出
    session.Stdout = &out   //这种会把输出自动写到Builder里，不用再循环输出

    modes := ssh.TerminalModes{
        ssh.ECHO:          0, // 0：不回显
        ssh.TTY_OP_ISPEED: 28800,
        ssh.TTY_OP_OSPEED: 28800,
    }
    //建立伪终端
    if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
        log.Fatal("ssh session requestpty error.")
        return nil, err
    }
    //登录到shell，此步必须在StdinPipe之后，否则产生错误：StdinPipe after process started
    if err := session.Shell(); err != nil {
        fmt.Printf("shell Error: %s\n", err)
        return nil, err
    }

    sshConnect := &sshConnect{client: client, session: session, in: in, out: out}
    return sshConnect, nil

}

func (c *sshConnect) Banner() {
    time.Sleep(time.Second * 1)
    fmt.Print(c.out.String()) //这里包含了prompt
    c.out.Reset()             //清除out，这样下一次string()时就不会有上次的输出了
}

// 获取prompt，这样每次获取字符串以后，通过strings.ReplaceAll(string, oldsub, newsub)替换为空
func (c *sshConnect) Prompt() string {
    c.in.Write([]byte("\n"))
    time.Sleep(time.Second * 1)
    prompt := c.out.String()
    c.out.Reset()
    return prompt
}

func (c *sshConnect) SendCommand(cmd string) string {
    cmd = strings.Trim(cmd, " \t")
    l := len(cmd)
    if l < 1 {
        return ""
    }
    if strings.Compare(cmd[l-1:], "\n") == -1 {
        cmd += "\n"
    }
    c.in.Write([]byte(cmd))
    time.Sleep(time.Second * 1)
    result := c.out.String()
    c.out.Reset()
    return result

}
