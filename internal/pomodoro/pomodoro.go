package pomodoro

import (
	"fmt"
	"image/color"
	"pomogoro/internal/music"
	"pomogoro/internal/pomoapp"
	"time"

	// Gui imports
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Structure to represent a Pomodoro instance
type PomodoroSettings struct {
	StartFocusTime int // Start time for focus in seconds
	StartRelaxTime int // Starting point for the relax timer in seconds
	Iterations     int // Number of times the Focus/Relax combination should be repeated
	IterationCount int // The current count of the number of iterations completed

	// TODO(map) I will need to include something like focus playist and relaxing playlist or something like that
	// eventually. For now I can just do pause music during break
	PauseDuringBreak bool
}

// All things related to the Pomodoro including canvas to draw timer, the settings, and current status
type PomodoroTimer struct {
	CurrentTimer int  // The current time on the timer in seconds
	IsRunning    bool // Flag for if the timer is running
	InBreakMode  bool // Flag for whether we are in the relax portion or focus portion of the timer

	PomodoroSettings    PomodoroSettings    // The settings of the particular timer
	PomodoroTimerCanvas PomodoroTimerCanvas // Canvas to draw the timer and access all components
}

func NewPomodoroTimer(library *music.Library, settings *pomoapp.Settings) *PomodoroTimer {
	pt := &PomodoroTimer{
		IsRunning:   false,
		InBreakMode: false,
	}

	pt.CreateDefaultCanvas(library, settings)

	return pt
}

func (pt *PomodoroTimer) StartTimer() {
	pt.IsRunning = true
	// TODO(map) Good enough for now but we should really count down the final break period too
	for pt.IsRunning && pt.PomodoroSettings.IterationCount < pt.PomodoroSettings.Iterations {
		if pt.CurrentTimer > 0 {
			// Update timer
			fmt.Printf("Current timer = %d", pt.CurrentTimer)
			time.Sleep(time.Second * 1)
			pt.CurrentTimer -= 1
			pt.UpdateTimerText()
		} else {
			// Conditionally increment the counter only when finishing a focus period and refresh the text
			if !pt.InBreakMode {
				pt.PomodoroSettings.IterationCount += 1
				pt.UpdateIterationText()
			}

			// Switch modes
			pt.InBreakMode = !pt.InBreakMode

			// Update the timer with the appropriate value after the previous timer finishes
			if pt.InBreakMode {
				pt.CurrentTimer = pt.PomodoroSettings.StartRelaxTime
				pt.PomodoroTimerCanvas.ModeText.Text = "Relax"
				pt.PomodoroTimerCanvas.ModeText.Refresh()
			} else {
				pt.CurrentTimer = pt.PomodoroSettings.StartFocusTime
				pt.PomodoroTimerCanvas.ModeText.Text = "Focus"
				pt.PomodoroTimerCanvas.ModeText.Refresh()
			}
			pt.UpdateTimerText()

		}
	}
}

func (pt *PomodoroTimer) PauseTimer() {
	pt.IsRunning = false
}

func (pt *PomodoroTimer) RestartTimer() {
	pt.CurrentTimer = pt.PomodoroSettings.StartFocusTime
	pt.PomodoroSettings.IterationCount = 0
	pt.UpdateTimerText()
	pt.UpdateIterationText()
	pt.PomodoroTimerCanvas.ModeText.Text = "Focus"
	pt.PomodoroTimerCanvas.ModeText.Refresh()
}

func (pt *PomodoroTimer) SetSettings(startFocusTime int, startRelaxTime int, iterations int) {
	// NOTE(map) Multiply by 60 for the focus and relax time because the input units is in minutes but we track in seconds
	// so the math is easier and so we can do one second increments on the timer itself.
	pomodoroSettings := PomodoroSettings{
		// TODO(map) Uncomment me after testing
		// StartFocusTime: startFocusTime * 60,
		StartFocusTime: startFocusTime,
		// StartRelaxTime:   startRelaxTime * 60,
		StartRelaxTime:   startRelaxTime,
		Iterations:       iterations,
		PauseDuringBreak: false,
	}
	pt.PomodoroSettings = pomodoroSettings

	// Refresh the text to display to the user
	pt.UpdateTimerText()
	pt.UpdateIterationText()
}

func (pt *PomodoroTimer) UpdateTimerText() {
	// TODO(map) Render the time correctly based on the timer running or if reset is hit.
	pt.PomodoroTimerCanvas.TimerText.Text = fmt.Sprintf(
		"%d min %d sec",
		int(pt.CurrentTimer/60),
		int(pt.CurrentTimer%60),
	)
	pt.PomodoroTimerCanvas.TimerText.Refresh()
}

func (pt *PomodoroTimer) UpdateIterationText() {
	pt.PomodoroTimerCanvas.IterationText.Text = fmt.Sprintf(
		"Completed %d of %d Iterations",
		pt.PomodoroSettings.IterationCount,
		pt.PomodoroSettings.Iterations,
	)
	pt.PomodoroTimerCanvas.IterationText.Refresh()
}

func (pt *PomodoroTimer) CreateDefaultCanvas(library *music.Library, settings *pomoapp.Settings) {
	circleContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(300, 300)),
		canvas.NewCircle(color.RGBA{0, 0, 0, 255}),
	)
	modeText := canvas.NewText("Focus", color.RGBA{255, 255, 255, 255})
	modeTextInnerContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(100, 50)),
		modeText,
	)
	modeTextContainer := container.New(layout.NewCenterLayout(), modeTextInnerContainer)
	timerText := canvas.NewText("No Timer Created", color.RGBA{255, 255, 255, 255})
	timerTextInnerContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(100, 50)),
		timerText,
	)
	timerTextContainer := container.New(layout.NewCenterLayout(), timerTextInnerContainer)
	iterationText := canvas.NewText("", color.RGBA{255, 255, 255, 255})
	iterationTextInnerContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(100, 50)),
		iterationText,
	)
	iterationTextContainer := container.New(layout.NewCenterLayout(), iterationTextInnerContainer)
	playButton := widget.NewButton("Play", func() {
		if !pt.IsRunning {
			go pt.StartTimer()
			if settings.LinkPlayers {
				// TODO(map) Maybe this is bad to assume that there is no player if the song is not Paused because it
				// can't have started then
				if library.CurrentSong.IsPaused {
					library.CurrentSong.Resume()
					library.CurrentSong.IsPaused = false
				} else {
					go library.CurrentSong.Play(settings.LibraryPath)
				}
			}
		} else {
			pt.PauseTimer()
			if settings.LinkPlayers {
				library.CurrentSong.Pause()
				library.CurrentSong.IsPaused = true
			}
		}
	})
	playButtonContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(100, 50)),
		playButton,
	)
	resetButton := widget.NewButton("Restart", func() {
		pt.PauseTimer()
		pt.RestartTimer()
	})
	resetButtonContainer := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(100, 50)),
		resetButton,
	)
	textContainer := container.New(
		layout.NewVBoxLayout(),
		modeTextContainer,
		timerTextContainer,
		iterationTextContainer,
	)
	controlContainer := container.New(
		layout.NewHBoxLayout(),
		playButtonContainer,
		resetButtonContainer,
	)

	textPlusControlContainer := container.NewVBox(textContainer, controlContainer)
	topLevelContainer := container.New(
		layout.NewCenterLayout(),
		circleContainer,
		textPlusControlContainer,
	)

	pt.PomodoroTimerCanvas = PomodoroTimerCanvas{
		CircleContainer:             circleContainer,
		ModeText:                    modeText,
		ModeTextInnerContainer:      modeTextInnerContainer,
		ModeTextContainer:           modeTextContainer,
		TimerText:                   timerText,
		TimerTextInnerContainer:     timerTextInnerContainer,
		TimerTextContainer:          timerTextContainer,
		IterationText:               iterationText,
		IterationTextInnerContainer: iterationTextInnerContainer,
		IterationTextContainer:      iterationTextContainer,
		PlayButton:                  playButton,
		PlayButtonContainer:         playButtonContainer,
		ResetButton:                 resetButton,
		ResetButtonContainer:        resetButtonContainer,
		TextContainer:               textContainer,
		ControlContainer:            controlContainer,
		TextPlusControlContainer:    textPlusControlContainer,
		TopLevelContainer:           topLevelContainer,
	}
}

type PomodoroTimerCanvas struct {
	// Circle container to hold all the data and controls
	CircleContainer *fyne.Container

	// Text to display mode
	ModeText               *canvas.Text
	ModeTextInnerContainer *fyne.Container
	ModeTextContainer      *fyne.Container

	// Timer related components
	TimerText               *canvas.Text
	TimerTextInnerContainer *fyne.Container
	TimerTextContainer      *fyne.Container

	// Iterations related components
	IterationText               *canvas.Text
	IterationTextInnerContainer *fyne.Container
	IterationTextContainer      *fyne.Container

	// Controls
	PlayButton           *widget.Button
	PlayButtonContainer  *fyne.Container
	ResetButton          *widget.Button
	ResetButtonContainer *fyne.Container

	// Parent containers
	TextContainer    *fyne.Container
	ControlContainer *fyne.Container

	// Wrapper
	TextPlusControlContainer *fyne.Container

	// Top level container
	TopLevelContainer *fyne.Container
}
