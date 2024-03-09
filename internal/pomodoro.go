package pomodoro

import (
    "image/color"
	"fmt"
	"time"

	// Gui imports
	"fyne.io/fyne/v2/canvas"
)

// Structure to represent a Pomodoro instance
type PomodoroSettings struct {
    StartFocusTime int // Start time for focus in seconds
    StartRelaxTime int // Starting point for the relax timer in seconds
    Iterations int // Number of times the Focus/Relax combination should be repeated
    IterationCount int // The current count of the number of iterations completed

    // TODO(map) I will need to include something like focus playist and relaxing playlist or something like that
    // eventually. For now I can just do pause music during break
    PauseDuringBreak bool
}

// All things related to the Pomodoro including canvas to draw timer, the settings, and current status
type PomodoroTimer struct {
    CurrentTimer int // The current time on the timer in seconds
    IsRunning bool // Flag for if the timer is running
    InBreakMode bool // Flag for whether we are in the relax portion or focus portion of the timer

    PomodoroSettings PomodoroSettings // The settings of the particular timer
    PomodoroTimerCanvas PomodoroTimerCanvas // Canvas to draw the timer and access all components
}

func (pt *PomodoroTimer) StartTimer() {
    pt.IsRunning = true
   // TODO(map) This is not ideal because it will always wait one additional second before actually pausing the timer.
    for pt.IsRunning && pt.CurrentTimer > 0 {
        fmt.Printf("Current timer = %d", pt.CurrentTimer)
        time.Sleep(time.Second * 1)
        pt.CurrentTimer -= 1
        pt.UpdateTimerText()
    }

    // If this is reached then the timer has finished and additional work is to be done.
    if pt.CurrentTimer <= 0 {
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
        } else {
            pt.CurrentTimer = pt.PomodoroSettings.StartFocusTime
        }
        pt.UpdateTimerText()

        // Reset the timer flag
        pt.IsRunning = false
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
}

func (pt *PomodoroTimer) SetSettings(startFocusTime int, startRelaxTime int, iterations int) {
    // NOTE(map) Multiply by 60 for the focus and relax time because the input units is in minutes but we track in seconds
    // so the math is easier and so we can do one second increments on the timer itself.
    pomodoroSettings := PomodoroSettings{
        // TODO(map) Uncomment me after testing
        // StartFocusTime: startFocusTime * 60,
        StartFocusTime: startFocusTime,
        StartRelaxTime: startRelaxTime * 60,
        Iterations: iterations,
        PauseDuringBreak: false,
    }
    pt.PomodoroSettings = pomodoroSettings

    // Refresh the text to display to the user
    pt.UpdateTimerText()
    pt.UpdateIterationText()
}

func (pt *PomodoroTimer) UpdateTimerText() {
    // TODO(map) Render the time correctly based on the timer running or if reset is hit.
    pt.PomodoroTimerCanvas.TimerText.Text = fmt.Sprintf("%d min %d sec", int(pt.CurrentTimer/60), int(pt.CurrentTimer%60))
    pt.PomodoroTimerCanvas.TimerText.Refresh()
}

func (pt *PomodoroTimer) UpdateIterationText() {
    pt.PomodoroTimerCanvas.IterationText.Text = fmt.Sprintf("Completed %d of %d Iterations", pt.PomodoroSettings.IterationCount, pt.PomodoroSettings.Iterations)
    pt.PomodoroTimerCanvas.IterationText.Refresh()
}

func (pt *PomodoroTimer) CreateDefaultCanvas() {
    circleContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(300, 300)), canvas.NewCircle(color.RGBA{0, 0, 0, 255}))
    timerText := canvas.NewText("No Timer Created", color.RGBA{255, 255, 255, 255})
    timerTextInnerContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 50)), timerText)
    timerTextContainer := container.New(layout.NewCenterLayout(), timerTextInnerContainer)
    iterationText := canvas.NewText("", color.RGBA{255, 255, 255, 255})
    iterationTextInnerContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 50)), iterationText)
    iterationTextContainer := container.New(layout.NewCenterLayout(), iterationTextInnerContainer)
    playButton := widget.NewButton("Play", func() {
        if !pt.IsRunning {
            go pt.StartTimer()
        } else {
            pt.PauseTimer()
        }
    })
    playButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 50)), playButton)
    resetButton := widget.NewButton("Restart", func() {
        pt.PauseTimer()
        pt.RestartTimer()
    })
    resetButtonContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(100, 50)), resetButton)
    textContainer := container.New(layout.NewVBoxLayout(), timerTextContainer, iterationTextContainer)
    controlContainer := container.New(layout.NewHBoxLayout(), playButtonContainer, resetButtonContainer)

    textPlusControlContainer := container.NewVBox(textContainer, controlContainer)
    topLevelContainer := container.New(layout.NewCenterLayout(), circleContainer, textPlusControlContainer)

    pt.PomodoroTimerCanvas = PomodoroTimerCanvas{
        CircleContainer: circleContainer,
        TimerText: timerText,
        TimerTextInnerContainer: timerTextInnerContainer,
        TimerTextContainer: timerTextContainer,
        IterationText: iterationText,
        IterationTextInnerContainer: iterationTextInnerContainer,
        IterationTextContainer: iterationTextContainer,
        PlayButton: playButton,
        PlayButtonContainer: playButtonContainer,
        ResetButton: resetButton,
        ResetButtonContainer: resetButtonContainer,
        TextContainer: textContainer,
        ControlContainer: controlContainer,
        TextPlusControlContainer: textPlusControlContainer,
        TopLevelContainer: topLevelContainer,
    }
}

type PomodoroTimerCanvas struct {
    // Circle container to hold all the data and controls
    CircleContainer *fyne.Container

    // Timer related components
    TimerText *canvas.Text
    TimerTextInnerContainer *fyne.Container
    TimerTextContainer *fyne.Container

    // Iterations related components
    IterationText *canvas.Text
    IterationTextInnerContainer *fyne.Container
    IterationTextContainer *fyne.Container

    // Controls
    PlayButton *widget.Button
    PlayButtonContainer *fyne.Container
    ResetButton * widget.Button
    ResetButtonContainer *fyne.Container

    // Parent containers
    TextContainer *fyne.Container
    ControlContainer *fyne.Container

    // Wrapper
    TextPlusControlContainer *fyne.Container

    // Top level container
    TopLevelContainer *fyne.Container
}


