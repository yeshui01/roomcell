/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief:gate 节点
 */
package thallgate

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// gate 节点
type FrameNodeHallGate struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeHallGate {
	return &FrameNodeHallGate{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeHallGate) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeHallGate) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectServer()
}
func (frameNode *FrameNodeHallGate) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeHallGate) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeHallGate) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeHallGate) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeHallGate) NodeType() int32 {
	return trnode.ETRNodeTypeHallGate
}
func (frameNode *FrameNodeHallGate) NodeIndex() int32 {
	return frameNode.nodeIndex
}

// func (frameNode *FrameNodeGate) SetUserFrameRun(func(curTimeMs int64)) {
// }
