package components

import (
	"gioui.org/unit"
	"gioui.org/widget"
	"go-monzo-wallet/ui/values"
	"image/color"
	_ "image/png"
)

type Icon struct {
	*widget.Icon
	Color color.NRGBA
}

// NewIcon returns a new Icon from IconVG data.
func NewIcon(icon *widget.Icon) *Icon {
	return &Icon{
		Icon: icon,
	}
}

func (icon *Icon) Layout(gtx values.C, iconSize unit.Dp) values.D {
	cl := color.NRGBA{A: 0xff}
	if icon.Color != (color.NRGBA{}) {
		cl = icon.Color
	}
	gtx.Constraints.Min.X = gtx.Dp(iconSize)
	return icon.Icon.Layout(gtx, cl)
}
