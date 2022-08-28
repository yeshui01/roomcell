/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief:gate 节点
 */
package tgate

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// root 节点
type FrameNodeGate struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeGate {
	return &FrameNodeGate{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (frameNode *FrameNodeGate) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (frameNode *FrameNodeGate) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return frameNode.InitConnectRoot()
}
func (frameNode *FrameNodeGate) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (frameNode *FrameNodeGate) RunStepRun(curTimeMs int64) bool {

	return true
}
func (frameNode *FrameNodeGate) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (frameNode *FrameNodeGate) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (frameNode *FrameNodeGate) NodeType() int32 {
	return trnode.ETRNodeTypeRoot
}
func (frameNode *FrameNodeGate) NodeIndex() int32 {
	return frameNode.nodeIndex
}

// func (frameNode *FrameNodeGate) SetUserFrameRun(func(curTimeMs int64)) {
// }
