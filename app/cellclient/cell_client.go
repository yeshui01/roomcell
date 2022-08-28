package cellclient

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"roomcell/app/account/accrouter"
	"roomcell/pkg/appconfig"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/webreq"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
	"github.com/sirupsen/logrus"
)

// 引用中文
func init() {
	fontPaths := findfont.List()
	for i, fpath := range fontPaths {
		// 楷体: simkai
		// 黑体: simhei
		if strings.Contains(fpath, "simkai.ttf") {
			os.Setenv("FYNE_FONT", fpath)
			fmt.Println(i)
			break
		}
	}
}

type IGameRoomWindow interface {
	RoomShow()
	DelRoomPlayer(roleID int64)
	AddRoomPlayer(p *pbclient.RoomPlayer)
}

type CellClient struct {
	clientCfg *appconfig.CellClientCfg
	fyneApp   fyne.App
	// logWidget *widget.TextGrid
	logWidget *widget.Label
	logScroll *container.Scroll
	totalLog  string
	loginForm *widget.Form

	// 右边状态栏
	curAccountLabel *widget.Label
	// 子窗口显示
	subWindowShow *widget.Label
	accInput      *widget.Entry
	pswdInput     *widget.Entry
	// 大厅客户端
	hallClient   *HallClient
	loginToken   string
	hallConnAddr string
	SeqID        int32
	SeqlLock     sync.Mutex

	roleData *pbclient.RoleInfo

	// 子窗口
	gameListWindow fyne.Window
	joinRoomWindow fyne.Window
	curGameChoose  *GameItemData
	gameRoom       IGameRoomWindow
}

func NewCellClient(configPath string) *CellClient {
	c := &CellClient{}
	cfg, err := appconfig.ReadCellClientConfigFromFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("read config error:%s", err.Error()))
	}
	c.clientCfg = cfg
	c.hallClient = NewHallClient()
	//c.gameListWindow = fyneApp.NewWindow()
	return c
}

func (c *CellClient) refreshLogContent(newContent string) {
	c.totalLog = c.totalLog + newContent + "\n"
	for {
		if strings.Count(c.totalLog, "\n") >= 30 {
			idx := strings.Index(c.totalLog, "\n")
			if idx != -1 {
				c.totalLog = c.totalLog[idx+1:]
			}
		} else {
			break
		}
	}

	c.logWidget.SetText(c.totalLog)
	c.logScroll.ScrollToBottom()
	//c.logWidget.SetText(newContent)

	// txtMulti.Refresh()
	c.logScroll.Refresh()
}

// 登录账号
func (c *CellClient) LoginAccount() {
	accHttpAddr := fmt.Sprintf("%s/account/login", c.clientCfg.AccountAddr)
	req := &accrouter.AccountLoginReq{
		UserName: c.accInput.Text,
		Pswd:     c.pswdInput.Text,
	}
	// reqData, err := json.Marshal(req)
	// if err != nil {
	// 	c.refreshLogContent(fmt.Sprintf("login account error:%s", err.Error()))
	// 	return
	// }
	err2, repData := webreq.PostJson(accHttpAddr, req)
	if err2 != nil {
		c.refreshLogContent(fmt.Sprintf("login account result error:%s", err2.Error()))
		return
	}
	rsp := &AccountLoginRsp{}
	strData := string(repData)
	c.refreshLogContent(fmt.Sprintf("login ret: %s", strData))
	err3 := json.Unmarshal(repData, rsp)
	if err3 != nil {
		c.refreshLogContent(fmt.Sprintf("unmarshal rsp error: %s", err3.Error()))
		return
	}
	if rsp.Code == 0 {
		c.loginToken = rsp.Data.Token
		c.hallConnAddr = rsp.Data.HallAddr
		// 连接
		errConn := c.hallClient.Connect(c.hallConnAddr)
		if errConn != nil {
			c.refreshLogContent(fmt.Sprintf("connect hall error:%s", errConn.Error()))
		} else {
			c.refreshLogContent(fmt.Sprintf("login hall sucess:%s", c.hallConnAddr))
			go func() {
				c.hallClient.RunWithWindow(c)
			}()
		}
		logResult, _ := json.MarshalIndent(rsp, "", "\t")
		c.refreshLogContent(fmt.Sprintf("login result: %s", string(logResult)))
	} else {
		c.refreshLogContent(fmt.Sprintf("login result error: %d", rsp.Code))
	}

}
func (c *CellClient) GenSeqID() int32 {
	c.SeqlLock.Lock()
	c.SeqID++
	if c.SeqID >= 999999999 {
		c.SeqID = 1
	}
	c.SeqlLock.Unlock()
	return c.SeqID
}

func (c *CellClient) Run() {
	loghlp.SetConsoleLogLevel(logrus.Level(c.clientCfg.LogLevel))
	myApp := fyneapp.New()
	c.fyneApp = myApp
	myWindow := myApp.NewWindow("模拟客户端")

	myApp.Settings().SetTheme(theme.DarkTheme())
	//c.InitGameListWindow()
	// subWindow := myApp.NewWindow("sub_window")
	// c.subWindowShow = widget.NewLabel("this is a sub window")
	// subWindow.SetContent(c.subWindowShow)
	// //subWindow.Show()
	// subWindow.Resize(fyne.NewSize(500, 400))
	//subWindow.CenterOnScreen()
	// 账号
	c.accInput = widget.NewEntry()
	c.accInput.SetPlaceHolder("请输入账号")
	// 密码
	c.pswdInput = widget.NewEntry()
	c.pswdInput.Password = true
	bottom := canvas.NewText("", color.White)

	//c.logWidget = widget.NewTextGrid()
	c.logWidget = widget.NewLabel("watch log output")
	c.logScroll = container.NewScroll(c.logWidget)
	c.logWidget.SetText("watch log output")
	c.loginForm = widget.NewForm(
		widget.NewFormItem("账号", c.accInput),
		widget.NewFormItem("密码", c.pswdInput),
	)
	c.loginForm.SubmitText = "登录"
	c.loginForm.OnSubmit = func() { // optional, handle form submission
		loghlp.Info("Form submitted:", c.accInput.Text)
		loghlp.Info("multiline:", c.pswdInput.Text)
		if c.accInput.Disabled() {
			loghlp.Info("accIput disabled!!")
			return
		}
		if len(c.accInput.Text) < 1 || len(c.pswdInput.Text) < 1 {
			loghlp.Error("账号密码格式错误!!!")
			c.refreshLogContent("account or pswd error!!!")
			return
		}
		// myWindow.Close()
		// 连接账号服务器TODO:
		c.LoginAccount()
		c.curAccountLabel.SetText(fmt.Sprintf("当前账号:%s", c.accInput.Text))
		//c.loginForm.Disable()
		c.accInput.Disable()
		c.pswdInput.Disable()
		//c.subWindowShow.SetText(fmt.Sprintf("account:%s,pswd:%s", c.accInput.Text, c.pswdInput.Text))
	}
	c.loginForm.CancelText = "退出登录"
	c.loginForm.OnCancel = func() {
		log.Println("退出登录")
		//c.loginForm.Enable()
		c.accInput.Enable()
		c.pswdInput.Enable()
		c.curAccountLabel.SetText("当前账号:未登录")
		c.refreshLogContent(fmt.Sprintf("[%s]exitGame", c.accInput.Text))
		//c.subWindowShow.SetText(fmt.Sprintf("[%s]exitGame", c.accInput.Text))
	}

	left := canvas.NewText("", color.White)

	c.curAccountLabel = widget.NewLabel("当前账号:未登录")
	// btnGame := widget.NewButton("进入游戏", func() {
	// 	c.refreshLogContent(fmt.Sprintf("[%s]enterGame", c.accInput.Text))
	// })
	btnEnterHall := widget.NewButton("进入大厅", func() {
		c.refreshLogContent(fmt.Sprintf("[%s]send enter hall", c.accInput.Text))
		//c.SendLoginHall()
		c.hallClient.LoginHallCh <- true
	})
	btnExitHall := widget.NewButton("退出大厅", func() {
		c.refreshLogContent("exit hall")
		c.hallClient.Stop()
	})
	btnExitRoom := widget.NewButton("退出房间", func() {
		c.refreshLogContent("exit room")
		doCmd := &ClientCmd{
			CmdType: EClientCmdLeaveRoom,
			CmdData: nil,
		}
		c.hallClient.ClientCmdCh <- doCmd
	})
	btnJoinRoom := widget.NewButton("加入房间", func() {
		c.refreshLogContent("ready join room")

		// 打开加入房间界面
		c.InitJoinRoomWindow()
		c.joinRoomWindow.Show()
	})
	btnGameList := widget.NewButton("游戏列表", func() {
		c.refreshLogContent(fmt.Sprintf("game list"))
		c.InitGameListWindow()
		c.gameListWindow.Show()
		// createRoomCmd := &ClientCmd{
		// 	CmdType: EClientCmdCreateRoom,
		// 	CmdData: &ClientCmdCreateRoomData{},
		// }
		// c.hallClient.ClientCmdCh <- createRoomCmd
	})
	rights := container.NewVBox(c.curAccountLabel, btnEnterHall, btnExitHall, btnGameList, btnJoinRoom, btnExitRoom)

	mainLayout := container.New(layout.NewBorderLayout(c.loginForm, bottom, left, rights),
		c.loginForm,
		bottom,
		left,
		rights,
		c.logScroll)
	myWindow.SetContent(mainLayout)
	myWindow.Resize(fyne.NewSize(float32(c.clientCfg.WindowWidth), float32(c.clientCfg.WindowHeight)))
	myWindow.CenterOnScreen()
	myWindow.SetMaster()

	myWindow.ShowAndRun()
	os.Unsetenv("FYNE_FONT")
}
