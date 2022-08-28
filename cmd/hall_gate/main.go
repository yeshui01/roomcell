package main

import (
	"os"
	"os/signal"
	"reflect"
	"roomcell/app/hallgate"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/trnode"
	"roomcell/pkg/wsserve"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetReportCaller(true)
	loghlp.ActiveConsoleLog()
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.Debugf("trframe hall_gate start")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "0", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	cfgPath := viper.GetString("configPath")
	wsServe := wsserve.NewWSServer()
	servIdx := cast.ToInt32(viper.GetString("index"))
	wsServe.SetupWSRouter("/ws", func(wsConn *websocket.Conn, err error) {
		trframe.PostUserCommand(protocol.CellCmdClassWebsocket, protocol.CmdTypeWebsocketConnect, wsConn)
	})
	stopSig := make(chan os.Signal)
	stopCh := make(chan bool)
	signal.Notify(stopSig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stopSig
		stopCh <- true
	}()

	trframe.Init(cfgPath, trnode.ETRNodeTypeHallGate, servIdx)
	// 运行websock监听
	wsAddr := trframe.GetFrameConfig().HallGateCfgs[servIdx].WsListenAddr
	go func() {
		wsServe.Run(wsAddr, stopCh, 0)
		trframe.Stop()
	}()
	hallGate := hallgate.NewHallGate()
	trframe.RegUserCommandHandler(func(frameCmd *trframe.TRFrameCommand) {
		loghlp.Infof("recv usercmd(%d_%d),type:%+v",
			frameCmd.UserCmd.GetCmdClass(),
			frameCmd.UserCmd.GetCmdType(),
			reflect.TypeOf(frameCmd.UserCmd.GetCmdData()).String(),
		)
		// 处理websocket命令
		hallGate.HandleCommand(frameCmd)
	})
	trframe.RegisterUserFrameRun(func(curTimeMs int64) {
		hallGate.FrameRun(curTimeMs)
	})
	loghlp.ActiveFileLogReportCaller(true)
	trframe.Start()
}
