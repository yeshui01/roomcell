package cellclient

import (
	"encoding/json"
	"fmt"
	"image/color"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PersonInfo struct {
	ID    int32
	Name  string
	Age   int32
	Sex   int32
	Extra string
}

func runBorderLayout(t *testing.T) {
	myApp := fyneapp.New()
	myWindow := myApp.NewWindow("Border Layout")
	myApp.Settings().SetTheme(theme.DarkTheme())
	top := canvas.NewText("top", color.White)
	left := canvas.NewText("left", color.White)
	bottom := canvas.NewText("bottom", color.White)
	right := canvas.NewText("right", color.White)

	// middle := canvas.NewText("middle", color.White)
	// 多行文本
	txtMulti := widget.NewTextGrid()
	//txtResults.ShowLineNumbers = true
	//txtResults.SetText(strings.TrimPrefix("msg1"+"\n"+"msg2", "\n"))
	ctnScroll := container.NewScroll(txtMulti)
	go func() {
		personId := int32(1)
		age := int32(10)
		sex := int32(1)
		secTick := time.NewTicker(time.Second * 1)
		totalText := "" // 所有的文本,\n分割
		for {
			select {
			case <-secTick.C:
				{
					// 添加一行内容
					personInfo := &PersonInfo{
						ID:    personId,
						Name:  fmt.Sprintf("person:%d", personId),
						Age:   age,
						Sex:   sex,
						Extra: "xtraaaaaaaaaaaaaaaaaaaa",
					}
					jd, errj := json.Marshal(personInfo)
					// jd, errj := json.MarshalIndent(personInfo, "", " ")
					if errj == nil {
						totalText = totalText + string(jd) + "\n"
						if len(txtMulti.Rows) >= 10 {
							idx := strings.Index(totalText, "\n")
							if idx != -1 {
								totalText = totalText[idx+1:]
							}
						}
						txtMulti.SetText(totalText)
					}

					ctnScroll.ScrollToBottom()
					// txtMulti.Refresh()
					ctnScroll.Refresh()
					personId++
					break
				}
			}
		}
	}()

	content := container.New(layout.NewBorderLayout(top, bottom, left, right),
		top, left, bottom, right, ctnScroll)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(600, 480))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
