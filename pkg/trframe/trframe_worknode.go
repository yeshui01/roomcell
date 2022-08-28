/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief:当前工作节点
 */
package trframe

type ITRFrameWorkNode interface {
	RunStepCheck(curTimeMs int64) bool
	RunStepInit(curTimeMs int64) bool
	RunStepPreRun(curTimeMs int64) bool
	RunStepRun(curTimeMs int64) bool
	RunStepStop(curTimeMs int64) bool
	RunStepEnd(curTimeMs int64) bool
	NodeType() int32
	NodeIndex() int32
	//SetUserFrameRun(func(curTimeMs int64))
}
