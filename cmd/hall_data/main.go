package main

import (
	"fmt"
	"os"
	"os/signal"
	"roomcell/app/halldata"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/trnode"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetReportCaller(true)
	loghlp.ActiveConsoleLog()
	loghlp.SetConsoleLogLevel(logrus.DebugLevel)
	loghlp.Debugf("trframe hall_data start")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "0", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	cfgPath := viper.GetString("configPath")
	servIdx := cast.ToInt32(viper.GetString("index"))
	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt)
	go func() {
		<-stopSig
		trframe.Stop()
	}()
	trframe.Init(cfgPath, trnode.ETRNodeTypeHallData, servIdx)
	hallData := halldata.NewHallData()
	if len(trframe.GetFrameConfig().HallDataCfgs[servIdx].LocalDBFile) > 0 {
		hallData.OpenLocalDB(trframe.GetFrameConfig().HallDataCfgs[servIdx].LocalDBFile)
	} else {
		mainCfg := trframe.GetFrameConfig().HallDataCfgs[servIdx]
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mainCfg.GameDb.User, mainCfg.GameDb.Pswd, mainCfg.GameDb.Host, mainCfg.GameDb.Port, mainCfg.GameDb.DbName)
		loghlp.Infof("open gamemysql:%s", dsn)
		hallData.OpenMysqlDB(dsn)
		//hallData.OpenMysqlDB(dsn)
	}

	// trframe.RegUserCommandHandler(func(frameCmd *trframe.TRFrameCommand) {
	// 	loghlp.Infof("recv usercmd(%d_%d),type:%+v",
	// 		frameCmd.UserCmd.GetCmdClass(),
	// 		frameCmd.UserCmd.GetCmdType(),
	// 		reflect.TypeOf(frameCmd.UserCmd.GetCmdData()).String(),
	// 	)
	// 	// 处理websocket命令
	// 	hallGate.HandleCommand(frameCmd)
	// })
	trframe.RegisterUserFrameRun(func(curTimeMs int64) {
		hallData.FrameRun(curTimeMs)
	})
	trframe.Start()

	hallData.HallDataGlobal().StopDBRun()
}
