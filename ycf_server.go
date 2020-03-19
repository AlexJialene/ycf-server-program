package main

import "C"
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
			return C.CString(string(marshal))
		}
	}

	return C.CString("")
}

func Put(monitorId *C.char, size C.int) {
	id := C.GoBytes(unsafe.Pointer(monitorId), size)
	s := string(id)
	if s != "" {
		//todo
	}
}

const initUrl = "http://127.0.0.1:9004/df/sys/systemMonitorList"
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
	/*s := post("ycf-home", "adcd7048512e64b48da55b027577886ee5a36350")
	response := BizResponse{}
	err := json.Unmarshal(s, &response)

	if err == nil {
		fmt.Println(response.Data[0].Path)

	}*/

	/*s := "ycf-home"
	i := "adcd7048512e64b48da55b027577886ee5a36350"

	bytes := []byte(s)
	i2 := []byte(i)

	systemCode := (*C.char)(unsafe.Pointer(&bytes))
	secret := (*C.char)(unsafe.Pointer(&i2))

	fmt.Println("123")
	system := InitSystem(systemCode ,secret, C.int(len(bytes)), C.int(len(i2)))
	fmt.Println(system)*/

}
