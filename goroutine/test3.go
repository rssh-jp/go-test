package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileWalk chan string

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.IsDir() {
		f <- path
	}

	return nil
}

func measure(t time.Time) {
	log.Println(time.Now().Sub(t))
}

func doSomething(filepath string) error {
	time.Sleep(time.Second)
	fd, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer fd.Close()

	buf := make([]byte, 1024)
	n, err := fd.Read(buf)
	if err != nil {
		return err
	}

	log.Println(n, filepath, string(buf))

	return nil
}

func main() {
	log.Println("START")
	defer log.Println("END")

	defer measure(time.Now())

	var targetPath string
	var semaphoreCount int

	flag.StringVar(&targetPath, "t", "./", "target path")
	flag.IntVar(&semaphoreCount, "s", 1, "semaphore count")
	flag.Parse()

	walker := make(fileWalk)

	go func() {
		if err := filepath.Walk(targetPath, walker.Walk); err != nil {
			log.Fatal(err)
		}

		close(walker)
	}()

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, semaphoreCount)

	for path := range walker {
		wg.Add(1)
		semaphore <- struct{}{}

		go func() {
			defer func() {
				<-semaphore
				wg.Done()
			}()
			doSomething(path)
		}()
	}

	wg.Wait()
}
