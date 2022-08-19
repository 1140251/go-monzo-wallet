package handlers

import "gioui.org/io/key"

// KeyEventHandler is implemented by pages and modals that require key event
// notifications.
type KeyEventHandler interface {
	// KeysToHandle returns an expression that describes a set of key
	// combinations that the implementer of this interface wishes to capture.
	// The HandleKeyPress() method will only be called when any of these key
	// combinations is pressed.
	KeysToHandle() key.Set
	// HandleKeyPress is called when one or more keys are pressed on the current
	// window that match any of the key combinations returned by KeysToHandle().
	HandleKeyPress(*key.Event)
}
