package handlers

import (
	"go-monzo-wallet/internal"
	"go-monzo-wallet/ui/components"
	"golang.org/x/text/message"
)

type Load struct {
	Theme *components.Theme

	Printer         *message.Printer
	Network         string
	CurrentAppWidth int
	Toast           *components.Toast
	WL              *internal.Wallet

	ToggleSync             func()
	DarkModeSettingChanged func(bool)
	LanguageSettingChanged func()
	CurrencySettingChanged func()
}

func (l *Load) RefreshTheme(window *DefaultWindowNavigator) {
	l.LanguageSettingChanged()
	l.CurrencySettingChanged()
	window.Reload()
}

// GetCurrentAppWidth returns the current width of the app's window.
func (l *Load) GetCurrentAppWidth() int {
	return l.CurrentAppWidth
}
