package main

import (
    "fmt"
    "unsafe"
)

type stA struct {
    A []int
    B string
}

type stB struct {
    A stA
    B string
}

func main() {

    var s stB
    s.A.B = "Test"

    var t *stB
    t = new(stB)
    t.A.B = "Hello"

    fmt.Println(unsafe.Sizeof(s))
    fmt.Println(s.A.B)
    fmt.Println(unsafe.Sizeof(t))
    fmt.Println(t.A.B)

}
