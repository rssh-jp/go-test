package main

import (
	"log"
)

func main() {
	ch := make(chan int, 1)

	log.Printf("%+v\n", ch)

	ch <- 30

	num, ok := <-ch

	log.Println(num, ok)

	ch <- 10

	close(ch)

	num, ok = <-ch

	log.Println(num, ok)

	num, ok = <-ch

	log.Println(num, ok)
}
