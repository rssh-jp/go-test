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

func Go(f func()error, ef func(error)){
    go func(){
        ef(f())
    }()
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

        Go(func()error{
            defer sg.Done()
            return execute(item)
        }, func(err error){
			if err != nil {
				switch err.(type) {
				case *os.PathError:
					chErr <- err
				default:
				}
			}
        })
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
