package gonoco/ssh_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	//"net"

	"os"

	"golang.org/x/crypto/ssh"
)

var (
	server   = "10.1.6.192:22"
	username = "root"
	password = "default"
)

func authByPublicKey(keyPath string) ssh.AuthMethod {

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
func main() {
	config := &ssh.ClientConfig{
		Config: ssh.Config{
			Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128",
				"3des-cbc", "aes128-cbc", "aes192-cbc", "aes256-cbc"},
			KeyExchanges: []string{"diffie-hellman-group1-sha1"},
		},

		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            "cisco",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //和下个函数一样，直接返回nil，不够安全
		// HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {return nil},
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		// Auth: []ssh.AuthMethod{authByPublicKey("/root/.ssh/id_rsa")},
		// Auth: []ssh.AuthMethod{ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
		// 	// Just send the password back for all questions
		// 	answers := make([]string, len(questions))
		// 	for i, _ := range answers {
		// 		answers[i] = "cisco" // replace this
		// 	} //这种模式cisco_ios不支持

		// 	return answers, nil
		// })},
	}
	// config.Config.Ciphers = append(config.Config.Ciphers, "aes128-cbc", "aes192-cbc", "aes256-cbc", "3des-cbc")
	// config.Config.KeyExchanges = append(config.Config.KeyExchanges, "diffie-hellman-group1-sha1")

	fmt.Println(config)

	sshClient, err := ssh.Dial("tcp", server, config)
	if err != nil {
		log.Fatal("创建ssh client 失败", err)
	}
	defer sshClient.Close()

	//创建ssh-session
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatal("创建ssh session 失败", err)
	}
	defer session.Close()
	// //执行远程命令
	// combo1, err := session.CombinedOutput(`uptime; whoami`)
	// if err != nil {
	//     log.Fatal("远程执行cmd 失败", err)
	// }
	// log.Println("命令输出:", string(combo1))

	// //执行远程命令，第二次就会失败：远程执行cmd 失败ssh: Stdout already set，更换combo2也不行
	// combo2, err := session.CombinedOutput(`cd /; ls -al`)
	// if err != nil {
	//     log.Fatal("远程执行cmd 失败", err)
	// }
	// log.Println("命令输出:", string(combo2))

	//设置terminalmodes的方式
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 打开/关闭回显
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		ssh.VSTATUS:       1,
	}
	//建立伪终端
	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		log.Fatal("创建requestpty出错", err)
		os.Exit(1)
	}
	// //设置session的标准输入是stdin
	stdin, err := session.StdinPipe()
	if err != nil {
		log.Fatal("输入错误", err)
		os.Exit(1)
	}
	// //设置session的标准输出和错误输出分别是os.stdout,os,stderr.就是输出到后台
	// session.Stdout = os.Stdout
	var b bytes.Buffer
	session.Stdout = &b
	session.Stderr = os.Stderr

	// stdin, err := session.StdinPipe()
	// if err != nil {
	//     fmt.Errorf("Unable to setup stdin for session: %v", err)
	//     os.Exit(1)
	// }
	// go io.Copy(stdin, os.Stdin)

	// stdout, err := session.StdoutPipe()
	// if err != nil {
	//     fmt.Errorf("Unable to setup stdout for session: %v", err)
	//     os.Exit(1)
	// }
	// go io.Copy(os.Stdout, stdout)

	// stderr, err := session.StderrPipe()
	// if err != nil {
	//     fmt.Errorf("Unable to setup stderr for session: %v", err)
	//     os.Exit(1)
	// }
	// go io.Copy(os.Stderr, stderr)

	err = session.Shell()
	if err != nil {
		log.Fatal("创建shell出错", err)
		os.Exit(1)
	}
	//commands := []string{"configure terminal", "show history"}
	commands := []string{"tmsh", "list ltm virtutal"}
	//将命令依次执行
	for _, cmd := range commands {
		fmt.Println(cmd)
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal("写入stdin出错", err)
			os.Exit(1)
		}
	}

	//执行等待
	err = session.Wait()
	if err != nil {
		log.Fatal("等待session出错", err)
		os.Exit(1)
	}

	fmt.Println("结束") //如果前面有session.Wait()，则一直阻塞，该句无法执行
}
