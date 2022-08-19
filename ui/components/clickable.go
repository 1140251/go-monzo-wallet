package components

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"go-monzo-wallet/ui/values"
	"image"
)

type Clickable struct {
	button    *widget.Clickable
	style     *ClickableStyle
	Hoverable bool
	Radius    Radius
	isEnabled bool
}

func (t *Theme) NewClickable(hoverable bool) *Clickable {
	return &Clickable{
		button:    &widget.Clickable{},
		style:     t.Styles.ClickableStyle,
		Hoverable: hoverable,
		isEnabled: true,
	}
}

func (cl *Clickable) Style() ClickableStyle {
	return *cl.style
}

func (cl *Clickable) ChangeStyle(style *ClickableStyle) {
	cl.style = style
}

func (cl *Clickable) Clicked() bool {
	return cl.button.Clicked()
}

func (cl *Clickable) IsHovered() bool {
	return cl.button.Hovered()
}

// SetEnabled enables/disables the clickable.
func (cl *Clickable) SetEnabled(enable bool, gtx *layout.Context) layout.Context {
	var mGtx layout.Context
	if gtx != nil && !enable {
		mGtx = gtx.Disabled()
	}

	cl.isEnabled = enable
	return mGtx
}

// Return clickable enabled/disabled state.
func (cl *Clickable) Enabled() bool {
	return cl.isEnabled
}

func (cl *Clickable) Layout(gtx values.C, w layout.Widget) values.D {
	return cl.button.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				tr := gtx.Dp(unit.Dp(cl.Radius.TopRight))
				tl := gtx.Dp(unit.Dp(cl.Radius.TopLeft))
				br := gtx.Dp(unit.Dp(cl.Radius.BottomRight))
				bl := gtx.Dp(unit.Dp(cl.Radius.BottomLeft))
				defer clip.RRect{
					Rect: image.Rectangle{Max: image.Point{
						X: gtx.Constraints.Min.X,
						Y: gtx.Constraints.Min.Y,
					}},
					NW: tl, NE: tr, SE: br, SW: bl,
				}.Push(gtx.Ops).Pop()
				clip.Rect{Max: gtx.Constraints.Min}.Push(gtx.Ops).Pop()

				if cl.Hoverable && cl.button.Hovered() {
					paint.Fill(gtx.Ops, cl.style.HoverColor)
				}

				for _, c := range cl.button.History() {
					DrawInk(gtx, c, cl.style.Color)
				}
				return layout.Dimensions{Size: gtx.Constraints.Min}
			}),
			layout.Stacked(w),
		)
	})
}
