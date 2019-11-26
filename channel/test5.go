package main

import (
	"log"
	"sync"
)

func test(size int) {
	ch := make(chan int, size)
	ch <- 1
	log.Printf("size=%d, ch=%d\n", size, <-ch)
}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		test(0)
	}()
	go func() {
		defer wg.Done()
		test(1)
	}()

	wg.Wait()
}
