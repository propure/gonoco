package main

import (
    "fmt"
    "log"
    "os"
    "strconv"
    "time"

    "golang.org/x/crypto/ssh"
)

var (
    host     string = "localhost"
    port     int    = 10022
    username string = "root"
    password string = "admin123"
)

type ptyRequestMsg struct {
    Term     string
    Columns  uint32 //列字符数
    Rows     uint32 //行数
    Width    uint32 //宽
    Height   uint32
    Modelist string
}

func main() {

    config := &ssh.ClientConfig{
        Config: ssh.Config{
            Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", // 用于linux shell
                //  "arcfour128", "arcfour256", "arcfour",
                "3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc"}, // 用于交换机
            KeyExchanges: []string{"diffie-hellman-group1-sha1"}, //const kexAlgoDH1SHA1 = "diffie-hellman-group1-sha1" ssh/kex.go
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

    //config.Config.Ciphers = append(config.Config.Ciphers, "3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc")
    fmt.Println("Connecting to ", host+":"+strconv.Itoa(port))
    client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), config)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    defer client.Close()

    session, _ := client.NewSession()
    defer session.Close()

    channel, inRequests, err := client.OpenChannel("session", nil)
    if err != nil {
        log.Fatal(err)
        os.Exit(1)
    }
    defer channel.Close()

    go func() {
        for req := range inRequests {
            if req.WantReply {
                req.Reply(false, nil)
            }
        }
    }()
    modes := ssh.TerminalModes{
        ssh.ECHO:          1,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }
    var modeList []byte
    for k, v := range modes {
        kv := struct {
            Key byte
            Val uint32
        }{k, v}
        modeList = append(modeList, ssh.Marshal(&kv)...)
    }
    modeList = append(modeList, 0)
    req := ptyRequestMsg{
        Term:     "xterm",
        Columns:  150,
        Rows:     35,
        Width:    uint32(150 * 8),
        Height:   uint32(35 * 8),
        Modelist: string(modeList),
    }

    ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
    if !ok || err != nil {
        log.Println(err)
        os.Exit(1)
    }

    ok, err = channel.SendRequest("shell", true, nil)
    if !ok || err != nil {
        log.Println(err)
        os.Exit(1)
    }

    fmt.Println("Ok!!!")

}
