/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-07-06 14:58:03
 * @LastEditTime: 2022-09-27 11:25:30
 * @FilePath: \roomcell\cmd\wsgate\main.go
 */
package main

import (
	"net"
	"roomcell/pkg/loghlp"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func testUnixSockServer() {
	var unixAddr *net.UnixAddr
	unixAddr, _ = net.ResolveUnixAddr("unix", "/tmp/cell_unixsock")
	net.Listen("unix", "./tmp/cell_unixsock")
	unixLisner, _ := net.ListenUnix("unix", unixAddr)
	unixLisner.Accept()
	unixLisner.AcceptUnix()
}

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
