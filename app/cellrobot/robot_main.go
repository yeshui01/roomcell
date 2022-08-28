package cellrobot

import (
	"roomcell/app/cellrobot/robotai"
	"roomcell/app/cellrobot/robotcore"
	"roomcell/app/cellrobot/robothandler"
	"roomcell/pkg/protocol"
)

func NewAiRobot(robotName string, targetRoomID int64, aiType int32) *robotcore.CellRobot {
	r := robotcore.NewCellRobot(robotName)
	r.TargetRoomID = targetRoomID

	switch aiType {
	case robotai.ERobotAiUndercover:
		{
			r.AiInstance = robotai.NewAiUndercover(r)
			break
		}
	}
	initRegisterCommnHandle(r)
	return r
}

// 注册基本消息处理
func initRegisterCommnHandle(robotObj *robotcore.CellRobot) {
	// room
	robotObj.RegMsgHandle(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerEnter, robothandler.RobotHandleECMsgRoomPushPlayerEnterNotify)
	robotObj.RegMsgHandle(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushPlayerLeave, robothandler.RobotHandleECMsgRoomPushPlayerLeaveNotify)
	robotObj.RegMsgHandle(protocol.ECMsgClassRoom, protocol.ECMsgRoomPushChat, robothandler.RobotHandleECMsgRoomPushChatNotify)

	// game
	robotObj.RegMsgHandle(protocol.ECMsgClassGame, protocol.ECMsgGamePushPlayerReadyStatus, robothandler.RobotHandleECMsgGamePushPlayerReadyStatusNotify)

	// undercover
	robotObj.RegMsgHandle(protocol.ECMsgClassGame, protocol.ECMsgGamePushUndercoverRoomData, robothandler.RobotHandleECMsgGamePushUndercoverRoomDataNotify)
}
