package library

import (
	"log"
	"os"
	"time"

	// Internal imports
	"pomogoro/internal/song"
)

type Library struct {
	LibraryPath        string
	Songs              []*song.Song
	CurrentSong        *song.Song
	CurrIdx            int
	PlayingCurrentSong bool
	PlayNextSongChan   chan bool
	PlayNextSong       bool
	HasNextSong        bool
}

func (library *Library) LoadLibrary(pathToLibrary string) {
	library.PlayNextSongChan = make(chan bool)
	library.LibraryPath = pathToLibrary
	songs, err := os.ReadDir(pathToLibrary)
	if err != nil {
		panic(err)
	}

	library.CurrIdx = 0

	for _, singleSong := range songs {
		log.Printf("Adding song %s to queue", singleSong.Name())
		song := song.NewSong(pathToLibrary, singleSong.Name())
		song.ApplyTag(library.LibraryPath, singleSong.Name())
		library.Songs = append(
			library.Songs,
			song,
		)
	}

	log.Print("Finished loading library...")

	// Set the current song as the first one and set the flag that there is a next song
	library.CurrentSong = library.Songs[0]
	library.HasNextSong = true
}

func (library *Library) PlayLibrary() {
	library.PlayingCurrentSong = true
	go library.CurrentSong.Play(library.PlayNextSongChan)
	for library.HasNextSong {
		time.Sleep(time.Second)
		select {
		case <-library.PlayNextSongChan:
			library.CurrIdx = library.CurrIdx + 1
			if library.CurrIdx >= len(library.Songs) {
				library.HasNextSong = false
			} else {
				library.CurrentSong = library.Songs[library.CurrIdx]
				go library.CurrentSong.Play(library.PlayNextSongChan)
			}
		}
	}
}

func (library *Library) DecIndex() {
	// TODO(map) Implement me
}

func (library *Library) IncIndex() {
	// TODO(map) Implement me
}
