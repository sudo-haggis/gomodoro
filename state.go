package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// AppState represents the complete application state for persistence
type AppState struct {
	SessionManagerState *SessionManager  `json:"session_manager"`
	CurrentState        string           `json:"current_state"`
	TimeRemaining       int              `json:"time_remaining"`
	LastSaved           time.Time        `json:"last_saved"`
	Settings            GoModoroSettings `json:"settings"`
}

// getStateFilePath returns the path where state should be saved
func getStateFilePath() string {
	// Use XDG_CONFIG_HOME or fallback to ~/.config
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "./gomodoro_state.json" // Fallback to current directory
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	// Create gomodoro config directory
	goModoroDir := filepath.Join(configDir, "gomodoro")
	os.MkdirAll(goModoroDir, 0755)

	return filepath.Join(goModoroDir, "session_state.json")
}

// saveAppState saves the current application state to disk
func saveAppState() error {
	state := AppState{
		SessionManagerState: sessionManager,
		CurrentState:        currentState,
		TimeRemaining:       timeRemaining,
		LastSaved:           time.Now(),
		Settings:            DefaultSettings,
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	statePath := getStateFilePath()
	return os.WriteFile(statePath, data, 0644)
}

// loadAppState loads the application state from disk
func loadAppState() error {
	statePath := getStateFilePath()

	// Check if state file exists
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		return nil // No state file, start fresh
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		return err
	}

	var state AppState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	// Restore application state
	if state.SessionManagerState != nil {
		sessionManager = state.SessionManagerState
	}
	currentState = state.CurrentState
	timeRemaining = state.TimeRemaining
	DefaultSettings = state.Settings

	// If we were in a running state, pause instead to avoid confusion
	if currentState == TimerRunning {
		currentState = TimerPaused
	}

	return nil
}

// clearAppState removes the saved state file
func clearAppState() error {
	statePath := getStateFilePath()
	return os.Remove(statePath)
}

// autoSaveState saves state periodically while running
func autoSaveState() {
	// Save state every 30 seconds when timer is running
	autoSaveTicker := time.NewTicker(30 * time.Second)
	defer autoSaveTicker.Stop()

	for {
		select {
		case <-autoSaveTicker.C:
			if currentState == TimerRunning || currentState == TimerPaused {
				saveAppState() // Ignore errors in auto-save
			}
		case <-stopChannel:
			// Save final state before exit
			saveAppState()
			return
		}
	}
}
