/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief:root 节点
 */
package troot

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// root 节点
type FrameNodeRoot struct {
	tframeObj iframe.ITRFrame
	nodeIndex int32
}

func New(frameObj iframe.ITRFrame, index int32) *FrameNodeRoot {
	return &FrameNodeRoot{
		tframeObj: frameObj,
		nodeIndex: index,
	}
}

func (fnRoot *FrameNodeRoot) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("frame node run step check")
	return true
}

func (fnRoot *FrameNodeRoot) RunStepInit(curTimeMs int64) bool {
	loghlp.Info("frame node run step init")
	return true
}
func (fnRoot *FrameNodeRoot) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("frame node run step preRun")
	return true
}
func (fnRoot *FrameNodeRoot) RunStepRun(curTimeMs int64) bool {

	return true
}
func (fnRoot *FrameNodeRoot) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("frame node run step stop")
	return true
}
func (fnRoot *FrameNodeRoot) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("frame node run step end")
	return true
}
func (fnRoot *FrameNodeRoot) NodeType() int32 {
	return trnode.ETRNodeTypeRoot
}
func (fnRoot *FrameNodeRoot) NodeIndex() int32 {
	return fnRoot.nodeIndex
}

func (fnRoot *FrameNodeRoot) SetUserFrameRun(func(curTimeMs int64)) {

}
