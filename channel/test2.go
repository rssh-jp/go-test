package main

import (
	"log"
)

func main() {
	const count = 10

	ch := make(chan int, count+10)

	for i := 0; i < count; i++ {
		ch <- i
	}

	log.Println("len", len(ch))
	for n := range ch {
		log.Println(n)
	}
	for {
		select {
		case n := <-ch:
			log.Println(n)
		}
	}
}
