package main

import(
    "log"
    "os"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

const(
    //filepath = "/opt/rsrc/giga1"
    filepath = "/opt/rsrc/mega1"
    //filepath = "/opt/rsrc/byte100"
)


type ReadSeeker struct{
    file *os.File
    limit int
    count int
    bytes int
}

func NewReadSeeker(fd *os.File)*ReadSeeker{
    return &ReadSeeker{
        file: fd,
        limit: 10 * 1024,
    }
}

func (r *ReadSeeker)Read(p []byte)(n int, err error){
    log.Println("IN", len(p))
    var retn int
    var index int
    for{
        t := time.Now()

        size := r.limit
        if len(p[index:]) < r.limit{
            size = len(p[index:])
        }
        if size == 0{
            break
        }

        log.Println(index, size)

        b := p[index:index + size]
        n, err := r.file.Read(b)
        if err != nil{
            log.Println("+++++++++++++++++++++++", n, err)
            return retn, err
        }

        index += size

        retn += n

        if retn >= len(p){
            break
        }

        diff := time.Second - time.Now().Sub(t)
        if diff > 0{
            time.Sleep(diff)
        }
    }

    log.Println("OUT", retn)
    return retn, nil
}

func (r *ReadSeeker)Read2(p []byte)(n int, err error){
    log.Println("###", r.count, len(p))
    defer func(){
        r.count += 1
    }()
    var index int
    var retn int
    for range time.Tick(time.Second){
        size := r.limit
        if len(p[index:]) < r.limit{
            size = len(p[index:])
        }
        if size == 0{
            break
        }
        log.Println(size)
        b := p[index:index + size]
        n, err := r.file.Read(b)
        if err != nil{
            return -1, err
        }

        retn += n

        index = index + size
    }

    return retn, nil
}

func (r *ReadSeeker)Read3(p []byte)(n int, err error){
    size := r.limit
    if len(p) < r.limit{
        size = len(p)
    }
    p = p[:size]
    if len(p) > 0{
        return r.file.Read(p)
    }
    return
}

func (r *ReadSeeker)Seek(offset int64, whence int)(n int64, err error){
    return r.file.Seek(offset, whence)
}

func main(){
    log.Println("START")
    defer log.Println("END")

    log.Println("+++++++++", 1)
    sess, err := session.NewSession()
    if err != nil{
        log.Fatal(err)
    }

    log.Println("+++++++++", 2)
    svc := s3.New(sess, &aws.Config{
        Region: aws.String("ap-northeast-1"),
    })

    log.Println("+++++++++", 3)
    fd, err := os.OpenFile(filepath, os.O_RDONLY, 0755)
    if err != nil{
        log.Fatal(err)
    }

    log.Println("+++++++++", 4)
    defer fd.Close()

    r := NewReadSeeker(fd)
    
    log.Println("+++++++++", 5)
    res, err := svc.PutObject(&s3.PutObjectInput{
        Body: r,
        Bucket: aws.String("test-araumi"),
        Key: aws.String("giga1"),
    })
    if err != nil{
        log.Fatal(err)
    }

    log.Println("##################RES", res)
}
