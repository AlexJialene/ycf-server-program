package main

import (
	/*"fmt"
	"github.com/robfig/cron"*/
	"fmt"
	"sync"
)

//计数器
type Counter struct {
	systemId int64
	path     string
	count    int
	name     string
}

type MonitorHandle struct {
	mu         sync.RWMutex
	monitorMap sync.Map
}

func (m *MonitorHandle) putMonitor() {

}

func main() {
	te := &Te{}
	fmt.Println(te.name)
	te.t1()
	fmt.Println(te.name)

	fmt.Println(te.t3())
	te.t2()
	fmt.Println(te.name)
}

type Te struct {
	name string
}

func (t Te) t1() {
	t.name = "123"
}

func (t Te) t3() string {
	return t.name

}

func (t *Te) t2() {
	t.name = "456"
}

func collect() {
}

//
