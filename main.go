package main

import (
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"log"
	"os"
)

var (
	otoContext *oto.Context
	ap         *oto.Player
)

func play(input []byte) {
	if _, err := ap.Write(input); err != nil {
		log.Panicf("%v", err)
	}
}

func main() {
	var otoErr error
	otoContext, otoErr = oto.NewContext(48000, 2, 2, 1000)
	if otoErr != nil {
		log.Panicf("Error creating oto.NewContext %v", otoErr)
	}
	ap = otoContext.NewPlayer()
	dir, err := os.ReadDir("audio/")
	if err != nil {
		return
	}
	sm := map[int][]byte{}
	for i, entry := range dir {
		fmt.Println(i, entry.Name())
		playFile, err := os.Open("audio/" + entry.Name())
		if err != nil {
			log.Panicf("Error opening file %s: %v", entry.Name(), err)
		}
		decodedFile, err := mp3.NewDecoder(playFile)
		if err != nil {
			log.Panicf("Error decoding file %s: %v", entry.Name(), err)
		}
		if sm[i], err = io.ReadAll(decodedFile); err != nil {
			log.Panicf("Error reading decodedFile %s: %v", entry.Name(), err)
		}
	}
	for {
		var t int
		if _, err := fmt.Scan(&t); err != nil {
			log.Panicf("%v", err)
		}
		// If close player and make new, will it interrupt playback
		play(sm[t])
	}
}
