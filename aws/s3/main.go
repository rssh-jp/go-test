package main

import (
	"io"
	"log"
	"os"
	"time"
	//"net"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	//"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/rssh-jp/go-bandwidth"
)

const (
	filepath = "/opt/rsrc/giga1"
	//filepath = "/opt/rsrc/halfgiga1"
	//filepath = "/opt/rsrc/mega1"
	//filepath = "/opt/rsrc/byte100"
)

func measuretime(t time.Time) {
	log.Println(time.Now().Sub(t))
}

type MyTransport struct {
	http.Transport
}

type ReadCloser struct {
	*bandwidth.ReadWriter
	io.Closer
}

func (r *ReadCloser) Read(p []byte) (int, error) {
	return r.ReadWriter.Read(p)
}
func (r *ReadCloser) Close() error {
	return r.Closer.Close()
}

func (t *MyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	log.Println("++++++++++++++++++++++++")
	log.Printf("%+v\n", req.Body)
	if req.Body == nil {
		return t.Transport.RoundTrip(req)
	}
	req.Body = &ReadCloser{
		ReadWriter: bandwidth.NewReader(req.Body, 10*1024*1024, time.Second),
		Closer:     req.Body,
	}
	return t.Transport.RoundTrip(req)
}

func main() {
	log.Println("START")
	defer log.Println("END")
	defer measuretime(time.Now())

	log.Println("+++++++++", 1)
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"),
		HTTPClient: &http.Client{
			Transport: &MyTransport{
				//Transport: http.Transport{
				//    Proxy: http.ProxyFromEnvironment,
				//    DialContext: (&net.Dialer{
				//        Timeout: 30 * time.Second,
				//        KeepAlive: 30 * time.Second,
				//        DualStack: true,
				//    }).DialContext,
				//    ForceAttemptHTTP2: true,
				//    MaxIdleConns: 100,
				//    IdleConnTimeout: 90 * time.Second,
				//    TLSHandshakeTimeout: 10 * time.Second,
				//    ExpectContinueTimeout: 1 * time.Second,

				//    //ResponseHeaderTimeout: 1 * time.Second,
				//    //WriteBufferSize: 1 * 1024 * 1024,
				//},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("+++++++++", 2)
	s3Inst := s3.New(sess)
	//uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader){
	//    //u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(1 * 1024 * 1024)
	//    //u.MaxUploadParts = 1000
	//    ////u.PartSize = 100 * 1024 * 1024
	//    //u.Concurrency = 1
	//})

	log.Println("+++++++++", 3)
	fd, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("+++++++++", 4)
	defer fd.Close()

	log.Println("+++++++++", 5)
	res, err := s3Inst.PutObject(&s3.PutObjectInput{
		Bucket: aws.String("test-araumi"),
		Key:    aws.String("giga1"),
		Body:   fd,
	})
	//res, err := uploader.Upload(&s3manager.UploadInput{
	//	Bucket: aws.String("test-rssh"),
	//	Key:    aws.String("giga1"),
	//	//Body: bandwidth.NewReader(fd, 1 * 1024 * 1024, time.Second),
	//	Body: fd,
	//})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("##################RES", res)
}
