package cellclient

import (
	"roomcell/pkg/sconst"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type GameItemData struct {
	GameName string
	RoomType int32
}

var gameItemList []*GameItemData = []*GameItemData{
	{GameName: "你画我猜", RoomType: sconst.EGameRoomTypeDrawGuess},
	{GameName: "谁是卧底", RoomType: sconst.EGameRoomTypeUndercover},
	// {GameName: "你画我猜3", RoomType: sconst.EGameRoomTypeDrawGuess},
	// {GameName: "你画我猜4", RoomType: sconst.EGameRoomTypeDrawGuess},
}

func (c *CellClient) InitGameListWindow() {
	//gameItemList = make([]*GameItemData, 0)
	if c.gameListWindow != nil {
		return
	}
	c.gameListWindow = c.fyneApp.NewWindow("游戏列表")
	c.gameListWindow.Resize(fyne.NewSize(400, 400))
	c.gameListWindow.CenterOnScreen()
	//gameListWindow.SetCon
	chooseList := widget.NewList(
		func() int {
			return len(gameItemList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(gameItemList[i].GameName)
		})
	chooseList.OnSelected = func(id int) {
		c.curGameChoose = gameItemList[id]
	}
	chooseList.OnUnselected = func(id int) {
		c.curGameChoose = nil
	}
	btnCreateRoom := widget.NewButton("创建房间", func() {
		c.SendCreateRoom()
	})
	btnCancelList := widget.NewButton("关闭", func() {
		c.gameListWindow.Hide()
	})
	c.gameListWindow.SetOnClosed(func() {
		c.gameListWindow = nil
	})
	// thisContent := container.NewVBox(chooseList, btnCreateRoom)
	thisContent := container.New(layout.NewBorderLayout(btnCancelList, btnCreateRoom, nil, nil), btnCancelList, btnCreateRoom, chooseList)
	c.gameListWindow.SetContent(thisContent)
}
