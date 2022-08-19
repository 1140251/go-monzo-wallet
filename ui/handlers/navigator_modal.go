package handlers

import "gioui.org/layout"

type Modal interface {
	// ID is a unique string that identifies the modal and may be used
	// to differentiate this modal from other modals.
	ID() string
	// OnAttachedToNavigator is called when navigation occurs; i.e. when a page
	// or modal is pushed into the window's display. The navigator parameter is
	// the PageNavigator or WindowNavigator object that is used to display the
	// content. This is called just before OnResume() is called.
	OnAttachedToNavigator(navigator PageNavigator)
	// OnResume is called to initialize data and get UI elements ready to be
	// displayed. This is called just before Handle() and Layout() are called (in
	// that order).
	OnResume()
	// Handle is called just before Layout() to determine if any user
	// interaction recently occurred on the modal and may be used to update the
	// page's UI components shortly before they are displayed.
	Handle()
	// Layout draws the modal's UI components into the provided layout context
	// to be eventually drawn on screen.
	Layout(gtx layout.Context) layout.Dimensions
	// OnDismiss is called after the modal is dismissed.
	// NOTE: The modal may be re-displayed on the app's window, in which case
	// OnResume() will be called again. This method should not destroy UI
	// components unless they'll be recreated in the OnResume() method.
	OnDismiss()
}
