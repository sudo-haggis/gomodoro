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
	timeRemaining = 30 // Will be set by session manager
	currentState  = TimerReady
	ticker        *time.Ticker // Go's built-in timer that fires every interval

	// UI references
	timeDisplay         *widget.Label
	currentSessionLabel *widget.Label
	completedLabel      *widget.Label
	completedList       *widget.Label
	remainingLabel      *widget.Label
	remainingList       *widget.Label
	startPauseBtn       *widget.Button
	resetBtn            *widget.Button
	skipBtn             *widget.Button
	myWindow            fyne.Window // Need reference for notifications
	myApp               fyne.App    // Need app reference for thread-safe UI updates

	// Channels for goroutine communication (like ship-to-ship signals!)
	timerChannel   = make(chan int, 1)    // Sends time updates (buffered to prevent blocking)
	controlChannel = make(chan string, 1) // Sends control commands (buffered to prevent blocking)
	stopChannel    = make(chan bool, 1)   // Tells goroutine to stop (buffered to prevent blocking)
)

// Settings structure for future configuration
type GoModoroSettings struct {
	Sessions             int // Number of work sessions (default: 6)
	ShortBreak           int // Short break minutes (default: 4)
	LongBreak            int // Long break minutes (default: 30)
	LongBreakFrequency   int // 1=/2, 2=/3, 3=/4, 4=/5 (default: 1)
	Surprises            int // Max surprise tasks per cycle (default: 3)
	SurpriseMinutes      int // Duration of surprise tasks (default: 2)
}

// Default settings
var DefaultSettings = GoModoroSettings{
	Sessions:           6,
	ShortBreak:         4,
	LongBreak:          30,
	LongBreakFrequency: 1,  // One long break in the middle
	Surprises:          3,
	SurpriseMinutes:    2,
}