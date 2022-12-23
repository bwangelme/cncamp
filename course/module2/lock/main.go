package main

import (
	"fmt"
	"sync"
	"time"
)

/**
RWMutex 的上锁规则

Lock 的优先级高于 RLock

Lock 上锁时， RLock 无法上锁
RLock 上锁时，RLock 可以上锁
RLock 上锁时，Lock 无法上锁
*/

func wLock(l *sync.RWMutex) {
	fmt.Println("Lock")
	l.Lock()
	time.Sleep(5 * time.Second)
	l.Unlock()
	fmt.Println("Unlock")
}

func rLock(l *sync.RWMutex) {
	time.Sleep(1 * time.Second)
	fmt.Println("Try RLock")
	l.RLock()
	fmt.Println("RLock Success")
}

func main() {
	var l = &sync.RWMutex{}

	go wLock(l)
	go rLock(l)

	time.Sleep(6 * time.Second)
}
