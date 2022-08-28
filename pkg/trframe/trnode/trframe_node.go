package trnode

import "roomcell/pkg/evhub"

// 节点类型
const (
	ETRNodeTypeNone = 0
	ETRNodeTypeRoot = 1
	ETRNodeTypeGate = 2

	// 大厅类型服务节点
	ETRNodeTypeHallGate    = 1000 // 大厅网关
	ETRNodeTypeHallServer  = 1001 // 大厅主服务器
	ETRNodeTypeHallData    = 1002 // 大厅用户服务器
	ETRNodeTypeHallRoomMgr = 1003 // 大厅房间管理器
	ETRNodeTypeHallRoom    = 1004 // 大厅房间

)

// 节点
type TRNodeInfo struct {
	ZoneID    int32
	NodeType  int32
	NodeIndex int32
	DesInfo   string
}

// 节点实体
type ITRNodeEntity interface {
	GetNodeInfo() *TRNodeInfo
	Equal(zoneID int32, nodeType int32, nodeIndex int32) bool
	SendMsg(msg *evhub.NetMessage) bool
	LastHeartTime() int64
	SetHeartTime(int64)
	IsServerClient() bool
	GetSessionID() int32
}
