// -*- mode:go;mode:go-playground -*-
// snippet of code @ 2020-02-12 10:52:46

package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"io/ioutil"
	"net/http"
	"strings"

	//"encoding/json"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/theme"

	//	"github.com/bregydoc/gtranslate"
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

// {"translation":["窗口"],"basic":{"us-phonetic":"ˈwɪndoʊ","phonetic":"ˈwɪndəʊ","uk-phonetic":"ˈwɪndəʊ","explains":["n. 窗；窗口；窗户"]},"query":"window","errorCode":0,"web":[{"value":["视窗","窗","窗口","窗户"],"key":"window"},{"value":["后窗","后窗玻璃"],"key":"Rear Window"},{"value":["凸窗","飘窗","窗台","八角窗"],"key":"bay window"}]}

type Youdao struct {
}

func doTranslate(src string) string {

	if src == "" {
		return ""
	}
	// curl http://fanyi.youdao.com/openapi.do\?keyfrom\=YouDaoCV\&key\=659600698\&type\=data\&doctype\=json\&version\=1.1\&q\=compose
	uuu := fmt.Sprintf("http://fanyi.youdao.com/openapi.do?keyfrom=YouDaoCV&key=659600698&type=data&doctype=json&version=1.1&q=%s", src)
	resp, err := http.Get(uuu)
	if err != nil {
		return ""
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
		// handle error
	}
	println(string(body))
	return string(body)

	// translated, err := gtranslate.TranslateWithParams(

	// 	src,
	// 	gtranslate.TranslationParams{
	// 		From: "en",
	// 		To:   "zh-CN",
	// 	},
	// )
	// if err != nil {
	// 	panic(err)
	// }

	//return translated

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
	// output := widget.NewEntry()
	output := widget.NewCard("Title", "sub", widget.NewLabel("Content"))
	w.SetContent(widget.NewVBox(
		ii,
		output,
		widget.NewHBox(

			widget.NewButton("翻译输入框", func() {
				translated := doTranslate(ii.Text)
				//output.SetText(translated)
				// translated[70]= "\n"
				translated = strings.Replace(translated, ",", "\n", -1)
				newContent := widget.NewLabel(translated)
				output.SetContent(newContent)
				//output.k
			}),
			widget.NewButton("翻译粘贴板", func() {
				translated := doTranslate(w.Clipboard().Content())
				//	output.SetText(translated)

				print(translated)
			}),
		),
	))
	//ctrlTab := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: desktop.ControlModifier}
	ctrlTab := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlTab, func(shortcut fyne.Shortcut) {
		println("got key")
		clipData := w.Clipboard().Content()
		ii.SetText(clipData)
		translated := doTranslate(clipData)
		if translated == "" {
			return
		}
		// print(translated)
		//output.SetText(translated)
		translated = strings.Replace(translated, ",", "\n", -1)
		newContent := widget.NewLabel(translated)
		output.SetContent(newContent)
	})

	w.ShowAndRun()
}
