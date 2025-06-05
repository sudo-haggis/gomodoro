package main

import (
	"testing"
	"time"
)

func TestFullPomodoroFlow(t *testing.T) {
	// Skip in short mode as this test takes time
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup
	DefaultSettings = GoModoroSettings{
		Sessions:           2,
		ShortBreak:         1, // 1 minute for testing
		LongBreak:          2,
		LongBreakFrequency: 0,
		Surprises:          1,
		SurpriseMinutes:    1,
	}

	sessionManager = NewSessionManager()
	currentState = TimerReady

	// Verify initial state
	if len(sessionManager.Sessions) < 3 {
		t.Fatal("Should have at least 3 sessions (2 work + 1 break)")
	}

	// Start first work session
	current := sessionManager.GetCurrentSession()
	if current.Type != SessionWork {
		t.Error("Should start with work session")
	}

	// Simulate timer start
	currentState = TimerRunning
	
	// Simulate timer finish
	currentState = TimerFinished

	// Move to next session
	nextSession()
	
	// Should now be on a break (possibly with surprise)
	current = sessionManager.GetCurrentSession()
	if current.Type != SessionShortBreak && current.Type != SessionSurprise {
		t.Errorf("Expected break or surprise, got %v", current.Type)
	}
}

func TestSettingsUpdate(t *testing.T) {
	// Initial settings
	DefaultSettings = GoModoroSettings{
		Sessions:           4,
		ShortBreak:         5,
		LongBreak:          15,
		LongBreakFrequency: 1,
		Surprises:          2,
		SurpriseMinutes:    3,
	}

	// Create initial session manager
	sessionManager = NewSessionManager()
	initialSessionCount := len(sessionManager.Sessions)

	// Update settings
	DefaultSettings.Sessions = 6
	DefaultSettings.LongBreakFrequency = 2

	// Recreate session manager (simulating save settings)
	sessionManager = NewSessionManager()
	newSessionCount := len(sessionManager.Sessions)

	// Verify sessions were updated
	if newSessionCount <= initialSessionCount {
		t.Error("Session count should increase when settings change")
	}

	// Verify long breaks were added
	longBreakCount := 0
	for _, session := range sessionManager.Sessions {
		if session.Type == SessionLongBreak {
			longBreakCount++
		}
	}

	if longBreakCount != 2 {
		t.Errorf("Expected 2 long breaks with frequency 2, got %d", longBreakCount)
	}
}

func TestConcurrentTimerOperations(t *testing.T) {
	// Setup channels
	controlChannel = make(chan string, 10)
	stopChannel = make(chan bool)

	// Start timer goroutine
	go func() {
		// Simplified timer loop for testing
		for {
			select {
			case cmd := <-controlChannel:
				switch cmd {
				case "start":
					currentState = TimerRunning
				case "pause":
					currentState = TimerPaused
				case "stop":
					return
				}
			case <-stopChannel:
				return
			}
		}
	}()

	// Send multiple commands rapidly
	controlChannel <- "start"
	controlChannel <- "pause"
	controlChannel <- "start"
	controlChannel <- "pause"

	// Give goroutine time to process
	time.Sleep(100 * time.Millisecond)

	// Clean up
	controlChannel <- "stop"
	time.Sleep(100 * time.Millisecond)
}

func TestSessionCompletionFlow(t *testing.T) {
	// Create a minimal session list
	DefaultSettings = GoModoroSettings{
		Sessions:           2,
		ShortBreak:         5,
		LongBreak:          15,
		LongBreakFrequency: 0,
		Surprises:          0,
	}

	sessionManager = NewSessionManager()

	// Complete all sessions
	for sessionManager.CurrentIndex < len(sessionManager.Sessions) {
		current := sessionManager.GetCurrentSession()
		if current == nil {
			break
		}

		// Mark as completed and move next
		current.Completed = true
		sessionManager.CurrentIndex++
	}

	// Verify all completed
	completed := sessionManager.GetCompletedSessions()
	if len(completed) != len(sessionManager.Sessions) {
		t.Error("All sessions should be marked as completed")
	}

	// Test cycle restart
	sessionManager = NewSessionManager()
	if sessionManager.CurrentIndex != 0 {
		t.Error("New cycle should start at index 0")
	}
}