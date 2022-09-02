package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"roomcell/app/account"
	"roomcell/pkg/appconfig"
	"roomcell/pkg/loghlp"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	loghlp.ActiveConsoleLog()
	loghlp.Debugf("roomcell_account start")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "100", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
	// go func() {
	// 	log.Println(http.ListenAndServe("127.0.0.1:6060", nil))
	// }()
	roomCellCfg := appconfig.NewRoomCellCfg()
	errConfig := appconfig.ReadRoomCellConfigFromFile(viper.GetString("configPath"), roomCellCfg)
	if errConfig != nil {
		loghlp.Errorf("read config error:%s", errConfig.Error())
		return
	}
	jvData, err := json.Marshal(roomCellCfg)
	if err != nil {
		loghlp.Debugf("err:%s", err.Error())
	} else {
		loghlp.Debugf("configJv:%s", string(jvData))
	}
	// runtime.SetMutexProfileFraction(1)
	// runtime.SetBlockProfileRate(1)
	loghlp.Debugf("roomCellCfg:%+v", roomCellCfg)
	loghlp.SetConsoleLogLevel(logrus.Level(roomCellCfg.AccountCfgs[0].LogLevel))
	accApp := account.NewAccount(roomCellCfg)
	if roomCellCfg.AccountLocalDB != nil {
		if len(roomCellCfg.AccountLocalDB.DbFile) > 0 {
			loghlp.Infof("open localdb:%s", roomCellCfg.AccountLocalDB.DbFile)
			accApp.OpenLocalDB(roomCellCfg.AccountLocalDB.DbFile)
		} else {
			panic("local db config error")
		}
	} else {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", roomCellCfg.AccountDB.User, roomCellCfg.AccountDB.Pswd, roomCellCfg.AccountDB.Host, roomCellCfg.AccountDB.Port, roomCellCfg.AccountDB.DbName)
		loghlp.Infof("open mysql:%s", dsn)
		accApp.OpenMysqlDB(dsn)
	}

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt)
	go func() {
		<-stopCh
		accApp.Stop()
	}()
	accApp.Run(roomCellCfg.AccountCfgs[0].ListenAddr, 0)
}
