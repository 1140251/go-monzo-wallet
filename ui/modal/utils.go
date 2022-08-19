package modal

import (
	"gioui.org/layout"
	"gioui.org/widget"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func editorsNotEmpty(editors ...*widget.Editor) bool {
	for _, e := range editors {
		if e.Text() == "" {
			return false
		}
	}
	return true
}
