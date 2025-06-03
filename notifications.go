package main

import (
	"os/exec"
	"time"

	"fyne.io/fyne/v2/dialog"
)

// showSystemNotification sends a system notification (Ubuntu/Linux)
func showSystemNotification(title, message string) {
	// Use notify-send command on Ubuntu
	cmd := exec.Command("notify-send", title, message, "-i", "alarm-clock", "-t", "5000")
	cmd.Run() // Fire and forget - don't block if it fails
}

// showAnnoyingPopup creates a modal dialog that ye MUST click
func showAnnoyingPopup(title, message string) {
	// Create an information dialog - must be called from UI thread
	dialog.ShowInformation(title, message, myWindow)
}

// triggerAllAlerts - the full annoying experience!
func triggerAllAlerts() {
	// 1. System notification (safe from any thread)
	showSystemNotification("🏴‍☠️ GoModoro Complete!", "Avast! Yer pomodoro session be finished!")

	// 2. Pop-up dialog (must run on UI thread)
	go func() {
		time.Sleep(100 * time.Millisecond)
		showAnnoyingPopup("🍅 Session Complete!",
			"Ahoy, captain! Yer 25-minute session be done!\n\nTime to take a break and stretch yer sea legs!")
	}()

	// 3. Request window focus (safe from any thread)
	myWindow.RequestFocus()
}
