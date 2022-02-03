//递归方式解析json串转换的json

package main

import (
    "encoding/json"
    "fmt"
    "strings"
)

// 如果碰到的json是：[{...}]
func jsonArrayParse(vv []interface{}) {
    for i, u := range vv {
        switch vv1 := u.(type) {
        case string:
            fmt.Println(i, "[string_] :", u)
        case float64:
            fmt.Println(i, "[float64_]:", u)
        case bool:
            fmt.Println(i, "[bool_]:", u)
        case nil:
            fmt.Println(i, "[nil_]:", u)
        case []interface{}:
            fmt.Println(i, "[array_] :", u)
            jsonArrayParse(vv1)
        case interface{}:
            fmt.Println(i, "[interface_]:", u)
            m1 := u.(map[string]interface{})
            jsonObjectParse(m1)
        default:
            fmt.Println("  ", i, "[type?_]", u, ", ", vv1)
        }
    }
}

// 如果碰到的json是：{...}
func jsonObjectParse(f interface{}) {
    m := f.(map[string]interface{})
    for k, v := range m {
        switch vv := v.(type) {
        case string:
            fmt.Println(k, "[string] :", vv)
        case float64:
            fmt.Println(k, "[float64]:", vv)
        case bool:
            fmt.Println(k, "[bool]:", vv)
        case nil:
            fmt.Println(k, "[nil]:", vv)
        case []interface{}:
            fmt.Println(k, "[array]:")
            jsonArrayParse(vv)
        case interface{}:
            fmt.Println(k, "[interface]:", vv)
            m1 := v.(map[string]interface{})
            jsonObjectParse(m1)
        default:
            fmt.Println(k, "[type?]", vv)
        }
    }
}
func main() {
    jsonStr := []byte(`{"Name":"aree", "Age":18,"From": [ "SZ", "GD" ],"data":[{"a":"aa","b":null},{"c":[]},{"list":["dd",1,"650827..."]}]}`)

    if strings.Index(string(jsonStr[:]), "[") == 0 { //如果第一个字符是"["，就用jsonArrayParse解析，否则用jsonObjectParse解析
        var f []interface{}
        fmt.Println("第一种")
        err := json.Unmarshal(jsonStr, &f)
        if err != nil {
            fmt.Println(err)
        }
        jsonArrayParse(f)
    } else {
        var f interface{}
        fmt.Println("第二种")
        err := json.Unmarshal(jsonStr, &f)
        if err != nil {
            fmt.Println(err)
        }
        jsonObjectParse(f)
    }

}
