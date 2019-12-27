package main

import (
	"log"
)

func main() {
	ch := make(chan int, 1)
	//chEnd := make(chan struct{})
	//go func(){
	//    loop:
	//    for{
	//        select{
	//        case num := <-ch:
	//            log.Println(len(ch))
	//            switch num{
	//            case 1:
	//                log.Println("++++++++++++++++")
	//                ch <- 2
	//            default:
	//                log.Println("--------------")
	//                break loop
	//            }
	//        }
	//    }
	//    chEnd <- struct{}{}
	//}()

	ch <- 1

	log.Println(<-ch)

	//<-chEnd
}
