package main

import (
	// Internal imports
	"pomogoro/internal/gui"
	"pomogoro/internal/library"
	"pomogoro/internal/player"
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
// * Ship with the settings stored within the application itself? This would make them not persistant though
// * Figure out a nice way to introduce playing music
// * Fill circle based on percentage of time ran
// * Don't allow for going over the total number of iterations
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

var settings = pomoapp.NewSettings(settingsFilePath, "", false, false, false)

func main() {
	// Load the settings for the application
	settings.Load()

	// Load library
	library := library.Library{}
	library.LoadLibrary("/home/michael/Desktop/programming/pomogoro/library", settings)

	// Load the player
	player := player.Player{IsPlaying: false, IsPaused: false}

	myApp := app.New()
	window := myApp.NewWindow(titleText)
	pomodoroTimer := pomodoro.NewPomodoroTimer(&library, settings)

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
	libraryView := gui.NewLibraryView(libraryListLabelText, &library, songDetailsView, settings, &player)

	// Info
	descriptionRow := container.New(
		layout.NewHBoxLayout(),
		descriptionLabelContainer,
		currentSongPlayingContainer,
	)
	// Control
	controls := gui.NewMusicControls(&library, &player, settings, pomodoroTimer)

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
