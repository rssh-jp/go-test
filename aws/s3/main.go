package main

import (
	"log"
	"os"
	"time"
    "net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

    //"github.com/rssh-jp/go-bandwidth"
)

const (
	filepath = "/extdisk1/tmp/giga1"
	//filepath = "/extdisk1/tmp/mega1"
	//filepath = "/extdisk1/tmp/byte100"
)


func measuretime(t time.Time) {
	log.Println(time.Now().Sub(t))
}

func main() {
	log.Println("START")
	defer log.Println("END")
	defer measuretime(time.Now())

	log.Println("+++++++++", 1)
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("ap-northeast-1"),
        HTTPClient: &http.Client{
            Transport: &http.Transport{
                ResponseHeaderTimeout: 1 * time.Second,
                WriteBufferSize: 1024 * 1024,
            },
        },
    })
	if err != nil {
		log.Fatal(err)
	}

	log.Println("+++++++++", 2)
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader){
        u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(1 * 1024 * 1024)
        u.MaxUploadParts = 1000
        //u.PartSize = 100 * 1024 * 1024
        u.Concurrency = 1
	})

	log.Println("+++++++++", 3)
	fd, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("+++++++++", 4)
	defer fd.Close()

	log.Println("+++++++++", 5)
	res, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("test-rssh"),
		Key:    aws.String("giga1"),
		//Body: bandwidth.NewReader(fd, 1 * 1024 * 1024, time.Second),
		Body: fd,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("##################RES", res)
}
