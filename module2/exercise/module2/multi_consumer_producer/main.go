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

func Producer(ctx context.Context, id int, queue chan<- int, wg *sync.WaitGroup) {
	t := time.NewTicker(time.Millisecond * 500)
	wg.Add(1)

	fmt.Printf("Start Producer %d\n", id)
LOOP:
	for {
		select {
		case <-t.C:
			queue <- rand.Intn(1000)
		case <-ctx.Done():
			fmt.Printf("Producer %d get cancel\n", id)
			break LOOP
		}
	}

	t.Stop()

	fmt.Printf("Stop the Producer %d\n", id)
	wg.Done()
}

func Consumer(ctx context.Context, id int, queue <-chan int, wg *sync.WaitGroup) {
	t := time.NewTicker(time.Second * 1)
	wg.Add(1)

	fmt.Printf("Start Consumer %d\n", id)
LOOP:
	for {
		select {
		case <-t.C:
			v := <-queue
			fmt.Printf("Consumer %d get value %v\n", id, v)
		case <-ctx.Done():
			fmt.Printf("Consumer %d get cancel\n", id)
			break LOOP
		}
	}

	t.Stop()
	fmt.Printf("Stop the Consumer %d\n", id)
	wg.Done()
}

func main() {
	var wg = new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	queue := make(chan int, 10)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < 4; i++ {
		go Consumer(ctx, i, queue, wg)
	}
	for i := 0; i < 2; i++ {
		go Producer(ctx, i, queue, wg)
	}

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
