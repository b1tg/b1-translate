// -*- mode:go;mode:go-playground -*-
// snippet of code @ 2020-02-12 10:52:46

package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"

	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"

	"github.com/bregydoc/gtranslate"
	"os"
)

var res *theme.ThemedResource
var translateIcon = &fyne.StaticResource{
	StaticName:    "res.png",
	StaticContent: []byte{},
}

func init() {
	res = theme.NewThemedResource(translateIcon, nil)
	bi, _ := os.Open("translate.png")
	bis, _ := bi.Stat()
	bisize := bis.Size()
	sc := make([]byte, bisize)
	bi.Read(sc)
	translateIcon = &fyne.StaticResource{
		StaticName:    "translate.png",
		StaticContent: sc,
	}

}
func doTranslate(src string) string {

	translated, err := gtranslate.TranslateWithParams(

		src,
		gtranslate.TranslationParams{
			From: "en",
			To:   "zh-CN",
		},
	)
	if err != nil {
		panic(err)
	}

	return translated

}
func main() {
	app := app.New()
	w := app.NewWindow("b1-translate")
	w.Resize(fyne.Size{
		Width:  500,
		Height: 40,
	})
	w.SetIcon(translateIcon)
	ii := widget.NewEntry()
	output := widget.NewEntry()
	w.SetContent(widget.NewVBox(
		ii,
		output,
		widget.NewHBox(

			widget.NewButton("翻译输入框", func() {
				translated := doTranslate(ii.Text)
				output.SetText(translated)
			}),
			widget.NewButton("翻译粘贴板", func() {
				translated := doTranslate(w.Clipboard().Content())
				output.SetText(translated)
			}),
		),
	))
	ctrlTab := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: desktop.ControlModifier}
	w.Canvas().AddShortcut(&ctrlTab, func(shortcut fyne.Shortcut) {
		clipData := w.Clipboard().Content()
		ii.SetText(clipData)
		translated := doTranslate(clipData)
		output.SetText(translated)
	})

	w.ShowAndRun()
}
