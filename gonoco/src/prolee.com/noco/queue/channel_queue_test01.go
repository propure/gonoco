package main

import (
    "fmt"
    "time"
)

type DataContainer struct {
    Queue chan interface{}
}

func NewDataContainer(max_queue_len int) (dc *DataContainer) {
    dc = &DataContainer{}
    dc.Queue = make(chan interface{}, max_queue_len)
    return dc
}

//非阻塞push
func (dc *DataContainer) Push(data interface{}, waittime time.Duration) bool {
    click := time.After(waittime)
    select {
    case dc.Queue <- data:
        return true
    case <-click:
        return false
    }
}

//非阻塞pop
func (dc *DataContainer) Pop(waittime time.Duration) (data interface{}) {
    click := time.After(waittime)
    select {
    case data = <-dc.Queue:
        return data
    case <-click:
        return nil
    }
}

//test
var MAX_WAIT_TIME = 10 * time.Millisecond

func main() {
    type dataItem struct {
        name string
        age  int
    }

    datacotainer := NewDataContainer(2)
    //add
    fmt.Printf("res=%v\n", datacotainer.Push(&dataItem{"zhangsan", 25}, MAX_WAIT_TIME))
    fmt.Printf("res=%v\n", datacotainer.Push(&dataItem{"lisi", 30}, MAX_WAIT_TIME))
    fmt.Printf("res=%v\n", datacotainer.Push(&dataItem{"wangwu", 28}, MAX_WAIT_TIME))

    //get
    var item interface{}
    item = datacotainer.Pop(MAX_WAIT_TIME)
    if item != nil {
        if tmp, ok := item.(*dataItem); ok { //interface转为具体类型，断言也是这样用，_, ok := element.(*type)
            fmt.Printf("item name:%v, age:%v\n", tmp.name, tmp.age)
        }
    }
}
