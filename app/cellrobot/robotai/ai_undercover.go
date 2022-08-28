package robotai

import (
	"fmt"
	"math/rand"
	"roomcell/app/cellrobot/robotcore"
	"roomcell/pkg/evhub"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
)

type AiUndercover struct {
	RobotObj *robotcore.CellRobot
	// game data
	RoomGameData *pbclient.RoomUndercoverDetail
	//
	IsOut    bool // 是否出局
	IsTalked bool // 本次是否已经发言
	IsVoted  bool // 本次是否投票了
	// TotedRoleID     int64 // 本次投票的角色
	// TotedPlayNumber int32 // 本次投票的角色序号
	SelfNumber int32
}

func NewAiUndercover(robotObj *robotcore.CellRobot) *AiUndercover {
	return &AiUndercover{
		RobotObj: robotObj,
	}
}

func (robotAi *AiUndercover) Update(curTime int64) {
	robotObj := robotAi.RobotObj
	switch robotObj.RobotStatus {
	case robotcore.ERobotStatusNone:
		{
			robotObj.LoginHall()
			break
		}
	case robotcore.ERobotStatusLoginHall:
		{
			if robotObj.StatusStep == robotcore.ERobotStatusStepIng {
				break
			}
			// 登录大厅完成
			if robotObj.TargetRoomID > 0 {
				robotObj.SendEnterRoom(robotObj.TargetRoomID)
			}
			break
		}
	case robotcore.ERobotStatusEnterRoom:
		{
			if robotObj.StatusStep == robotcore.ERobotStatusStepIng {
				break
			}
			robotAi.UpdatePlay(curTime)

			break
		}
	}
}
func (robotAi *AiUndercover) UpdatePlay(curTime int64) {
	if robotAi.RoomGameData == nil {
		robotAi.RoomGameData = robotAi.RobotObj.RoomDetail.UndercoverRoomData
		robotAi.RobotObj.Debugf("robot(%s) init undercove room data", robotAi.RobotObj.RobotName)
		return
	}
	if robotAi.SelfNumber == 0 {
		if robotAi.RoomGameData != nil {
			if selfData, ok := robotAi.RoomGameData.PlayersGameData[robotAi.RobotObj.UserID]; ok {
				robotAi.SelfNumber = selfData.PlayerNumber
			}
		} else {
			return
		}
	}
	robotObj := robotAi.RobotObj
	switch robotAi.RoomGameData.RoomStep {
	case sconst.EUndercoverStepReady:
		{
			if !robotObj.IsReady {
				robotObj.DoActGameReady(1)
			}
			break
		}
	case sconst.EUndercoverStepGenWords:
		{
			robotAi.RobotObj.Debugf("robot(%s) wait gen words ...", robotAi.RobotObj.GetRobotName())
			break
		}
	case sconst.EUndercoverStepTalk:
		{
			// 是否自己发言
			if robotAi.RoomGameData.TalkRoleID != robotAi.RobotObj.UserID {
				break
			}
			if robotAi.IsTalked {
				break
			}
			// 发言
			robotAi.DoActTalk()
			break
		}
	case sconst.EUndercoverStepVote:
		{
			if robotAi.IsVoted {
				break
			}
			robotAi.DoActVote()
			break
		}
	case sconst.EUndercoverStepVoteSummary:
		{
			robotAi.RobotObj.Debugf("robot(%s) wait EUndercoverStepVoteSummary ...", robotAi.RobotObj.GetRobotName())
			break
		}
	case sconst.EUndercoverStepEnd:
		{
			robotAi.RobotObj.Debugf("robot(%s) wait EUndercoverStepEnd ...", robotAi.RobotObj.GetRobotName())
			break
		}
	}
}

// 新一轮重置
func (robotAi *AiUndercover) ResetNewTurn() {
	robotAi.RobotObj.Debugf("robot(%s) ResetNewTurn", robotAi.RobotObj.GetRobotName())
	robotAi.IsTalked = false
	robotAi.IsVoted = false
	robotAi.IsOut = false
	robotAi.SelfNumber = 0
}

// 新一轮发言重置
func (robotAi *AiUndercover) ResetNewTalk() {
	robotAi.RobotObj.Debugf("robot(%s) ResetNewTalk", robotAi.RobotObj.GetRobotName())
	robotAi.IsTalked = false
	robotAi.IsVoted = false
}

// robot action
func (robotAi *AiUndercover) DoActTalk() {
	robotAi.RobotObj.DoAction(robotcore.ActUndercoverTalk, func() {
		reqMsg := &pbclient.ECMsgRoomChatReq{
			TalkContent: fmt.Sprintf("hello, Im robot(%s)!!!", robotAi.RobotObj.GetRobotName()),
		}
		robotAi.RobotObj.RemoteCall(protocol.ECMsgClassRoom,
			protocol.ECMsgRoomChat,
			reqMsg,
			func(sMsg *evhub.NetMessage) {
				rsp := &pbclient.ECMsgRoomChatRsp{}
				if !trframe.DecodePBMessage(sMsg, rsp) {
					return
				}
				robotAi.RobotObj.LogCbMsgInfo(sMsg, rsp)
				if sMsg.Head.Result == protocol.ECodeSuccess {
					robotAi.IsTalked = true
				}
				robotAi.RobotObj.EndAction()
			},
		)
	})
}

// 投票
func (robotAi *AiUndercover) DoActVote() {
	// 随机选一个玩家
	var roleList []int64 = make([]int64, 0)
	for _, v := range robotAi.RoomGameData.PlayersGameData {
		if v.IsOut {
			continue
		}
		if v.RoleID == robotAi.RobotObj.UserID {
			continue
		}
		roleList = append(roleList, v.RoleID)
	}
	if len(roleList) < 1 {
		return
	}
	idx := rand.Intn(len(roleList))
	voteRoleID := roleList[idx]
	robotAi.RobotObj.DoAction(robotcore.ActUndercoverVote, func() {
		reqMsg := &pbclient.ECMsgGameUndercoverVoteReq{
			TargetRoleID: voteRoleID,
		}
		robotAi.RobotObj.RemoteCall(protocol.ECMsgClassGame,
			protocol.ECMsgGameUndercoverVote,
			reqMsg,
			func(sMsg *evhub.NetMessage) {
				rsp := &pbclient.ECMsgGameUndercoverVoteRsp{}
				if !trframe.DecodePBMessage(sMsg, rsp) {
					return
				}
				robotAi.RobotObj.LogCbMsgInfo(sMsg, rsp)
				if sMsg.Head.Result == protocol.ECodeSuccess {
					robotAi.IsVoted = true
				}
				robotAi.RobotObj.EndAction()
			},
		)
	})
}
