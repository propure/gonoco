package ssh

import (
    "bytes"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "strconv"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
)

type SSHClient struct {
    Client  *ssh.Client
    Session *ssh.Session
}

func (this *SSHClient) SSHConnect(host string, port int, username string, password string, idkeyfilepath string) (*ssh.Session, error) {

    config := ssh.ClientConfig{
        Config: ssh.Config{
            /**
             * //golang.org/x/crypto/ssh/common.go
             * var supportedCiphers = []string{
             *     "aes128-ctr", "aes192-ctr", "aes256-ctr",
             *     "aes128-gcm@openssh.com",
             *     chacha20Poly1305ID,
             *     "arcfour256", "arcfour128", "arcfour",
             *     aes128cbcID,
             *     tripledescbcID,
             * }
             *
             * const chacha20Poly1305ID = "chacha20-poly1305@openssh.com"
             * const (
             *     gcmCipherID    = "aes128-gcm@openssh.com"
             *     aes128cbcID    = "aes128-cbc"
             *     tripledescbcID = "3des-cbc"
             * )
             **/
            Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "arcfour128", "arcfour256", "arcfour",
                "3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc"},

            /**
             * //golang.org/x/crypto/ssh/kex.go
             * const (
             *     kexAlgoDH1SHA1          = "diffie-hellman-group1-sha1"
             *     kexAlgoDH14SHA1         = "diffie-hellman-group14-sha1"
             *     kexAlgoECDH256          = "ecdh-sha2-nistp256"
             *     kexAlgoECDH384          = "ecdh-sha2-nistp384"
             *     kexAlgoECDH521          = "ecdh-sha2-nistp521"
             *     kexAlgoCurve25519SHA256 = "curve25519-sha256@libssh.org"
             *
             *     // For the following kex only the client half contains a production
             *     // ready implementation. The server half only consists of a minimal
             *     // implementation to satisfy the automated tests.
             *     kexAlgoDHGEXSHA1        = "diffie-hellman-group-exchange-sha1"
             *     kexAlgoDHGEXSHA256      = "diffie-hellman-group-exchange-sha256"
             * )
             */
            KeyExchanges: []string{"diffie-hellman-group1-sha1", "diffie-hellman-group14-sha1", "ecdh-sha2-nistp256",
                "ecdh-sha2-nistp384", "ecdh-sha2-nistp521", "curve25519-sha256@libssh.org", "diffie-hellman-group-exchange-sha1",
                "diffie-hellman-group-exchange-sha256"},
        },
        User: username,
        Auth: []ssh.AuthMethod{
            ssh.Password(password), //普通用户密码登录，ssh.Password()内置的函数，可以用PasswordCallback自定义
            ssh.KeyboardInteractive( // f5一般是交互式键盘登录，回调函数模拟这个过程
                func(user, instruction string, questions []string, echos []bool) ([]string, error) {
                    answers := make([]string, len(questions))
                    for i, _ := range answers {
                        answers[i] = password
                    }
                    return answers, nil
                },
            ),
            ssh.PublicKeysCallback( //通过id_rsa文件登录，用户名root
                func() ([]ssh.Signer, error) {
                    if idkeyfilepath == "" {
                        return nil, fmt.Errorf("argument 'idkeyfile' is empty string")
                    }

                    content, err := ioutil.ReadFile(idkeyfilepath) //客户端保存私钥，服务端保存公钥,私钥就是id_rsa，公钥就是id_rsa.pub
                    if err != nil {
                        err = fmt.Errorf("unable to read id key file: %v", err)
                        return nil, err
                    }

                    signer, err := ssh.ParsePrivateKey(content)
                    if err != nil {
                        err = fmt.Errorf("unable to parse private key: %v", err)
                        return nil, err
                    }

                    return []ssh.Signer{signer}, nil //返回的是个ssh.Signer切片
                },
            ),
            // ssh.PublicKeys(func(keyPath string) ssh.Signer {
            //     pemCodePrivateKey, err := ioutil.ReadFile(keyPath)
            //     if err != nil {
            //         return nil
            //     }
            //     //fmt.Println(keyData)
            //     signer, err := ssh.ParsePrivateKey(pemCodePrivateKey)
            //     if err != nil {
            //         log.Fatal("ParsePrivateKeyWithPassphrase")
            //         return nil
            //     }
            //     return signer
            // }(idkeyfilepath)),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        //HostKeyCallback: hostKeyCallBackFunc(h.Host),
        BannerCallback: func(message string) error {
            return nil
        },
        ClientVersion: "",
        HostKeyAlgorithms: []string{"ssh-rsa", "rsa-sha2-256", "rsa-sha2-512", "ssh-ed25519", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384",
            "ecdsa-sha2-nistp521"},
        Timeout: time.Second * 5,
    }

    client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), &config) //第三个参数是一个指针，返回的client也是指针
    if err != nil {
        fmt.Printf("Dial Error: %s\n", err)
        return nil, err
    }
    defer client.Close()

    session, err := client.NewSession() //NewSession返回的是指针
    if err != nil {
        fmt.Printf("NewSession Error: %s\n", err)
        return nil, err
    }
    defer session.Close()

    return session, nil
}

func (this SSHClient) sessionTermimal(session *ssh.Session) (*ssh.Session, error) {
    in, err := session.StdinPipe() //in io.WriteCloser: 用来接收输入
    if err != nil {
        fmt.Printf("StdinPipe Error: %s\n", err)
        return nil, err
    }
    defer in.Close()

    var out bytes.Buffer //Buffer类型变量out用于存储输出
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
        return nil, err
    }
    //登录到shell，此步必须在StdinPipe之后，否则产生错误：StdinPipe after process started
    if err := session.Shell(); err != nil {
        fmt.Printf("shell Error: %s\n", err)
        return nil, err
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
    s = strings.Trim(s, " \n\r\t")
    if len(s) < 1 {
        fmt.Printf("err2 : %s\n", errors.New("no banner"))
    } else {
        fmt.Printf("banner: %s\n\n", s)
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
