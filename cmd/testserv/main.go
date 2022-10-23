/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-07-18 18:38:20
 * @LastEditTime: 2022-09-27 11:11:59
 * @FilePath: \roomcell\cmd\testserv\main.go
 */
package main

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"sync"
	"time"
)

// type TestServer struct {
// 	tserver.TServer
// 	appConfig *appconfig.RoomCellConfig
// }

// func NewTestServer() *TestServer {
// 	s := tserver.NewTServer(0, 0, 0)
// 	ts := &TestServer{}
// 	ts.TServer = *s
// 	return ts
// }

// func (ts *TestServer) RunStepCheck() bool {
// 	ts.TServer.RunStepCheck(time.Now().Unix())
// 	loghlp.Infof("TestServer RunStepCheck")
// 	return true
// }
// func (ts *TestServer) RunStepInit() bool {
// 	ts.TServer.RunStepInit(time.Now().Unix())
// 	loghlp.Infof("TestServer RunStepInit")
// 	return true
// }

var servAddr string = ":15000"
var hubServ *evhub.EventHub = nil

func runServer() {
	hubServ = evhub.NewHub()
	hubServ.Listen(evhub.ListenModeTcp, servAddr)
	hubServ.OnMessage(func(netSession *evhub.NetSession, netMsg *evhub.NetMessage) {
		loghlp.Infof("server recv message:%s", string(netMsg.Data))
		netSession.Send(netMsg)
	})
	hubServ.Run()
}
func runClient() {
	evClient := evhub.NewEvhubClient(servAddr)
	msg := evhub.MakeEmptyMessage()
	msg.Data = []byte("hello world!")
	evClient.SendMsg(msg)
	recvMsg, err := evClient.RecvMsg(3000)
	if err == nil {
		loghlp.Infof("client recv message:%s", string(recvMsg.Data))
	} else {
		loghlp.Errorf("recv msg error:%s", err.Error())
	}
	evClient.Close()
}

func main() {
	loghlp.ActiveConsoleLog()
	// 测试一下protobuf的反射
	// pbName := pbtools.GetFullNameByMessage(&pbclient.RoleInfo{})
	// loghlp.Infof("pbName:%s", pbName)
	// pbMsg := pbtools.GetNewMessageObjByName(pbName)
	// roleInfo := pbMsg.(*pbclient.RoleInfo)
	// roleInfo.Name = "xxx"
	// baseServ := NewTestServer()
	// baseServ.RunStepCheck()
	// baseServ.RunStepInit()
	//baseServ.Run()
	var gw sync.WaitGroup
	gw.Add(1)
	go func() {
		runServer()
		gw.Done()
	}()
	time.Sleep(time.Second * 2)
	gw.Add(1)
	go func() {
		runClient()
		gw.Done()
	}()
	gw.Wait()
}
