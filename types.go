package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Timer states - using constants (like enums in other languages)
const (
	TimerReady    = "ready"
	TimerRunning  = "running"
	TimerPaused   = "paused"
	TimerFinished = "finished"
)

// Global application state
var (
	// Timer state
	timeRemaining = 30 // 60*25 = minutes in seconds
	currentState  = TimerReady
	ticker        *time.Ticker // Go's built-in timer that fires every interval

	// UI references
	timeDisplay   *widget.Label
	startPauseBtn *widget.Button
	resetBtn      *widget.Button
	myWindow      fyne.Window // Need reference for notifications
	myApp         fyne.App    // Need app reference for thread-safe UI updates

	// Channels for goroutine communication (like ship-to-ship signals!)
	timerChannel   = make(chan int)    // Sends time updates
	controlChannel = make(chan string) // Sends control commands
	stopChannel    = make(chan bool)   // Tells goroutine to stop
)

// Settings structure for future configuration
type GoModoroSettings struct {
	Sessions   int // Number of work sessions (default: 6)
	ShortBreak int // Short break minutes (default: 4)
	LongBreak  int // Long break minutes (default: 1)
	Surprises  int // Surprise tasks per session (default: 3)
}

// Default settings
var DefaultSettings = GoModoroSettings{
	Sessions:   6,
	ShortBreak: 4,
	LongBreak:  1,
	Surprises:  3,
}
