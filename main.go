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
	audioPlayer *oto.Context
)

func play(decodeFile *mp3.Decoder) {
	log.Printf("Loading %d bytes", decodeFile.Length())
	ap := audioPlayer.NewPlayer()
	if _, err := io.Copy(ap, decodeFile); err != nil {
		log.Printf("Failed playing file: %v", err)

	}
	defer ap.Close()
}

func main() {
	var otoErr error
	audioPlayer, otoErr = oto.NewContext(48000, 2, 2, 4096)
	if otoErr != nil {
		log.Panicf("Error oto.NewContext %v", otoErr)

	}
	dir, err := os.ReadDir("audio/")
	if err != nil {
		return
	}
	sm := map[int]*mp3.Decoder{}
	for i, entry := range dir {
		fmt.Println(i, entry.Name())

		playFile, err := os.Open("audio/" + entry.Name())
		if err != nil {
			log.Panicf("Error opening file %s: %v", entry.Name(), err)
		}
		sm[i], err = mp3.NewDecoder(playFile)
		if err != nil {
			log.Panicf("Error decoding file %s: %v", entry.Name(), err)
		}
	}
	for {
		var t int
		_, err := fmt.Scan(&t)
		if err != nil {
			log.Panicf("%v", err)
		}
		fmt.Println(t)
		play(sm[t])
	}
}

/*
package main


import (
	"embed"
	"fmt"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed audio/*
var content embed.FS

func main() {
	dir, err := content.ReadDir("audio")
	if err != nil {
		return
	}
	var ftp string
	for i, entry := range dir {
		fmt.Printf("%d: %v\n", i, entry.Name())
		ftp = "audio/" + entry.Name()
		break
	}
	f, err := content.Open(ftp)
	if err != nil {
		log.Fatal(err)
	}
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	println(format.SampleRate.N(time.Second / 60))
	speaker.Init(48000, format.SampleRate.N(time.Second/60))
	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)
	streamer.Close()

	for {
		fmt.Print("Press [ENTER] to play sound! ")
		fmt.Scanln()
		shot := buffer.Streamer(0, buffer.Len())
		speaker.Play(shot)

	}
}
*/
