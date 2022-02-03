package main

import (
    "reflect"
    "runtime"
    "strings"
)

type BaseVar interface {
    //对外暴露的接口
    Test()
}

//方便传递
type VarDecriptor func(params ...interface{}) BaseVar

var gFormatCode = map[string]int{
    "VarBinary": 1,
}

func FormatCode(f VarDecriptor) int {
    //取方法名
    fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
    names := strings.Split(fname, ".")
    funcName := names[1]
    //从全局变量中取值
    return gFormatCode[funcName]
}

func main() {

}
