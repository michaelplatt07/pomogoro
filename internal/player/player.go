package player

import (
	"fmt"
	"time"

	// Internal imports
	"pomogoro/internal/library"
	"pomogoro/internal/messages"
	"pomogoro/internal/pomoapp"
)

type Player struct {
	SongControlChan chan messages.ChannelMessage
	IsPlaying       bool
	IsPaused        bool
}

func (player *Player) Play(library *library.Library, settings *pomoapp.Settings) {
	fmt.Println("Initializing channel...")
	player.SongControlChan = make(chan messages.ChannelMessage)
	go library.CurrentSong.Play(player.SongControlChan)
	player.IsPlaying = true
	for {
		time.Sleep(time.Second)
		message, ok := <-player.SongControlChan
		if !ok {
			return
		} else {
			fmt.Println("Got message = ", message)
			if message.SongFinished == true {
				fmt.Println("Song finished message received")
				if library.HasNextSong && settings.AutoPlay && !settings.Shuffle {
					library.IncIndex()
					fmt.Println("Starting next song...")
					go library.CurrentSong.Play(player.SongControlChan)
					fmt.Println("Started next song...")
				} else if library.HasNextSong && settings.AutoPlay && settings.Shuffle {
					library.NextShuffle()
					fmt.Println("Starting next song...")
					go library.CurrentSong.Play(player.SongControlChan)
					fmt.Println("Started next song...")
				} else {
					fmt.Println("Stopping player...")
					library.CurrentSong.Stop(false)
					break
				}
			} else if message.SongStopped == true {
				fmt.Println("Song stopped message received")
				break
			} else if message.SongPaused == true {
				fmt.Println("Song paused message received")
				player.IsPlaying = false
				player.IsPaused = true
			} else if message.SongResumed == true {
				fmt.Println("Song resumed message received")
				player.IsPlaying = true
				player.IsPaused = false
			} else if message.SongSkipped == true {
				fmt.Println("Song skipped message received")
				// Start playing the next song if the stage is not paused
				if player.IsPlaying {
					fmt.Println("Playing next song")
					go library.CurrentSong.Play(player.SongControlChan)
				}
			}
		}
	}
	fmt.Println("Closing channel...")
	close(player.SongControlChan)
	fmt.Println("Closed...")
}
