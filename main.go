package main

import (
	"flag"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"log"
	"os"
	"strings"
)

var (
	otoContext *oto.Context
	content    map[string][]byte
)

func play(input []byte) {
	ap := otoContext.NewPlayer()
	if _, err := ap.Write(input); err != nil {
		log.Panicf("failed writing to player: %v", err)
	}
	if err := ap.Close(); err != nil {
		log.Panicf("failed to clone player: %v", err)
	}
}

func gui() {
	a := app.New()
	w := a.NewWindow("soundboard")
	ng := container.NewGridWithColumns(func() int {
		var col int
		c := len(content)
		ipc := 10 // items per column
		for {
			if c > ipc {
				col++
				c = c - ipc
				continue
			}
			if c < ipc && c > 0 {
				col++
				break
			}
		}
		return col
	}())
	for index, bytes := range content {
		sl := index
		nb := bytes
		ng.Add(widget.NewButton(sl, func() {
			go play(nb)
		}))
	}
	w.SetContent(ng)
	w.ShowAndRun()
}

func main() {
	var otoErr error
	otoContext, otoErr = oto.NewContext(48000, 2, 2, 256)
	if otoErr != nil {
		log.Panicf("Error creating oto.NewContext %v", otoErr)
	}
	var path string
	flag.StringVar(&path, "d", "audio/", "dir of 48khz mp3 files")
	flag.Parse()
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}
	content = map[string][]byte{}
	for _, entry := range dir {
		en := entry.Name()
		if en[len(en)-4:] != ".mp3" {
			continue
		}
		file, err := os.Open("audio/" + entry.Name())
		if err != nil {
			log.Panicf("Error opening file %s: %v", entry.Name(), err)
		}
		decodedFile, err := mp3.NewDecoder(file)
		if err != nil {
			log.Panicf("Error decoding file %s: %v", entry.Name(), err)
		}
		if content[strings.Replace(en, ".mp3", "", 1)], err = io.ReadAll(decodedFile); err != nil {
			log.Panicf("Error reading decodedFile %s: %v", entry.Name(), err)
		}
	}
	gui()
}
