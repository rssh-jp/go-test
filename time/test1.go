package main

import(
    "log"
    "time"
)

func main(){
    t := time.Now()
    t = t.AddDate(0, 0, -1)
    t = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
    log.Println(t)
}
