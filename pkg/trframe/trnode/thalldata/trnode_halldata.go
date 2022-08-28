package thalldata

/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-07-29 14:14:17
 * @LastEditTime: 2022-07-29 14:14:17
 * @Brief:hall 节点
 */

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// gate 节点
type FrameNodeHallData struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeHallData {
	return &FrameNodeHallData{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeHallData) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeHallData) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeHallData) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeHallData) RunStepRun(curTimeMs int64) bool {
	return true
}
func (frameNode *FrameNodeHallData) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeHallData) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeHallData) NodeType() int32 {
	return trnode.ETRNodeTypeHallServer
}
func (frameNode *FrameNodeHallData) NodeIndex() int32 {
	return frameNode.nodeIndex
}

// func (frameNode *FrameNodeGate) SetUserFrameRun(func(curTimeMs int64)) {
// }
