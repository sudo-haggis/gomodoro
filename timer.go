package main

import (
	"fmt"
	"time"
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

// startTimer begins a new timer session
func startTimer() {
	if currentState == TimerReady || currentState == TimerPaused {
		currentState = TimerRunning
		// Create a new ticker that fires every second
		ticker = time.NewTicker(1 * time.Second)
		updateUI()
	}
}

// pauseTimer pauses the current session
func pauseTimer() {
	if currentState == TimerRunning {
		currentState = TimerPaused
		if ticker != nil {
			ticker.Stop() // Stop the countdown
		}
		updateUI()
	}
}

// resetTimer resets to a fresh 25-minute session
func resetTimer() {
	currentState = TimerReady
	timeRemaining = 25 * 60 // Reset to 25 minutes
	if ticker != nil {
		ticker.Stop()
	}
	updateUI()
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
				startTimer()
			case "pause":
				pauseTimer()
			case "reset":
				resetTimer()
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
					// Timer finished - UNLEASH THE KRAKEN OF NOTIFICATIONS!
					currentState = TimerFinished
					timeRemaining = 0
					ticker.Stop()
					// Fire all the alerts!
					triggerAllAlerts()
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
