package re_test

import (
    "fmt"
    "re"
    "testing"
)

func TestSearch(t *testing.T) {
    a, err := re.Search("([a-z]+) (?P<a>\\d+)", "test 123456")
    if err == nil {
        fmt.Println(a.NamedGroup("a"))
        fmt.Println(a.NamedGroup("b"))
        fmt.Println(a.Group(1))
        fmt.Println(a.Group(9))
    } else {
        fmt.Println(err)
    }
}
