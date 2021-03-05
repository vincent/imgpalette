package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ColorButton struct {
	widget.BaseWidget
	onTap func(color color.Color)
	Rect  *canvas.Rectangle
}

func (g *ColorButton) Tapped(a *fyne.PointEvent) {
	g.onTap(g.Rect.FillColor)
}

func NewColorButton(color color.Color, onTap func(color.Color)) *ColorButton {
	ret := &ColorButton{onTap: onTap}
	ret.ExtendBaseWidget(ret)
	ret.Rect = canvas.NewRectangle(color)
	ret.Rect.SetMinSize(fyne.NewSize(40, 40))
	return ret
}

func (g *ColorButton) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

func (g *ColorButton) CreateRenderer() fyne.WidgetRenderer {
	return &colorButtonRenderer{colorButton: g}
}

type colorButtonRenderer struct {
	colorButton *ColorButton
}

func (g *colorButtonRenderer) BackgroundColor() color.Color {
	return g.colorButton.Rect.FillColor
}

func (g *colorButtonRenderer) Destroy() {
}

func (g *colorButtonRenderer) Layout(size fyne.Size) {
	g.colorButton.Rect.Resize(size)
}

func (g *colorButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(40, 40)
}

func (g *colorButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{g.colorButton.Rect}
}

func (g *colorButtonRenderer) Refresh() {
	g.colorButton.Rect.Refresh()
}
