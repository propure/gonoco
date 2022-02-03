package main

import (
    "fmt"
    "reflect"
)

type Dark interface {
    Run(string)
    Speak(string) string
}

type bird struct {
}

type chick struct {
}

type myInt int

// func (c *chick) Run(msg string) {
//     fmt.Print(msg)
// }

// func (c *chick) Speak() string { //参数不存在，所以没有实现Speak()
//     return ""
// }

func (b *bird) Run(msg string) {
    fmt.Print(msg)
}

func (b *bird) Speak(string) string {
    return ""
}

func (b *myInt) Run(msg string) {
    fmt.Print(msg)
}

func (b *myInt) Speak(string) string {
    return ""
}

var (
    //判断结构体bird是否实现了Dark接口
    _ Dark = (*bird)(nil) // _ Interface = (*Type)(nil)是判断结构体Type是否实现了接口Interface，把nil转化成*bird类型 赋值后即丢弃。
    // _ Dark = (*chick)(nil) // 没有实现所有方法，所有编译报错:
    // cannot use (*chick)(nil) (value of type *chick) as Dark value in variable declaration:
    //  wrong type for method Speak (have func() string, want func(string) string)
    _ Dark = (*myInt)(nil)
)

func main() {
    var bird Dark = (*bird)(nil)
    // var chick Dark = (*chick)(nil) //这种是间接判断, 没有实现所有方法，所有编译报错:
    var myInt Dark = (*myInt)(nil)
    fmt.Println(reflect.TypeOf(bird))
    // fmt.Println(reflect.TypeOf(chick)*)
    fmt.Println(reflect.TypeOf(myInt))
}
