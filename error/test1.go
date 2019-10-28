package main

import (
	"errors"
	"log"
	"os"
	"strconv"
)

var (
	ErrNoo = errors.New("NOOOOOOOOOOOOOOOOOO")
)

func execute() error {
	fd, err := os.Open("unko")
	if err != nil {
		return err
	}

	defer fd.Close()

	buf := make([]byte, 8)

	_, err = fd.Read(buf)
	if err != nil {
		return err
	}

	num, err := strconv.Atoi(string(buf))
	if err != nil {
		return err
	}

	log.Println(num)

	return nil
}

func main() {
	err := execute()

	switch e := err.(type) {
	case *os.PathError:
		log.Println("PathError", e)
	default:
		log.Println("Other", e)
	}
}
