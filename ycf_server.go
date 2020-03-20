package main

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

var isInit = false

//string  == id
//int  == path
var pattens map[string]string

//export InitSystem
func InitSystem(sc, s *C.char, size1, size2 C.int) *C.char {
	if !isInit {
		systemCode := C.GoBytes(unsafe.Pointer(sc), size1)
		secret := C.GoBytes(unsafe.Pointer(s), size2)
		bytes := post(string(systemCode), string(secret))

		response := BizResponse{}
		json.Unmarshal(bytes, &response)
		if response.Ret == "success" {
			marshal, err := json.Marshal(response.Data)
			if err == nil {
				//start config
				isInit = true
				pattens = make(map[string]string, len(response.Data))
				for _, v := range response.Data {
					monitor := v
					pattens[monitor.Id] = monitor.Path
				}
				//init cron
				s := "0/10 * * * * ? "
				i := cron.New()
				addFunc := i.AddFunc(s, Statistics)
				if addFunc == nil {
					i.Start()
				}

				return C.CString(string(marshal))
			}
		}
	}

	return C.CString("")
}

//export PushUrl
func PushUrl(str *C.char, size C.int) {
	if isInit {
		s := C.GoBytes(unsafe.Pointer(str), size)
		url := string(s)
		if url == "" {
			//if eq url
			for k, v := range pattens {
				matched, err := regexp.MatchString(v, url)
				if err == nil {
					if matched {
						Add(k, 1)
						return
					}
				}
			}
		}
	}
}

//export Put
func Put(monitorId *C.char, size, count C.int) {
	if isInit {
		id := C.GoBytes(unsafe.Pointer(monitorId), size)
		s := string(id)
		if s != "" {
			Add(s, int(count))
		}
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
