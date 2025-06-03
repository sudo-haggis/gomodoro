// Every Go file starts with a package declaration
// "main" is special - it creates an executable program
package main

import (
	"fmt"  // For string formatting
	"time" // For time operations and ticker

	"fyne.io/fyne/v2"           // Base fyne package for types like Size
	"fyne.io/fyne/v2/app"       // Creates the application
	"fyne.io/fyne/v2/container" // Layout containers (VBox, HBox, etc)
	"fyne.io/fyne/v2/widget"    // UI widgets (buttons, labels, etc)
)

// Timer states - using constants (like enums in other languages)
const (
	TimerReady    = "ready"
	TimerRunning  = "running"
	TimerPaused   = "paused"
	TimerFinished = "finished"
)

// Global variables for timer state
var (
	timeRemaining = 25 * 60 // 25 minutes in seconds
	currentState  = TimerReady
	ticker        *time.Ticker // Go's built-in timer that fires every interval
	timeDisplay   *widget.Label
	startPauseBtn *widget.Button
	resetBtn      *widget.Button

	// Channels for goroutine communication (like ship-to-ship signals!)
	timerChannel   = make(chan int)    // Sends time updates
	controlChannel = make(chan string) // Sends control commands
	stopChannel    = make(chan bool)   // Tells goroutine to stop
)

// formatTime converts seconds to MM:SS format
func formatTime(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}

// updateUI changes the display based on current timer state
func updateUI() {
	timeDisplay.SetText(formatTime(timeRemaining))

	// Change button text based on state
	switch currentState {
	case TimerReady:
		startPauseBtn.SetText("🏴‍☠️ Start Timer!")
		resetBtn.SetText("Reset")
	case TimerRunning:
		startPauseBtn.SetText("⏸️ Pause")
		resetBtn.SetText("Reset")
	case TimerPaused:
		startPauseBtn.SetText("▶️ Resume")
		resetBtn.SetText("Reset")
	case TimerFinished:
		startPauseBtn.SetText("🎉 Session Complete!")
		resetBtn.SetText("New Session")
		timeDisplay.SetText("00:00 - Ahoy! Well done!")
	}
}

// timerGoroutine runs in the background and handles the countdown
// This is where the Go concurrency magic happens!
func timerGoroutine() {
	for {
		select {
		// Listen for control commands from UI
		case command := <-controlChannel:
			switch command {
			case "start":
				if currentState == TimerReady || currentState == TimerPaused {
					currentState = TimerRunning
					// Create a new ticker that fires every second
					ticker = time.NewTicker(1 * time.Second)
					updateUI()
				}
			case "pause":
				if currentState == TimerRunning {
					currentState = TimerPaused
					if ticker != nil {
						ticker.Stop() // Stop the countdown
					}
					updateUI()
				}
			case "reset":
				currentState = TimerReady
				timeRemaining = 25 * 60 // Reset to 25 minutes
				if ticker != nil {
					ticker.Stop()
				}
				updateUI()
			}

		// Listen for ticker events (every second when running)
		case <-func() <-chan time.Time {
			if ticker != nil {
				return ticker.C
			}
			// Return a channel that never sends if ticker is nil
			return make(<-chan time.Time)
		}():
			if currentState == TimerRunning {
				timeRemaining--
				if timeRemaining <= 0 {
					// Timer finished!
					currentState = TimerFinished
					timeRemaining = 0
					ticker.Stop()
				}
				updateUI()
			}

		// Listen for stop signal (when app closes)
		case <-stopChannel:
			if ticker != nil {
				ticker.Stop()
			}
			return
		}
	}
}

func main() {
	// Create the app
	myApp := app.New()
	myApp.SetIcon(nil)

	// Create the main window
	myWindow := myApp.NewWindow("GoModoro - Pomodoro Timer")
	myWindow.Resize(fyne.NewSize(400, 300))

	// Create UI elements
	title := widget.NewLabel("🍅 GoModoro Timer")
	title.Alignment = fyne.TextAlignCenter

	// Big time display
	timeDisplay = widget.NewLabel(formatTime(timeRemaining))
	timeDisplay.Alignment = fyne.TextAlignCenter
	// Make the time display bigger (this is a Fyne-specific way)
	timeDisplay.TextStyle = fyne.TextStyle{Bold: true}

	// Start/Pause button - uses anonymous function (closure)
	startPauseBtn = widget.NewButton("🏴‍☠️ Start Timer!", func() {
		switch currentState {
		case TimerReady, TimerPaused:
			controlChannel <- "start" // Send command to goroutine
		case TimerRunning:
			controlChannel <- "pause"
		case TimerFinished:
			controlChannel <- "reset" // Start new session
		}
	})

	// Reset button
	resetBtn = widget.NewButton("Reset", func() {
		controlChannel <- "reset"
	})

	// Layout buttons horizontally
	buttonContainer := container.NewHBox(startPauseBtn, resetBtn)

	// Main layout - arrange everything vertically
	content := container.NewVBox(
		widget.NewLabel(""), // Spacer
		title,
		widget.NewLabel(""), // Spacer
		timeDisplay,
		widget.NewLabel(""), // Spacer
		buttonContainer,
		widget.NewLabel(""), // Spacer
	)

	myWindow.SetContent(content)
	myWindow.CenterOnScreen()

	// Start the timer goroutine BEFORE showing the window
	go timerGoroutine()

	// Handle window closing - send stop signal to goroutine
	myWindow.SetCloseIntercept(func() {
		stopChannel <- true
		myApp.Quit()
	})

	// Show the window and run (this blocks until window closes)
	myWindow.ShowAndRun()
}
