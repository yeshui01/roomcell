package hallclient

import (
	"log"
	"net/url"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type HallClient struct {
	// evHub *evhub.EventHub
	WsConn *websocket.Conn
	ChDone chan struct{}
}

func NewHallClient() *HallClient {
	return &HallClient{
		ChDone: make(chan struct{}),
	}
}

func (hclient *HallClient) Connect(servAddr string) error {
	u := url.URL{Scheme: "ws", Host: servAddr, Path: "/ws"}
	loghlp.Infof("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
		return err
	}
	hclient.WsConn = c
	return nil
}

func (hclient *HallClient) Run() {
	defer hclient.WsConn.Close()
	go func() {
		defer close(hclient.ChDone)
		for {
			msgType, message, err := hclient.WsConn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			switch msgType {
			case websocket.TextMessage:
				{
					loghlp.Debugf("recv txtmessage:%s", string(message))
					break
				}
			case websocket.BinaryMessage:
				{
					// sMsg := evhub.MakeEmptyMessage()
					// sMsg.Decode(message)
					sMsg := evhub.DecodeClientMsgToServerMsg(message)
					loghlp.Debugf("recv BinaryMessage(%d_%d)SeqID:%d", sMsg.Head.MsgClass, sMsg.Head.MsgType, sMsg.SecondHead.ReqID)
					if sMsg.Head.MsgClass == protocol.ECMsgClassPlayer && sMsg.Head.MsgType == protocol.ECMsgPlayerLoginHall {
						pbRep := &pbclient.ECMsgPlayerLoginHallRsp{}
						errParse := proto.Unmarshal(sMsg.Data, pbRep)
						if trframe.DecodePBMessage(sMsg, pbRep) {
							loghlp.Debugf("recv ECMsgPlayerLoginHallRsp:%+v", pbRep)
						} else {
							loghlp.Errorf("decode ECMsgPlayerLoginHallRsp error:%s", errParse.Error())
						}
					}
					break
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-hclient.ChDone:
			return
		case t := <-ticker.C:
			err := hclient.WsConn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				loghlp.Errorf("write error:%s", err.Error())
				return
			}
			cToken, _ := genClientJwtToken(1, "testAcc1", 1)
			reqMsg := &pbclient.ECMsgPlayerLoginHallReq{
				Token: cToken,
			}
			sMsg := trframe.MakePBMessage(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerLoginHall, reqMsg, protocol.ECodeSuccess)
			sMsg.SecondHead = &evhub.NetMsgSecondHead{
				ReqID: 1,
			}
			sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
			err = hclient.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
			if err != nil {
				loghlp.Errorf("write binary error:%s", err.Error())
			}
		}
	}
}
