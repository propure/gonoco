package main

import (
    "fmt"
    "reflect"

    "github.com/jakecoffman/cron"
    //"prolee.com/noco/socket"
)

func main() {
    //socket.Tcpserver()
    fmt.Println("Hello, World!")
    c := cron.New()
    fmt.Println(reflect.TypeOf(c))
    c.AddFunc("0 5 * * * *", func() { fmt.Println("Every 5 minutes") }, "Often")
    c.AddFunc("@hourly", func() { fmt.Println("Every hour") }, "Frequent")
    c.AddFunc("@every 1h30m", func() { fmt.Println("Every hour thirty") }, "Less Frequent")
    c.Start()

    var a interface{} = 10
    b, ok := a.(int)
    if ok {
        fmt.Println("a(", reflect.TypeOf(a), ")", a)
        fmt.Println("b(", reflect.TypeOf(b), ")", b)
    }
}
