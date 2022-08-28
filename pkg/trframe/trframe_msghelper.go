/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-07-15 14:14:17
 * @Brief:当前工作节点
 */
package trframe

import (
	"reflect"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbtools"
	"roomcell/pkg/trframe/iframe"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// 发送消息(本区内)
func (tf *TRFrame) forwardZoneMessage(msgClass int32, msgType int32, msgData []byte, nodeTpye int32, nodeIndex int32, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv) bool {
	sendMsg := MakeInnerMsg(msgClass, msgType, msgData)
	// 找到节点
	frameNode, err := tf.frameNodeMgr.FindNode(tf.frameConfig.ZoneID, nodeTpye, nodeIndex)
	if frameNode == nil || err != nil {
		loghlp.Errorf("forward message fail, not find node(%d,%d)", nodeTpye, nodeIndex)
		return false
	}
	callInfo := tf.remoteMsgMgr.makeCallInfo(msgClass, msgType, cb, env)
	setSecondHead(sendMsg, 0, callInfo.reqID, 0)
	frameNode.SendMsg(sendMsg)
	loghlp.Debugf("forwardMessage,toNode(%s),msg(%d_%d),reqID:%d",
		frameNode.GetNodeInfo().DesInfo,
		msgClass,
		msgType,
		callInfo.reqID,
	)
	return true
}
func (tf *TRFrame) ForwardMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, nodeTpye int32, nodeIndex int32, cb iframe.MsgCallbackFunc, env *iframe.TRRemoteMsgEnv) bool {
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		loghlp.Errorf("forward message fail, proto marsh error:%s", err.Error())
		return false
	}
	return tf.forwardZoneMessage(msgClass, msgType, data, nodeTpye, nodeIndex, cb, env)
}
func (tf *TRFrame) pushMessage(msgClass int32, msgType int32, msgData []byte, nodeTpye int32, nodeIndex int32) bool {
	sendMsg := MakeInnerMsg(msgClass, msgType, msgData)
	// 找到节点
	frameNode, err := tf.frameNodeMgr.FindNode(tf.frameConfig.ZoneID, nodeTpye, nodeIndex)
	if frameNode == nil || err != nil {
		loghlp.Errorf("push message fail, not find node(%d,%d)", nodeTpye, nodeIndex)
		return false
	}
	frameNode.SendMsg(sendMsg)
	loghlp.Debugf("push message,toNode(%s),msg(%d_%d),reqID:%d",
		frameNode.GetNodeInfo().DesInfo,
		msgClass,
		msgType,
	)
	return true
}
func (tf *TRFrame) PushZoneMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, nodeTpye int32, nodeIndex int32) bool {
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		loghlp.Errorf("PushZoneMessage fail, proto marsh error:%s", err.Error())
		return false
	}
	return tf.pushMessage(msgClass, msgType, data, nodeTpye, nodeIndex)
}

// 发送回复消息
func (tf *TRFrame) SendReplyMessage(okCode int32, pbRep proto.Message, env *iframe.TRRemoteMsgEnv) {
	env.SrcMessage.Head.Result = uint16(okCode)
	if pbRep != nil {
		env.SrcMessage.Data, _ = EncodePBMessage(pbRep)
	}
	if env.SrcMessage.SecondHead != nil {
		env.SrcMessage.SecondHead.RepID = env.SrcMessage.SecondHead.ReqID
	}
	session := tf.evHub.GetSession(env.SrcSessionID)
	if session != nil {
		session.Send(env.SrcMessage)
		loghlp.Debugf("SendReplyMessage(%d_%d)[%s]:%+v",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
			pbtools.GetFullNameByMessage(pbRep),
			pbRep,
		)
	} else {
		loghlp.Errorf("SendReplyMessage(%d_%d),but not find session(%d)",
			env.SrcMessage.Head.MsgClass,
			env.SrcMessage.Head.MsgType,
			env.SrcSessionID,
		)
	}
}

func MakePBMessage(msgClass int32, msgType int32, pbMsg protoreflect.ProtoMessage, okCode int32) *evhub.NetMessage {
	data, err := proto.Marshal(pbMsg)
	if err == nil {
		evMsg := evhub.MakeMessage(msgClass, msgType, data)
		return evMsg
	} else {
		loghlp.Errorf("make pb message error:%s", err.Error())
	}
	return nil
}

func DecodePBMessage(sMsg *evhub.NetMessage, pbMsg protoreflect.ProtoMessage) bool {
	errParse := proto.Unmarshal(sMsg.Data, pbMsg)
	if errParse != nil {
		loghlp.Errorf("decode PB(%s) error:%s", reflect.TypeOf(pbMsg).String(), errParse.Error())
		return false
	}
	return true
}
func DecodePBMessage2(msgData []byte, pbMsg protoreflect.ProtoMessage) bool {
	errParse := proto.Unmarshal(msgData, pbMsg)
	if errParse != nil {
		loghlp.Errorf("decode PB(%s) error:%s", reflect.TypeOf(pbMsg).String(), errParse.Error())
		return false
	}
	return true
}
func EncodePBMessage(pbMsg proto.Message) ([]byte, error) {
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		loghlp.Errorf("EncodePBMessageError:%s", err.Error())
	}

	return data, err
}
func LogMsgInfo(netMsg *evhub.NetMessage, pbMsg proto.Message) {
	loghlp.Debugf("recv reqmsg(%d_%d)[%s]:%+v",
		netMsg.Head.MsgClass,
		netMsg.Head.MsgType,
		pbtools.GetFullNameByMessage(pbMsg),
		pbMsg,
	)
}
func LogCbMsgInfo(netMsg *evhub.NetMessage, pbMsg proto.Message) {
	loghlp.Debugf("recv cbmsg(%d_%d)[%s]:%+v",
		netMsg.Head.MsgClass,
		netMsg.Head.MsgType,
		pbtools.GetFullNameByMessage(pbMsg),
		pbMsg,
	)
}
