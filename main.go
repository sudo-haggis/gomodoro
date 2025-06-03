// GoModoro - A Pirate's Pomodoro Timer
// Built with Go and Fyne for desktop and mobile
package main

import (
	"fyne.io/fyne/v2"           // Base fyne package for types like Size
	"fyne.io/fyne/v2/app"       // Creates the application
	"fyne.io/fyne/v2/container" // Layout containers (VBox, HBox, etc)
	"fyne.io/fyne/v2/widget"    // UI widgets (buttons, labels, etc)
)

// createTimerUI builds the main timer interface
func createTimerUI() *fyne.Container {
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

	// Settings button
	settingsBtn := widget.NewButton("⚙️ Settings", func() {
		showSettingsWindow()
	})

	// Layout buttons horizontally
	buttonContainer := container.NewHBox(startPauseBtn, resetBtn, settingsBtn)

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

	return content
}

func main() {
	// Create the app - store in global variable
	myApp = app.New()
	myApp.SetIcon(nil)

	// Create the main window - store in global variable for notifications
	myWindow = myApp.NewWindow("GoModoro - Pomodoro Timer")
	myWindow.Resize(fyne.NewSize(400, 300))

	// Create and set the main UI
	content := createTimerUI()
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
