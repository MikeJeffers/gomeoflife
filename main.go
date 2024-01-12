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
	"golang.org/x/mobile/event/paint"
)

func main() {
	width := 1000
	height := 600
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{Width: width, Height: height})
		if err != nil {
			log.Fatal(err)
			return
		}
		defer w.Release()

		drift := 10
		ticker := time.NewTicker(500 * time.Millisecond)
		doneChan := make(chan bool)

		go func() {
			for {
				select {
				case t := <-ticker.C:
					fmt.Println("Tick at", t)
					drift += 10
					w.Fill(image.Rect(0, 0, width, height), color.RGBA{0x00, 0x00, 0x00, 0xff}, screen.Src)
					w.Fill(image.Rect(drift, 25, drift+10, 10), color.RGBA{0x00, 0xff, 0xff, 0x22}, screen.Over)
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
			case paint.Event:
				w.Publish()
			}
		}
	})
}
