package trframe

import (
	"roomcell/pkg/loghlp"
)

func (tf *TRFrame) RunStepCheck(curTimeMs int64) bool {
	loghlp.Info("TRFrame::RunStepCheck")
	if !tf.curWorkNode.RunStepCheck(curTimeMs) {
		return false
	}
	if tf.userStepRun[ETRFrameStepCheck] != nil {
		return tf.userStepRun[ETRFrameStepCheck](curTimeMs)
	}
	return true
}
func (tf *TRFrame) RunStepInit(curTimeMs int64) bool {
	if !tf.curWorkNode.RunStepInit(curTimeMs) {
		return false
	}
	tf.registerFrameHandler()
	if tf.userStepRun[ETRFrameStepInit] != nil {
		return tf.userStepRun[ETRFrameStepInit](curTimeMs)
	}
	return true
}
func (tf *TRFrame) RunStepPreRun(curTimeMs int64) bool {
	loghlp.Info("TRFrame::RunStepPreRun")

	if !tf.curWorkNode.RunStepPreRun(curTimeMs) {
		return false
	}
	if tf.userStepRun[ETRFrameStepPreRun] != nil {
		return tf.userStepRun[ETRFrameStepPreRun](curTimeMs)
	}
	return true
}
func (tf *TRFrame) RunStepRun(curTimeMs int64) bool {
	tf.curWorkNode.RunStepRun(curTimeMs)
	for _, fn := range tf.loopFuncList {
		fn(curTimeMs)
	}
	tf.updateKeepNodeAlive(curTimeMs)
	return true
}
func (tf *TRFrame) RunStepStop(curTimeMs int64) bool {
	loghlp.Info("TRFrame::RunStepStop")
	if !tf.curWorkNode.RunStepStop(curTimeMs) {
		return false
	}
	if tf.userStepRun[ETRFrameStepStop] != nil {
		return tf.userStepRun[ETRFrameStepStop](curTimeMs)
	}
	return true
}
func (tf *TRFrame) RunStepEnd(curTimeMs int64) bool {
	loghlp.Info("TRFrame::RunStepEnd")
	if !tf.curWorkNode.RunStepEnd(curTimeMs) {
		return false
	}
	if tf.userStepRun[ETRFrameStepEnd] != nil {
		return tf.userStepRun[ETRFrameStepEnd](curTimeMs)
	}
	return true
}
