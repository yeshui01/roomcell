package hallgate

import (
	"roomcell/app/hallgate/hallgatehandler"
	"roomcell/app/hallgate/hallgatemain"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"

	"google.golang.org/protobuf/proto"
)

type HallGate struct {
	UserMgr       *hallgatemain.HGateUserManager
	ConnMgr       *hallgatemain.HGateClientManager
	CmdHandlerMap map[int64]func(frameCmd *trframe.TRFrameCommand)
}

func NewHallGate() *HallGate {
	hg := &HallGate{
		CmdHandlerMap: make(map[int64]func(frameCmd *trframe.TRFrameCommand)),
	}
	hg.UserMgr = hallgatemain.NewHGateUserManager()
	hg.ConnMgr = hallgatemain.NewHGateClientManager()
	hallgatehandler.InitHallGateServe(hg)
	hg.RegisterMsgHandler()
	hg.InitCmdHandle()
	return hg
}

func GetCmdHandlerKey(cmdClass int32, cmdType int32) int64 {
	k := int64(cmdClass)<<32 + int64(cmdType)
	return k
}
func (hg *HallGate) InitCmdHandle() {
	hg.RegisterCmdHandler(protocol.CellCmdClassWebsocket,
		protocol.CmdTypeWebsocketConnect,
		hg.HandleCmdWSConnect)
	hg.RegisterCmdHandler(protocol.CellCmdClassWebsocket,
		protocol.CmdTypeWebsocketClosed,
		hg.HandleCmdWSClose)
	hg.RegisterCmdHandler(protocol.CellCmdClassWebsocket,
		protocol.CmdTypeWebsocketMessage,
		hg.HandleCmdWSMessage)

}

func (hg *HallGate) RegisterCmdHandler(cmdClass int32, cmdType int32, handle func(frameCmd *trframe.TRFrameCommand)) {
	hKey := GetCmdHandlerKey(cmdClass, cmdType)
	if _, ok := hg.CmdHandlerMap[hKey]; ok {
		loghlp.Errorf("repeated register cmd(%d_%d) handle", cmdClass, cmdType)
	} else {
		hg.CmdHandlerMap[hKey] = handle
	}
}

func (hg *HallGate) FrameRun(curTimeMs int64) {

}

func (hg *HallGate) SendWSClientReplyMessage(okCode int32, cltRep proto.Message, env *iframe.TRRemoteMsgEnv) {
	msgData, _ := proto.Marshal(cltRep)
	hg.SendWSClientReplyMessage2(okCode, msgData, env)
}
func (hg *HallGate) SendWSClientReplyMessage2(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
	hgsession, ok := env.UserData.(*hallgatemain.HGateConnction)
	if !ok {
		hgsession = nil
		loghlp.Errorf("SendWSClientReplyMessage2,hgsession convert fail")
		// 根据玩家ID查找
		if env.SrcMessage != nil {
			if env.SrcMessage.SecondHead.ID > 0 {
				gateUser := hg.GetUserManager().GetGateUser(env.SrcMessage.SecondHead.ID)
				if gateUser != nil {
					hgsession = gateUser.GetGateConnect()
				}
			}
		}
		return
	}
	if hgsession == nil {
		return
	}
	env.SrcMessage.Data = msgData
	env.SrcMessage.Head.Result = uint16(okCode)
	clientReqID := uint64(0)
	if env.SrcMessage.SecondHead != nil {
		if env.SrcMessage.SecondHead.ReqID > 0 {
			env.SrcMessage.SecondHead.RepID = env.SrcMessage.SecondHead.ReqID
			clientReqID = env.SrcMessage.SecondHead.ReqID
		}
	}
	hgsession.SendMsg(env.SrcMessage)
	loghlp.Debugf("SendWSClientReplyMessage2(%d_%d)(clientReqNo:%d)",
		env.SrcMessage.Head.MsgClass,
		env.SrcMessage.Head.MsgType,
		clientReqID,
	)
}
func (hg *HallGate) GetGateConnMgr() *hallgatemain.HGateClientManager {
	return hg.ConnMgr
}
