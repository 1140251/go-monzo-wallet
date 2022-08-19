package components

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/gen2brain/beeep"
	"go-monzo-wallet/ui/values"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

type duration int32

const (
	Short duration = iota
	Long
)
const (
	icon  = "ui/assets/decredicons/qrcodeSymbol.png"
	title = "Godcr"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

type Notification interface {
	Notify(message string, d ...duration)
	NotifyError(message string, d ...duration)
}

type Toast struct {
	sync.Mutex
	theme   *Theme
	success bool
	message string
	timer   *time.Timer
}

func NewToast(th *Theme) *Toast {
	return &Toast{
		theme: th,
	}
}

// Notify is called to display a message indicating a successful action. It updates the toast object with the toast message
// and duration. The duration parameter is optional.
func (t *Toast) Notify(message string, d ...duration) {
	t.notify(message, true, d...)
}

// NotifyError is called to display a message indicating a failed action. It updates the toast object with the toast message
// and duration. The duration parameter is optional.
func (t *Toast) NotifyError(message string, d ...duration) {
	t.notify(message, false, d...)
}

// notify updates notification parameters on the toast object. It takes the message, success and duration
// parameters.
func (t *Toast) notify(message string, success bool, d ...duration) {
	var notificationDelay duration
	if len(d) > 0 {
		notificationDelay = d[0]
	}

	t.Lock()
	t.message = message
	t.success = success
	t.timer = time.NewTimer(getDurationFromDelay(notificationDelay))
	t.Unlock()
}

func getDurationFromDelay(d duration) time.Duration {
	if d == Long {
		return 5 * time.Second
	}
	return 2 * time.Second
}

func (t *Toast) Layout(gtx layout.Context) layout.Dimensions {
	t.handleToastDisplay(gtx)
	if t.timer == nil {
		return layout.Dimensions{}
	}

	color := t.theme.Color.Success
	if !t.success {
		color = t.theme.Color.Danger
	}

	card := t.theme.Card()
	card.Color = color
	return layout.Center.Layout(gtx, func(gtx C) D {
		return layout.Inset{Top: values.MarginPadding65}.Layout(gtx, func(gtx C) D {
			return card.Layout(gtx, func(gtx C) D {
				return layout.Inset{
					Top: values.MarginPadding7, Bottom: values.MarginPadding7,
					Left: values.MarginPadding15, Right: values.MarginPadding15,
				}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					msg := t.theme.Body1(t.message)
					msg.Color = t.theme.Color.Surface
					return msg.Layout(gtx)
				})
			})
		})
	})
}

func (t *Toast) handleToastDisplay(gtx layout.Context) {
	if t.timer == nil {
		return
	}

	select {
	case <-t.timer.C:
		t.timer = nil
		op.InvalidateOp{}.Add(gtx.Ops)
	default:
	}
}

type SystemNotification struct {
	iconPath string
	message  string
}

func NewSystemNotification() (*SystemNotification, error) {
	absolutePath, err := getAbsolutePath()
	if err != nil {
		return nil, err
	}

	return &SystemNotification{
		iconPath: filepath.Join(absolutePath, icon),
	}, nil
}

func (s *SystemNotification) Notify(message string) error {
	err := beeep.Notify(title, message, s.iconPath)
	if err != nil {
		return err
	}

	return nil
}

func getAbsolutePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %s", err.Error())
	}

	exSym, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", fmt.Errorf("error getting filepath after evaluating sym links")
	}

	return path.Dir(exSym), nil
}
