package re

import (
    "fmt"
    "regexp"
)

type match struct { // 这里居然不用大写首字母
    rePtr  *regexp.Regexp
    matchs []string
}

func (m *match) NamedGroup(name string) string {
    i := m.rePtr.SubexpIndex(name)
    if i >= 0 {
        return m.matchs[i]
    } else {
        return ""
    }

}

func (m *match) Group(index int) string {
    count := len(m.matchs)
    if index < 0 || index >= count {
        return ""
    }
    return m.matchs[index]

}

func IsMatched(pattern string, s string) (bool, error) {
    ret, err := regexp.MatchString(pattern, s) // 这种方法不用编译正则表达式
    return ret, err
}

func Search(pattern string, s string) (*match, error) {

    reptr, err := regexp.Compile(pattern)
    if err != nil {
        return nil, err
    }

    matchs := reptr.FindStringSubmatch(s) // 返回[]string
    if len(matchs) < 1 {
        return nil, fmt.Errorf("no matching string")
    }

    var m = &match{reptr, matchs}

    return m, nil
}
