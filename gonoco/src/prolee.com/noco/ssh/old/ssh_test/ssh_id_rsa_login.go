package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "strconv"
    "strings"
    "time"

    "golang.org/x/crypto/ssh"
    //"time"
)

func main() {
    var (
        // command       = "ls /"
        host          = "localhost"
        port          = 10022
        username      = "root"
        idkeyfilepath = "C:\\users\\prole\\.ssh\\id_rsa" //ssh服务器的/root/.ssh/authorized_keys保存了公钥文件id_rsa.pub内容
    )
    config := ssh.ClientConfig{
        User: username,
        Auth: []ssh.AuthMethod{
            ssh.PublicKeysCallback( //只有一个Auth方法，就是通过id_rsa登录，用户名root
                func() ([]ssh.Signer, error) {
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
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), //必须，否则
    }

    client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), &config)
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

    // output, err := session.Output("ls /") //只有一个命令
    // if err != nil {
    //     //err = fmt.Errorf("run command %s on host %s error %v", command, host, err)
    //     fmt.Printf("run command %s on host %s error %v", "ls /", host, err)
    //     return
    // }
    // fmt.Println(string(output)) //output是[]byte所以打印值是：

    in, err := session.StdinPipe() //in io.WriteCloser: 用来接收输入
    defer in.Close()
    if err != nil {
        log.Fatal("ssh session StdinPipe error")
        return
    }

    /**
     * strings.Builder实现了Write()方法，也就实现了io.Writer接口
     * 而Session.Stdout的类型就是io.Writer
     * 可以把string.builder赋值给Stdout
     **/
    var out strings.Builder //bytes.Buffer或string.Builder类型变量out用于存储输出
    session.Stdout = &out
    // out, err := session.StdoutPipe() //out io.Reader: 输出，实际上时session的ch属性，也就是一个通道
    // if err != nil {
    //     log.Fatal("ssh session StdinPipe error")
    //     return
    // }

    //这个返回的时一个channel
    // out, err := session.StdoutPipe() //in io.WriteCloser: 用来接收输入 //这种方式需要有一个b []byte接收out.Read(b)，这个最大长度一般固定
    // if err != nil {
    //     log.Fatal("ssh session StdinPipe error")
    //     return nil, err
    // }

    modes := ssh.TerminalModes{
        ssh.ECHO:          0, // 0：不回显
        ssh.TTY_OP_ISPEED: 28800,
        ssh.TTY_OP_OSPEED: 28800,
    }
    //建立伪终端
    if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
        log.Fatal("ssh session requestpty error.")
        return
    }
    //登录到shell，此步必须在StdinPipe之后，否则产生错误：StdinPipe after process started
    if err := session.Shell(); err != nil {
        fmt.Printf("shell Error: %s\n", err)
        return
    }

    // var s strings.Builder
    // b := make([]byte, 1024)
    // for i := 0; i < 3; i++ { //100将来替换为banner_timeout
    //     time.Sleep(time.Microsecond * 10)
    //     n, err := out.Read(b)
    //     if err != nil {
    //         if err == io.EOF {
    //             break
    //         }
    //     }
    //     if n == 0 {
    //         break
    //     }
    //     // fmt.Println("xxx")
    //     fmt.Print(b)
    //     // s.Write(b)
    // }
    time.Sleep(time.Second * 2)
    fmt.Println(out.String()) //打印banner
    out.Reset()               //清空out

    in.Write([]byte("\n"))
    time.Sleep(time.Second * 2)
    prompt := out.String() //获取prompt

    in.Write([]byte("ls /\n"))
    time.Sleep(time.Second * 2)
    buf := out.String()
    buf = strings.ReplaceAll(buf, prompt, "")
    fmt.Println(buf)
    out.Reset()

    in.Write([]byte("ls\n"))
    fmt.Printf("%s", out.String())

    in.Write([]byte("pwd\n"))
    fmt.Printf("%s", out.String())
    return
}
