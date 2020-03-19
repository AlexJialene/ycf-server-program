package main

import (
	"sync"
)

type Monitor struct {
	id    int64
	count int
}

//系统是否认证成功
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
