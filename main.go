// -*- mode:go;mode:go-playground -*-
// snippet of code @ 2020-02-12 10:52:46

package main

import (
	"errors"
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
	"github.com/bregydoc/gtranslate"
	"github.com/m7shapan/njson"
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

type WebExplain struct {
	Value []string `njson:"value"`
	Key   string   `njson:"key"`
}
type Youdao struct {
	BasicExplains []string     `njson:"basic.explains"`
	WebExplains   []WebExplain `njson:"web"`
	ErrorCode     int          `njson:"errorCode"`
}

func doTranslate(src string, engine string) (string, error) {
	if engine == EngineYoudao {
		yd, err := youdaoTranlsate(src)
		return yd, err
	} else if engine == EngineGoogle {
		gg, err := googleTranlsate(src)
		return gg, err
	} else {
		//return fmt.Sprintf("google:\n%s\n\nyoudao:\n%s\n", gg, yd)
		return "", errors.New("Engine not found")
	}
}
func youdaoTranlsate(src string) (string, error) {

	if src == "" {
		return "", errors.New("Empty input")
	}
	// curl http://fanyi.youdao.com/openapi.do\?keyfrom\=YouDaoCV\&key\=659600698\&type\=data\&doctype\=json\&version\=1.1\&q\=compose
	uuu := fmt.Sprintf("http://fanyi.youdao.com/openapi.do?keyfrom=YouDaoCV&key=659600698&type=data&doctype=json&version=1.1&q=%s", src)
	resp, err := http.Get(uuu)
	if err != nil {
		return "", fmt.Errorf("http get error: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("http read body error: %s", err)
	}

	println(string(body))
	youdao := Youdao{}
	err = njson.Unmarshal(body, &youdao)
	if err != nil {
		return "", fmt.Errorf("json unmarshal error: %s", err)
	}
	if youdao.ErrorCode > 0 {
		return "", fmt.Errorf("server error: %d", youdao.ErrorCode)
	}

	translateResult := youdao.BasicExplains
	if len(youdao.BasicExplains) == 0 && len(youdao.WebExplains) > 0 {
		translateResult = youdao.WebExplains[0].Value
	}

	return strings.Join(translateResult, "\n"), nil

}

func googleTranlsate(input string) (string, error) {
	translated, err := gtranslate.TranslateWithParams(
		input,
		gtranslate.TranslationParams{
			From: "en",
			To:   "zh-CN",
		},
	)
	if err != nil {
		return "", fmt.Errorf("gtranslate error: %s", err)
	}
	return translated, nil
}

const EngineYoudao = "Youdao"
const EngineGoogle = "Google"

func main() {
	app := app.New()
	w := app.NewWindow("b1-translate")
	w.Resize(fyne.Size{
		Width:  500,
		Height: 40,
	})
	w.SetIcon(translateIcon)
	radio := widget.NewRadioGroup([]string{EngineYoudao, EngineGoogle}, nil)
	radio.Selected = EngineYoudao
	ii := widget.NewEntry()
	output := widget.NewEntry()
	status := widget.NewLabel("...")
	w.SetContent(widget.NewVBox(
		radio,
		ii,
		output,
		widget.NewHBox(

			widget.NewButton("翻译输入框", func() {
				//translated := doTranslate(ii.Text)
				////output.SetText(translated)
				//// translated[70]= "\n"
				//translated1 := strings.Join(translated, "\n")
				//newContent := widget.NewLabel(translated1)
				//output.SetContent(newContent)
				//output.k
			}),
			widget.NewButton("翻译粘贴板", func() {
				//translated := doTranslate(w.Clipboard().Content())
				////	output.SetText(translated)
				//translated1 := strings.Join(translated, "\n")
				//newContent := widget.NewLabel(translated1)
				//output.SetContent(newContent)

				// print(translated)
			}),
			status,
		),
	))
	//ctrlTab := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: desktop.ControlModifier}
	ctrlTab := desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: desktop.SuperModifier}
	w.Canvas().AddShortcut(&ctrlTab, func(shortcut fyne.Shortcut) {
		println("get key")
		clipData := w.Clipboard().Content()
		ii.SetText(clipData)
		translated, err := doTranslate(clipData, radio.Selected)
		if err != nil {
			errMsg := err.Error()
			limit := 20
			if len(errMsg) < limit {
				limit = len(errMsg)
			}
			status.SetText(err.Error()[0:limit])

			output.SetText("")
		} else {
			output.SetText(translated)
			status.SetText("ok")
		}
	})

	w.ShowAndRun()
}
