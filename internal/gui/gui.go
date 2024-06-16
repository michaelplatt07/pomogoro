package gui

import (
	"fmt"
	"log"
	"strconv"

	// Internal imports
	"pomogoro/internal/library"
	"pomogoro/internal/music"
	"pomogoro/internal/player"
	"pomogoro/internal/pomoapp"
	"pomogoro/internal/pomodoro"

	// Gui imports
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Gui struct {
	Toolbar *widget.Toolbar
}

func NewGui(
	app fyne.App,
	pomodoroTimer *pomodoro.PomodoroTimer,
	appSettings *pomoapp.Settings,
	library *music.Library,
) *Gui {
	toolbar := CreateNewToolbar(app, pomodoroTimer, appSettings)
	return &Gui{
		Toolbar: toolbar,
	}
}

type PomodoroCreationWindow struct {
	Window             fyne.Window
	Container          *fyne.Container
	FocusTimeInput     *widget.Entry
	RelaxTimeInput     *widget.Entry
	IterationTimeInput *widget.Entry
}

func NewPomodoroCreationWindow(app fyne.App, p *pomodoro.PomodoroTimer) *PomodoroCreationWindow {
	// This will initialize and build the window and provide links to the fields that would be used to retrieve input
	// or modify values elsewhere
	pomodoroWindow := app.NewWindow("New Pomodoro")
	focusTimeLabel := widget.NewLabel("Enter focus time in minutes: ")
	focusTimeLabelContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(200, 40)),
		focusTimeLabel,
	)
	focusTimeText := widget.NewEntry()
	focusTimeTextContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(200, 40)),
		focusTimeText,
	)

	relaxTimeLabel := widget.NewLabel("Enter relax time in minutes: ")
	relaxTimeLabelContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(200, 40)),
		relaxTimeLabel,
	)
	relaxTimeText := widget.NewEntry()
	relaxTimeTextContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(200, 40)),
		relaxTimeText,
	)

	iterationTimeLabel := widget.NewLabel("Enter the number of iterations to complete: ")
	iterationTimeLabelContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(200, 40)),
		iterationTimeLabel,
	)
	iterationTimeText := widget.NewEntry()
	iterationTimeTextContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(200, 40)),
		iterationTimeText,
	)

	createTimerButton := widget.NewButton("Create Timer", func() {
		dialog.ShowConfirm(
			"Confirm",
			"Do you want to create this timer? Current timer will be overridden",
			func(confirm bool) {
				fmt.Println("Starting new Pomodoro Timer")
				// TODO(map) Error handling here
				focusTime, _ := strconv.Atoi(focusTimeText.Text)
				relaxTime, _ := strconv.Atoi(relaxTimeText.Text)
				iterationTime, _ := strconv.Atoi(iterationTimeText.Text)
				p.SetSettings(focusTime, relaxTime, iterationTime)
				p.RestartTimer()
				p.UpdateTimerText()
				p.UpdateIterationText()
				pomodoroWindow.Close()
			},
			pomodoroWindow,
		)
	})

	timerButtonContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(200, 40)),
		createTimerButton,
	)
	textContainer := container.NewVBox(
		focusTimeLabelContainer,
		relaxTimeLabelContainer,
		iterationTimeLabelContainer,
	)
	inputContainer := container.NewVBox(
		focusTimeTextContainer,
		relaxTimeTextContainer,
		iterationTimeTextContainer,
	)
	pomodoroInfoContainer := container.New(layout.NewHBoxLayout(), textContainer, inputContainer)
	content := container.New(layout.NewVBoxLayout(), pomodoroInfoContainer, timerButtonContainer)

	return &PomodoroCreationWindow{
		Window:             pomodoroWindow,
		Container:          content,
		FocusTimeInput:     focusTimeText,
		RelaxTimeInput:     relaxTimeText,
		IterationTimeInput: iterationTimeText,
	}
}

func (p *PomodoroCreationWindow) Render() {
	p.Window.SetContent(p.Container)
	p.Window.Resize(fyne.NewSize(400, 400))
	p.Window.Show()
}

type SettingsWindow struct {
	Window    fyne.Window
	Container *fyne.Container
}

func NewSettingsWindow(app fyne.App, s *pomoapp.Settings) *SettingsWindow {
	settingsWindow := app.NewWindow("Settings")

	// Widget creation
	libraryPathLabel := widget.NewLabel("Library path: ")
	libraryPath := widget.NewEntry()
	libraryPath.SetText(s.LibraryPath)
	autoPlayCheckBox := widget.NewCheck("Autoplay next song", func(checked bool) {
		s.AutoPlay = checked
	})
	autoPlayCheckBox.Checked = s.AutoPlay
	shuffleCheckBox := widget.NewCheck("Shuffle", func(checked bool) {
		s.Shuffle = checked
	})
	shuffleCheckBox.Checked = s.Shuffle
	linkPlayersCheckBox := widget.NewCheck(
		"Link Players (Pausing timer pauses music and vice versa)",
		func(checked bool) {
			s.Shuffle = checked
		},
	)
	linkPlayersCheckBox.Checked = s.LinkPlayers

	saveButton := widget.NewButton("Save", func() {
		dialog.ShowConfirm(
			"Confirm",
			"Are you sure you want to save these settings?",
			func(confirm bool) {
				s.Save(
					libraryPath.Text,
					autoPlayCheckBox.Checked,
					shuffleCheckBox.Checked,
					linkPlayersCheckBox.Checked,
				)
				settingsWindow.Close()
			},
			settingsWindow,
		)
	})

	// Create the rows
	libraySettingsRow := container.New(layout.NewVBoxLayout(), libraryPathLabel, libraryPath)
	playSettingsRow := container.New(
		layout.NewHBoxLayout(),
		autoPlayCheckBox,
		shuffleCheckBox,
		linkPlayersCheckBox,
	)
	saveRow := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), saveButton)

	content := container.New(layout.NewVBoxLayout(), libraySettingsRow, playSettingsRow, saveRow)

	return &SettingsWindow{
		Window:    settingsWindow,
		Container: content,
	}
}

func (p *SettingsWindow) Render() {
	p.Window.SetContent(p.Container)
	p.Window.Resize(fyne.NewSize(400, 400))
	p.Window.Show()
}

type LibraryView struct {
	LibraryList     *widget.List
	Container       *fyne.Container
	Library         *library.Library
	SongDetailsView *SongDetailsView
}

func NewLibraryView(
	labelText string,
	library *library.Library,
	songDetailsView *SongDetailsView,
	settings *pomoapp.Settings,
	player *player.Player,
) *LibraryView {
	libraryListLabel := widget.NewLabel(labelText)
	libraryList := widget.NewList(
		func() int {
			return len(library.Songs)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel(library.CurrentSong.Name)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(library.Songs[i].Name)
		})
	libraryListLabelContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(50, 50)),
		libraryListLabel,
	)
	libraryListContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(300, 400)),
		libraryList,
	)

	// Initialize and refresh right away because the first song should be selected
	l := LibraryView{
		LibraryList: libraryList,
		Container: container.New(
			layout.NewHBoxLayout(),
			libraryListLabelContainer,
			libraryListContainer,
			songDetailsView.Container,
		),
		Library:         library,
		SongDetailsView: songDetailsView,
	}
	l.UpdateSelected()

	libraryList.OnSelected = func(index int) {
		// Kill the currently playing song
		if library.CurrentSong.Player != nil && library.CurrentSong.Player.IsPlaying() {
			library.CurrentSong.Stop(false) // TODO(map) Is this right?
		}

		// Update the index and new song and start playing
		library.CurrIdx = index
		library.CurrentSong = library.Songs[library.CurrIdx]
		l.UpdateSelected()
		go library.CurrentSong.Play(player.SongControlChan)
	}

	return &l
}

func (l *LibraryView) UpdateSelected() {
	l.LibraryList.Select(l.Library.CurrIdx)
	l.LibraryList.Refresh()

	currentSong := l.Library.CurrentSong
	l.SongDetailsView.TitleInput.SetText(currentSong.Tag.Title())
	l.SongDetailsView.ArtistInput.SetText(currentSong.Tag.Artist())
	l.SongDetailsView.AlbumInput.SetText(currentSong.Tag.Album())
	l.SongDetailsView.GenreInput.SetText(currentSong.Tag.Genre())
	l.SongDetailsView.TitleInput.Refresh()
	l.SongDetailsView.ArtistInput.Refresh()
	l.SongDetailsView.AlbumInput.Refresh()
	l.SongDetailsView.GenreInput.Refresh()
}

type SongDetailsView struct {
	Container *fyne.Container

	// TODO(map) Include labels for these guys
	TitleInput  *widget.Entry
	ArtistInput *widget.Entry
	AlbumInput  *widget.Entry
	GenreInput  *widget.Entry

	CurrentSong *music.Song
}

func NewSongDetailsView(labelText string) *SongDetailsView {
	detailsLabel := widget.NewLabel(labelText)
	titleInput := widget.NewEntry()
	artistInput := widget.NewEntry()
	albumInput := widget.NewEntry()
	genreInput := widget.NewEntry()

	saveId3DataButton := widget.NewButton("Save", func() {
		fmt.Println("TODO(map) Implement the song save here")
	})

	songDetailsContainer := container.New(
		layout.NewVBoxLayout(),
		detailsLabel,
		titleInput,
		artistInput,
		albumInput,
		genreInput,
		saveId3DataButton,
	)

	return &SongDetailsView{
		Container: container.New(
			layout.NewGridWrapLayout(fyne.NewSize(300, 400)),
			songDetailsContainer,
		),
		TitleInput:  titleInput,
		ArtistInput: artistInput,
		AlbumInput:  albumInput,
		GenreInput:  genreInput,
	}
}

type MusicControls struct {
	Container  *fyne.Container
	PrevButton *widget.Button
	PlayButton *widget.Button
	StopButton *widget.Button
	NextButton *widget.Button
}

func NewMusicControls(
	library *library.Library,
	player *player.Player,
	settings *pomoapp.Settings,
	pomodoroTimer *pomodoro.PomodoroTimer,
) *MusicControls {
	prevButton := widget.NewButton("Prev", func() {
		log.Println("Prev clicked")
		if library.CurrIdx == 0 {
			// Do nothing because we can't decrement
			fmt.Println("Cannot go to previous song")
		} else {
			library.CurrentSong.Stop(true)
			library.DecIndex()
		}
	})
	playButton := widget.NewButton("Play", func() {
		log.Println("Play clicked")
		// Case where there is no Player set because the initial launch of the MP3 hasn't happened
		if library.CurrentSong.Player == nil {
			log.Println("No song set, playing song...")

			// Start the player
			go player.Play(library, settings)

			// Start the pomodoro timer if the timer and music controls are linked
			if settings.LinkPlayers {
				go pomodoroTimer.StartTimer()
			}
		} else if library.CurrentSong.Player.IsPlaying() { // Case of song is currently playing
			log.Println("Pausing song...")
			library.CurrentSong.Pause(player.SongControlChan)

			// Pause the pomodoro timer if the timer and music controls are linked
			if settings.LinkPlayers {
				pomodoroTimer.PauseTimer()
			}
		} else if !library.CurrentSong.Player.IsPlaying() { // Case where song is paused
			log.Println("Resuming song...")
			library.CurrentSong.Resume(player.SongControlChan)

			// Resume the pomodoro timer if the timer and music controls are linked
			if settings.LinkPlayers {
				go pomodoroTimer.StartTimer()
			}
		}
	})
	stopButton := widget.NewButton("Stop", func() {
		log.Println("Stop clicked")
		library.CurrentSong.Stop(false)

		// Pause the pomodoro timer if the timer and music controls are linked
		if settings.LinkPlayers {
			pomodoroTimer.PauseTimer()
		}
	})
	nextButton := widget.NewButton("Next", func() {
		log.Println("Next clicked")
		if !library.HasNextSong {
			// Do nothing because we can't decrement
			fmt.Println("Cannot go to next song")
		} else {
			library.CurrentSong.Stop(true)
			library.IncIndex()
		}
	})

	// Containers around the buttons to ensure their size doesn't grow beyond what is desired
	prevButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), prevButton)
	playButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), playButton)
	stopButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), stopButton)
	nextButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(50, 50)), nextButton)

	// Control
	controlsRow := container.New(
		layout.NewHBoxLayout(),
		prevButtonContainer,
		playButtonContainer,
		stopButtonContainer,
		nextButtonContainer,
	)

	return &MusicControls{
		Container:  container.New(layout.NewCenterLayout(), controlsRow),
		PrevButton: prevButton,
	}
}

func CreateNewToolbar(
	app fyne.App,
	pomodoroTimer *pomodoro.PomodoroTimer,
	appSettings *pomoapp.Settings,
) *widget.Toolbar {
	return widget.NewToolbar(
		// TODO(map) What's a good icon to use here? Maybe explore the idea of making my own resource
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			pomodoroCreationWindow := NewPomodoroCreationWindow(app, pomodoroTimer)
			pomodoroCreationWindow.Render()
		}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SettingsIcon(), func() {
			pomodoroSettingsWindow := NewSettingsWindow(app, appSettings)
			pomodoroSettingsWindow.Render()
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)
}
