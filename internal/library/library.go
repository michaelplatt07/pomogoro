package library

import (
	"log"
	"os"

	// Internal imports
	"pomogoro/internal/song"
)

type Library struct {
	LibraryPath        string
	Songs              []*song.Song
	CurrentSong        *song.Song
	CurrIdx            int
	PlayingCurrentSong bool
	PlayNextSong       bool
	HasNextSong        bool
}

func (library *Library) LoadLibrary(pathToLibrary string) {
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

func (library *Library) DecIndex() {
	if library.CurrIdx-1 > 0 {
		library.CurrIdx = library.CurrIdx - 1
		library.CurrentSong = library.Songs[library.CurrIdx]
		library.HasNextSong = true
	}
}

func (library *Library) IncIndex() {
	library.CurrIdx = library.CurrIdx + 1
	library.CurrentSong = library.Songs[library.CurrIdx]
	if library.CurrIdx >= len(library.Songs)-1 {
		library.HasNextSong = false
	} else {
		library.HasNextSong = true
	}
}
