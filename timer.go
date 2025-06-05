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
	updateUIWithSession()
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

// resetTimer resets the current session
func resetTimer() {
	currentState = TimerReady
	if current := sessionManager.GetCurrentSession(); current != nil {
		timeRemaining = current.Duration
	} else {
		timeRemaining = 25 * 60 // Default fallback
	}
	// Safely stop ticker if it exists
	if ticker != nil {
		ticker.Stop()
		ticker = nil
	}
	updateUI()
}

// nextSession moves to the next session in the list
func nextSession() {
	if sessionManager.NextSession() {
		currentState = TimerReady
		if current := sessionManager.GetCurrentSession(); current != nil {
			timeRemaining = current.Duration
		}
		updateUI()
	} else {
		// All sessions complete - start new cycle
		sessionManager = NewSessionManager()
		currentState = TimerReady
		if current := sessionManager.GetCurrentSession(); current != nil {
			timeRemaining = current.Duration
		}
		updateUI()
	}
}

// skipSession skips the current session
func skipSession() {
	if sessionManager.SkipCurrentSession() {
		currentState = TimerReady
		if current := sessionManager.GetCurrentSession(); current != nil {
			timeRemaining = current.Duration
		}
		updateUI()
	} else {
		// No more sessions - start new cycle
		sessionManager = NewSessionManager()
		currentState = TimerReady
		if current := sessionManager.GetCurrentSession(); current != nil {
			timeRemaining = current.Duration
		}
		updateUI()
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
				startTimer()
			case "pause":
				pauseTimer()
			case "reset":
				resetTimer()
			case "next":
				nextSession()
			case "skip":
				skipSession()
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
