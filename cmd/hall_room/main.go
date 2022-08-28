package main

import (
	"fmt"
	"os"
	"os/signal"
	"roomcell/app/hallroom"
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
	loghlp.Debugf("trframe hall_room start")
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
	trframe.Init(cfgPath, trnode.ETRNodeTypeHallRoom, servIdx)
	hallServ := hallroom.NewHallRoom()
	// 敏感词初始化
	loghlp.Infof("init sensitive")
	hallServ.GetGlobalData().InitSensitiveWords(fmt.Sprintf("%s/sensitive_words.txt", cfgPath))
	trframe.RegisterUserFrameRun(func(curTimeMs int64) {
		hallServ.FrameRun(curTimeMs)
	})
	loghlp.ActiveFileLogReportCaller(true)
	trframe.Start()
}
