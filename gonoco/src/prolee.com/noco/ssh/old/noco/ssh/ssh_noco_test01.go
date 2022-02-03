package ssh

import (
  "bytes"
  "errors"
  "fmt"
  "io"
  "regexp"
  "strconv"
  "strings"
  "time"

  "golang.org/x/crypto/ssh"
)

type SSHCollector struct {
  device_type string
  username    string
  password    string
  host        string
  port        int
  base_prompt string
  client      *ssh.Client
  session     *ssh.Session
  in          io.WriteCloser
  out         *bytes.Buffer //不能是非指针
}

var device_type_list := []string{
  "cisco_ios",
  "cisco_nxos",
  "cisco_asa",
  "huawei_vrpv8",
  "h3c_comware",
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

func NewCollector(host string, port int, username string, password string, device_type string) (*SSHCollector, error) {
  collect := new(SSHCollector)
  collect.device_type = device_type
  err := collect.SSHConnect(host, port, username, password)
  if err != nil {
    return nil, err
  }
  return collect, nil
}

func (collect *SSHCollector) SSHConnect(host string, port int, username string, password string) error {

  // Dial的config参数
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

  // 创建ssh连接
  client, err := ssh.Dial("tcp", host+":"+strconv.Itoa(port), &config)
  if err != nil {
    fmt.Printf("Dial Error: %s\n", err)
    return err
  }
  //defer client.Close()

  //创建ssh session
  session, err := client.NewSession()
  if err != nil {
    fmt.Printf("NewSession Error: %s\n", err)
    return err
  }
  //defer session.Close()

  // 设置ssh session shell的输入
  in, err := session.StdinPipe() //io.WriterCloser
  if err != nil {
    fmt.Printf("StdinPipe Error: %s\n", err)
    return err
  }
  //defer in.Close()

  // 设置ssh session shell的输出
  //var out bytes.Buffer
  var out io.ReadCloser
  session.Stdout = &out

  //设置terminalmodes的方式
  modes := ssh.TerminalModes{
    ssh.ECHO:          1, // 0：不回显
    ssh.TTY_OP_ISPEED: 14400,
    ssh.TTY_OP_OSPEED: 14400,
  }
  //建立伪终端
  if err = session.RequestPty("xterm", 80, 40, modes); err != nil {
    fmt.Printf("session requestpty Error: %s\n", err)
    return err
  }
  //登录到shell，此步必须在StdinPipe之后，否则产生错误：StdinPipe after process started
  if err := session.Shell(); err != nil {
    fmt.Printf("session shell Error: %s\n", err)
    return err
  }

  collect.client = client
  collect.session = session
  collect.in = in
  collect.out = &out

  collect.bannerWait()
  collect.setPrompt()

  return nil

}

func (collect *SSHCollector) Close() {
  if collect.in != nil {
    collect.in.Close()
  }

  if collect.session != nil {
    collect.session.Close()
  }

  if collect.client != nil {
    collect.client.Close()
  }
}

func (collect *SSHCollector) SendCommand(command string) (string, error) {

  if _, err := io.WriteString(collect.in, command+"\n"); err != nil {
    fmt.Printf("WriteString error: %s\n", err)
    return "", err
  }

  time.Sleep(time.Millisecond * 2000)

  s, err := collect.out.ReadString(0)
  if err != nil && err != io.EOF {
    fmt.Printf("sendcommand Read error: %s\n", err)
    return "", err
  }

  return s, nil
}

func (collect *SSHCollector) SendMultiCommand(commands []string) (map[string]string, error) {

  //output := make(map[string]string, len(commands)) //指定的长度没有实际意义，但可以减少内存重新分配。
  output := map[string]string{}
  for _, cmd := range commands { //第一个值是索引

    if _, err := io.WriteString(collect.in, cmd+"\n"); err != nil {
      fmt.Printf("WriteString error: %s\n", err)
      return nil, err
    }

    time.Sleep(time.Millisecond * 2000)

    s, err := collect.out.ReadString(0)
    if err != nil && err != io.EOF {
      fmt.Printf("sendcommand Read error: %s\n", err)
      return nil, err
    }
    output[cmd] = s
  }

  return output, nil
}

func (collect *SSHCollector) bannerWait() {
  //处理banner
  time.Sleep(time.Millisecond * 1000)

  s, err := collect.out.ReadString(0)
  if err != nil && err != io.EOF {
    fmt.Printf("bannerWait Read error: %s\n", err)
  }

  if len(s) < 1 {
    fmt.Printf("err2 : %s\n", errors.New("no banner"))
  } else {
    fmt.Printf("banner :`%s`\n\n", s)
  }
}

func (collect *SSHCollector) GetPrompt() string {
  return collect.base_prompt

}

func (collect *SSHCollector) setPrompt() {
  if _, err := io.WriteString(collect.in, "\n"); err != nil {
    fmt.Printf("WriteString error: %s\n", err)
    return
  }

  time.Sleep(time.Millisecond * 1000)

  s, err := collect.out.ReadString(0)
  if err != nil && err != io.EOF {
    fmt.Printf("setPrompt Read error: %s\n", err)
  }

  s = strings.Trim(s, " \n\r\t")
  if len(s) < 1 {
    fmt.Printf("err2 : %s\n", errors.New("no Prompt"))
  } else {
    collect.base_prompt = s
  }
}

func RegexESCAP(s string) string {
  r := regexp.QuoteMeta(strconv.Quote(s))
  i := strings.Replace(r, `\\`, `\`, -1)
  i = strings.Replace(i, `"`, ``, -1)
  return i
}
