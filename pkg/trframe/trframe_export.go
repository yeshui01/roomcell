/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief:接口导出
 */
package trframe

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/tframeconfig"
	"roomcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	frameCore *TRFrame = nil
)

func Init(frameConfigPath string, nodeType int32, nodeIndex int32) {
	if frameCore == nil {
		frameCore = newTRFrame(frameConfigPath, nodeType, nodeIndex)
	}
}

func Start() {
	frameCore.Run()
}
func Stop() {
	frameCore.evHub.Stop()
}
func GetFrameConfig() *tframeconfig.FrameConfig {
	return frameCore.frameConfig
}
func GetCurNodeIndex() int32 {
	return frameCore.curWorkNode.NodeIndex()
}

// 注册用户层帧函数
func RegisterUserFrameRun(frameRun func(curTimeMs int64)) {
	frameCore.loopFuncList = append(frameCore.loopFuncList, frameRun)
}

// 注册阶段处理函数
func RegisterUserStepRun(step int32, stepRun func(curTimeMs int64) bool) {
	if step < 0 || step >= int32(ETRFrameStepFinal) {
		loghlp.Errorf("step out of range for register user step run!!!")
		return
	}
	frameCore.userStepRun[step] = stepRun
}

// 注册业务消息处理
func RegWorkMsgHandler(msgClass int32, msgType int32, msgHander iframe.FrameMsgHandler) {
	frameCore.msgDispatcher.RegisterMsgHandler(msgClass, msgType, msgHander)
}

// 注册用户命令回调
func RegUserCommandHandler(callback func(userCmd *TRFrameCommand)) {
	frameCore.userCmdHandle = callback
}

// 发送用户命令
func PostUserCommand(cmdClass int32, cmdType int32, cmdData interface{}) {
	frameCore.evHub.PostCommand(cmdClass, cmdType, cmdData)
}

// 发送消息
func ForwardZoneMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, nodeTpye int32, nodeIndex int32, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv) bool {
	return frameCore.ForwardMessage(msgClass, msgType, pbMsg, nodeTpye, nodeIndex, cb, env)
}

// 推送消息
func PushZoneMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, nodeTpye int32, nodeIndex int32) bool {
	return frameCore.PushZoneMessage(msgClass, msgType, pbMsg, nodeTpye, nodeIndex)
}
func ForwardZoneClientPBMessage(msgClass int32, msgType int32, pbMsg proto.Message, nodeTpye int32, nodeIndex int32, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv, roleID int64) bool {
	msgData, _ := proto.Marshal(pbMsg)
	return ForwardZoneClientMessage(msgClass, msgType, msgData, nodeTpye, nodeIndex, cb, env, roleID)
}
func ForwardZoneClientMessage(msgClass int32, msgType int32, msgData []byte, nodeTpye int32, nodeIndex int32, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv, roleID int64) bool {
	sendMsg := MakeInnerMsg(msgClass, msgType, msgData)
	// 找到节点
	frameNode, err := frameCore.frameNodeMgr.FindNode(frameCore.frameConfig.ZoneID, nodeTpye, nodeIndex)
	if frameNode == nil || err != nil {
		loghlp.Errorf("forward message fail, not find node(%d,%d)", nodeTpye, nodeIndex)
		return false
	}
	callInfo := frameCore.remoteMsgMgr.makeCallInfo(msgClass, msgType, cb, env)
	setSecondHead(sendMsg, roleID, callInfo.reqID, 0)
	frameNode.SendMsg(sendMsg)
	return true
}
func pushZoneClientMessage(msgClass int32, msgType int32, msgData []byte, nodeTpye int32, nodeIndex int32, roleID int64) bool {
	sendMsg := MakeInnerMsg(msgClass, msgType, msgData)
	// 找到节点
	frameNode, err := frameCore.frameNodeMgr.FindNode(frameCore.frameConfig.ZoneID, nodeTpye, nodeIndex)
	if frameNode == nil || err != nil {
		loghlp.Errorf("forward message fail, not find node(%d,%d)", nodeTpye, nodeIndex)
		return false
	}
	setSecondHead(sendMsg, roleID, 0, 0)
	frameNode.SendMsg(sendMsg)
	return true
}
func PushZoneClientPBMessage(msgClass int32, msgType int32, pbMsg proto.Message, nodeTpye int32, nodeIndex int32, roleID int64) bool {
	msgData, _ := proto.Marshal(pbMsg)
	return pushZoneClientMessage(msgClass, msgType, msgData, nodeTpye, nodeIndex, roleID)
}
func BroadcastMessage(msgClass int32, msgType int32, pbMsg proto.Message, nodeTpye int32) bool {
	msgData, err := proto.Marshal(pbMsg)
	if err != nil {
		return false
	}
	nodeList := frameCore.frameNodeMgr.GetNodeListByType(nodeTpye)
	for _, oneNode := range nodeList {
		curMsgData := make([]byte, len(msgData))
		copy(curMsgData, msgData)
		sendMsg := MakeInnerMsg(msgClass, msgType, curMsgData)
		oneNode.SendMsg(sendMsg)
	}
	return true
}
func SendReplyMessage(okCode int32, pbRep proto.Message, env *iframe.TRRemoteMsgEnv) {

	frameCore.SendReplyMessage(okCode, pbRep, env)
}

func DispatchMsg(session iframe.ISession, msg *evhub.NetMessage, customData interface{}) bool {
	return frameCore.msgDispatcher.Dispatch(session, msg, customData)
}

func GetNodeListByType(nodeType int32) []trnode.ITRNodeEntity {
	return frameCore.frameNodeMgr.GetNodeListByType(nodeType)
}

// 当前系统时间,非即时时间
func GetFrameSysNowTime() int64 {
	return frameCore.nowFrameTimeMs / 1000
}
