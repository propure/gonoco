package ssh_test

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

func NewConfig(username string, password string, idkeyfilepath string) *ssh.ClientConfig {
    config := ssh.ClientConfig{
        User: "root",
        Auth: []ssh.AuthMethod{
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
            }(idkeyfilepath)), //一般就是"/root/.ssh/.id_rsa"文件
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), //func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
        BannerCallback: func(message string) error {
            return nil
        },
    }
    return &config
}

func main() {
    config := ssh.ClientConfig{
        User: "root",
        Auth: []ssh.AuthMethod{
            ssh.PublicKeys(func(keyPath string) ssh.Signer {
                keyData, err := ioutil.ReadFile(keyPath)
                if err != nil {
                    return nil
                }
                //fmt.Println(keyData)
                signer, err := ssh.ParsePrivateKey(keyData)
                if err != nil {
                    log.Fatal("ParsePrivateKeyWithPassphrase")
                    return nil
                }
                return signer
            }("/root/.ssh/id_rsa")),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), //func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }
        BannerCallback: func(message string) error {
            return nil
        },
    }

    client, err := ssh.Dial("tcp", "127.0.0.1"+":"+strconv.Itoa(22), &config)
    if err != nil {
        fmt.Printf("Dial Error: %s\n", err)
        return
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        fmt.Printf("NewSession Error: %s\n", err)
        return
    }
    defer session.Close()

    in, err := session.StdinPipe()
    if err != nil {
        fmt.Printf("StdinPipe Error: %s\n", err)
        return
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
        return
    }
    //登录到shell，此步必须在StdinPipe之后，否则产生错误：StdinPipe after process started
    if err := session.Shell(); err != nil {
        fmt.Printf("shell Error: %s\n", err)
        return
    }

    //处理banner
    var s string
    b := make([]byte, 1024)
    for i := 0; i < 100; i++ { //100将来替换为banner_timeout
        time.Sleep(time.Second)
        n, err := out.Read(b)
        if err != nil {
            if err == io.EOF {
                //fmt.Println("EOF")
                break
            }
            fmt.Print(err)
        }
        s += string(b[:n])

    }

    s = strings.Trim(s, "\n\r\t ")
    p := strings.Split(s, "\n")
    t := p[len(p)-1]
    if len(t) < 1 {
        fmt.Printf("err2 : %s\n", errors.New("no banner"))
    } else {
        fmt.Printf("banner: `%s`\n\n", t)
    }

}
