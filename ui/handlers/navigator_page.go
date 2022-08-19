package handlers

// PageNavigator defines methods for navigating between pages in a window or a
// MasterPage.
type PageNavigator interface {
	// CurrentPage returns the page that is at the top of the stack. Returns nil
	// if the stack is empty.
	CurrentPage() Page
	// CurrentPageID returns the ID of the current page or an empty string if no
	// page is displayed.
	CurrentPageID() string
	// Display causes the specified page to be displayed on the parent window or
	// page. All other instances of this same page will be closed and removed
	// from the backstack.
	Display(page Page)
	// CloseCurrentPage dismisses the page at the top of the stack and gets the
	// next page ready for display.
	CloseCurrentPage()
	// ClosePagesAfter dismisses all pages from the top of the stack until (and
	// excluding) the page with the specified ID. If no page is found with the
	// provided ID, no page will be popped. The page with the specified ID will
	// be displayed after the other pages are popped.
	ClosePagesAfter(keepPageID string)
	// ClearStackAndDisplay dismisses all pages in the stack and displays the
	// specified page.
	ClearStackAndDisplay(page Page)
	// CloseAllPages dismisses all pages in the stack.
	CloseAllPages()
}
