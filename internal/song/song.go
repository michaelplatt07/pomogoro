package song

import (
	"fmt"
	"log"
	"os"
	"time"

	// Internal imports
	"pomogoro/internal/messages"

	// MP3 imports
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"

	// ID3
	"github.com/bogem/id3v2"
)

type Song struct {
	Name     string
	FilePath string
	Skipped  bool
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
	}
}

func (song *Song) Play(songControlsChan chan messages.ChannelMessage) {
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
		if song.Skipped {
			// Reset skipped flag and published the skipped message
			song.Skipped = false
			songControlsChan <- messages.ChannelMessage{SongSkipped: true}
			break
		} else if song.Player == nil {
			// Song was stopped so we published stopped message
			fmt.Println("Publishing stopped message")
			songControlsChan <- messages.ChannelMessage{SongStopped: true}
			fmt.Println("Published stopped message")
			break
		} else if !song.Player.IsPlaying() {
			// The song was ended by playing to completion and we should publish the message and break
			if song.Player.UnplayedBufferSize() == 0 {
				fmt.Println("Publishing completed message")
				songControlsChan <- messages.ChannelMessage{SongFinished: true}
				fmt.Println("Published completed message")
				break
			}
		}
	}
}

func (song *Song) Resume(songControlsChan chan messages.ChannelMessage) {
	song.Player.Play()
	fmt.Println("Publishing resume message...")
	songControlsChan <- messages.ChannelMessage{SongResumed: true}
	fmt.Println("Published resume message...")
}

func (song *Song) Pause(songControlsChan chan messages.ChannelMessage) {
	song.Player.Pause()
	fmt.Println("Publishing pause message...")
	songControlsChan <- messages.ChannelMessage{SongPaused: true}
	fmt.Println("Published pause message...")
}

func (song *Song) Stop(skipped bool) {
	if song.Player != nil {
		song.Player.Close()
	}
	song.Player = nil
	song.Skipped = skipped
}
