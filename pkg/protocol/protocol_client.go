package protocol

// 客户端的从2开始
const (
	ECMsgClassPlayer = 2 // 玩家
	ECMsgClassRoom   = 3 // 房间
	ECMsgClassGame   = 4 // 玩游戏
)

// ECMsgClassPlayer = 2
const (
	ECMsgPlayerLoginHall = 1 // 登录大厅
	ECMsgPlayerKeepHeart = 2 // 心跳
)

// ECMsgClassRoom   = 3 // 房间
const (
	ECMsgRoomCreate            = 1 // 创建房间
	ECMsgRoomQuery             = 2 // 房间查询
	ECMsgRoomEnter             = 3 // 进入房间
	ECMsgRoomPushPlayerEnter   = 4 // 推送有玩家进入房间
	ECMsgRoomLeave             = 5 // 离开房间
	ECMsgRoomPushPlayerLeave   = 6 // 推送有玩家离开房间
	ECMsgRoomPushPlayerOffline = 7 // 推送有玩家离线
	ECMsgRoomChat              = 8 // 房间聊天
	ECMsgRoomPushChat          = 9 // 房间聊天-推送
)

// ECMsgClassGame   = 4 // 玩游戏
const (
	ECMsgGameReadyOpt              = 1 // 玩家准备操作
	ECMsgGamePushPlayerReadyStatus = 2 // 推送玩家准备状态

	ECMsgGameDrawPaint        = 10 // 你画我猜-画图
	ECMsgGamePushDrawPaint    = 11 // 你画我猜-同步推送画图数据
	ECMsgGameDrawGuessWords   = 12 // 你画我猜-猜词
	ECMsgGamePushDrawGuess    = 13 // 你画我猜-推送玩家的猜词
	ECMsgGameGrawSetting      = 14 // 你画我猜-房主设定游戏规则
	ECMsgGamePushDrawRoomData = 15 // 你画我猜-推送画图房间游戏数据更新
	ECMsgGameDrawSelectWords  = 16 // 选择词语

	ECMsgGamePushUndercoverRoomData     = 30 // 谁是卧底-推送房间游戏数据
	ECMsgGamePushPlayerUnderWords       = 31 // 谁是卧底-推送玩家卧底词汇更新
	ECMsgGameUndercoverVote             = 32 // 谁是卧底-投票
	ECMsgGamePushUndercoverTalkerChange = 33 // 谁是卧底-发言人变更
	ECMsgGamePushUndercoverVote         = 34 // 谁是卧底-推送投票
	ECMsgGamePushUndercoverOut          = 35 // 谁是卧底-推送玩家出局
	ECMsgGameUndercoverTalk             = 36 // 谁是卧底-发言
	ECMsgGamePushUndercoverTalk         = 37 // 谁是卧底-推送发言

	ECMsgGamePushNumberBombRoomData      = 50 // 数字炸弹-推送房间游戏数据
	ECMsgGamePushNumberBombGuesserChange = 51 // 数字炸弹-猜数字玩家变更
	ECMsgGamePushNumberBombRangeChange   = 52 // 数字炸弹-数字范围变更
	ECMsgGameNumberBombGuess             = 53 // 数字炸弹-猜数字
	ECMsgGamePushNumberBombGuess         = 54 // 数字炸弹-猜数字推送
	ECMsgGameNumberBombChoosePunishment  = 55 // 数字炸弹-选择惩罚
	ECMsgGameNumberBombSetting           = 56 // 数字炸弹-设定

	ECMsgGameRescueSetting      = 70 // 拯救玩家-设定
	ECMsgGamePushRescueRoomData = 71 // 拯救玩家-推送房间游戏数据
	ECMsgGamePushRescueSetting  = 72 // 拯救玩家-设定推送
	ECMsgGameRescueRecvGift     = 73 // 拯救玩家-收到礼物
	ECMsgGamePushRescueRecvGift = 74 // 拯救玩家-收到礼物推送
)
