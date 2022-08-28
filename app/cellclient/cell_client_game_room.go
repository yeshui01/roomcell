package cellclient

import (
	"fmt"
	"image/color"
	"roomcell/pkg/pb/pbclient"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 通用游戏房间
// 房间
type GameRoomBase struct {
	gameWindow fyne.Window
	mainClient *CellClient
	listWidget *widget.List
	playerList []*pbclient.RoomPlayer
	chatEntry  *widget.Entry
}

func NewGameRoomBase(c *CellClient, roomID int64, gameName string) *GameRoomBase {
	room := &GameRoomBase{
		mainClient: c,
	}
	room.mainClient = c
	room.gameWindow = c.fyneApp.NewWindow(fmt.Sprintf("%s(%s)的游戏房间(%d)", gameName, c.roleData.Name, roomID))
	btnReady := widget.NewButton("准备", func() {
		room.mainClient.refreshLogContent("ready game")
		readyCmd := &ClientCmd{
			CmdType: EClientCmdDrawReady,
			CmdData: true,
		}
		room.mainClient.hallClient.ClientCmdCh <- readyCmd
	})
	btnReadyCancel := widget.NewButton("取消准备", func() {
		room.mainClient.refreshLogContent("ready game cancel")
		readyCmd := &ClientCmd{
			CmdType: EClientCmdDrawReady,
			CmdData: false,
		}
		room.mainClient.hallClient.ClientCmdCh <- readyCmd
	})
	optList := container.NewVBox(btnReady, btnReadyCancel)
	optList.Move(fyne.NewPos(600, 0))
	// 加一条分割线
	splitLine := canvas.NewLine(color.White)
	splitLine.Position1 = fyne.NewPos(550, 0)
	splitLine.Position2 = fyne.NewPos(550, 700)
	// 玩家列表
	playerLabel := widget.NewLabel("房间玩家:")
	playerLabel.Move(fyne.NewPos(560, 70))
	room.listWidget = widget.NewList(
		func() int {
			return len(room.playerList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(room.playerList[i].Nickname)
		},
	)
	room.listWidget.Move(fyne.NewPos(550, 100))
	room.listWidget.Resize(fyne.NewSize(100, 500))
	//content := container.New(layout.NewBorderLayout(nil, nil, nil, optList), roomName, optList, contentDraw)
	// 加一个聊天框
	room.chatEntry = widget.NewEntry()
	room.chatEntry.SetPlaceHolder("input chat message")
	room.chatEntry.Move(fyne.NewPos(0, 0))
	room.chatEntry.Resize(fyne.NewSize(350, 50))
	// 发送按钮
	btnChat := widget.NewButton("发送消息", func() {
		if len(room.chatEntry.Text) > 0 {
			room.mainClient.refreshLogContent("chat message")
			readyCmd := &ClientCmd{
				CmdType: EClientCmdChatMessage,
				CmdData: room.chatEntry.Text,
			}
			room.mainClient.hallClient.ClientCmdCh <- readyCmd
			room.chatEntry.SetText("")
		}
	})
	btnChat.Move(fyne.NewPos(360, 0))
	btnChat.Resize(fyne.NewSize(40, 30))
	content := container.NewWithoutLayout(room.chatEntry, splitLine, optList, playerLabel, room.listWidget)
	room.gameWindow.SetContent(content)

	room.gameWindow.Resize(fyne.NewSize(700, 600))
	room.gameWindow.SetFixedSize(true) // 固定大小,禁止拉伸
	room.gameWindow.CenterOnScreen()
	room.gameWindow.Show()
	return room
}

func (gameRoom *GameRoomBase) RoomShow() {

}

func (gameRoom *GameRoomBase) AddRoomPlayer(roomPlayer *pbclient.RoomPlayer) {
	gameRoom.playerList = append(gameRoom.playerList, roomPlayer)
	gameRoom.listWidget.Refresh()
	//gameRoom.gameWindow.Content().Refresh()
}
func (gameRoom *GameRoomBase) DelRoomPlayer(roleID int64) {
	lenN := len(gameRoom.playerList)
	for k, v := range gameRoom.playerList {
		if v.RoleID == roleID {
			if k != lenN-1 {
				for j := k; j < lenN-1; j++ {
					gameRoom.playerList[j] = gameRoom.playerList[j+1]
				}
			}
			break
		}
	}
	gameRoom.playerList = gameRoom.playerList[0 : lenN-1]
	gameRoom.listWidget.Refresh()
}
