package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/mccutchen/palettor"
	"github.com/nfnt/resize"
)

const (
	imgWidth   = 180
	imgHeight  = 130
	cellWidth  = 200
	cellHeight = 260
)

type ImageBox struct {
	fyne.Container
	Image   *canvas.Image
	Palette *fyne.Container
}

func (g *ImageBox) Cursor() desktop.Cursor {
	return desktop.CrosshairCursor
}

func NewImageBox(path string, onColorSelect func(color color.Color)) *ImageBox {

	source, err := getImageFromFilePath(path)

	// Reduce it to a manageable size
	i := resize.Thumbnail(imgWidth, imgHeight, source, resize.Lanczos3)

	img := NewImage(i, onColorSelect)

	// Extract the 6 most dominant colors, halting the clustering algorithm
	// after 100 iterations if the clusters have not yet converged.
	palette, err := palettor.Extract(5, 100, i)
	if err != nil {
		log.Println("cannot extract palette:", err)
	}

	size := fyne.NewSize(imgWidth, imgHeight)
	palbox := createPaletteBox(palette.Colors(), onColorSelect)
	margin := widget.NewLabel("")

	center := container.NewVBox(img, palbox, margin)
	center.Refresh()

	container := &ImageBox{
		Image:   &img.Image,
		Palette: palbox,
		Container: fyne.Container{
			Objects: []fyne.CanvasObject{layout.NewSpacer(), center, layout.NewSpacer()},
			Layout:  layout.NewHBoxLayout(),
		},
	}
	container.Resize(size)
	return container
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func createPaletteBox(colors []color.Color, onTap func(color.Color)) *fyne.Container {
	palbox := container.NewHBox()

	palbox.Add(layout.NewSpacer())
	for _, color := range colors {
		c := NewColorButton(color, onTap)
		palbox.Add(c)
	}
	palbox.Add(layout.NewSpacer())

	return palbox
}
