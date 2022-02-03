package main

import (
    "fmt"
    "regexp"
)

type Match struct {
    rePtr       *regexp.Regexp
    matchString []string
}
type MyInt int

func (m *Match) NamedGroup(name string) (string, error) {

    i := m.rePtr.SubexpIndex(name)
    if i >= 0 {
        return m.matchString[i], nil
    } else {
        return "", fmt.Errorf("name '%s' not in NameGroup", name)
    }

}

func (m *Match) Group(index int) (string, error) {
    count := len(m.matchString)
    if index < 0 || index >= count {
        return "", fmt.Errorf("index out of matchs range")
    }
    return m.matchString[index], nil

}

// 是否匹配正则表达式, 正确返回true
func Consistent(pattern string, s string) (bool, error) {
    ret, err := regexp.MatchString(pattern, s) //这种方法不用编译正则表达式
    return ret, err
}

func _find(pattern string, s string) (*Match, error) {

    reptr, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }

    matchingString := reptr.FindStringSubmatch(s)
    if len(matchingString) < 1 {
        return nil, fmt.Errorf("no matching string")
    }

    var m = &Match{reptr, matchingString}

    return m, nil

}

func Find(pattern string, s string) (*Match, error) {
    return _find(pattern, s)
}

func find(pattern string, text string) (map[string]string, error) {
    re, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }
    groupNames := re.SubexpNames() //获取所有命名组的名称 [] string
    groupCount := len(groupNames)  //命名组数量
    if groupCount < 1 {
        return nil, nil
    }

    matchs := re.FindStringSubmatch(text)
    if len(matchs) < 1 {
        return nil, nil
    }
    ret := make(map[string]string, groupCount)
    for i := 0; i < groupCount; i++ {
        name := groupNames[i]
        if name == "" { //跳过完整匹配，index=0
            continue
        }
        ret[name] = matchs[re.SubexpIndex(name)]

    }
    return ret, nil
}

func test01() {
    // pattern := `^\s*(?P<hello>Hello)\s*,\s*(?P<world>World).*`
    // re, err := regexp.Compile(pattern)
    // if err != nil {
    // 	fmt.Println("compile error: ", err)
    // 	os.Exit(1)
    // }
    // fmt.Println(re.FindString("Hello,World"))
    // fmt.Println(re.FindStringSubmatch("Hello,World"))
    // fmt.Println(re.FindStringIndex("Hello,World"))
    // fmt.Println(re.FindStringSubmatchIndex("Hello,World"))
    // fmt.Println(re.FindAllString("Hello,World", -1))
    // fmt.Println(re.FindAllStringSubmatch("Hello,World", -1))
    // fmt.Println(re.FindAllStringIndex("Hello,World", -1))
    // fmt.Println(re.FindAllStringSubmatchIndex("Hello,World", -1))

    // matchs := re.FindStringSubmatch("Hello,World")

    // fmt.Println(re.SubexpNames())
    // fmt.Println(re.SubexpIndex("*"))
    // fmt.Println(re.SubexpIndex(""))
    // fmt.Println(re.SubexpIndex("hello"))
    // fmt.Println(re.SubexpIndex("world"))
    // fmt.Println(matchs[re.SubexpIndex("hello")])
    // fmt.Println(matchs[re.SubexpIndex("world")])

    m, err := find(`^(?P<ethernet_interface>Eth\S+)\s+(?P<vlan>(\d+|--|monitor))\s+(?P<type>eth)\s+`+
        `(?P<mode>(trunk|access|routed))\s+(?P<status>(up|down))\s+`+
        `(?P<reason>(none|SFP not inserted|Administratively down|channel admin down|Link not connected|suspended\(no LACP PDUs\)|XCVR not inserted))\s+`+
        `(?P<speed>(\d+G?|auto))(\(D\))\s+`+
        `(?P<port_channel>(\d+|--))`,
        `Eth1/28       461     eth  access down    Administratively down      auto(D) 99`)
    if err != nil {
        fmt.Println("err")
    }

    fmt.Println(m)

}

func main() {

    //var map1 map[string]string //如果仅仅这样，这时并没有分配内存空间
    //map1 := make(map[string]string) //用:=就无需var map1
    // map1 := map[string]string{}
    // map1["hello"] = "world"
    // buf := make([]byte, 4096)
    test01()

}
