package main

import (
	"image"
	"image/color"
	"log"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

func runGrid(send, recv chan string, grid *Grid) {
	go func() {
		for msg := range recv {
			if msg == "end" {
				return
			}
			grid.nextState()
			send <- "updated"
		}
	}()
}

func main() {
	width := 1000
	height := 600
	size := 25
	var grid Grid
	grid.init(width/size, height/size)
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{Width: width, Height: height, Title: "GOme of Life"})
		if err != nil {
			log.Fatal(err)
			return
		}
		defer w.Release()

		send := make(chan string)
		recv := make(chan string)
		runGrid(recv, send, &grid)
		send <- "go!"

		go func() {
			for {
				_, ok := <-recv
				if !ok {
					return
				}
				w.Fill(image.Rect(0, 0, width, height), color.RGBA{0x00, 0x00, 0x00, 0xff}, screen.Over)
				for _, cell := range grid.cells {
					x, y, state := cell.drawState()
					if state == 1 {
						w.Fill(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size), color.RGBA{0x00, 0xff, 0xff, 0x22}, screen.Over)
					}
				}
				w.Publish()
				time.Sleep(500 * time.Millisecond)
				send <- "ready"
			}
		}()

		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					send <- "end"
					return
				}
			}
		}
	})
}
