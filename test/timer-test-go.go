package main

import (
	"testing"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		seconds  int
		expected string
	}{
		{0, "00:00"},
		{59, "00:59"},
		{60, "01:00"},
		{90, "01:30"},
		{1500, "25:00"},
		{3600, "60:00"},
	}

	for _, tt := range tests {
		result := formatTime(tt.seconds)
		if result != tt.expected {
			t.Errorf("formatTime(%d) = %s; want %s", tt.seconds, result, tt.expected)
		}
	}
}

func TestTimerStateTransitions(t *testing.T) {
	// Initialize for testing
	sessionManager = NewSessionManager()
	currentState = TimerReady
	timeRemaining = 25 * 60

	// Test start timer
	startTimer()
	if currentState != TimerRunning {
		t.Error("Timer should be in running state after start")
	}

	// Test pause timer
	pauseTimer()
	if currentState != TimerPaused {
		t.Error("Timer should be in paused state after pause")
	}

	// Test resume (start from paused)
	startTimer()
	if currentState != TimerRunning {
		t.Error("Timer should resume to running state")
	}

	// Clean up ticker
	if ticker != nil {
		ticker.Stop()
		ticker = nil
	}
}

func TestResetTimer(t *testing.T) {
	// Setup
	sessionManager = NewSessionManager()
	currentState = TimerRunning
	timeRemaining = 100

	// Perform reset
	resetTimer()

	// Verify state
	if currentState != TimerReady {
		t.Error("Timer should be in ready state after reset")
	}

	// Check time is reset to current session duration
	current := sessionManager.GetCurrentSession()
	if current != nil && timeRemaining != current.Duration {
		t.Errorf("Time remaining should be %d, got %d", current.Duration, timeRemaining)
	}
}

func TestNextSession(t *testing.T) {
	// Create a small session list
	DefaultSettings = GoModoroSettings{
		Sessions:           2,
		ShortBreak:         5,
		LongBreak:          15,
		LongBreakFrequency: 0,
		Surprises:          0,
	}

	sessionManager = NewSessionManager()
	initialSession := sessionManager.GetCurrentSession()

	// Move to next session
	nextSession()

	newSession := sessionManager.GetCurrentSession()
	if newSession == initialSession {
		t.Error("Should have moved to a different session")
	}

	if currentState != TimerReady {
		t.Error("Timer should be in ready state after moving to next session")
	}
}

func TestSkipSession(t *testing.T) {
	// Setup
	sessionManager = NewSessionManager()
	initialIndex := sessionManager.CurrentIndex

	// Skip current session
	skipSession()

	if sessionManager.CurrentIndex <= initialIndex {
		t.Error("Session index should have increased after skip")
	}

	if currentState != TimerReady {
		t.Error("Timer should be in ready state after skip")
	}
}

func TestTimerEdgeCases(t *testing.T) {
	// Test with nil session manager
	sessionManager = nil
	
	// These should not panic
	resetTimer()
	if timeRemaining != 25*60 {
		t.Error("Should use default time when session manager is nil")
	}

	// Test state transitions from finished state
	currentState = TimerFinished
	startTimer()
	// Should not start from finished state
	if currentState == TimerRunning {
		t.Error("Should not start timer from finished state")
	}
}