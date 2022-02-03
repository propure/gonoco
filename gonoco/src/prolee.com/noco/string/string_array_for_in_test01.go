package main

import "fmt"

func main() {
    s := []string{"Hello", "World"}
    // s := make([]string, 10)
    // s = append(s, "Hello")
    // s = append(s, "World")

    for _, t, := range s { //range 切片时，第一个返回的时index
        fmt.Println(t)
    }
}
