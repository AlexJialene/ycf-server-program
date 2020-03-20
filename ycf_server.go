package main

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unsafe"
)

//export InitSystem
func InitSystem(sc, s *C.char, size1, size2 C.int) *C.char {
	systemCode := C.GoBytes(unsafe.Pointer(sc), size1)
	secret := C.GoBytes(unsafe.Pointer(s), size2)
	bytes := post(string(systemCode), string(secret))

	response := BizResponse{}
	json.Unmarshal(bytes, &response)
	if response.Ret == "success" {
		marshal, err := json.Marshal(response.Data)
		if err == nil {

			s := "0/10 * * * * ? "
			i := cron.New()
			addFunc := i.AddFunc(s, Statistics)
			if addFunc == nil {
				i.Start()
			}

			return C.CString(string(marshal))
		}
	}

	return C.CString("")
}

//export Put
func Put(monitorId *C.char, size, count C.int) {
	id := C.GoBytes(unsafe.Pointer(monitorId), size)
	s := string(id)
	if s != "" {
		Add(s, int(count))
	}
}

const initUrl = "http://127.0.0.1:9004/df/sys/systemMonitorList"
const pushUrl = "http://127.0.0.1:9004/df/sys/systemMonitorPush"
const contentType = "application/x-www-form-urlencoded"

func post(sc, s string) []byte {
	param := "systemCode=" + sc + "&secret=" + s
	resp, err := http.Post(initUrl, contentType, strings.NewReader(param))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	return bytes

}

func get() {
	resp, err := http.Get("http://127.0.0.1:9004/df/sys/systemMonitorPush")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
}

func postJson(s string) {
	request, e := http.NewRequest("POST", pushUrl, strings.NewReader(s))
	if e != nil {
		fmt.Println(e)
	} else {
		request.Header.Add("content-type", "application/json")
		client := &http.Client{Timeout: 60 * time.Second}
		resp, e := client.Do(request)
		if e == nil {
			request.Body.Close()
			resp.Body.Close()
		}
	}
}

type BizResponse struct {
	Ret  string          `json:"ret"`
	Data []SystemMonitor `json:"data"`
}

type SystemMonitor struct {
	Id       string `json:"id"`
	SystemId string `json:"systemId"`
	Path     string `json:"path"`
}

func main() {

}
