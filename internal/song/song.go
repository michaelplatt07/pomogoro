package song

import (
	"log"
	"os"
	"time"

	// MP3 imports
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"

	// ID3
	"github.com/bogem/id3v2"
)

type Song struct {
	Name     string
	FilePath string
	IsPaused bool
	Tag      *id3v2.Tag

	// Used just for accessing the player for functionality. The file has to be opened and be streamed so having an
	// instance of the player actually being initialized doesn't work too well unless I wanted to keep the bytes of the
	// file in memory.
	Player oto.Player
}

func (song *Song) ApplyTag(libraryPath string, songName string) {
	tag, _ := id3v2.Open(libraryPath+"/"+songName, id3v2.Options{Parse: true})
	song.Tag = tag
}

func NewSong(libraryPath string, songName string) *Song {
	return &Song{
		Name:     songName,
		FilePath: libraryPath + "/" + songName,
		IsPaused: false,
	}
}

func (song *Song) Play(playNextSongChan chan bool) {
	// Open the file to stream the contents
	f, err := os.Open(song.FilePath)
	if err != nil {
		log.Println("Err opening file")
		panic(err)
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		log.Println("Err setting up decoder")
		panic(err)
	}

	c, ready, err := oto.NewContext(d.SampleRate(), 2, 2)
	if err != nil {
		panic(err)
	}
	<-ready

	song.Player = c.NewPlayer(d)
	song.Player.SetVolume(0.1)
	song.Player.Play()

	// Wait for the song to finish playing
	for {
		time.Sleep(time.Second)
		if !song.Player.IsPlaying() && !song.IsPaused {
			break
		}
	}

	// Write to channel that song ended
	playNextSongChan <- true
}

func (song *Song) Resume() {
	song.Player.Play()
	song.IsPaused = false
}

func (song *Song) Pause() {
	song.Player.Pause()
	song.IsPaused = true
}

func (song *Song) Stop() {
	// TODO(map) Implement me
}
