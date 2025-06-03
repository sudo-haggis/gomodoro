// Every Go file starts with a package declaration
// "main" is special - it creates an executable program
package main

// Import section - like #include in C or imports in Python
// Parentheses let us import multiple packages cleanly
import (
	"fyne.io/fyne/v2/app"       // Creates the application
	"fyne.io/fyne/v2/container" // Layout containers (VBox, HBox, etc)
	"fyne.io/fyne/v2/widget"    // UI widgets (buttons, labels, etc)
	"fyne.io/fyne/v2"           // Base fyne package for types like Size
)

// main() is the entry point - Go runs this function first
func main() {
	// := is Go's "short variable declaration"
	// It declares AND assigns in one line (type inferred automatically)
	// Equivalent to: var myApp fyne.App = app.New()
	myApp := app.New()
	myApp.SetIcon(nil) // We'll add a proper pirate icon later!
	
	// Create the main window - string parameter is the window title
	myWindow := myApp.NewWindow("GoModoro - Ahoy!")
	// fyne.NewSize(width, height) creates a size object
	myWindow.Resize(fyne.NewSize(400, 300))
	
	// widget.NewLabel() creates a text label
	// Go strings can contain emojis directly!
	greeting := widget.NewLabel("🏴‍☠️ Ahoy, matey! Welcome aboard the GoModoro!")
	// Dot notation accesses struct fields (like object properties)
	greeting.Alignment = fyne.TextAlignCenter
	
	subtitle := widget.NewLabel("Yer productivity ship be ready to sail!")
	subtitle.Alignment = fyne.TextAlignCenter
	
	// Go functions can be passed as parameters (like JavaScript callbacks)
	// func() { ... } is an anonymous function (lambda/closure)
	testButton := widget.NewButton("Hoist the Colors! 🏴‍☠️", func() {
		// This function runs when button is clicked
		greeting.SetText("🦜 Avast! The timer be ready for action, captain!")
	})
	
	// container.NewVBox arranges widgets vertically
	// Go allows trailing commas in function calls (handy for multiline)
	content := container.NewVBox(
		widget.NewLabel(""), // Empty label acts as spacer
		greeting,
		widget.NewLabel(""), // Another spacer
		subtitle,
		widget.NewLabel(""), 
		testButton,
		widget.NewLabel(""), 
	)
	
	// Method chaining - set content then center window
	myWindow.SetContent(content)
	myWindow.CenterOnScreen()
	
	// ShowAndRun() displays window and starts the GUI event loop
	// This blocks until the window is closed
	myWindow.ShowAndRun()
}
