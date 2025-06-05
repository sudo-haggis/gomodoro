# 🍅 GoModoro - A Pirate's Pomodoro Timer

A Pomodoro timer built with Go and Fyne, designed for both desktop (Ubuntu) and mobile (Android) use.

## Features

- **Flexible Sessions**: Configure work sessions with breaks
- **Smart Breaks**: Short breaks after each session, with configurable long breaks
- **Surprise Tasks**: Random mini-tasks to keep things interesting
- **Session Tracking**: See completed and upcoming sessions
- **Annoying Notifications**: System notifications and pop-ups when sessions complete
- **Pirate Theme**: Because why not? 🏴‍☠️

## Installation

### Prerequisites

1. **Go 1.19+** installed
2. **Fyne dependencies**:
   ```bash
   # Ubuntu/Debian
   sudo apt-get install gcc libgl1-mesa-dev xorg-dev
   ```

### Setup

1. **Clone/Create project structure**:
   ```bash
   mkdir gomodoro
   cd gomodoro
   go mod init gomodoro
   ```

2. **Install Fyne**:
   ```bash
   go get fyne.io/fyne/v2@latest
   ```

3. **Add all the Go files**:
   - `main.go`
   - `types.go`
   - `timer.go`
   - `settings.go`
   - `session.go`
   - `ui_updates.go`
   - `notifications.go`

4. **Build and run**:
   ```bash
   go build .
   ./gomodoro
   ```

## Building for Android

1. **Install fyne command tool**:
   ```bash
   go install fyne.io/fyne/v2/cmd/fyne@latest
   ```

2. **Package for Android**:
   ```bash
   fyne package -os android -appID com.pirate.gomodoro
   ```

3. **Install on your Pixel 6**:
   ```bash
   adb install GoModoro.apk
   ```

## Settings Configuration

- **Sessions per cycle**: Number of work sessions (default: 6)
- **Short break**: Minutes for short breaks (default: 4)
- **Long break**: Minutes for long breaks (default: 30)
- **Long break frequency**:
  - 0 = No long breaks
  - 1 = One break in the middle
  - 2 = Two breaks (at thirds)
  - 3 = Three breaks (at quarters)
  - 4 = Four breaks (at fifths)
- **Max surprises**: Maximum surprise tasks per cycle (default: 3)
- **Surprise duration**: Minutes per surprise task (default: 2)

## Session Flow

1. Always starts with a work session
2. After each work session (except the last):
   - 50% chance of a surprise task
   - Followed by a break (short or long based on position)
3. Always ends with a work session

## Troubleshooting

- **Notifications not working on Ubuntu**: Ensure `notify-send` is installed
- **Window too small on mobile**: The app is designed for Pixel 6 aspect ratio
- **Timer not updating**: Check that the goroutine is running properly

## Future Features

- Settings persistence
- Custom surprise task list
- Statistics tracking
- Sound notifications
- Theme customization

## License

This be free software, matey! Use it as ye please! 🏴‍☠️