package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Monitor struct {
	Time int64             `json:"time"`
	List []*MonitorControl `json:"list"`
}

type MonitorControl struct {
	//monitor id
	Id string `json:"id"`
	//count
	Count int `json:"count"`
}

var monitorMap sync.Map

//statistical map read-write lock
var m = new(sync.RWMutex)

func Add(id string, count int) {
	m.RLock()
	defer m.RUnlock()
	value, ok := monitorMap.Load(id)
	if ok {
		control := value.(*MonitorControl)
		control.Count += count

	} else {
		monitorMap.Store(id, &MonitorControl{id, count})
	}
}

//statistics data
func Statistics() {
	fmt.Println("ticker run ")
	m.Lock()
	//slices
	controls := make([]*MonitorControl, 0)

	monitorMap.Range(func(key, value interface{}) bool {
		control := value.(*MonitorControl)
		controls = append(controls, control)
		monitorMap.Delete(key)
		return true
	})
	m.Unlock()

	if len(controls) > 0 {
		fmt.Println("synchronization monitor data")
		monitor := Monitor{time.Now().UnixNano() / 1e6, controls}
		bytes, e := json.Marshal(monitor)
		if e == nil {
			postJson(string(bytes))
		}
	}

}

//Whether the system has been certified successfully
var sysStatus sync.Map

func validate(sc string) bool {
	value, ok := sysStatus.Load(sc)
	if ok {
		return value.(bool)
	}
	return false
}

func putSysStatus(sc string, flag bool) {
	sysStatus.Store(sc, flag)
}
