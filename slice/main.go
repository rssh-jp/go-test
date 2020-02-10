package main

import(
    "log"
)

func main(){
    const size = 1024 * 1024 * 1024 * 1024
    s := make([]byte, size)

    log.Println(len(s))
}
