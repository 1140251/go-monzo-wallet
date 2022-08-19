package components

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"go-monzo-wallet/ui/values"
)

func UniformPadding(gtx layout.Context, body layout.Widget) layout.Dimensions {
	width := gtx.Constraints.Max.X

	padding := values.MarginPadding24

	if (width - 2*gtx.Dp(padding)) > gtx.Dp(MaxWidth) {
		paddingValue := float32(width-gtx.Dp(MaxWidth)) / 2
		padding = unit.Dp(paddingValue)
	}

	return layout.Inset{
		Top:    values.MarginPadding24,
		Right:  padding,
		Bottom: values.MarginPadding24,
		Left:   padding,
	}.Layout(gtx, body)
}

func UniformMobile(gtx layout.Context, isHorizontal, withList bool, body layout.Widget) layout.Dimensions {
	insetRight := values.MarginPadding10
	if withList {
		insetRight = values.MarginPadding0
	}

	insetTop := values.MarginPadding24
	if isHorizontal {
		insetTop = values.MarginPadding0
	}

	return layout.Inset{
		Top:   insetTop,
		Right: insetRight,
		Left:  values.MarginPadding10,
	}.Layout(gtx, body)
}
