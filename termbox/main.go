package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

const c = termbox.ColorDefault

func drawString(x, y int, str string) {
	for index, r := range str {
		termbox.SetCell(x+index, y, r, c, c)
	}

	termbox.Flush()
}

func main() {
	log.Println("START")
	defer log.Println("END")

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer termbox.Close()

	termbox.Clear(c, c)

	width, height := termbox.Size()

	drawString(0, 0, "+")
	drawString(width-1, 0, "+")
	drawString(0, height-1, "+")
	drawString(width-1, height-1, "+")

	drawString(1, 1, "abcあいうそんdef")
	drawString(1, 10, fmt.Sprintf("%d, %d", width, height))

	var word string

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break loop
			case termbox.KeyEnter:
				if word == "exit" {
					break loop
				}

				word = ""
				w := ""
				for i := 0; i < width; i++ {
					w += " "
				}
				drawString(1, 3, w)
			default:
				word += string(ev.Ch)
				drawString(1, 3, word)
				drawString(1, 2, fmt.Sprintf("%+v", ev))
			}
		}
	}
}
