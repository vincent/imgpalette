package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
)

type Image struct {
	canvas.Image
	onTap func(color color.Color)
}

func (g *Image) Tapped(a *fyne.PointEvent) {
	col := g.Image.Image.At(int(a.Position.X), int(a.Position.Y))
	r, _, _, _ := col.RGBA()
	if r > 200 {
		r = r * 1
	}
	g.onTap(col)
}

func NewImage(i image.Image, onTap func(color.Color)) *Image {
	ret := &Image{
		onTap: onTap,
		Image: *canvas.NewImageFromImage(i),
	}

	ret.Image.SetMinSize(fyne.NewSize(200, 150))
	ret.Image.FillMode = canvas.ImageFillContain

	return ret
}

func (g *Image) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

func (g *Image) CreateRenderer() fyne.WidgetRenderer {
	return &ImageRenderer{colorButton: g}
}

type ImageRenderer struct {
	colorButton *Image
}

func (g *ImageRenderer) BackgroundColor() color.Color {
	return color.Transparent
}

func (g *ImageRenderer) Destroy() {
}

func (g *ImageRenderer) Layout(size fyne.Size) {
	g.colorButton.Image.Resize(size)
}

func (g *ImageRenderer) MinSize() fyne.Size {
	return g.colorButton.MinSize()
}

func (g *ImageRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{&g.colorButton.Image}
}

func (g *ImageRenderer) Refresh() {
	g.colorButton.Image.Refresh()
}
