package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
    "strconv"
    
    // Internal imports
    "pomogoro/internal/pomodoro"
    "pomogoro/internal/music"

	// Gui imports
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// TODO(map) List of things to correct
// * Move everything to separate modules to make code nicer and better
// * Figure out a nice way to introduce playing music
// *Fill circle based on percentage of time ran
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
    // currSong music.Song
    // currSongDetails music.SongDetails

    // Declare an empty queue
    // songQueue = music.Queue{CurrIdx: 0}
)

func main() {
    // Load the settings for the application
    loadSettings(settingsFilePath)

    // Read the library to load the songs into the application
    library := music.Library{}
    library.LoadLibrary(settings.LibraryPath)

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
			return len(library.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(library.GetCurrentSong().Name)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(library.Songs[i].Name)
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
    // TODO(map) Move this out somewhere that makes sense
    currentSong := library.GetCurrentSong()
	titleInput.SetText(currentSong.Title())
    artistInput.SetText(currentSong.Artist())
    albumInput.SetText(currentSong.Album())
    genreInput.SetText(currentSong.Genre())
    titleInput.Refresh()
    artistInput.Refresh()
    albumInput.Refresh()
    genreInput.Refresh()

    saveId3DataButton := widget.NewButton("Save", func() {
        currentSong.SaveDetails(
            titleInput.Text,
            artistInput.Text,
            albumInput.Text,
            genreInput.Text,
        )
	})
    songDetailsContainer := container.New(layout.NewVBoxLayout(), detailsLabel, titleInput, artistInput, albumInput, genreInput, saveId3DataButton)
    songDetailsParentContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 400)), songDetailsContainer)

    // TODO(map) Make this its own component to encapsulate
    libraryList.Select(library.CurrIdx)
    libraryList.Refresh()

    // Music controls
    prevButton := widget.NewButton(prevButtonText, func() {
        log.Println("Prev clicked")
        // TODO(map) This is totally not safe since it can go out of bounds. Temp measure
        library.CurrIdx = library.CurrIdx - 1
        libraryList.Select(library.CurrIdx)
        libraryList.Refresh()
        // TODO(map) This shouldn't start autoplaying if it wasn't playing before hand
        // A few conditions to consider:
        // If music is playing and button clicked, start playing after updating current song
        // If music is paused and button clicked, don't start playing
        // If music is stopped and button clicked, don't start playing
        if currentSong.Player != nil { // Close the current player
            currentSong.Player.Close()
            currentSong = library.GetCurrentSong()

            // Start the song because the previous song was playing
            go currentSong.Play(settings.LibraryPath)
        } else {
            currentSong = library.GetCurrentSong()
        }

        // Set the title to be displayed
        currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currentSong.Title(), currentSong.Artist()))
        currentSongPlaying.Refresh()

        // Set the details and refresh
        titleInput.SetText(currentSong.Title())
        artistInput.SetText(currentSong.Artist())
        albumInput.SetText(currentSong.Album())
        genreInput.SetText(currentSong.Genre())
        titleInput.Refresh()
        artistInput.Refresh()
        albumInput.Refresh()
        genreInput.Refresh()
    })
    playButton := widget.NewButton(playButtonText, func() {
        log.Println("Play clicked")
        if currentSong.Player == nil { // There is no Player set because the initial launch of the MP3 hasn't happened
            log.Println("No song set, playing song...")
            libraryList.Select(library.CurrIdx)
            libraryList.Refresh()
            go currentSong.Play(settings.LibraryPath)

            // Set the title to be displayed
            currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currentSong.Title(), currentSong.Artist()))
            currentSongPlaying.Refresh()

            // Set the details and refresh
            titleInput.SetText(currentSong.Title())
            artistInput.SetText(currentSong.Artist())
            albumInput.SetText(currentSong.Album())
            genreInput.SetText(currentSong.Genre())
            titleInput.Refresh()
            artistInput.Refresh()
            albumInput.Refresh()
            genreInput.Refresh()
        } else if currentSong.Player.IsPlaying() {
            log.Println("Pausing song...")
            currentSong.Player.Pause()
            currentSong.IsPaused = true
        } else if !currentSong.Player.IsPlaying() && currentSong.IsPaused {
            log.Println("Resuming song...")
            currentSong.Player.Play()
            currentSong.IsPaused = false
        }
    })
    stopButton := widget.NewButton(stopButtonText, func() {
        log.Println("Stop clicked")
        currentSong.Player.Close()
        currentSong = library.GetCurrentSong()
        currentSongPlaying.SetText("Currently Playing: ")
    })
    nextButton := widget.NewButton(nextButtonText, func() {
        log.Println("Next clicked")
        // TODO(map) This is totally not safe since it can go out of bounds. Temp measure
        library.CurrIdx = library.CurrIdx + 1
        libraryList.Select(library.CurrIdx)
        libraryList.Refresh()
        // TODO(map) This shouldn't start autoplaying if it wasn't playing before hand
        // A few conditions to consider:
        // If music is playing and button clicked, start playing after updating current song
        // If music is paused and button clicked, don't start playing
        // If music is stopped and button clicked, don't start playing
        if currentSong.Player != nil { // Close the current player
            currentSong.Player.Close()
            currentSong = library.GetCurrentSong()

            // Start the song because the previous song was playing
            go currentSong.Play(settings.LibraryPath)
        } else {
            currentSong = library.GetCurrentSong()
        }

        
        // Set the title to be displayed
        currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currentSong.Title(), currentSong.Artist()))
        currentSongPlaying.Refresh()
        // Set the details and refresh
        titleInput.SetText(currentSong.Title())
        artistInput.SetText(currentSong.Artist())
        albumInput.SetText(currentSong.Album())
        genreInput.SetText(currentSong.Genre())
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
    // content := container.New(layout.NewVBoxLayout(), toolbar, pomodoroTimer.PomodoroTimerCanvas.TopLevelContainer, descriptionRow, libraryRow)

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

func createPomodoroTimer() *pomodoro.PomodoroTimer {
    pt := &pomodoro.PomodoroTimer{
        IsRunning: false,
        InBreakMode: false,
        
        // Settings are initialized later in the SetSettings method. This is the standard for a blank timer
    }

    pt.CreateDefaultCanvas()

    return pt
}
