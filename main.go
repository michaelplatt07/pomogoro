package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
    "strconv"
    
    // Internal imports
    "pomogoro/internal/pomodoro"

	// Gui imports
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	// MP3 imports
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto/v2"

	// ID3
	"github.com/bogem/id3v2"
)
// TODO(map) List of things to correct
// * Move everything to separate modules to make code nicer and better
// * Fill circle based on percentage of time ran
// * Don't allow for going over the total number of iterations
// * Automatically start the timer when switching between focus and relax periods
// * Include text to show if it is focus time or relax time
// * Add ability to save timer and link playlist
// * Toggle text of the button between play and pause
// * Link playlists to the focus and relax timer
// * Save setting to store whether music should pause during the relax timer

// Holds the settings for the program
type Settings struct {
    LibraryPath string
    AutoPlay bool
    Shuffle bool
}

// Represents a song and its current state and player
type Song struct {
    Name string
    IsPaused bool
    Player oto.Player
    Tag *id3v2.Tag
}

// Queue for a list of songs and an index to track
type Queue struct {
    Songs []Song
    CurrIdx int
}

// Struct for displaying and altering values of the Song detials
type SongDetails struct {
    Title string
    Artist string
    Album string
    Genre string
}

const (
    settingsFilePath = "/home/michael/Desktop/programming/pomogoro/settings.json"
    // TODO(map) Implement the saved Pomodoros
    // savedPomodoros = "/home/michael/Desktop/programming/pomogoro/saved_pomodoros.json"

    // Sizes
    width = 800
    height = 600

    // Text
    titleText = "Pomo-Go-ro"
    descriptionText = "Welcome to your all-in-one focus partner"
    browseFileText = "Browse"
    prevButtonText = "Prev"
    playButtonText = "Play"
    stopButtonText = "Stop"
    pauseButtonText = "Pause"
    nextButtonText = "Next"
    libraryListLabelText = "Library:"
    detailsLabelText = "Song Details"
)

var (
    settings Settings     

    // Delcare a Song instance that can be referenced through out the program as needed
    currSong Song
    currSongDetails SongDetails

    // Declare an empty queue
    songQueue = Queue{CurrIdx: 0}
)

func main() {
    // Load the settings for the application
    loadSettings(settingsFilePath)

    // Read the library to load the songs into the application
    readLibrary()

	myApp := app.New()
	window := myApp.NewWindow(titleText)
    pomodoroTimer := createPomodoroTimer()

    // Toolbar
    toolbar := widget.NewToolbar(
        // TODO(map) What's a good icon to use here? Maybe explore the idea of making my own resource
        widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			createPomodoroWindow(myApp, pomodoroTimer)
		}),
        widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			createSettingsWindow(myApp)
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

    // About info
    descriptionLabel := widget.NewLabel(descriptionText)
    currentSongPlaying := widget.NewLabel("Currently Playing:")
    descriptionLabelContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 50)), descriptionLabel)
    currentSongPlayingContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(400, 50)), currentSongPlaying)

    // Library View
    libraryListLabel := widget.NewLabel(libraryListLabelText)
    libraryList := widget.NewList(
		func() int {
			return len(songQueue.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(currSong.Name)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(songQueue.Songs[i].Name)
		})
    libraryListLabelContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), libraryListLabel)
    libraryListContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 400)), libraryList)

    // Song details view
    detailsLabel := widget.NewLabel(detailsLabelText)
    // TODO(map) Include labels for these guys
    titleInput := widget.NewEntry()
    artistInput := widget.NewEntry()
    albumInput := widget.NewEntry()
    genreInput := widget.NewEntry()
	titleInput.SetPlaceHolder(currSongDetails.Title)
	artistInput.SetPlaceHolder(currSongDetails.Artist)
	albumInput.SetPlaceHolder(currSongDetails.Album)
	genreInput.SetPlaceHolder(currSongDetails.Genre)
    saveId3DataButton := widget.NewButton("Save", func() {
		log.Println("Title was:", titleInput.Text)
		log.Println("Artist was:", artistInput.Text)
		log.Println("Album was:", albumInput.Text)
		log.Println("Genre was:", genreInput.Text)
        currSong.Tag.SetTitle(titleInput.Text)
        currSong.Tag.SetArtist(artistInput.Text)
        currSong.Tag.SetAlbum(albumInput.Text)
        currSong.Tag.SetGenre(genreInput.Text)
        if err := currSong.Tag.Save(); err != nil {
		    log.Fatal("Error while saving a tag: ", err)
	    }
	})
    songDetailsContainer := container.New(layout.NewVBoxLayout(), detailsLabel, titleInput, artistInput, albumInput, genreInput, saveId3DataButton)
    songDetailsParentContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 400)), songDetailsContainer)


    // Music controls
    prevButton := widget.NewButton(prevButtonText, func() {
        log.Println("Prev clicked")
        // TODO(map) This is totally not safe since it can go out of bounds. Temp measure
        songQueue.CurrIdx = songQueue.CurrIdx - 1
        libraryList.Select(songQueue.CurrIdx)
        libraryList.Refresh()
        if currSong != (Song{}) {
            currSong.Player.Close()
            currSong = Song{}
        }
        go playSong()
        
        // Update the struct for the song details so it's ready to be referenced
        updateCurrentSongDetails()
        // Set the title to be displayed
        currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currSongDetails.Title, currSongDetails.Artist))
        currentSongPlaying.Refresh()
        // Set the details and refresh
        titleInput.SetText(currSongDetails.Title)
        artistInput.SetText(currSongDetails.Artist)
        albumInput.SetText(currSongDetails.Album)
        genreInput.SetText(currSongDetails.Genre)
        titleInput.Refresh()
        artistInput.Refresh()
        albumInput.Refresh()
        genreInput.Refresh()
    })
    playButton := widget.NewButton(playButtonText, func() {
        log.Println("Play clicked")
        if currSong == (Song{}) {
            log.Println("No song set, playing song...")
            libraryList.Select(songQueue.CurrIdx)
            libraryList.Refresh()
            go playSong()

            // Update the struct for the song details so it's ready to be referenced
            updateCurrentSongDetails()
            // Set the title to be displayed
            currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currSongDetails.Title, currSongDetails.Artist))
            currentSongPlaying.Refresh()
            // Set the details and refresh
            titleInput.SetText(currSongDetails.Title)
            artistInput.SetText(currSongDetails.Artist)
            albumInput.SetText(currSongDetails.Album)
            genreInput.SetText(currSongDetails.Genre)
            titleInput.Refresh()
            artistInput.Refresh()
            albumInput.Refresh()
            genreInput.Refresh()
        } else if currSong.Player.IsPlaying() {
            log.Println("Pausing song...")
            currSong.Player.Pause()
            currSong.IsPaused = true
        } else if !currSong.Player.IsPlaying() && currSong.IsPaused {
            log.Println("Resuming song...")
            currSong.Player.Play()
            currSong.IsPaused = false
        }
        
    })
    stopButton := widget.NewButton(stopButtonText, func() {
        log.Println("Stop clicked")
        currSong.Player.Close()
        currSong = Song{}
        currentSongPlaying.SetText("Currently Playing: ")
    })
    nextButton := widget.NewButton(nextButtonText, func() {
        log.Println("Next clicked")
        // TODO(map) This is totally not safe since it can go out of bounds. Temp measure
        songQueue.CurrIdx = songQueue.CurrIdx + 1
        libraryList.Select(songQueue.CurrIdx)
        libraryList.Refresh()
        if currSong != (Song{}) {
            currSong.Player.Close()
            currSong = Song{}
        }
        go playSong()

        // Update the struct for the song details so it's ready to be referenced
        updateCurrentSongDetails()
        // Set the title to be displayed
        currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currSongDetails.Title, currSongDetails.Artist))
        currentSongPlaying.Refresh()
        // Set the details and refresh
        titleInput.SetText(currSongDetails.Title)
        artistInput.SetText(currSongDetails.Artist)
        albumInput.SetText(currSongDetails.Album)
        genreInput.SetText(currSongDetails.Genre)
        titleInput.Refresh()
        artistInput.Refresh()
        albumInput.Refresh()
        genreInput.Refresh()
    })

    // Containers around the buttons to ensure their size doesn't grow beyond what is desired
    prevButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), prevButton)
    playButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), playButton)
    stopButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), stopButton)
    nextButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), nextButton)

    // Rows
    // Info
    descriptionRow := container.New(layout.NewHBoxLayout(), descriptionLabelContainer, currentSongPlayingContainer)
    // Library
    libraryRow := container.New(layout.NewHBoxLayout(), libraryListLabelContainer, libraryListContainer, songDetailsParentContainer)
    // Control
    controlsRow := container.New(layout.NewHBoxLayout(), prevButtonContainer, playButtonContainer, stopButtonContainer, nextButtonContainer)
    controlsRowParent := container.New(layout.NewCenterLayout(), controlsRow)

    // Parent container
    content := container.New(layout.NewVBoxLayout(), toolbar, pomodoroTimer.PomodoroTimerCanvas.TopLevelContainer, descriptionRow, libraryRow, controlsRowParent)

	window.SetContent(content)
	window.Resize(fyne.NewSize(width, height))
	window.ShowAndRun()

}

func createPomodoroWindow(app fyne.App, pt *pomodoro.PomodoroTimer) {
    newPomodoroWindow := app.NewWindow("New Pomodoro")

    pomodoroNameLabel := widget.NewLabel("Pomodoro name: ")
    pomodoroNameLabelContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), pomodoroNameLabel)
    pomodoroNameText := widget.NewEntry()
    pomodoroNameTextContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), pomodoroNameText)
    
    focusTimeLabel := widget.NewLabel("Enter focus time in minutes: ")
    focusTimeLabelContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), focusTimeLabel)
    focusTimeText := widget.NewEntry()
    focusTimeTextContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), focusTimeText)
    
    relaxTimeLabel := widget.NewLabel("Enter relax time in minutes: ")
    relaxTimeLabelContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), relaxTimeLabel)
    relaxTimeText := widget.NewEntry()
    relaxTimeTextContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), relaxTimeText)
    
    iterationTimeLabel := widget.NewLabel("Enter the number of iterations to complete: ")
    iterationTimeLabelContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), iterationTimeLabel)
    iterationTimeText := widget.NewEntry()
    iterationTimeTextContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), iterationTimeText)

    createTimerButton := widget.NewButton("Create new Timer", func() {
        fmt.Println("Starting new Pomodoro Timer")
        // TODO(map) Error handling here
        focusTime, _ := strconv.Atoi(focusTimeText.Text)
        relaxTime, _ := strconv.Atoi(relaxTimeText.Text)
        iterationTime, _ := strconv.Atoi(iterationTimeText.Text)
        pt.SetSettings(focusTime, relaxTime, iterationTime)
        pt.RestartTimer()

    })
    timerButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(200, 40)), createTimerButton)
    textContainer := container.NewVBox(pomodoroNameLabelContainer, focusTimeLabelContainer, relaxTimeLabelContainer, iterationTimeLabelContainer)
    inputContainer := container.NewVBox(pomodoroNameTextContainer, focusTimeTextContainer, relaxTimeTextContainer, iterationTimeTextContainer)
    pomodoroInfoContainer := container.New(layout.NewHBoxLayout(), textContainer, inputContainer)
    content := container.New(layout.NewVBoxLayout(), pomodoroInfoContainer, timerButtonContainer)

    newPomodoroWindow.SetContent(content)
	newPomodoroWindow.Resize(fyne.NewSize(400, 400))
	newPomodoroWindow.Show()
}

func createSettingsWindow(app fyne.App) {
    settingsWindow := app.NewWindow("Setttings")

    // Widget creation
    libraryPathLabel := widget.NewLabel("Library path: ")
    libraryPath := widget.NewEntry()
    libraryPath.SetText(settings.LibraryPath)
    autoPlayCheckBox := widget.NewCheck("Autoplay next song", func (checked bool) {
        settings.AutoPlay = checked
    })
    autoPlayCheckBox.Checked = settings.AutoPlay
    shuffleCheckBox := widget.NewCheck("Shuffle", func (checked bool) {
        settings.Shuffle = checked
    })
    shuffleCheckBox.Checked = settings.Shuffle
    saveButton := widget.NewButton("Save", func() {
        settings.LibraryPath = libraryPath.Text
        settings.AutoPlay = autoPlayCheckBox.Checked
        settings.Shuffle = shuffleCheckBox.Checked
        dialog.ShowConfirm("Confirm", "Are you sure you want to save these settings?", func(confirm bool) {
            saveSettings(settingsFilePath, settings)
            settingsWindow.Close()
        }, settingsWindow)
	})

    // Create the rows
    libraySettingsRow := container.New(layout.NewVBoxLayout(), libraryPathLabel, libraryPath)
    playSettingsRow := container.New(layout.NewHBoxLayout(), autoPlayCheckBox, shuffleCheckBox)
    saveRow := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), saveButton)

    content := container.New(layout.NewVBoxLayout(), libraySettingsRow, playSettingsRow, saveRow)
    settingsWindow.SetContent(content)
	settingsWindow.Resize(fyne.NewSize(400, 100))
	settingsWindow.Show()
}

func loadSettings(pathToSettings string) {
    settingsFile, _ := os.ReadFile(pathToSettings)
    err := json.Unmarshal(settingsFile, &settings)
    if err != nil {
        log.Print("Failure in unmarshelling the settings data")
    }
}

func saveSettings(pathToSettings string, settings Settings) {
    // TODO(map) Handle errors gracefully
    file, _ := json.MarshalIndent(settings, "", "    ")
	_ = os.WriteFile(pathToSettings, file, 0644)
}

func readLibrary() {
    songs, err := os.ReadDir(settings.LibraryPath)
    if err != nil {
        panic(err)
    }
    
   songQueue = Queue{
        Songs: []Song{},
        CurrIdx: 0,
    }

    for _, song := range songs {
        log.Printf("Adding song %s to queue", song.Name())
        songQueue.Songs = append(songQueue.Songs, Song{Name: song.Name(), IsPaused: false})
    }
    
}

func createPomodoroTimer() *pomodoro.PomodoroTimer {
    pt := &pomodoro.PomodoroTimer{
        IsRunning: false,
        InBreakMode: false,
        
        // Settings are initialized later in the SetSettings method. This is the standard for a blank timer
    }

    pt.CreateDefaultCanvas()

    return pt
}

func updateCurrentSongDetails() {
    // Update the currSong variable so the program can reference it when needed for controls
    currSong = songQueue.Songs[songQueue.CurrIdx]

    // Update the current song tag details
    tag, err := id3v2.Open(settings.LibraryPath + "/" + songQueue.Songs[songQueue.CurrIdx].Name, id3v2.Options{Parse: true})
        if err != nil {
            log.Println("Error reading ID3 tag")
            panic(err)
        }
    currSong.Tag = tag

    currSongDetails = SongDetails{
        Title: tag.Title(),
        Artist: tag.Artist(),
        Album: tag.Album(),
        Genre: tag.Genre(),
    }
}

func playSong() {
    // Open the file that is associated with the currently selected song in the queue
    f, err := os.Open(settings.LibraryPath + "/" + songQueue.Songs[songQueue.CurrIdx].Name)

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
    currSong.Player = p

    // TODO(map) Apply volume adjust
    p.SetVolume(0.2)

    for {
		time.Sleep(time.Second)
		if !p.IsPlaying() && !currSong.IsPaused {
			break
		}
	}

}

