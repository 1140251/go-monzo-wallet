package ui

import (
	"errors"
	giouiApp "gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/sirupsen/logrus"
	"go-monzo-wallet/internal"
	"go-monzo-wallet/ui/assets"
	"go-monzo-wallet/ui/components"
	"go-monzo-wallet/ui/handlers"
	"go-monzo-wallet/ui/pages"
	"go-monzo-wallet/ui/values"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Window struct {
	*giouiApp.Window
	load      *handlers.Load
	navigator handlers.WindowNavigator
}

func CreateWindow() (*Window, error) {

	giouiWindow := giouiApp.NewWindow(giouiApp.MinSize(values.AppWidth, values.AppHeight), giouiApp.Title(values.StrAppTitle))

	win := &Window{
		Window:    giouiWindow,
		navigator: handlers.NewNavigator(giouiWindow.Invalidate),
	}

	l, err := win.NewLoad()
	if err != nil {
		return nil, err
	}
	win.load = l

	return win, nil

}

func (win *Window) NewLoad() (*handlers.Load, error) {
	th := components.NewTheme(assets.FontCollection(), assets.Icons, false)
	if th == nil {
		return nil, errors.New("unexpected error while loading theme")
	}

	l := &handlers.Load{
		Theme:   th,
		Toast:   components.NewToast(th),
		Printer: message.NewPrinter(language.English),
		WL:      &internal.Wallet{},
	}

	return l, nil

}

// HandleEvents runs main event handling and page rendering loop.
func (win *Window) HandleEvents() {

	for {
		e := <-win.Events()
		switch evt := e.(type) {

		case system.DestroyEvent:
			win.navigator.CloseAllPages()
			return // exits the loop, caller will exit the program.

		case system.FrameEvent:
			ops := win.handleFrameEvent(evt)
			evt.Frame(ops)

		default:
			logrus.Info("Unhandled window event %v\n", e)
		}
	}
}

// handleFrameEvent is called when a FrameEvent is received by the active
// window. It expects a new frame in the form of a list of operations that
// describes what to display and how to handle input. This operations list
// is returned to the caller for displaying on screen.
func (win *Window) handleFrameEvent(evt system.FrameEvent) *op.Ops {
	switch {
	case win.navigator.CurrentPage() == nil:
		// Prepare to display the StartPage if no page is currently displayed.
		win.navigator.Display(pages.NewStartPage(win.load))

	default:
		// The app window may have received some user interaction such as key
		// presses, a button click, etc which triggered this FrameEvent. Handle
		// such interactions before re-displaying the UI components. This
		// ensures that the proper interface is displayed to the user based on
		// the action(s) they just performed.
		win.handleRelevantKeyPresses(evt)
		win.navigator.CurrentPage().HandleUserInteractions()
		if modal := win.navigator.TopModal(); modal != nil {
			modal.Handle()
		}
	}

	// Generate an operations list with instructions for drawing the window's UI
	// components onto the screen. Use the generated ops to request key handlers.
	ops := win.prepareToDisplayUI(evt)
	win.addKeyEventRequestsToOps(ops)

	return ops
}

// handleRelevantKeyPresses checks if any open modal or the current page is a
// load.KeyEventHandler AND if the provided system.FrameEvent contains key press
// handlers for the modal or page.
func (win *Window) handleRelevantKeyPresses(evt system.FrameEvent) {
	handleKeyPressFor := func(tag string, maybeHandler interface{}) {
		handler, ok := maybeHandler.(handlers.KeyEventHandler)
		if !ok {
			return
		}
		for _, event := range evt.Queue.Events(tag) {
			if keyEvent, isKeyEvent := event.(key.Event); isKeyEvent && keyEvent.State == key.Press {
				handler.HandleKeyPress(&keyEvent)
			}
		}
	}

	// Handle key handlers on the top modal first, if there's one.
	// Only handle key handlers on the current page if no modal is displayed.
	if modal := win.navigator.TopModal(); modal != nil {
		handleKeyPressFor(modal.ID(), modal)
	} else {
		handleKeyPressFor(win.navigator.CurrentPageID(), win.navigator.CurrentPage())
	}
}

// prepareToDisplayUI creates an operation list and writes the layout of all the
// window UI components into it. The created ops is returned and may be used to
// record further operations before finally being rendered on screen via
// system.FrameEvent.Frame(ops).
func (win *Window) prepareToDisplayUI(evt system.FrameEvent) *op.Ops {
	backgroundWidget := layout.Expanded(func(gtx C) D {
		return components.Fill(gtx, win.load.Theme.Color.Gray4)
	})

	currentPageWidget := layout.Stacked(func(gtx C) D {
		if modal := win.navigator.TopModal(); modal != nil {
			gtx = gtx.Disabled()
		}
		return win.navigator.CurrentPage().Layout(gtx)
	})

	topModalLayout := layout.Stacked(func(gtx C) D {
		modal := win.navigator.TopModal()
		if modal == nil {
			return layout.Dimensions{}
		}
		return modal.Layout(gtx)
	})

	// Use a StackLayout to write the above UI components into an operations
	// list via a graphical context that is linked to the ops.
	ops := &op.Ops{}
	gtx := layout.NewContext(ops, evt)
	layout.Stack{Alignment: layout.N}.Layout(
		gtx,
		backgroundWidget,
		currentPageWidget,
		topModalLayout,
		layout.Stacked(win.load.Toast.Layout),
	)

	return ops
}

// addKeyEventRequestsToOps checks if the current page or any modal has
// registered to be notified of certain key handlers and updates the provided
// operations list with instructions to generate a FrameEvent if any of the
// desired keys is pressed on the window.
func (win *Window) addKeyEventRequestsToOps(ops *op.Ops) {
	requestKeyEvents := func(tag string, desiredKeys key.Set) {
		if desiredKeys == "" {
			return
		}

		// Execute the key.InputOP{}.Add operation after all other operations.
		// This is particularly important because some pages call op.Defer to
		// signfiy that some operations should be executed after all other
		// operations, which has an undesirable effect of discarding this key
		// operation unless it's done last, after all other defers are done.
		m := op.Record(ops)
		key.InputOp{Tag: tag, Keys: desiredKeys}.Add(ops)
		op.Defer(ops, m.Stop())
	}

	// Request key handlers on the top modal, if necessary.
	// Only request key handlers on the current page if no modal is displayed.
	if modal := win.navigator.TopModal(); modal != nil {
		if handler, ok := modal.(handlers.KeyEventHandler); ok {
			requestKeyEvents(modal.ID(), handler.KeysToHandle())
		}
	} else {
		if handler, ok := win.navigator.CurrentPage().(handlers.KeyEventHandler); ok {
			requestKeyEvents(win.navigator.CurrentPageID(), handler.KeysToHandle())
		}
	}
}
