package main

import (
    "encoding/json"
    "fmt"
)

type Server struct {
    ServerName string
    ServerIP   []string
}

type Serverslice struct {
    Servers []Server
}

func main() {
    var s Serverslice
    str := `{
                "servers": [
                    {
                        "serverName": "Guangzhou_Base",
                        "serverIP": [
                            "127.0.0.1",
                            "192.168.0.1"
                        ]
                    }, {
                        "serverName": "Beijing_Base",
                        "serverIP": [
                            "127.0.0.2",
                            "192.168.0.2"
                        ]
                    }
                ]
            }`

    err := json.Unmarshal([]byte(str), &s) //解析会自动嵌套，无需一级级

    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(s)
    fmt.Println(s.Servers[0].ServerName)
    fmt.Println(s.Servers[0].ServerIP[1])
}
