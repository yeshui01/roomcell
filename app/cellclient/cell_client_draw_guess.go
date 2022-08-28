package cellclient

import (
	"fmt"
	"image/color"
	"roomcell/pkg/pb/pbclient"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type DrawBoard struct {
	*canvas.Rectangle
	mainClient *CellClient
	gameWindow fyne.Window
	drawing    bool
	painData   *pbclient.DrawPainData
	playerList []*pbclient.RoomPlayer
	painObject []fyne.CanvasObject
	listWidget *widget.List
}

func (d *DrawBoard) MouseIn(me *desktop.MouseEvent) {
	fmt.Println("MouseIn")
}
func (d *DrawBoard) MouseMoved(me *desktop.MouseEvent) {
	if d.drawing {
		fmt.Printf("MouseMoved:%+v", me)
		oneCircle := canvas.NewCircle(color.White)
		oneCircle.StrokeWidth = 1
		oneCircle.Resize(fyne.NewSize(5, 5))
		oneCircle.Move(me.Position)
		mainContent := d.gameWindow.Content().(*fyne.Container)
		// newLine := canvas.NewLine(color.White)
		// newLine.Move(me.Position)
		mainContent.AddObject(oneCircle)
		//mainContent.AddObject(newLine)
		oneCircle.Refresh()
		mainContent.Refresh()
		d.painData.DrawPoints = append(d.painData.DrawPoints, &pbclient.PainPoint{
			PosX: int32(me.PointEvent.Position.X) * 1000,
			PosY: int32(me.PointEvent.Position.Y) * 1000,
		})
		d.painObject = append(d.painObject, oneCircle)
	}
	//oneCircle := canvas.NewCircle(color.Black)
	//oneCircle.StrokeColor = color.Gray{0x99}
}
func (d *DrawBoard) MouseOut() {
	fmt.Printf("MouseOut")
	//mainContent := d.gameWindow.Content().(*fyne.Container)
	//mainContent.RemoveAll()
	//mainContent.Refresh()
	d.drawing = false
}
func (d *DrawBoard) MouseDown(me *desktop.MouseEvent) {
	fmt.Printf("MouseDown")
	if me.Button == desktop.RightMouseButton {
		// 清除
		if len(d.painObject) > 0 {
			mainContent := d.gameWindow.Content().(*fyne.Container)
			for _, o := range d.painObject {
				mainContent.Remove(o)
			}
			mainContent.Refresh()
			d.painObject = nil
		}
	} else {
		d.drawing = true
		d.painData = &pbclient.DrawPainData{
			DrawPoints: make([]*pbclient.PainPoint, 0),
		}
	}
}
func (d *DrawBoard) MouseUp(me *desktop.MouseEvent) {
	fmt.Printf("MouseUp")
	d.drawing = false
	//d.mainClient.SendDrawPaint(d.painData)
	optCmd := &ClientCmd{
		CmdType: EClientCmdSyncPainData,
		CmdData: d.painData,
	}
	d.mainClient.hallClient.ClientCmdCh <- optCmd
	d.painData = nil
}

func (d *DrawBoard) DrawByPaindata(data *pbclient.DrawPainData) {
	fmt.Println("DrawByPaindata")
	for _, point := range data.DrawPoints {
		x := float32(point.PosX) / 1000
		y := float32(point.PosY) / 1000

		oneCircle := canvas.NewCircle(color.White)
		oneCircle.StrokeWidth = 1
		oneCircle.Resize(fyne.NewSize(5, 5))
		oneCircle.Move(fyne.NewPos(x, y))
		mainContent := d.gameWindow.Content().(*fyne.Container)
		mainContent.AddObject(oneCircle)
		oneCircle.Refresh()
		mainContent.Refresh()
	}
}

func NewDrawBoard() *DrawBoard {
	d := &DrawBoard{}
	d.Rectangle = canvas.NewRectangle(color.Transparent)
	return d
}

// 房间
type RoomDrawGuess struct {
	gameWindow fyne.Window
	mainClient *CellClient
	drawBoard  *DrawBoard
}

func NewRoomDrawGuess(c *CellClient) *RoomDrawGuess {
	room := &RoomDrawGuess{
		mainClient: c,
		drawBoard:  NewDrawBoard(),
	}
	room.drawBoard.mainClient = c
	room.gameWindow = c.fyneApp.NewWindow(fmt.Sprintf("你画我猜(%s)的游戏房间", c.roleData.Name))
	room.drawBoard.gameWindow = room.gameWindow
	//room.drawBoard.drawWindow = c.fyneApp.NewWindow("画板")
	room.drawBoard.Resize(fyne.NewSize(400, 600))
	//contentDraw := container.NewWithoutLayout(room.drawBoard)
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
	splitLine.Position1 = fyne.NewPos(590, 0)
	splitLine.Position2 = fyne.NewPos(590, 700)
	// 玩家列表
	playerLabel := widget.NewLabel("房间玩家:")
	playerLabel.Move(fyne.NewPos(600, 70))
	room.drawBoard.listWidget = widget.NewList(
		func() int {
			return len(room.drawBoard.playerList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(room.drawBoard.playerList[i].Nickname)
		},
	)
	room.drawBoard.listWidget.Move(fyne.NewPos(600, 100))
	room.drawBoard.listWidget.Resize(fyne.NewSize(100, 500))
	//content := container.New(layout.NewBorderLayout(nil, nil, nil, optList), roomName, optList, contentDraw)
	content := container.NewWithoutLayout(room.drawBoard, splitLine, optList, playerLabel, room.drawBoard.listWidget)
	room.gameWindow.SetContent(content)

	room.gameWindow.Resize(fyne.NewSize(700, 600))
	room.gameWindow.SetFixedSize(true) // 固定大小,禁止拉伸
	room.gameWindow.CenterOnScreen()
	room.gameWindow.Show()
	return room
}

func (gameRoom *RoomDrawGuess) RoomShow() {

}

func (gameRoom *RoomDrawGuess) AddRoomPlayer(roomPlayer *pbclient.RoomPlayer) {
	gameRoom.drawBoard.playerList = append(gameRoom.drawBoard.playerList, roomPlayer)
	gameRoom.drawBoard.listWidget.Refresh()
	//gameRoom.gameWindow.Content().Refresh()
}
func (gameRoom *RoomDrawGuess) DelRoomPlayer(roleID int64) {
	lenN := len(gameRoom.drawBoard.playerList)
	for k, v := range gameRoom.drawBoard.playerList {
		if v.RoleID == roleID {
			if k != lenN-1 {
				for j := k; j < lenN-1; j++ {
					gameRoom.drawBoard.playerList[j] = gameRoom.drawBoard.playerList[j+1]
				}
			}
			break
		}
	}
	gameRoom.drawBoard.playerList = gameRoom.drawBoard.playerList[0 : lenN-1]
	gameRoom.drawBoard.listWidget.Refresh()
}
