package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	timeoutCtx, _ := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	go func(ctx context.Context) {
		ticker := time.NewTicker(1 * time.Second)
		for _ = range ticker.C {
			select {
			case <-ctx.Done():
				fmt.Println("child process interrupt")
				return
			default:
				fmt.Println("enter default")
			}
		}
	}(timeoutCtx)
	time.Sleep(1 * time.Second)
	select {
	case <-timeoutCtx.Done():
		time.Sleep(1 * time.Second)
		fmt.Println("main process exit")
	}
}
