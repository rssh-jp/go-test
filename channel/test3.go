package main

import (
	"log"
	"sync"
)

var (
	ch = make(chan int)
	wg sync.WaitGroup
)

func main() {
	chEnd := make(chan struct{})
	go func() {
	loop:
		for {
			select {
			case num := <-ch:
				switch num {
				case 1:
					log.Println("++++++++++++++++")
					execute(1)
				default:
					log.Println("--------------")
					break loop
				}
			}
		}
		chEnd <- struct{}{}
	}()

	execute(0)

	wg.Wait()

	log.Println("+")

	<-chEnd
}

func execute(count int) {
	if count == 1 {
		ch <- 2
		return
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		ch <- 1
	}()
}
