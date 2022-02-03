package main

import (
    "fmt"
    "strings"
)

func main() {
    s := "Hello, World"
    fmt.Printf("%s\n", strings.Trim(s, "deH "))
}
