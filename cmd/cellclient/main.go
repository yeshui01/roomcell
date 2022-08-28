package main

import (
	"roomcell/app/cellclient"
	"roomcell/pkg/loghlp"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("configPath", "./", "config file path")
	pflag.String("index", "100", "server index")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	loghlp.ActiveConsoleLog()

	loghlp.Debugf("configPath=%s", viper.GetString("configPath"))
	loghlp.Debugf("serverIndex=%s", viper.GetString("index"))

	configPath := viper.GetString("configPath")
	cellClient := cellclient.NewCellClient(configPath)

	cellClient.Run()
	//os.Unsetenv("FYNE_FONT")
}
