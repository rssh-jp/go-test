package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"reflect"
	"time"
)

const (
	filepath = "/opt/rsrc/giga1"
	//filepath = "tmp/killo1"
	bufferSize = 1 * 1024 * 1024
)

func main() {
	s := time.Now()

	defer func() {
		log.Println(time.Now().Sub(s))
	}()

	log.Println("+++++++++++++", 1)
	fd, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("+++++++++++++", 2)

	defer fd.Close()

	log.Println(fd.Seek(0, os.SEEK_END))

	log.Println("+++++++++++++", 3)
	for _, item := range []func(io.Reader) error{
		readBufioReader,
		//readBufioScanner,
		//readReader,
	} {
		func(f func(r io.Reader) error) {
			log.Println("------------", 1)
			s := time.Now()

			defer func() {
				log.Println(getFunctionName(f), time.Now().Sub(s))
			}()

			log.Println("------------", 2)
			_, err := fd.Seek(0, os.SEEK_SET)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("------------", 3)

			err = f(fd)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("------------", 4)
		}(item)
	}

}

func getFunctionName(f func(r io.Reader) error) string {
	typ := reflect.TypeOf(f)
	return typ.Name()
}

func readBufioReader(r io.Reader) error {
	log.Println("####", 1)
	b := bufio.NewReader(r)

	buf := make([]byte, bufferSize)

	log.Println("####", 2)
	count := 0
	for i := 0; ; i++ {
		if i%100 == 0 {
			log.Println("     ##", i)
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

func readBufioScanner(r io.Reader) error {
	b := bufio.NewScanner(r)

	count := 0

	for b.Scan() {
		buf := b.Bytes()
		count += len(buf) + 1
	}

	log.Println(count)

	return nil
}

func readReader(r io.Reader) error {
	buf := make([]byte, bufferSize)

	count := 0
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		count += n
	}

	log.Println(count)

	return nil
}
