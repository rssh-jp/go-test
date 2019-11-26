package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/rssh-jp/go-syncgroup"
)

func execute(filepath string) error {
	fd, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer fd.Close()

	buf := make([]byte, 4)

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
		"un",
		"un",
		"un",
		"un",
		"un",
		"un",
		"un",
		"un",
		"un",
		"un",
		"aaa",
	}

	var (
		chErr = make(chan error, len(list))
		sg    = syncgroup.New(semaphoreCount)
	)

	defer sg.Close()

	for _, item := range list {
		sg.Add()

		go func(chErr chan error, filepath string) {
			defer sg.Done()

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

	sg.Wait()

	for i := 0; i < len(list); i++ {
		select {
		case err := <-chErr:
			log.Println(err)
		default:
			break
		}
	}
}
