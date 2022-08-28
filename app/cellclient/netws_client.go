package cellclient

import (
	"log"
	"net/url"
	"roomcell/pkg/crossdef"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func genClientJwtToken(userID int64, userName string, dataZone int32) (string, error) {
	j := &crossdef.JWT{
		[]byte(crossdef.SignKey),
	}
	tokenTime := int64(3600) // 默认1小时
	claims := crossdef.CustomClaims{
		Account:       userName,
		LastLoginTime: 0,
		CchId:         "",
		DataZone:      dataZone,
		UserID:        userID,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,      // 签名生效时间
			ExpiresAt: time.Now().Unix() + tokenTime, // 过期时间 一小时
			Issuer:    crossdef.SignKey,              //签名的发行者
		},
	}

	token, err := j.CreateToken(claims)
	return token, err
}

type HallClient struct {
	// evHub *evhub.EventHub
	WsConn      *websocket.Conn
	ChDone      chan struct{}
	LoginHallCh chan bool
	ClientCmdCh chan *ClientCmd
}

func NewHallClient() *HallClient {
	return &HallClient{
		ChDone:      make(chan struct{}),
		LoginHallCh: make(chan bool),
		ClientCmdCh: make(chan *ClientCmd, 100),
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

func (hclient *HallClient) Stop() {
	hclient.ChDone <- struct{}{}
	//close(hclient.ChDone)
}
