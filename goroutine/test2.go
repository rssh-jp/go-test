package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	errInvalid = errors.New("Invelid")
)

func measure(t time.Time) {
	log.Println(time.Now().Sub(t))
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

	chSrc := make(chan string)
	chDone := make(chan struct{})
	chErr := make(chan error)
	chIsEnd := make(chan bool)

	defer close(chSrc)
	defer close(chDone)
	defer close(chErr)
	defer close(chIsEnd)

	var wg sync.WaitGroup

	// 実行するやつ
	go func(chSrc chan string, chErr chan error, wg *sync.WaitGroup, semaphoreCount int) {
		chSemaphore := make(chan struct{}, semaphoreCount)

		for loop := true; loop; {
			select {
			case p, ok := <-chSrc:
				if !ok {
					loop = false
					break
				}
				chSemaphore <- struct{}{}
				go func() {
					defer func() {
						wg.Done()
						<-chSemaphore
					}()

					err := doSomething(p)
					if err != nil {
						chErr <- err
						loop = false
					}
				}()
			}
		}
	}(chSrc, chErr, &wg, semaphoreCount)

	// 追加するやつ
	go func(chSrc chan string, chErr chan error, chIsEnd chan bool, wg *sync.WaitGroup, targetPath string) {
		readdir(chSrc, chErr, wg, targetPath)

		chIsEnd <- true
	}(chSrc, chErr, chIsEnd, &wg, targetPath)

	// 終了を検知するやつ
	go func(chDone chan struct{}, chIsEnd chan bool, wg *sync.WaitGroup) {
		defer func() {
			chDone <- struct{}{}
		}()

		isEnd := <-chIsEnd
		if !isEnd {
			return
		}

		wg.Wait()
	}(chDone, chIsEnd, &wg)

	// 終了を待つ
	select {
	case err := <-chErr:
		log.Fatal(err)
	case <-chDone:
		break
	}
}

// ディレクトリを読み取り
func readdir(chSrc chan string, chErr chan error, wg *sync.WaitGroup, targetPath string) {
	fileinfos, err := ioutil.ReadDir(targetPath)
	if err != nil {
		chErr <- err
		return
	}

	for _, fileinfo := range fileinfos {
		srcpath := filepath.Join(targetPath, fileinfo.Name())

		if fileinfo.IsDir() {
			readdir(chSrc, chErr, wg, srcpath)
		} else {
			wg.Add(1)

			log.Println(srcpath)
			chSrc <- srcpath
		}
	}
}

// 実行したい関数
func doSomething(srcpath string) error {
	time.Sleep(time.Second)

	info, err := os.Stat(srcpath)
	if err != nil {
		return err
	}

	if info.Name() == "d" {
		return errInvalid
	}

	return nil
}
