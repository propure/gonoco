package main

import (
	"fmt"
	"noco/ssh"
	"os"
)

func main() {
	collect, err := ssh.NewCollector("localhost", 22, "root", "admin123", "cisco_ios")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer collect.Close()
	fmt.Printf("Prompt: %s\n", collect.GetPrompt())
	// s1, err := collect.SendCommand("terminal length 0")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("s1 = `%s`\n\n", s1)

	// s2, err := collect.SendCommand("show running-config")
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Printf("s2 = `%s`", s2)

	s3, err := collect.SendMultiCommand([]string{"uptime", "cat /etc/passwd"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for c, o := range s3 {
		fmt.Printf("%s : `%s`\n", c, o)
	}

}
