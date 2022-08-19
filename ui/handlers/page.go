package handlers

import (
	"gioui.org/layout"
)

// Page defines methods that control the appearance and functionality of
// UI components displayed on a window.
type Page interface {
	// ID is a unique string that identifies the page and may be used
	// to differentiate this page from other pages.
	ID() string
	// OnAttachedToNavigator is called when navigation occurs; i.e. when a page
	// or modal is pushed into the window's display. The navigator parameter is
	// the PageNavigator or DefaultWindowNavigator object that is used to display the
	// content. This is called just before OnNavigatedTo() is called.
	OnAttachedToNavigator(navigator PageNavigator)
	// OnNavigatedTo is called when the page is about to be displayed and may be
	// used to initialize page features that are only relevant when the page is
	// displayed. This is called just before HandleUserInteractions() and
	// Layout() are called (in that order).
	OnNavigatedTo()
	// HandleUserInteractions is called just before Layout() to determine
	// if any user interaction recently occurred on the page and may be
	// used to update the page's UI components shortly before they are
	// displayed.
	HandleUserInteractions()
	// Layout draws the page UI components into the provided layout context
	// to be eventually drawn on screen.
	Layout(layout.Context) layout.Dimensions
	// OnNavigatedFrom is called when the page is about to be removed from
	// the displayed window. This method should ideally be used to disable
	// features that are irrelevant when the page is NOT displayed.
	// NOTE: The page may be re-displayed on the app's window, in which case
	// OnNavigatedTo() will be called again. This method should not destroy UI
	// components unless they'll be recreated in the OnNavigatedTo() method.
	OnNavigatedFrom()
}

// Closable should be implemented by pages and modals that want to know when
// they are closed in order to perform some cleanup actions.
type Closable interface {
	// OnClosed is called to indicate that a specific instance of a page or
	// modal has been dismissed and will no longer be displayed.
	OnClosed()
}
