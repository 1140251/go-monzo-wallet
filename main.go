package main

import (
	"gioui.org/app"
	"go-monzo-wallet/ui"
	"os"
)

func main() {

	win, err := ui.CreateWindow()
	if err != nil {
		os.Exit(1)
	}

	go func() {
		win.HandleEvents() // blocks until the app window is closed
		os.Exit(0)
	}()

	// Start the GUI frontend.
	app.Main()

}
