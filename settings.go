package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Settings UI widgets (global so we can read their values)
var (
	sessionsEntry          *widget.Entry
	shortBreakEntry        *widget.Entry
	longBreakEntry         *widget.Entry
	longBreakFreqSelect    *widget.Select
	surprisesEntry         *widget.Entry
	surpriseDurationEntry  *widget.Entry
	settingsWindow         fyne.Window
)

// createSettingsUI builds the settings configuration page
func createSettingsUI() *fyne.Container {
	// Title
	title := widget.NewLabel("⚙️ GoModoro Settings")
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	// Instructions
	instructions := widget.NewLabel("Configure yer Pomodoro session parameters, matey!")
	instructions.Alignment = fyne.TextAlignCenter

	// Sessions setting
	sessionsLabel := widget.NewLabel("Sessions per cycle:")
	sessionsEntry = widget.NewEntry()
	sessionsEntry.SetText(strconv.Itoa(DefaultSettings.Sessions))
	sessionsEntry.SetPlaceHolder("6")

	// Short break setting
	shortBreakLabel := widget.NewLabel("Short break (minutes):")
	shortBreakEntry = widget.NewEntry()
	shortBreakEntry.SetText(strconv.Itoa(DefaultSettings.ShortBreak))
	shortBreakEntry.SetPlaceHolder("5")

	// Long break setting
	longBreakLabel := widget.NewLabel("Long break (minutes):")
	longBreakEntry = widget.NewEntry()
	longBreakEntry.SetText(strconv.Itoa(DefaultSettings.LongBreak))
	longBreakEntry.SetPlaceHolder("30")

	// Long break frequency setting
	longBreakFreqLabel := widget.NewLabel("Long break frequency:")
	longBreakFreqSelect = widget.NewSelect([]string{
		"None",
		"1 (middle)",
		"2 (thirds)",
		"3 (quarters)",
		"4 (fifths)",
	}, func(value string) {})
	longBreakFreqSelect.SetSelectedIndex(DefaultSettings.LongBreakFrequency)

	// Surprises setting
	surprisesLabel := widget.NewLabel("Max surprise tasks per cycle:")
	surprisesEntry = widget.NewEntry()
	surprisesEntry.SetText(strconv.Itoa(DefaultSettings.Surprises))
	surprisesEntry.SetPlaceHolder("3")

	// Surprise duration setting
	surpriseDurationLabel := widget.NewLabel("Surprise task duration (minutes):")
	surpriseDurationEntry = widget.NewEntry()
	surpriseDurationEntry.SetText(strconv.Itoa(DefaultSettings.SurpriseMinutes))
	surpriseDurationEntry.SetPlaceHolder("2")

	// Save button
	saveBtn := widget.NewButton("💾 Save & Close", func() {
		saveSettings()
		settingsWindow.Close()
	})

	// Cancel button
	cancelBtn := widget.NewButton("❌ Cancel", func() {
		settingsWindow.Close()
	})

	// Button container
	buttonContainer := container.NewHBox(saveBtn, cancelBtn)

	// Form-like layout using a border container for better mobile experience
	form := container.NewVBox(
		widget.NewSeparator(),
		sessionsLabel,
		sessionsEntry,
		widget.NewSeparator(),
		shortBreakLabel,
		shortBreakEntry,
		widget.NewSeparator(),
		longBreakLabel,
		longBreakEntry,
		longBreakFreqLabel,
		longBreakFreqSelect,
		widget.NewSeparator(),
		surprisesLabel,
		surprisesEntry,
		surpriseDurationLabel,
		surpriseDurationEntry,
		widget.NewSeparator(),
		buttonContainer,
	)

	// Main container
	content := container.NewVBox(
		widget.NewLabel(""), // Spacer
		title,
		instructions,
		widget.NewLabel(""), // Spacer
		form,
		widget.NewLabel(""), // Spacer
	)

	return content
}

// saveSettings reads the form values and updates the global settings
func saveSettings() {
	// Parse sessions (with error handling)
	if sessions, err := strconv.Atoi(sessionsEntry.Text); err == nil && sessions > 0 {
		DefaultSettings.Sessions = sessions
	}

	// Parse short break
	if shortBreak, err := strconv.Atoi(shortBreakEntry.Text); err == nil && shortBreak > 0 {
		DefaultSettings.ShortBreak = shortBreak
	}

	// Parse long break
	if longBreak, err := strconv.Atoi(longBreakEntry.Text); err == nil && longBreak > 0 {
		DefaultSettings.LongBreak = longBreak
	}

	// Parse long break frequency
	DefaultSettings.LongBreakFrequency = longBreakFreqSelect.SelectedIndex()

	// Parse surprises
	if surprises, err := strconv.Atoi(surprisesEntry.Text); err == nil && surprises >= 0 {
		DefaultSettings.Surprises = surprises
	}

	// Parse surprise duration
	if surpriseDuration, err := strconv.Atoi(surpriseDurationEntry.Text); err == nil && surpriseDuration > 0 {
		DefaultSettings.SurpriseMinutes = surpriseDuration
	}

	// Rebuild session list with new settings
	if sessionManager != nil {
		sessionManager = NewSessionManager()
		updateSessionDisplay()
		
		// Reset timer to use new session duration
		if currentState == TimerReady {
			if current := sessionManager.GetCurrentSession(); current != nil {
				timeRemaining = current.Duration
				updateUI()
			}
		}
	}

	// TODO: Save settings to file for persistence
	// For now, they just update the in-memory defaults
}

// showSettingsWindow creates and displays the settings window
func showSettingsWindow() {
	settingsWindow = myApp.NewWindow("GoModoro Settings")
	settingsWindow.Resize(fyne.NewSize(350, 600)) // Taller for more settings

	content := createSettingsUI()
	settingsWindow.SetContent(content)
	settingsWindow.CenterOnScreen()
	settingsWindow.Show()
}