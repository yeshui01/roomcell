package cellclient

import (
	"fmt"
	"log"
	"roomcell/pkg/evhub"
	"time"

	"github.com/gorilla/websocket"
)

func (hclient *HallClient) RunWithWindow(c *CellClient) {
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
					c.refreshLogContent(fmt.Sprintf("recv txtmessage:%s", string(message)))
					break
				}
			case websocket.BinaryMessage:
				{
					// sMsg := evhub.MakeEmptyMessage()
					// sMsg.Decode(message)
					sMsg := evhub.DecodeClientMsgToServerMsg(message)
					c.refreshLogContent(fmt.Sprintf("recv BinaryMessage(%d_%d),ok:%d,SeqID:%d", sMsg.Head.MsgClass, sMsg.Head.MsgType, sMsg.Head.Result, sMsg.SecondHead.ReqID))
					doCmd := &ClientCmd{
						CmdType: EClientCmdRecvMessage,
						CmdData: sMsg,
					}
					c.hallClient.ClientCmdCh <- doCmd

					// if sMsg.Head.MsgClass == protocol.ECMsgClassPlayer && sMsg.Head.MsgType == protocol.ECMsgPlayerLoginHall {
					// 	pbRep := &pbclient.ECMsgPlayerLoginHallRsp{}
					// 	errParse := proto.Unmarshal(sMsg.Data, pbRep)
					// 	if trframe.DecodePBMessage(sMsg, pbRep) {
					// 		c.refreshLogContent(fmt.Sprintf("recv ECMsgPlayerLoginHallRsp:%+v", pbRep))
					// 	} else {
					// 		c.refreshLogContent(fmt.Sprintf("decode ECMsgPlayerLoginHallRsp error:%s", errParse.Error()))
					// 	}
					// }
					break
				}
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 10)
	tickerHeart := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	defer tickerHeart.Stop()
	for {
		select {
		case <-hclient.ChDone:
			hclient.WsConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			return
		case <-tickerHeart.C:
			// reqMsg := &pbclient.ECMsgPlayerKeepHeartReq{}
			// sMsg := trframe.MakePBMessage(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerKeepHeart, reqMsg, protocol.ECodeSuccess)
			// sMsg.SecondHead = &evhub.NetMsgSecondHead{
			// 	ReqID: 1,
			// }
			// sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
			// err := hclient.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
			// if err != nil {
			// 	loghlp.Errorf("write heart error:%s", err.Error())
			// }
			doCmd := &ClientCmd{
				CmdType: EClientCmdKeepHeart,
				CmdData: nil,
			}
			c.hallClient.ClientCmdCh <- doCmd
		case <-ticker.C:
			// err := hclient.WsConn.WriteMessage(websocket.TextMessage, []byte("this is test heart"))
			// if err != nil {
			// 	loghlp.Errorf("write error:%s", err.Error())
			// 	return
			// }
			//c.SendLoginHall()
			//cToken, _ := genClientJwtToken(1, "testAcc1", 1)
			// reqMsg := &pbclient.ECMsgPlayerLoginHallReq{
			// 	Token: cToken,
			// }
			// sMsg := trframe.MakePBMessage(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerLoginHall, reqMsg, protocol.ECodeSuccess)
			// sMsg.SecondHead = &evhub.NetMsgSecondHead{
			// 	ReqID: 1,
			// }
			// sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
			// err = hclient.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
			// if err != nil {
			// 	loghlp.Errorf("write binary error:%s", err.Error())
			// }
		case <-hclient.LoginHallCh:
			{
				c.SendLoginHall()
			}
		case doCmd := <-hclient.ClientCmdCh:
			{
				c.HandleCmd(doCmd)
			}
		}
	}
}
