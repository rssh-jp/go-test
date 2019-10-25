package main

import (
	"log"
)

func main() {
	ch := make(chan int)

	for i := 0; i < 10; i++ {
		ch <- i
	}

	select {
	case n := <-ch:
		log.Println(n)
	}
}
