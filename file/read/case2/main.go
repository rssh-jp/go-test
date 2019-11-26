package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"time"
)

const (
	filepath   = "giga1"
	bufferSize = 1 * 1024 * 1024
)

func readBufioReader(r io.Reader) error {
	log.Println("-------", 1)
	b := bufio.NewReader(r)

	buf := make([]byte, bufferSize)

	log.Println("-------", 2)
	count := 0
	for i := 0; ; i++ {
		if i%100 == 0 {
			log.Println("##", i, "MB")
		}
		n, err := b.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		count += n
	}

	log.Println("####", 3)
	log.Println(count)

	return nil
}

func main() {
	s := time.Now()

	defer func() {
		log.Println(time.Now().Sub(s))
	}()

	var srcpath string
	flag.StringVar(&srcpath, "s", filepath, "source file path")
	flag.Parse()

	fd, err := os.Open(srcpath)
	if err != nil {
		log.Fatal(err)
	}

	defer fd.Close()

	log.Println(fd.Seek(0, os.SEEK_END))

	_, err = fd.Seek(0, os.SEEK_SET)
	if err != nil {
		log.Fatal(err)
	}

	err = readBufioReader(fd)
	if err != nil {
		log.Fatal(err)
	}
}
