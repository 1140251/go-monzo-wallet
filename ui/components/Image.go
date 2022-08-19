package components

import (
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	values2 "go-monzo-wallet/ui/values"
	"image"
)

type Image struct {
	*widget.Image
}

func NewImage(src image.Image) *Image {
	return &Image{
		Image: &widget.Image{
			Src: paint.NewImageOp(src),
		},
	}
}

func (img *Image) Layout(gtx values2.C) values2.D {
	return img.Image.Layout(gtx)
}

func (img *Image) Layout12dp(gtx values2.C) values2.D {
	return img.LayoutSize(gtx, values2.MarginPadding12)
}

func (img *Image) Layout16dp(gtx values2.C) values2.D {
	return img.LayoutSize(gtx, values2.MarginPadding16)
}

func (img *Image) Layout24dp(gtx values2.C) values2.D {
	return img.LayoutSize(gtx, values2.MarginPadding24)
}

func (img *Image) Layout36dp(gtx values2.C) values2.D {
	return img.LayoutSize(gtx, values2.MarginPadding36)
}

func (img *Image) Layout48dp(gtx values2.C) values2.D {
	return img.LayoutSize(gtx, values2.MarginPadding48)
}

func (img *Image) LayoutSize(gtx values2.C, size unit.Dp) values2.D {
	width := img.Src.Size().X
	scale := float32(size) / float32(width)
	img.Scale = scale
	return img.Layout(gtx)
}
