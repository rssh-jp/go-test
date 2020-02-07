package main

import (
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	//filepath = "/extdisk1/tmp/giga1"
	filepath = "/extdisk1/tmp/mega1"
	//filepath = "/extdisk1/tmp/byte100"
)

type ReadSeeker struct {
	file     *os.File
	limit    int
	bytes    int
	t        time.Time
	duration time.Duration
	count    int
}

func NewReadSeeker(fd *os.File) *ReadSeeker {
	return &ReadSeeker{
		file:     fd,
		limit:    100 * 1024,
		t:        time.Now(),
		duration: time.Second,
	}
}

const (
	taskCheck = iota + 1
	taskExec
	taskSleep
	taskEnd
)

func (r *ReadSeeker) Read(p []byte) (n int, err error) {
	if r.count < 7 {
		return r.file.Read(p)
	}

	task := taskCheck
	var index int
	var retn int
	for isLoop := true; isLoop; {
		switch task {
		case taskCheck:
			if index >= len(p) {
				task = taskEnd
				break
			}

			if r.bytes >= r.limit {
				task = taskSleep
			} else {
				task = taskExec
			}
		case taskExec:
			size := r.limit - r.bytes
			if size > len(p[index:]) {
				size = len(p[index:])
			}

			b := p[index : index+size]

			n, err = r.file.Read(b)

			index += n
			retn += n
			r.bytes += n

			if err != nil {
				task = taskEnd
				break
			}

			task = taskCheck

		case taskSleep:
			diff := r.duration - time.Now().Sub(r.t)
			if diff > 0 {
				log.Println("SLEEP", diff, r.bytes)
				time.Sleep(diff)
			}

			r.t = time.Now()

			r.bytes -= r.limit

			task = taskCheck
		case taskEnd:
			isLoop = false
		}
	}

	return retn, err
}

func (r *ReadSeeker) Seek(offset int64, whence int) (n int64, err error) {
	r.count++
	log.Println("SEEK", offset, whence)
	return r.file.Seek(offset, whence)
}

func measuretime(t time.Time) {
	log.Println(time.Now().Sub(t))
}

func main() {
	log.Println("START")
	defer log.Println("END")
	defer measuretime(time.Now())

	log.Println("+++++++++", 1)
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("+++++++++", 2)
	svc := s3.New(sess, &aws.Config{
		Region: aws.String("ap-northeast-1"),
	})

	log.Println("+++++++++", 3)
	fd, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("+++++++++", 4)
	defer fd.Close()

	r := NewReadSeeker(fd)

	log.Println("+++++++++", 5)
	res, err := svc.PutObject(&s3.PutObjectInput{
		Body: r,
		//Body: fd,
		Bucket: aws.String("test-araumi"),
		Key:    aws.String("giga1"),
		Metadata: map[string]*string{
			"max_bandwidth": aws.String("10KB/S"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("##################RES", res)
}
