package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Producer(ctx context.Context, queue chan<- int, wg *sync.WaitGroup) {
	t := time.NewTicker(time.Second * 1)
	wg.Add(1)

LOOP:
	for {
		select {
		case <-t.C:
			queue <- rand.Intn(1000)
		case <-ctx.Done():
			fmt.Println("Producer get cancel")
			break LOOP
		}
	}

	t.Stop()

	fmt.Println("Stop the Producer")
	wg.Done()
}

func Consumer(ctx context.Context, queue <-chan int, wg *sync.WaitGroup) {
	t := time.NewTicker(time.Second * 1)
	wg.Add(1)

LOOP:
	for {
		select {
		case <-t.C:
			v := <-queue
			fmt.Println("Consumer get value", v)
		case <-ctx.Done():
			fmt.Println("Consumer get cancel")
			break LOOP
		}
	}

	t.Stop()
	fmt.Println("Stop the Consumer")
	wg.Done()
}

func main() {
	var wg = new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan int, 10)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go Consumer(ctx, queue, wg)
	go Producer(ctx, queue, wg)

	for {
		select {
		case sig := <-sigs:
			fmt.Println("Receive ", sig)
			cancel()
			wg.Wait()
			fmt.Println("Stop the main worker")
			return
		}
	}
}
