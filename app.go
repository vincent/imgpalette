package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vincent/impallette/internal"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

const (
	startCellWidth = 200
	startWinWidth  = 600
	startWinHeight = 250
)

type App struct {
	colorFormat string
	app         fyne.App
	win         fyne.Window
	grid        *fyne.Container
	lastColors  []color.Color
	last        *fyne.Container
	bar         *widget.ProgressBar
	sel         *widget.Select
	knownFiles  map[string]bool
}

func NewApp() *App {
	a := app.NewWithID("com.impallette")
	w := a.NewWindow("ColorPick")

	grid := container.New(layout.NewGridLayout(3))
	self := &App{
		app:        a,
		grid:       grid,
		win:        w,
		knownFiles: make(map[string]bool, 100),
	}

	toolbar := self.toolbar(self.onLoadImage)
	scroll := container.NewVScroll(container.NewVBox(grid, layout.NewSpacer()))
	scroll.SetMinSize(fyne.NewSize(680, 200))

	c := container.New(layout.NewBorderLayout(nil, toolbar, nil, nil), scroll, toolbar)
	w.SetContent(c)

	for _, f := range strings.Split(a.Preferences().StringWithFallback("files", ""), "$") {
		if f != "" {
			self.knownFiles[f] = true
		}
	}

	return self
}

func (a *App) toolbar(onAdd func()) *fyne.Container {
	openButton := widget.NewButton("Ajouter", onAdd)

	a.bar = widget.NewProgressBar()
	a.bar.Min = 0
	a.bar.Max = 100
	a.setBarProgress(0)

	a.last = container.NewHBox()

	a.sel = widget.NewSelect([]string{"rgb", "hex", "hsl"}, a.onFormatSelect)
	a.sel.SetSelected("rgb")

	hbox := container.NewHBox(openButton, a.bar, layout.NewSpacer(), a.last, layout.NewSpacer(), a.sel)
	hbox.Refresh()

	return hbox
}

func (a *App) onLoadImage() {
	fd := dialog.NewFileOpen(a.onFileSelect, a.win)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}))
	fd.Show()
}

func (a *App) onFileSelect(reader fyne.URIReadCloser, err error) {
	if err != nil || reader == nil || reader.URI() == nil {
		return
	}
	a.addImageFiles([]string{reader.URI().Path()})
}

func (a *App) onFormatSelect(value string) {
	a.colorFormat = value
}

func (a *App) run() {
	a.win.ShowAndRun()
}

func (a *App) loadSavedImageFiles() {
	a.addImageFiles(strings.Split(a.app.Preferences().StringWithFallback("files", ""), "$"))
}

func (a *App) saveKnownImageFiles() {
	known := []string{}
	for f, _ := range a.knownFiles {
		known = append(known, f)
	}
	a.app.Preferences().SetString("files", strings.Join(known, "$"))
}

func (a *App) addImageFiles(paths []string) {
	count := float64(len(paths))
	wg := sync.WaitGroup{}

	a.setBarProgress(0)
	for i, path := range paths {

		wg.Add(1)
		go func(i int, path string) {
			defer wg.Done()
			file, err := os.Open(path)
			if errors.Is(err, os.ErrNotExist) {
				log.Println(err.Error())
				return
			}
			defer file.Close()
			imgbox := NewImageBox(path, a.onColorSelect)
			a.grid.Add(&imgbox.Container)

			if _, ok := a.knownFiles[path]; !ok {
				a.knownFiles[path] = true
			}

			a.setBarProgress(100 / count)
			a.grid.Refresh()
		}(i, path)
	}

	wg.Wait()
	a.saveKnownImageFiles()
	a.setBarProgress(100)
}

func (a *App) saveColor(col color.Color) {
	rgb, _ := internal.ColorToRGB(col)
	html := rgb.ToHTML()
	for _, c := range a.lastColors {
		crgb, _ := internal.ColorToRGB(c)
		chtml := crgb.ToHTML()
		if chtml == html {
			return
		}
	}

	if len(a.lastColors) >= 10 {
		a.lastColors = a.lastColors[1:]
	}

	a.lastColors = append(a.lastColors, col)

	objects := []fyne.CanvasObject{}
	a.last.Add(layout.NewSpacer())
	for _, c := range a.lastColors {
		b := NewColorButton(c, a.onColorSelect)
		objects = append(objects, b)
	}

	a.last.Objects = objects
	a.last.Resize(fyne.NewSize(150, 35))
	a.last.Refresh()
}

func (a *App) onColorSelect(col color.Color) {
	a.saveColor(col)
	a.win.Clipboard().SetContent(a.toColorFormat(col))
}

func (a *App) toColorFormat(col color.Color) string {
	switch a.colorFormat {
	case "rgb":
		r, g, b, _ := col.RGBA()
		return fmt.Sprintf("%v, %v, %v", r>>8, g>>8, b>>8)
	case "hex":
		c, _ := internal.ColorToRGB(col)
		return c.ToHTML()
	case "hsl":
		c, _ := internal.ColorToRGB(col)
		h := c.ToHSL()
		return fmt.Sprintf("%f, %f, %f", h.H, h.S, h.L)
	}
	return ""
}

func (a *App) setBarProgress(value float64) {
	time.Sleep(20 * time.Millisecond)
	a.bar.SetValue(value + a.bar.Value)
	a.bar.Resize(fyne.NewSize(150, 35))
	a.bar.Show()

	if a.bar.Value >= 99.9 {
		a.bar.SetValue(0)
		a.bar.Hide()
	}
}
