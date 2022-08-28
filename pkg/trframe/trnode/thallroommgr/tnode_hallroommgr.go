package thallroommgr

/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-08-01 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief:hall 节点
 */

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// gate 节点
type FrameNodeRoomMgr struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeRoomMgr {
	return &FrameNodeRoomMgr{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeRoomMgr) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeRoomMgr) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	// return frameNode.InitConnectServer()
	return true
}
func (frameNode *FrameNodeRoomMgr) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeRoomMgr) RunStepRun(curTimeMs int64) bool {
	return true
}
func (frameNode *FrameNodeRoomMgr) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeRoomMgr) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeRoomMgr) NodeType() int32 {
	return trnode.ETRNodeTypeHallServer
}
func (frameNode *FrameNodeRoomMgr) NodeIndex() int32 {
	return frameNode.nodeIndex
}

// func (frameNode *FrameNodeGate) SetUserFrameRun(func(curTimeMs int64)) {
// }
