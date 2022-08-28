package trframe

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/trframe/trnode"
)

type FrameSession struct {
	netSession     *evhub.NetSession
	nodeType       int32
	nodeInfo       *trnode.TRNodeInfo
	heartTime      int64 // 最近的心跳时间
	isServerClient bool
}

func (fs *FrameSession) GetSessionID() int32 {
	return fs.netSession.GetSessionID()
}

func (fs *FrameSession) GetNodeInfo() *trnode.TRNodeInfo {
	return fs.nodeInfo
}
func (fs *FrameSession) Equal(zoneID int32, nodeType int32, nodeIndex int32) bool {
	if fs.nodeInfo == nil {
		return false
	}
	return fs.nodeInfo.ZoneID == zoneID && fs.nodeInfo.NodeType == nodeType && fs.nodeInfo.NodeIndex == nodeIndex
}
func (fs *FrameSession) SendMsg(msg *evhub.NetMessage) bool {
	if fs.netSession == nil {
		return false
	}
	fs.netSession.Send(msg)
	return true
}

func (fs *FrameSession) LastHeartTime() int64 {
	return fs.heartTime
}
func (fs *FrameSession) SetHeartTime(keepTime int64) {
	fs.heartTime = keepTime
}
func (fs *FrameSession) IsServerClient() bool {
	return fs.isServerClient
}
