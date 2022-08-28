package robotcore

import (
	"roomcell/pkg/evhub"
)

func genHandleKey(msgClass int32, msgType int32) int32 {
	return msgClass*1000 + msgType
}

func (robotObj *CellRobot) RegMsgHandle(msgClass int32, msgType int32, handleFunc RobotMsgHandler) {
	handleKey := genHandleKey(msgClass, msgType)
	robotObj.MsgHandlers[handleKey] = handleFunc
}

// // 注册基本消息处理
// func (robotObj *CellRobot) InitRegisterCommnHandle() {
// 	// room
// 	robotObj.RegMsgHandle(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerEnter, RobotHandleECMsgRoomPushPlayerEnterNotify)
// 	robotObj.RegMsgHandle(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerLeave, RobotHandleECMsgRoomPushPlayerLeaveNotify)
// 	robotObj.RegMsgHandle(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushChat, RobotHandleECMsgRoomPushChatNotify)

// 	// game
// 	robotObj.RegMsgHandle(protocol.ECMsgClassGame, protocol.ECMsgGamePushPlayerReadyStatus, RobotHandleECMsgGamePushPlayerReadyStatusNotify)
// }

func (robotObj *CellRobot) DispatchRobotMsg(sMsg *evhub.NetMessage) {
	handleKey := genHandleKey(int32(sMsg.Head.MsgClass), int32(sMsg.Head.MsgType))
	if sMsg.SecondHead != nil && sMsg.SecondHead.ReqID > 0 {
		if callEnv, ok := robotObj.AsyncCall[int32(sMsg.SecondHead.ReqID)]; ok {
			robotObj.Debugf("handle robot callback msg(%d_%d),seqid(%d)", sMsg.Head.MsgClass, sMsg.Head.MsgType, sMsg.SecondHead.ReqID)
			callEnv.CallbackFunc(sMsg)
			delete(robotObj.AsyncCall, int32(sMsg.SecondHead.ReqID))
		} else {
			robotObj.Errorf("not find robot callenv, callback msg(%d_%d),seqid(%d)", sMsg.Head.MsgClass, sMsg.Head.MsgType, sMsg.SecondHead.ReqID)
		}
		return
	}
	if handleFunc, ok := robotObj.MsgHandlers[handleKey]; ok {
		if robotObj.RealRobotIns != nil {
			handleFunc(robotObj.RealRobotIns, sMsg)
		} else {
			handleFunc(robotObj, sMsg)
		}
	} else {
		robotObj.Debugf("not find robot msg(%d_%d) handler", sMsg.Head.MsgClass, sMsg.Head.MsgType)
	}
}
