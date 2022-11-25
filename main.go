package main

import (
	"flag"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"log"
	"os"
)

type soundboard struct {
	otoContext *oto.Context
	ap         *oto.Player
	content    map[string][]byte
}

func (s *soundboard) play(input []byte) {
	if _, err := s.ap.Write(input); err != nil {
		log.Panicf("%v", err)
	}
}

func (s *soundboard) gui() {
	a := app.New()
	w := a.NewWindow("soundboard")
	vb := container.NewVBox()
	for s2, bytes := range s.content {
		nb := bytes
		vb.Add(widget.NewButton(s2, func() {
			s.play(nb)
		}))
	}
	w.SetContent(vb)
	w.ShowAndRun()
}

func main() {
	s := soundboard{}
	var otoErr error
	s.otoContext, otoErr = oto.NewContext(48000, 2, 2, 256)
	if otoErr != nil {
		log.Panicf("Error creating oto.NewContext %v", otoErr)
	}
	s.ap = s.otoContext.NewPlayer()
	var path string
	flag.StringVar(&path, "d", "audio/", "dir of 48khz mp3 files")
	flag.Parse()
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}
	s.content = map[string][]byte{}
	for i, entry := range dir {
		// TODO : only open mp3 files
		fmt.Println(i, entry.Name())
		playFile, err := os.Open("audio/" + entry.Name())
		if err != nil {
			log.Panicf("Error opening file %s: %v", entry.Name(), err)
		}
		decodedFile, err := mp3.NewDecoder(playFile)
		if err != nil {
			log.Panicf("Error decoding file %s: %v", entry.Name(), err)
		}
		if s.content[entry.Name()], err = io.ReadAll(decodedFile); err != nil {
			log.Panicf("Error reading decodedFile %s: %v", entry.Name(), err)
		}
	}
	s.gui()
}
