package main

import (
	"os"
	"os/signal"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/wsserve"

	"github.com/gorilla/websocket"
)

func main() {
	loghlp.ActiveConsoleLog()
	serv := wsserve.NewWSServer()
	stopSig := make(chan os.Signal)
	stopCh := make(chan bool)
	signal.Notify(stopSig, os.Interrupt)
	go func() {
		<-stopSig
		stopCh <- true
	}()
	// echo
	serv.SetupWSRouter("/ws", func(wsConn *websocket.Conn, err error) {
		if err == nil {
			for {
				msgType, msgData, err := wsConn.ReadMessage()
				if err != nil {
					loghlp.Errorf("ws error:%s", err.Error())
					break
				}
				switch msgType {
				case websocket.TextMessage:
					{
						loghlp.Infof("recv txt message:%s", string(msgData))
					}
				}
				// 回复消息
				wsConn.WriteMessage(msgType, msgData)
			}
		}
	})
	serv.Run("localhost:7200", stopCh, 0)
}
