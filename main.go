// GoModoro - A Pirate's Pomodoro Timer
// Built with Go and Fyne for desktop and mobile
package main

import (
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2"           // Base fyne package for types like Size
	"fyne.io/fyne/v2/app"       // Creates the application
	"fyne.io/fyne/v2/container" // Layout containers (VBox, HBox, etc)
	"fyne.io/fyne/v2/widget"    // UI widgets (buttons, labels, etc)
)

// createTimerUI builds the main timer interface
func createTimerUI() *fyne.Container {
	// Initialize session manager
	sessionManager = NewSessionManager()

	// Create UI elements
	title := widget.NewLabel("🍅 GoModoro Timer")
	title.Alignment = fyne.TextAlignCenter

	// Current session label
	currentSessionLabel = widget.NewLabel("")
	currentSessionLabel.Alignment = fyne.TextAlignCenter
	currentSessionLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Big time display
	timeDisplay = widget.NewLabel(formatTime(timeRemaining))
	timeDisplay.Alignment = fyne.TextAlignCenter
	timeDisplay.TextStyle = fyne.TextStyle{Bold: true}

	// Session progress display - smaller text
	completedLabel = widget.NewLabel("Previous:")
	completedList = widget.NewLabel("")
	// Note: Fyne doesn't support direct color styling in labels
	// We'll use RichText widgets for colored text in production

	remainingLabel = widget.NewLabel("Upcoming:")
	remainingList = widget.NewLabel("")

	// Start/Pause button
	startPauseBtn = widget.NewButton("🏴‍☠️ Start Timer!", func() {
		switch currentState {
		case TimerReady, TimerPaused:
			controlChannel <- "start"
		case TimerRunning:
			controlChannel <- "pause"
		case TimerFinished:
			controlChannel <- "next"
		}
	})

	// Reset button
	resetBtn = widget.NewButton("Reset", func() {
		controlChannel <- "reset"
	})

	// Skip button
	skipBtn = widget.NewButton("Skip →", func() {
		controlChannel <- "skip"
	})

	// Settings button
	settingsBtn := widget.NewButton("⚙️", func() {
		showSettingsWindow()
	})

	// Layout buttons more compactly
	mainButtonContainer := container.NewHBox(startPauseBtn)
	secondaryButtonContainer := container.NewHBox(resetBtn, skipBtn, settingsBtn)

	// Compact session lists
	sessionProgress := container.NewVBox(
		completedLabel,
		completedList,
		widget.NewSeparator(),
		remainingLabel,
		remainingList,
	)

	// Main layout - more compact
	content := container.NewVBox(
		title,
		currentSessionLabel,
		widget.NewLabel(""), // Small spacer
		timeDisplay,
		widget.NewLabel(""), // Small spacer
		mainButtonContainer,
		secondaryButtonContainer,
		widget.NewLabel(""), // Small spacer
		sessionProgress,
	)

	// Update the display
	updateSessionDisplay()

	return content
}

func main() {
	// Create the app - store in global variable
	myApp = app.New()
	myApp.SetIcon(nil)

	// Try to load previous state first
	if err := loadAppState(); err != nil {
		// If loading fails, start fresh but don't crash
		sessionManager = NewSessionManager()
		currentState = TimerReady
		if current := sessionManager.GetCurrentSession(); current != nil {
			timeRemaining = current.Duration
		}
	}

	// Create the main window - store in global variable for notifications
	myWindow = myApp.NewWindow("GoModoro - Pomodoro Timer")
	myWindow.Resize(fyne.NewSize(400, 500)) // Smaller height

	// Create and set the main UI
	content := createTimerUI()
	myWindow.SetContent(content)
	myWindow.CenterOnScreen()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		// Save state before exiting
		saveAppState()
		stopChannel <- true
		myApp.Quit()
	}()

	// Start the timer goroutine BEFORE showing the window
	go timerGoroutine()

	// Start auto-save goroutine
	go autoSaveState()

	// Handle window closing - send stop signal to goroutine
	myWindow.SetCloseIntercept(func() {
		// Save state on window close
		saveAppState()
		stopChannel <- true
		myApp.Quit()
	})

	// Show the window and run (this blocks until window closes)
	myWindow.ShowAndRun()
}

