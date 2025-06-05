package main

import (
	"fmt"
	"strings"
)

// updateSessionDisplay updates the session progress lists
func updateSessionDisplay() {
	if sessionManager == nil {
		return
	}

	// Update current session label
	if current := sessionManager.GetCurrentSession(); current != nil {
		currentSessionLabel.SetText(current.GetSessionLabel() + " (" + current.GetSessionDurationString() + ")")
		
		// Update time remaining based on current session
		if currentState == TimerReady {
			timeRemaining = current.Duration
			timeDisplay.SetText(formatTime(timeRemaining))
		}
	} else {
		currentSessionLabel.SetText("🎉 All Sessions Complete!")
	}

	// Show only 3 completed sessions (most recent)
	completed := sessionManager.GetCompletedSessions()
	if len(completed) > 0 {
		var completedText []string
		start := 0
		if len(completed) > 3 {
			start = len(completed) - 3
		}
		for i := start; i < len(completed); i++ {
			completedText = append(completedText, "✓ "+completed[i].GetSessionLabel())
		}
		completedList.SetText(strings.Join(completedText, "\n"))
	} else {
		completedList.SetText("None yet")
	}

	// Show only next 3 remaining sessions
	remaining := sessionManager.GetRemainingSessions()
	if len(remaining) > 0 {
		var remainingText []string
		limit := 3
		if len(remaining) < limit {
			limit = len(remaining)
		}
		for i := 0; i < limit; i++ {
			remainingText = append(remainingText, "• "+remaining[i].GetSessionLabel())
		}
		remainingList.SetText(strings.Join(remainingText, "\n"))
	} else {
		remainingList.SetText("Last one!")
	}
}

// updateUIWithSession updates UI based on current timer state and session
func updateUIWithSession() {
	timeDisplay.SetText(formatTime(timeRemaining))
	
	// Update session display
	updateSessionDisplay()

	// Change button text based on state
	switch currentState {
	case TimerReady:
		startPauseBtn.SetText("🏴‍☠️ Start Timer!")
		resetBtn.SetText("Reset Current")
		skipBtn.Enable()
	case TimerRunning:
		startPauseBtn.SetText("⏸️ Pause")
		resetBtn.SetText("Reset Current")
		skipBtn.Disable() // Can't skip while running
	case TimerPaused:
		startPauseBtn.SetText("▶️ Resume")
		resetBtn.SetText("Reset Current")
		skipBtn.Enable()
	case TimerFinished:
		if sessionManager != nil && sessionManager.CurrentIndex < len(sessionManager.Sessions)-1 {
			startPauseBtn.SetText("➡️ Next Session")
			resetBtn.SetText("Repeat Session")
		} else {
			startPauseBtn.SetText("🎊 All Done!")
			resetBtn.SetText("Start New Cycle")
		}
		skipBtn.Disable()
		
		// Special finish message
		// Safely get current session with nil check
		if sessionManager != nil {
			current := sessionManager.GetCurrentSession()
			if current != nil {
				timeDisplay.SetText(fmt.Sprintf("00:00 - %s Complete!", current.GetSessionLabel()))
			}
		}
	}
}