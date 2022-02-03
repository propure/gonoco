package main

import (
    "fmt"
    "net/http"
)

//HTTPHandlerHello  无参数
func HTTPHandlerHello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "hello world, \n URL %s!", r.URL.Path[1:])
}

func main() {
    http.HandleFunc("/test/usr/ddddd", HTTPHandlerHello)
    http.ListenAndServe(":8080", nil)
}
