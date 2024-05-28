package music

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

// Represents a song and its current state and player
type Song struct {
	Name     string
	IsPaused bool
	Player   oto.Player
	Tag      *id3v2.Tag
}

func (song *Song) Title() string {
	return song.Tag.Title()
}

func (song *Song) Artist() string {
	return song.Tag.Artist()
}

func (song *Song) Album() string {
	return song.Tag.Album()
}

func (song *Song) Genre() string {
	return song.Tag.Genre()
}

func (song *Song) SaveDetails(title string, artist string, album string, genre string) {
	song.Tag.SetTitle(title)
	song.Tag.SetArtist(artist)
	song.Tag.SetAlbum(album)
	song.Tag.SetGenre(genre)

	if err := song.Tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}

// TODO(map) Figure out wtf to do with this libraryPath param
func (song *Song) Play(libraryPath string) {
	// Open the file that is associated with the currently selected song in the queue
	f, err := os.Open(libraryPath + "/" + song.Name)
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

	p := c.NewPlayer(d)
	defer p.Close()
	log.Println("Playing song")
	p.Play()

	// Assign the player so the controls works
	song.Player = p

	// TODO(map) Apply volume adjust
	p.SetVolume(0.2)

	for {
		time.Sleep(time.Second)
		if !p.IsPlaying() && !song.IsPaused {
			break
		}
	}
}

func (song *Song) Pause() {
	song.Player.Pause()
}

func (song *Song) Resume() {
	song.Player.Play()
}

func (song *Song) Stop() {
	song.Player.Close()
	song.Player = nil
}

// Library that will load up all the songs available to be played. Differs from Queue in that a Queue is a user defined
// list of songs to play but Library is everything
type Library struct {
	Songs       []Song
	CurrIdx     int
	CurrentSong Song
}

func (library *Library) LoadLibrary(pathToLibrary string) {
	songs, err := os.ReadDir(pathToLibrary)
	if err != nil {
		panic(err)
	}

	library.Songs = []Song{}
	library.CurrIdx = 0

	for _, song := range songs {
		// TODO(map) song.Name() is actually just returning the file name. It makese sense because of how my stuff is named but really we should be reading the ID3 tags here.
		log.Printf("Adding song %s to queue", song.Name())
		// TODO(map) Error handling
		details, _ := id3v2.Open(pathToLibrary+"/"+song.Name(), id3v2.Options{Parse: true})
		library.Songs = append(
			library.Songs,
			Song{Name: song.Name(), Tag: details, IsPaused: false},
		)
	}

	log.Print("Finished loading library...")

	// Set the current song as the first on
	library.CurrentSong = library.Songs[library.CurrIdx]
}

func (library *Library) GetCurrentSong() Song {
	return library.Songs[library.CurrIdx]
}

func (library *Library) IncIndex() {
	library.CurrIdx += 1
	library.CurrentSong = library.Songs[library.CurrIdx]
}

func (library *Library) DecIndex() {
	library.CurrIdx -= 1
	library.CurrentSong = library.Songs[library.CurrIdx]
}
