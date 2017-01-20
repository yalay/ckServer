package controllers

import (
	"sync"
)

var gCount *Count

type Count struct {
	sync.RWMutex
	imCount map[string]int
	ckCount map[string]int
}

func init() {
	gCount = newCount()
}

func newCount() *Count {
	return &Count{
		imCount: make(map[string]int, 0),
		ckCount: make(map[string]int, 0),
	}
}

func GetImCount(articleKey string) int {
	gCount.RLock()
	imCount := gCount.imCount[articleKey]
	gCount.RUnlock()
	return imCount
}

func GetCkCount(articleKey string) int {
	gCount.RLock()
	ckCount := gCount.ckCount[articleKey]
	gCount.RUnlock()
	return ckCount
}

func IncImCount(articleKey string) {
	gCount.Lock()
	curCount := gCount.imCount[articleKey]
	gCount.imCount[articleKey] = curCount + 1
	gCount.Unlock()
}

func IncCkCount(articleKey string) {
	gCount.Lock()
	curCount := gCount.ckCount[articleKey]
	gCount.ckCount[articleKey] = curCount + 1
	gCount.Unlock()
}
