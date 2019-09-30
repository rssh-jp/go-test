package main

import (
	"log"
	"os"
	"time"
)

const (
	retryNumMax = 3
)

type execResult struct {
	duration   time.Duration
	retryCount int
	err        error
}

func failure(filename string, count int) (res execResult) {
	defer func() execResult {
		log.Printf("defer : %+v\n", res)
		if res.err != nil && retryNumMax > count {
			return failure(filename, count+1)
		}

		return res
	}()

	s := time.Now()

	fd, err := os.OpenFile(filename, os.O_APPEND, 0644)
	if err != nil {
		res.err = err
		return
	}

	defer fd.Close()

	res.duration = time.Now().Sub(s)
	res.retryCount = count

	log.Println("-----------")

	return
}

func recursiveExec1(filename string, count int) (res execResult) {
	res.retryCount = count
	res.duration, res.err = exec1(filename)

	log.Printf("      : %+v\n", res)

	if res.err != nil && retryNumMax > count {
		return recursiveExec1(filename, count+1)
	}

	return res
}

func exec1(filename string) (duration time.Duration, err error) {
	s := time.Now()

	fd, err := os.OpenFile(filename, os.O_APPEND, 0644)
	if err != nil {
		return
	}

	defer fd.Close()

	duration = time.Now().Sub(s)

	return
}

func main() {
	res := failure("nice.file", 0)
	log.Printf("res   : %+v\n", res)
	if res.err != nil {
		log.Fatal(res.err)
	}

	res = recursiveExec1("nice.file", 0)
	log.Printf("res   : %+v\n", res)
	if res.err != nil {
		log.Fatal(res.err)
	}
}
