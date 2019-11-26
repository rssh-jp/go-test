package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func execute(filepath string) error {
	fd, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer fd.Close()

	buf := make([]byte, 8)

	n, err := fd.Read(buf)
	if err != nil {
		return err
	}

	num, err := strconv.Atoi(strings.Trim(string(buf[:n]), "\n"))
	if err != nil {
		return err
	}

	log.Println(num)

	return nil
}
func main() {
	const (
		semaphoreCount = 3
		retryCount     = 3
	)

	list := []string{
		"unko",
		"aaa",
	}

	var (
		chSemaphore = make(chan struct{}, semaphoreCount)
		chErr       = make(chan error, len(list))
		wg          sync.WaitGroup
	)

	defer func() {
		close(chSemaphore)
		close(chErr)
	}()

	wg.Add(len(list))

	for _, item := range list {
		chSemaphore <- struct{}{}

		go func(chErr chan error, filepath string) {
			defer func() {
				wg.Done()
				<-chSemaphore
			}()

		loop:
			for i := 0; i < retryCount; i++ {
				err := execute(filepath)
				if err != nil {
					switch err.(type) {
					case *os.PathError:
						chErr <- err
						break loop
					default:
					}
				}
				break
			}
		}(chErr, item)
	}

	wg.Wait()

	for i := 0; i < len(list); i++ {
		select {
		case err := <-chErr:
			log.Println(err)
		default:
			break
		}
	}
}
