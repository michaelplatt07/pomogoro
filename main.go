package main

import (
	// "fmt"
	// "log"

	// Internal imports
	"pomogoro/internal/gui"
	"pomogoro/internal/music"
	"pomogoro/internal/pomoapp"
	"pomogoro/internal/pomodoro"

	// Gui imports
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// TODO(map) List of things to correct
// * Move everything to separate modules to make code nicer and better
// * Ship with the settings stored within the application itself? This would make them not persistant though
// * Figure out a nice way to introduce playing music
// * Fill circle based on percentage of time ran
// * Don't allow for going over the total number of iterations
// * Automatically start the timer when switching between focus and relax periods
// * Include text to show if it is focus time or relax time
// * Add ability to save timer and link playlist
// * Toggle text of the button between play and pause
// * Link playlists to the focus and relax timer
// * Save setting to store whether music should pause during the relax timer
// * Refresh library when changed

const (
	settingsFilePath = "/home/michael/Desktop/programming/pomogoro/settings.json"
	// TODO(map) Implement the saved Pomodoros
	// savedPomodoros = "/home/michael/Desktop/programming/pomogoro/saved_pomodoros.json"

	// Sizes
	width  = 800
	height = 600

	// Text
	titleText            = "Pomo-Go-ro"
	descriptionText      = "Welcome to your all-in-one focus partner"
	browseFileText       = "Browse"
	prevButtonText       = "Prev"
	playButtonText       = "Play"
	stopButtonText       = "Stop"
	pauseButtonText      = "Pause"
	nextButtonText       = "Next"
	libraryListLabelText = "Library:"
	detailsLabelText     = "Song Details"
)

var settings = pomoapp.NewSettings(settingsFilePath, "", false, false)

func main() {
	// Load the settings for the application
	settings.Load()

	// Read the library to load the songs into the application
	library := music.Library{}
	library.LoadLibrary(settings.LibraryPath)

	myApp := app.New()
	window := myApp.NewWindow(titleText)
	pomodoroTimer := pomodoro.NewPomodoroTimer()

	// Toolbar
	toolbar := gui.CreateNewToolbar(myApp, pomodoroTimer, settings)

	// About info
	descriptionLabel := widget.NewLabel(descriptionText)
	currentSongPlaying := widget.NewLabel("Currently Playing:")
	descriptionLabelContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(300, 50)),
		descriptionLabel,
	)
	currentSongPlayingContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(400, 50)),
		currentSongPlaying,
	)

	// Song details view
	songDetailsView := gui.NewSongDetailsView(detailsLabelText)

	// Library View
	libraryView := gui.NewLibraryView(libraryListLabelText, &library, songDetailsView)

	// Music controls
	// prevButton := widget.NewButton(prevButtonText, func() {
	//     log.Println("Prev clicked")
	//     // TODO(map) This is totally not safe since it can go out of bounds. Temp measure
	//     library.DecIndex()
	//     libraryView.UpdateSelected()
	//     // TODO(map) This shouldn't start autoplaying if it wasn't playing before hand
	//     // A few conditions to consider:
	//     // If music is playing and button clicked, start playing after updating current song
	//     // If music is paused and button clicked, don't start playing
	//     // If music is stopped and button clicked, don't start playing
	//     // if currentSong.Player != nil { // Close the current player
	//     //     currentSong.Player.Close()
	//     //     currentSong = library.GetCurrentSong()
	//
	//         // Start the song because the previous song was playing
	//         // go currentSong.Play(settings.LibraryPath)
	//     // } else {
	//     //     currentSong = library.GetCurrentSong()
	//     // }
	//
	//     // Set the title to be displayed
	//     // currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currentSong.Title(), currentSong.Artist()))
	//     // currentSongPlaying.Refresh()
	// })
	// playButton := widget.NewButton(playButtonText, func() {
	//     log.Println("Play clicked")
	//     // if currentSong.Player == nil { // There is no Player set because the initial launch of the MP3 hasn't happened
	//     //     log.Println("No song set, playing song...")
	//     //     // libraryList.Select(library.CurrIdx)
	//     //     // libraryList.Refresh()
	//     //     go currentSong.Play(settings.LibraryPath)
	//     //
	//     //     // Set the title to be displayed
	//     //     currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currentSong.Title(), currentSong.Artist()))
	//     //     currentSongPlaying.Refresh()
	//     //
	//     // } else if currentSong.Player.IsPlaying() {
	//     //     log.Println("Pausing song...")
	//     //     currentSong.Player.Pause()
	//     //     currentSong.IsPaused = true
	//     // } else if !currentSong.Player.IsPlaying() && currentSong.IsPaused {
	//     //     log.Println("Resuming song...")
	//     //     currentSong.Player.Play()
	//     //     currentSong.IsPaused = false
	//     // }
	// })
	// stopButton := widget.NewButton(stopButtonText, func() {
	//     log.Println("Stop clicked")
	//     // currentSong.Player.Close()
	//     // currentSong = library.GetCurrentSong()
	//     // currentSongPlaying.SetText("Currently Playing: ")
	// })
	// nextButton := widget.NewButton(nextButtonText, func() {
	//     log.Println("Next clicked")
	//     // TODO(map) This is totally not safe since it can go out of bounds. Temp measure
	//     library.IncIndex()
	//     libraryView.UpdateSelected()
	//     // TODO(map) This shouldn't start autoplaying if it wasn't playing before hand
	//     // A few conditions to consider:
	//     // If music is playing and button clicked, start playing after updating current song
	//     // If music is paused and button clicked, don't start playing
	//     // If music is stopped and button clicked, don't start playing
	//     // if currentSong.Player != nil { // Close the current player
	//     //     currentSong.Player.Close()
	//     //     currentSong = library.GetCurrentSong()
	//     //
	//     //     // Start the song because the previous song was playing
	//     //     go currentSong.Play(settings.LibraryPath)
	//     // } else {
	//     //     currentSong = library.GetCurrentSong()
	//     // }
	//
	//
	//     // Set the title to be displayed
	//     // currentSongPlaying.SetText(fmt.Sprintf("Currently Playing: %s by %s", currentSong.Title(), currentSong.Artist()))
	//     // currentSongPlaying.Refresh()
	// })

	// Containers around the buttons to ensure their size doesn't grow beyond what is desired
	// prevButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), prevButton)
	// playButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), playButton)
	// stopButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), stopButton)
	// nextButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), nextButton)

	// Rows
	// Info
	descriptionRow := container.New(
		layout.NewHBoxLayout(),
		descriptionLabelContainer,
		currentSongPlayingContainer,
	)
	// Control
	// controlsRow := container.New(layout.NewHBoxLayout(), prevButtonContainer, playButtonContainer, stopButtonContainer, nextButtonContainer)
	// controlsRowParent := container.New(layout.NewCenterLayout(), controlsRow)
	controls := gui.NewMusicControls(&library, libraryView, settings.LibraryPath)

	// Parent container
	content := container.New(
		layout.NewVBoxLayout(),
		toolbar,
		pomodoroTimer.PomodoroTimerCanvas.TopLevelContainer,
		descriptionRow,
		libraryView.Container,
		controls.Container,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(width, height))
	window.ShowAndRun()
}
