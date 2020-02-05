package main

import(
    "flag"
    "log"
    "strings"
    "math/rand"
    "time"
)

func main(){
    log.Println("START")
    defer log.Println("END")

    var argNames string
    flag.StringVar(&argNames, "n", "test1 test2 test3", "space separated name list")
    flag.Parse()

    names := strings.Split(argNames, " ")

    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    r.Shuffle(len(names), func(i, j int){
        names[i], names[j] = names[j], names[i]
    })

    for index, name := range names{
        log.Printf("%d番 : %s\n", index + 1, name)
    }
}
