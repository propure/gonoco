package main

import (
    "fmt"
)

type People interface {
    Speak(string) string
}

type Student struct{}

func (s *Student) Speak(msg string) string {
    return "hi, " + msg
}

func main() {
    var peo People = &Student{} //注意这里,Student{}不是指针，所以不是People指针类型，所以错误
    fmt.Println(peo.Speak("Tom"))
}
