package protocol

const (
	// 业务层面
	ESMsgClassPlayer int32 = 300 // 玩家类别
	ESMsgClassRoom   int32 = 301 // 房间类别
)

// ESMsgClassPlayer int32 = 300 // 玩家类别
const (
	ESMsgPlayerLoadRole   = 1 // 加载角色数据
	ESMsgPlayerSaveRole   = 2 // 保存角色数据
	ESMsgPlayerLoginHall  = 3 // 登录大厅
	ESMsgPlayerDisconnect = 4 // 玩家连接断开
	ESMsgPlayerKickOut    = 5 // 踢掉玩家
)

// ESMsgClassRoom int32 = 301 // 房间类别
const (
	ESMsgRoomCreate     = 1 // 创建房间
	ESMsgRoomAutoDelete = 2 // 房间删除
	ESMsgRoomEnter      = 3 // 进入房间
	ESMsgRoomLeave      = 4 // 离开房间
	ESMsgRoomFind       = 5 // 获取房间信息
)
