package main

import (
	_ "embed"
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

//go:embed sb.ico
var icon []byte

var (
	otoContext *oto.Context
	content    map[string][]byte
	ipc        = 10 // items per column
)

func play(input []byte) {
	ap := otoContext.NewPlayer()
	if _, err := ap.Write(input); err != nil {
		log.Panicf("failed writing to player: %v", err)
	}
	if err := ap.Close(); err != nil {
		log.Panicf("failed to close player: %v", err)
	}
}

func gui() {
	a := app.New()
	w := a.NewWindow("go-soundboard")
	w.SetIcon(fyne.NewStaticResource("", icon))
	ng := container.NewGridWithColumns(func() int {
		var col int
		c := len(content)
		if c%ipc == 0 {
			return c / ipc
		}
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
	keys := make([]string, 0, len(content))
	for k := range content {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sl := k
		nb := content[k]
		ng.Add(widget.NewButton(sl, func() {
			go play(nb)
		}))
	}
	w.SetContent(ng)
	w.ShowAndRun()
}

func main() {
	var otoErr error
	otoContext, otoErr = oto.NewContext(48000, 2, 2, 1024)
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
		file, err := os.Open(path + en)
		if err != nil {
			log.Panicf("Error opening file %s: %v", en, err)
		}
		decodedFile, err := mp3.NewDecoder(file)
		if err != nil {
			log.Panicf("Error decoding file %s: %v", en, err)
		}
		if content[strings.Replace(en, ".mp3", "", 1)], err = io.ReadAll(decodedFile); err != nil {
			log.Panicf("Error reading decodedFile %s: %v", en, err)
		}
	}
	gui()
}
