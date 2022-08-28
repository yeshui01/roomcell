package main

import (
	"roomcell/app/hallclient"
	"roomcell/pkg/loghlp"
)

func main() {
	loghlp.ActiveConsoleLog()
	hallClient := hallclient.NewHallClient()
	hallClient.Connect("localhost:7200")
	hallClient.Run()
}
