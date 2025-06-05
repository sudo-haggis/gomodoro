package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}

// SessionType represents the type of timer session
type SessionType string

const (
	SessionWork       SessionType = "work"
	SessionShortBreak SessionType = "short_break"
	SessionLongBreak  SessionType = "long_break"
	SessionSurprise   SessionType = "surprise"
)

// SessionSlot represents a single time slot in the Pomodoro cycle
type SessionSlot struct {
	Type       SessionType `json:"type"`
	Duration   int         `json:"duration"` // in seconds
	Completed  bool        `json:"completed"`
	Current    bool        `json:"current"`
	SessionNum int         `json:"session_num,omitempty"` // Only for work sessions
}

// SessionManager handles the current todo list and session progression
type SessionManager struct {
	Sessions         []SessionSlot `json:"sessions"`
	CurrentIndex     int           `json:"current_index"`
	WorkCount        int           `json:"work_count"`
	TotalWorkCount   int           `json:"total_work_count"`
	SurpriseCount    int           `json:"surprise_count"`
	MaxSurpriseCount int           `json:"max_surprise_count"`
}

// Global session manager
var sessionManager *SessionManager

// NewSessionManager creates a fresh session list based on settings
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		Sessions:         make([]SessionSlot, 0),
		CurrentIndex:     0,
		WorkCount:        0,
		TotalWorkCount:   DefaultSettings.Sessions,
		SurpriseCount:    0,
		MaxSurpriseCount: DefaultSettings.Surprises,
	}
	
	sm.buildSessionList()
	return sm
}

// calculateLongBreakPositions determines which sessions get long breaks
func (sm *SessionManager) calculateLongBreakPositions() []int {
	positions := make([]int, 0)
	
	if DefaultSettings.LongBreakFrequency == 0 {
		return positions // No long breaks
	}
	
	// Calculate how many long breaks we need
	// 1 = 1 break (middle), 2 = 2 breaks (thirds), 3 = 3 breaks (quarters), etc.
	numBreaks := DefaultSettings.LongBreakFrequency
	
	if numBreaks >= DefaultSettings.Sessions {
		// Too many breaks requested, cap it
		numBreaks = DefaultSettings.Sessions - 1
	}
	
	// Calculate interval between long breaks
	interval := float64(DefaultSettings.Sessions) / float64(numBreaks + 1)
	
	// Place long breaks at calculated intervals
	for i := 1; i <= numBreaks; i++ {
		position := int(float64(i) * interval)
		// Ensure we don't place a long break after the last session
		if position > 0 && position < DefaultSettings.Sessions {
			positions = append(positions, position)
		}
	}
	
	return positions
}

// buildSessionList creates the full cycle of sessions
func (sm *SessionManager) buildSessionList() {
	sm.Sessions = make([]SessionSlot, 0)
	longBreakPositions := sm.calculateLongBreakPositions()
	workNum := 1
	
	// Always start with a work session
	sm.Sessions = append(sm.Sessions, SessionSlot{
		Type:       SessionWork,
		Duration:   25 * 60, // 25 minutes default
		Completed:  false,
		Current:    true,
		SessionNum: workNum,
	})
	workNum++
	
	// Add remaining sessions with breaks
	for i := 1; i < DefaultSettings.Sessions; i++ {
		// Check if we should add a surprise task (50/50 chance)
		if sm.SurpriseCount < sm.MaxSurpriseCount && rand.Float32() < 0.5 {
			sm.Sessions = append(sm.Sessions, SessionSlot{
				Type:      SessionSurprise,
				Duration:  DefaultSettings.SurpriseMinutes * 60,
				Completed: false,
				Current:   false,
			})
			sm.SurpriseCount++
		}
		
		// Determine if this should be a long break
		isLongBreak := false
		for _, pos := range longBreakPositions {
			if pos == i {
				isLongBreak = true
				break
			}
		}
		
		// Add the appropriate break
		if isLongBreak {
			sm.Sessions = append(sm.Sessions, SessionSlot{
				Type:      SessionLongBreak,
				Duration:  DefaultSettings.LongBreak * 60,
				Completed: false,
				Current:   false,
			})
		} else {
			sm.Sessions = append(sm.Sessions, SessionSlot{
				Type:      SessionShortBreak,
				Duration:  DefaultSettings.ShortBreak * 60,
				Completed: false,
				Current:   false,
			})
		}
		
		// Add the next work session
		sm.Sessions = append(sm.Sessions, SessionSlot{
			Type:       SessionWork,
			Duration:   25 * 60,
			Completed:  false,
			Current:    false,
			SessionNum: workNum,
		})
		workNum++
	}
}

// GetCurrentSession returns the current active session
func (sm *SessionManager) GetCurrentSession() *SessionSlot {
	if sm.CurrentIndex < len(sm.Sessions) {
		return &sm.Sessions[sm.CurrentIndex]
	}
	return nil
}

// NextSession moves to the next session (polymorphic progression)
func (sm *SessionManager) NextSession() bool {
	if sm.CurrentIndex < len(sm.Sessions) {
		// Mark current as completed
		sm.Sessions[sm.CurrentIndex].Completed = true
		sm.Sessions[sm.CurrentIndex].Current = false
		
		// Move to next
		sm.CurrentIndex++
		
		// Mark new current
		if sm.CurrentIndex < len(sm.Sessions) {
			sm.Sessions[sm.CurrentIndex].Current = true
			return true
		}
	}
	return false // No more sessions
}

// SkipCurrentSession skips the current session without marking complete
func (sm *SessionManager) SkipCurrentSession() bool {
	if sm.CurrentIndex < len(sm.Sessions) {
		sm.Sessions[sm.CurrentIndex].Current = false
		sm.CurrentIndex++
		
		if sm.CurrentIndex < len(sm.Sessions) {
			sm.Sessions[sm.CurrentIndex].Current = true
			return true
		}
	}
	return false
}

// RestartCurrentSession resets the timer for current session
func (sm *SessionManager) RestartCurrentSession() {
	// This just resets the timer, the session manager doesn't need to do anything
	// The timer.go will handle resetting timeRemaining
}

// GetSessionLabel returns a formatted string for the session type
func (s *SessionSlot) GetSessionLabel() string {
	switch s.Type {
	case SessionWork:
		return fmt.Sprintf("🍅 Work Session %d", s.SessionNum)
	case SessionShortBreak:
		return "☕ Short Break"
	case SessionLongBreak:
		return "🏖️ Long Break"
	case SessionSurprise:
		return "⚡ Surprise Task!"
	default:
		return "Unknown"
	}
}

// GetSessionDurationString returns formatted duration
func (s *SessionSlot) GetSessionDurationString() string {
	return fmt.Sprintf("%d min", s.Duration/60)
}

// ToJSON exports the session manager state
func (sm *SessionManager) ToJSON() (string, error) {
	data, err := json.MarshalIndent(sm, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromJSON imports session manager state
func (sm *SessionManager) FromJSON(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), sm)
}

// GetCompletedSessions returns list of completed sessions
func (sm *SessionManager) GetCompletedSessions() []SessionSlot {
	completed := make([]SessionSlot, 0)
	for _, session := range sm.Sessions {
		if session.Completed {
			completed = append(completed, session)
		}
	}
	return completed
}

// GetRemainingSessions returns list of remaining sessions
func (sm *SessionManager) GetRemainingSessions() []SessionSlot {
	remaining := make([]SessionSlot, 0)
	foundCurrent := false
	
	for _, session := range sm.Sessions {
		if session.Current {
			foundCurrent = true
			continue // Don't include current in remaining
		}
		if foundCurrent && !session.Completed {
			remaining = append(remaining, session)
		}
	}
	return remaining
}