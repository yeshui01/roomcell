package robotcore

const (
	ERobotStatusNone      = 0
	ERobotStatusLoginHall = 1 // 登录大厅
	ERobotStatusEnterRoom = 2 // 进入房间
)

const (
	ERobotStatusStepIng    = 0 // 进行中
	ERobotStatusStepFinish = 1 // 完成
)

const (
	ActGameReady      = "act_gameready"
	ActUndercoverTalk = "act_undercover_talk" // 谁是卧底-发言
	ActUndercoverVote = "act_undercover_vote" // 谁是卧底-投票
)
