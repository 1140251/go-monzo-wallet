package handlers

import (
	"sync"
)

// WindowNavigator defines methods for page navigation, displaying modals and
// reloading the entire window display.
type WindowNavigator interface {
	PageNavigator
	// ShowModal displays a modal over the current page. Any previously
	// displayed modal will be hidden by this new modal.
	ShowModal(Modal)
	// DismissModal dismisses the modal with the specified ID, if it was
	// previously displayed by this WindowNavigator. If there are more than 1
	// modal with the specified ID, only the top-most instance is dismissed.
	DismissModal(modalID string)
	// TopModal returns the top-most modal in display or nil if there is no
	// modal in display.
	TopModal() Modal
	// Reload causes the entire window display to be reloaded. If a page is
	// currently displayed, this should call the page's HandleUserInteractions()
	// method. If a modal is displayed, the modal's Handle() method should also
	// be called. Finally, the current page and modal's Layout methods should be
	// called to render the entire window's display.
	Reload()
}

type DefaultWindowNavigator struct {
	reloadDisplayFn func()
	subPages        *PageStack
	modalMutex      sync.Mutex
	modals          []Modal
}

// NewNavigator creates an instance of a DefaultWindowNavigator.
func NewNavigator(reloadDisplayFn func()) WindowNavigator {
	w := &DefaultWindowNavigator{
		reloadDisplayFn: reloadDisplayFn,
		subPages:        NewPageStack("main window"),
	}
	return w
}

// CurrentPage returns the page that is at the top of the stack. Returns nil if
// the stack is empty.
// Part of the PageNavigator interface.
func (window *DefaultWindowNavigator) CurrentPage() Page {
	return window.subPages.Top()
}

// CurrentPageID returns the ID of the current page or an empty string if no
// page is displayed.
// Part of the PageNavigator interface.
func (window *DefaultWindowNavigator) CurrentPageID() string {
	if currentPage := window.CurrentPage(); currentPage != nil {
		return currentPage.ID()
	}
	return ""
}

// Display causes the specified page to be displayed on this window. All other
// instances of this same page will be closed and removed from the backstack.
// Part of the PageNavigator interface.
func (window *DefaultWindowNavigator) Display(newPage Page) {
	pushed := window.subPages.Push(newPage, window)
	if pushed {
		window.Reload()
	}
}

// CloseCurrentPage dismisses the page at the top of the stack and gets the next
// page ready for display.
// Part of the PageNavigator interface.
func (window *DefaultWindowNavigator) CloseCurrentPage() {
	popped := window.subPages.Pop()
	if popped {
		window.Reload()
	}
}

// ClosePagesAfter dismisses all pages from the top of the stack until (and
// excluding) the page with the specified ID. If no page is found with the
// provided ID, no page will be popped. The page with the specified ID will be
// displayed after the other pages are popped.
// Part of the PageNavigator interface.
func (window *DefaultWindowNavigator) ClosePagesAfter(keepPageID string) {
	popped := window.subPages.PopAfter(func(page Page) bool {
		return page.ID() == keepPageID
	})
	if popped {
		window.Reload()
	}
}

// ClearStackAndDisplay dismisses all pages in the stack and displays the
// specified page.
// Part of the PageNavigator interface.
func (window *DefaultWindowNavigator) ClearStackAndDisplay(newPage Page) {
	newPage.OnAttachedToNavigator(window)
	window.subPages.Reset(newPage)
	window.Reload()
}

// CloseAllPages dismisses all pages in the stack.
// Part of the PageNavigator interface.
func (window *DefaultWindowNavigator) CloseAllPages() {
	window.subPages.Reset()
	window.Reload()
}

// ShowModal displays a modal over the current page. Any previously displayed
// modal will be hidden by this new modal. NOTE: Allows displaying multiple
// instances of the same modal.
// Part of the DefaultWindowNavigator interface.
func (window *DefaultWindowNavigator) ShowModal(modal Modal) {
	window.modalMutex.Lock()
	window.modals = append(window.modals, modal)
	window.modalMutex.Unlock()

	modal.OnAttachedToNavigator(window)
	modal.OnResume()
	window.Reload()
}

// DismissModal dismisses the modal with the specified ID, if it was previously
// displayed by this DefaultWindowNavigator. If there are more than 1 modal with the
// specified ID, only the top-most instance is dismissed.
// Part of the DefaultWindowNavigator interface.
func (window *DefaultWindowNavigator) DismissModal(modalID string) {
	var modalToDismiss Modal

	window.modalMutex.Lock()
	for i := len(window.modals) - 1; i >= 0; i-- {
		modal := window.modals[i]
		if modal.ID() == modalID {
			modalToDismiss = modal
			window.modals = append(window.modals[:i], window.modals[i+1:]...)
			break
		}
	}
	window.modalMutex.Unlock()

	if modalToDismiss != nil {
		modalToDismiss.OnDismiss() // do garbage collection in modal
		window.Reload()
	}
}

// TopModal returns the top-most modal in display or nil if there is no modal in
// display.
// Part of the DefaultWindowNavigator interface.
func (window *DefaultWindowNavigator) TopModal() Modal {
	window.modalMutex.Lock()
	defer window.modalMutex.Unlock()
	if l := len(window.modals); l > 0 {
		return window.modals[l-1]
	}
	return nil
}

// Reload causes the entire window display to be reloaded. If a page is
// currently displayed, this will call the page's HandleUserInteractions()
// method. If a modal is displayed, the modal's Handle() method will also be
// called. Finally, the current page and modal's Layout() methods are called to
// render the entire window's display.
// Part of the DefaultWindowNavigator interface.
func (window *DefaultWindowNavigator) Reload() {
	window.reloadDisplayFn()
}