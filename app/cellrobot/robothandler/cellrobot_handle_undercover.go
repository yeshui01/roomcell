package robothandler

import (
	"roomcell/app/cellrobot/robotai"
	"roomcell/app/cellrobot/robotcore"
	"roomcell/pkg/evhub"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/sconst"
)

func RobotHandleECMsgGamePushUndercoverRoomDataNotify(robotIns robotcore.ICellRobot, sMsg *evhub.NetMessage) {
	notifyMsg := &pbclient.ECMsgGamePushUndercoverRoomDataNotify{}
	robotIns.LogRecvMsgInfo(sMsg, notifyMsg)
	robotObj := robotIns.(*robotcore.CellRobot)
	aiObj := robotObj.AiInstance.(*robotai.AiUndercover)
	aiObj.RoomGameData = notifyMsg.RoomGameData
	robotObj.Infof("recv undercove room data notify,step(%d)", aiObj.RoomGameData.GetRoomStep())
	switch aiObj.RoomGameData.RoomStep {
	case sconst.EUndercoverStepReady:
		{
			aiObj.ResetNewTurn()
			break
		}
	case sconst.EUndercoverStepTalk:
		{
			aiObj.ResetNewTalk()
			break
		}
	case sconst.EUndercoverStepGenWords:
		{
			break
		}
	case sconst.EUndercoverStepVote:
		{
			break
		}
	case sconst.EUndercoverStepVoteSummary:
		{
			break
		}
	case sconst.EUndercoverStepEnd:
		{
			break
		}
	}
}
