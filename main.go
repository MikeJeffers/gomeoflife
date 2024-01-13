package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

func main() {
	width := 1000
	height := 600
	size := 25
	var grid Grid
	grid.init(width/size, height/size)
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{Width: width, Height: height, Title: "REKT"})
		if err != nil {
			log.Fatal(err)
			return
		}
		defer w.Release()

		ticker := time.NewTicker(1 * time.Second)
		doneChan := make(chan bool)

		go func() {
			for {
				select {
				case t := <-ticker.C:
					fmt.Println("Tick at", t)
					grid.nextState()
					w.Fill(image.Rect(0, 0, width, height), color.RGBA{0x00, 0x00, 0x00, 0xff}, screen.Src)
					for _, cell := range grid.cells {
						x, y, state := cell.drawState()
						if state == 1 {
							w.Fill(image.Rect(x*size, y*size, (x+1)*size, (y+1)*size), color.RGBA{0x00, 0xff, 0xff, 0x22}, screen.Over)
						}
					}
					w.Publish()
				case <-doneChan:
					return
				}
			}
		}()

		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					doneChan <- true
					return
				}
				// case paint.Event:
				// 	w.Publish()
			}
		}
	})
}
