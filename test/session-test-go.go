package main

import (
	"testing"
)

func TestNewSessionManager(t *testing.T) {
	// Save original settings
	originalSettings := DefaultSettings
	defer func() {
		DefaultSettings = originalSettings
	}()

	// Test with specific settings
	DefaultSettings = GoModoroSettings{
		Sessions:           4,
		ShortBreak:         5,
		LongBreak:          15,
		LongBreakFrequency: 1,
		Surprises:          2,
		SurpriseMinutes:    3,
	}

	sm := NewSessionManager()

	// Check basic properties
	if sm.TotalWorkCount != 4 {
		t.Errorf("Expected 4 work sessions, got %d", sm.TotalWorkCount)
	}

	if sm.MaxSurpriseCount != 2 {
		t.Errorf("Expected max 2 surprises, got %d", sm.MaxSurpriseCount)
	}

	// Check that we start with a work session
	if len(sm.Sessions) == 0 || sm.Sessions[0].Type != SessionWork {
		t.Error("First session should be a work session")
	}

	// Check that we end with a work session
	lastSession := sm.Sessions[len(sm.Sessions)-1]
	if lastSession.Type != SessionWork {
		t.Error("Last session should be a work session")
	}
}

func TestCalculateLongBreakPositions(t *testing.T) {
	tests := []struct {
		name      string
		sessions  int
		frequency int
		expected  int // expected number of long breaks
	}{
		{"No long breaks", 6, 0, 0},
		{"One in middle", 6, 1, 1},
		{"Two thirds", 6, 2, 2},
		{"Three quarters", 8, 3, 3},
		{"Four fifths", 10, 4, 4},
		{"More breaks than sessions", 4, 5, 3}, // Should cap at sessions-1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalSettings := DefaultSettings
			defer func() {
				DefaultSettings = originalSettings
			}()

			DefaultSettings.Sessions = tt.sessions
			DefaultSettings.LongBreakFrequency = tt.frequency

			sm := NewSessionManager()
			positions := sm.calculateLongBreakPositions()

			if len(positions) != tt.expected {
				t.Errorf("Expected %d long break positions, got %d", tt.expected, len(positions))
			}
		})
	}
}

func TestSessionProgression(t *testing.T) {
	// Create a simple session manager
	DefaultSettings = GoModoroSettings{
		Sessions:           3,
		ShortBreak:         5,
		LongBreak:          15,
		LongBreakFrequency: 0,
		Surprises:          0,
		SurpriseMinutes:    2,
	}

	sm := NewSessionManager()

	// Test initial state
	current := sm.GetCurrentSession()
	if current == nil || current.Type != SessionWork {
		t.Error("Should start with a work session")
	}

	// Progress through sessions
	if !sm.NextSession() {
		t.Error("Should be able to progress to next session")
	}

	current = sm.GetCurrentSession()
	if current == nil || current.Type != SessionShortBreak {
		t.Error("Second session should be a short break")
	}

	// Test skip functionality
	initialIndex := sm.CurrentIndex
	if !sm.SkipCurrentSession() {
		t.Error("Should be able to skip current session")
	}

	if sm.CurrentIndex != initialIndex+1 {
		t.Error("Skip should advance to next session")
	}
}

func TestSessionLabels(t *testing.T) {
	tests := []struct {
		sessionType SessionType
		sessionNum  int
		expected    string
	}{
		{SessionWork, 1, "🍅 Work Session 1"},
		{SessionShortBreak, 0, "☕ Short Break"},
		{SessionLongBreak, 0, "🏖️ Long Break"},
		{SessionSurprise, 0, "⚡ Surprise Task!"},
	}

	for _, tt := range tests {
		session := SessionSlot{
			Type:       tt.sessionType,
			SessionNum: tt.sessionNum,
		}

		label := session.GetSessionLabel()
		if label != tt.expected {
			t.Errorf("Expected label '%s', got '%s'", tt.expected, label)
		}
	}
}

func TestSurpriseTaskLimit(t *testing.T) {
	// Set deterministic settings
	DefaultSettings = GoModoroSettings{
		Sessions:           10,
		ShortBreak:         5,
		LongBreak:          15,
		LongBreakFrequency: 0,
		Surprises:          3, // Max 3 surprises
		SurpriseMinutes:    2,
	}

	// Run multiple times to account for randomness
	for i := 0; i < 10; i++ {
		sm := NewSessionManager()

		// Count surprise tasks
		surpriseCount := 0
		for _, session := range sm.Sessions {
			if session.Type == SessionSurprise {
				surpriseCount++
			}
		}

		if surpriseCount > 3 {
			t.Errorf("Surprise count %d exceeds maximum of 3", surpriseCount)
		}
	}
}

func TestJSONSerialization(t *testing.T) {
	sm := NewSessionManager()
	
	// Export to JSON
	jsonData, err := sm.ToJSON()
	if err != nil {
		t.Fatalf("Failed to export to JSON: %v", err)
	}

	// Create new manager and import
	sm2 := &SessionManager{}
	err = sm2.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to import from JSON: %v", err)
	}

	// Compare key fields
	if sm2.CurrentIndex != sm.CurrentIndex {
		t.Error("CurrentIndex mismatch after JSON roundtrip")
	}

	if len(sm2.Sessions) != len(sm.Sessions) {
		t.Error("Session count mismatch after JSON roundtrip")
	}
}