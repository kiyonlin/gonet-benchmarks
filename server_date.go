package main

import (
	"sync/atomic"
	"time"
)

func init() {
	updateServerDate()
}

var (
	ServerDate atomic.Value
)

func updateServerDate() {
	refreshServerDate()
	go func() {
		for {
			time.Sleep(time.Second)
			refreshServerDate()
		}
	}()
}

func refreshServerDate() {
	ServerDate.Store(time.Now().In(time.UTC).AppendFormat(nil, "Mon, 02 Jan 2006 15:04:05 GMT"))
}
