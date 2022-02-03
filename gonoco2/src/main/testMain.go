package main

import (
  "fmt"
  "re"
)

func main() {

  a, _ := re.Search("(?P<a>\n\\d*)", "test \n123456")
  fmt.Println(a["a"])

}
