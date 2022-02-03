package http

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "strconv"
)

func HttpGet2(url string) (string, error) {
  client := &http.Client{}
  request, err := http.NewRequest("GET", url, nil)
  if err != nil {
    fmt.Println(err)
  }

  cookie := &http.Cookie{Name: "Tom", Value: strconv.Itoa(123)}
  request.AddCookie(cookie) //向request中添加cookie

  //设置request的header
  request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
  request.Header.Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
  request.Header.Set("Accept-Encoding", "gzip,deflate,sdch")
  request.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
  request.Header.Set("Cache-Control", "max-age=0")
  request.Header.Set("Connection", "keep-alive")

  response, err := client.Do(request)
  if err != nil {
    fmt.Println(err)
    return "", err
  }

  defer response.Body.Close()
  fmt.Println(response.StatusCode)

  var body []byte
  if response.StatusCode == 200 {
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
      fmt.Println(err)
      return "", err
    }
    fmt.Println(string(body))
  }

  return string(body), nil
}

func HttpGet(url string) (string, error) {
  response, err := http.Get(url)
  if err != nil {
    fmt.Println("http.Get error: ", err.Error())
    return "", err
  }

  defer response.Body.Close()

  body, err := ioutil.ReadAll(response.Body)
  if err != nil {
    fmt.Println("ioutil.ReadAll error: ", err.Error())
    return "", err

  }

  return string(body[:]), nil
}
