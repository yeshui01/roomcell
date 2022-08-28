package cellclient

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/cast"
)

func (c *CellClient) InitJoinRoomWindow() {
	//gameItemList = make([]*GameItemData, 0)
	if c.joinRoomWindow != nil {
		return
	}
	c.joinRoomWindow = c.fyneApp.NewWindow("加入房间")
	c.joinRoomWindow.Resize(fyne.NewSize(300, 150))
	c.joinRoomWindow.CenterOnScreen()
	// chooseList := widget.NewList(
	// 	func() int {
	// 		return len(gameItemList)
	// 	},
	// 	func() fyne.CanvasObject {
	// 		return widget.NewLabel("template")
	// 	},
	// 	func(i widget.ListItemID, o fyne.CanvasObject) {
	// 		o.(*widget.Label).SetText(gameItemList[i].GameName)
	// 	})
	// chooseList.OnSelected = func(id int) {
	// 	c.curGameChoose = gameItemList[id]
	// }
	// chooseList.OnUnselected = func(id int) {
	// 	c.curGameChoose = nil
	// }

	title := widget.NewLabel("输入房间号:")
	roomInput := widget.NewEntry()
	roomInput.SetPlaceHolder("这里输入房间号...")

	btnJoinRoom := widget.NewButton("加入", func() {
		roomID := cast.ToInt64(roomInput.Text)
		doCmd := &ClientCmd{
			CmdType: EClientCmdEnterRoom,
			CmdData: roomID,
		}
		c.hallClient.ClientCmdCh <- doCmd
	})
	btnCancelList := widget.NewButton("关闭", func() {
		roomInput.SetText("")
		c.joinRoomWindow.Hide()
	})
	c.joinRoomWindow.SetOnClosed(func() {
		c.joinRoomWindow = nil
	})
	thisContent := container.NewVBox(btnCancelList, title, roomInput, btnJoinRoom)
	//thisContent := container.New(layout.NewBorderLayout(btnCancelList, btnCreateRoom, nil, nil), btnCancelList, btnCreateRoom, chooseList)
	c.joinRoomWindow.SetContent(thisContent)
}
