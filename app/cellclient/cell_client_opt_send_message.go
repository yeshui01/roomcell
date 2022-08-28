package cellclient

import (
	"fmt"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func (c *CellClient) SendLoginHall() {
	// cToken, _ := genClientJwtToken(1, "testAcc1", 1)
	reqMsg := &pbclient.ECMsgPlayerLoginHallReq{
		Token: c.loginToken,
	}
	sMsg := trframe.MakePBMessage(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerLoginHall, reqMsg, protocol.ECodeSuccess)
	sMsg.SecondHead = &evhub.NetMsgSecondHead{
		ReqID: uint64(c.GenSeqID()),
	}
	sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
	err := c.hallClient.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
	if err != nil {
		loghlp.Errorf("write binary error:%s", err.Error())
	}
}
func (c *CellClient) SendCreateRoom() {
	// cToken, _ := genClientJwtToken(1, "testAcc1", 1)
	if c.curGameChoose == nil {
		return
	}
	reqMsg := &pbclient.ECMsgRoomCreateReq{
		GameType: c.curGameChoose.RoomType,
	}
	sMsg := trframe.MakePBMessage(protocol.ECMsgClassRoom, protocol.ECMsgRoomCreate, reqMsg, protocol.ECodeSuccess)
	sMsg.SecondHead = &evhub.NetMsgSecondHead{
		ReqID: uint64(c.GenSeqID()),
	}
	sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
	err := c.hallClient.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
	if err != nil {
		loghlp.Errorf("write binary error:%s", err.Error())
	}
}
func (c *CellClient) SendLeaveRoom() {
	// cToken, _ := genClientJwtToken(1, "testAcc1", 1)
	reqMsg := &pbclient.ECMsgRoomLeaveReq{}
	c.sendMsgToServer(protocol.ECMsgClassRoom, protocol.ECMsgRoomLeave, reqMsg)
}
func (c *CellClient) sendMsgToServer(msgClass int32, msgType int32, pbMsg proto.Message) {
	sMsg := trframe.MakePBMessage(msgClass, msgType, pbMsg, protocol.ECodeSuccess)
	sMsg.SecondHead = &evhub.NetMsgSecondHead{
		ReqID: uint64(c.GenSeqID()),
	}
	sendData := evhub.EncodeServerMsgToClientMsg(sMsg)
	err := c.hallClient.WsConn.WriteMessage(websocket.BinaryMessage, sendData)
	c.refreshLogContent(fmt.Sprintf("send msg(%d_%d)[%d] to server", msgClass, msgType, sMsg.SecondHead.ReqID))
	if err != nil {
		loghlp.Errorf("write binary error:%s", err.Error())
	}
}

func (c *CellClient) SendEnterRoom(roomID int64) {
	reqMsg := &pbclient.ECMsgRoomEnterReq{
		RoomID: roomID,
	}
	c.sendMsgToServer(protocol.ECMsgClassRoom, protocol.ECMsgRoomEnter, reqMsg)
}

func (c *CellClient) SendKeepHeart() {
	reqMsg := &pbclient.ECMsgPlayerKeepHeartReq{}
	c.sendMsgToServer(protocol.ECMsgClassPlayer, protocol.ECMsgPlayerKeepHeart, reqMsg)
}

func (c *CellClient) SendDrawReady(isReady int32) {
	reqMsg := &pbclient.ECMsgGameReadyOptReq{
		Ready: isReady,
	}
	c.sendMsgToServer(protocol.ECMsgClassGame, protocol.ECMsgGameReadyOpt, reqMsg)
}

func (c *CellClient) SendDrawPaint(painData *pbclient.DrawPainData) {
	reqMsg := &pbclient.ECMsgGameDrawPaintReq{
		CurPain: painData,
	}
	c.sendMsgToServer(protocol.ECMsgClassGame, protocol.ECMsgGameDrawPaint, reqMsg)
}

func (c *CellClient) SendChatMessage(talkContent string) {
	reqMsg := &pbclient.ECMsgRoomChatReq{
		TalkContent: talkContent,
	}
	c.sendMsgToServer(protocol.ECMsgClassRoom, protocol.ECMsgRoomChat, reqMsg)
}
