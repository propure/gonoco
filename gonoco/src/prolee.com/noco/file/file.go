package file

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
)

func GetFileStat(path string) {

    fileStat, err := os.Stat(path)

    if err != nil && os.IsNotExist(err) { //判断err是不是"文件不存在"
        fmt.Println("file is not Exist")
    }

    //获取文件名
    fmt.Println(fileStat.Name())

    //判断是否是目录，返回bool类型，等价于Mode().IsDir()
    fmt.Println(fileStat.IsDir())

    //获取文件修改时间
    fmt.Println(fileStat.ModTime())

    //文件的模式位
    fmt.Println(fileStat.Mode())

    //获取文件大小
    fmt.Println(fileStat.Size())

    //底层数据来源（可以返回nil
    fmt.Println(fileStat.Sys())
}

func IsExist(path string) (bool, error) {
    _, err := os.Stat(path)

    if err != nil || os.IsNotExist(err) { //判断err是不是"文件不存在"
        return false, fmt.Errorf("file %s is not exist", path)
    }

    return true, nil
}

func Dirname(s string) string {
    i := strings.LastIndex(s, "/")
    if i > 0 {
        return s[:i]
    } else {
        return ""
    }
}

func Basename(s string) string {
    i := strings.LastIndex(s, "/")
    if i > 0 {
        return s[i+1:]
    } else {
        return s
    }
}

func GetFileContent(path string) (string, error) {

    stat, err := os.Stat(path)
    if err != nil && os.IsNotExist(err) {
        return "", fmt.Errorf("file %s is not exist", path)
    }
    size := stat.Size()
    buf := make([]byte, size+1)

    //os.OpenFile(path, os.O_RDONLY, 0)
    file, err := os.Open(path)
    if err != nil {
        return "", fmt.Errorf("file %s open error", path)
    }
    defer file.Close()

    _, err = file.Read(buf)
    if err != nil {
        return "", fmt.Errorf("file %s read error", path)
    }
    buf[size] = 0

    return string(buf), nil
}

func GetFileContent2(path string) (string, error) {

    if _, err := os.Stat(path); os.IsNotExist(err) {
        fmt.Println("file no exist")
    }

    buf, err := ioutil.ReadFile(path)
    if err != nil {
        fmt.Println("ioutil.ReadFile: ", err)
    }

    return string(buf), nil
}
