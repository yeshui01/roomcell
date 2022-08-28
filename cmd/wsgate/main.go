package main

import (
	"roomcell/pkg/loghlp"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	loghlp.ActiveConsoleLog()
	loghlp.Debugf("roomcell_wsgate start")
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "100", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))
}
