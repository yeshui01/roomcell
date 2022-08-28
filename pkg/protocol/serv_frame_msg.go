package protocol

const (
	EFrameMsgRegisterServerInfo   int32 = 1 // 注册服务信息
	EFrameMsgKeepNodeHeart        int32 = 2 // 节点心跳
	EFrameMsgServerMsgConvert     int32 = 3 // 服务器消息转换
	EFrameMsgPushMsgToClient      int32 = 4 // 推送消息给玩家
	EFrameMsgBroadcastMsgToClient int32 = 5 // 广播消息给玩家
)
