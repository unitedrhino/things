package utils

import (
	"sync"
	"time"
)

//带过期时间的sync.map
type ExpMap struct {
	Map sync.Map
}

type CheckExp interface {
	IsTimeOut() bool
}

func NewExpMap(checkTime time.Duration) *ExpMap {
	expMap := ExpMap{}
	go func() {
		time.Sleep(checkTime)
		expMap.Map.Range(func(key, value any) bool {
			v := value.(CheckExp)
			if v.IsTimeOut() {
				expMap.Map.Delete(key)
			}
			return true
		})
	}()
	return &expMap
}
