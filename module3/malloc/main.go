package main

//#cgo LDFLAGS:
//char* allocMemory();
import "C"
import (
	"fmt"
	"time"
)

func main() {
	holder := []*C.char{}

	for i := 1; i < 10; i++ {
		fmt.Printf("Allocating %dMb memory, raw memory is %d\n", i*100, i*100*1024*1025)
		holder = append(holder, (*C.char)(C.allocMemory()))
		time.Sleep(1 * time.Minute)
	}
}
