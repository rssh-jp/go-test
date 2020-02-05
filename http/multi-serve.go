package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	log.Println("START")
	defer log.Println("END")

	http.HandleFunc("/", handler)

	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			err := http.ListenAndServe(":8080", nil)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	wg.Wait()
}
