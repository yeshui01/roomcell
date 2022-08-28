package sconst

const (
	EGameRoomTypeNone       = 0  // 空房间
	EGameRoomTypeDrawGuess  = 1  // 你画我猜
	EGameRoomTypeUndercover = 2  // 谁是卧底
	EGameRoomTypeChat       = 10 // 聊天房间

)

// 房间状态
const (
	ERoomStatusIdle = 0 // 空闲
	ERoomStatusGame = 1 // 游戏中
)

// 你画我猜房间阶段
const (
	EDrawGuessStepReady        = 0 // 准备阶段
	EDrawGuessStepSelectDrawer = 1 // 选择画图玩家
	EDrawGuessStepSelectWords  = 2 // 选词
	EDrawGuessStepDraw         = 3 // 画图阶段
	EDrawGuessStepEnd          = 4 // 画图阶段结束
	EDrawGuessStepGameEnd      = 5 // 游戏结束
)

// 断线原因
const (
	EPlayerOfflineReasonNormal       = 0 // 正常断线
	EPlayerOfflineReasonKickOut      = 1 // 踢人
	EPlayerOfflineReasonReplaceLogin = 2 // 顶号登录
)

// 谁是卧底房间阶段
const (
	EUndercoverStepReady       = 0 // 准备阶段
	EUndercoverStepGenWords    = 1 // 分配词语阶段
	EUndercoverStepTalk        = 2 // 发言阶段
	EUndercoverStepVote        = 3 // 投票阶段
	EUndercoverStepVoteSummary = 4 // 投票结束汇总阶段
	EUndercoverStepEnd         = 5 // 游戏结束阶段
)
