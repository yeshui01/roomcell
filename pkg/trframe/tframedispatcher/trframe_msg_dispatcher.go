package tframedispatcher

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbtools"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe/iframe"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type FrameMsgDispatcher struct {
	frameInstance iframe.ITRFrame
	handlerList   map[int32]iframe.FrameMsgHandler
}

func (dsp *FrameMsgDispatcher) RegisterMsgHandler(msgClass int32, msgType int32, msgHander iframe.FrameMsgHandler) {
	handleKey := dsp.genHandlerKey(msgClass, msgType)
	dsp.handlerList[handleKey] = msgHander
}

func (dsp *FrameMsgDispatcher) Dispatch(session iframe.ISession, msg *evhub.NetMessage, customData interface{}) bool {
	handleKey := dsp.genHandlerKey(int32(msg.Head.MsgClass), int32(msg.Head.MsgType))
	msgHandler, ok := dsp.handlerList[handleKey]
	if ok {
		loghlp.Debugf("DispatchMsg(%d_%d)",
			msg.Head.MsgClass,
			msg.Head.MsgType,
		)
		tmsgCtx := &iframe.TMsgContext{
			FrameInstance: dsp.frameInstance,
			Session:       session,
			NetMessage:    msg,
			CustomData:    customData,
		}
		okCode, retData, rt := msgHandler(tmsgCtx)
		if rt == iframe.EHandleContent {
			pbRep, ok := retData.(protoreflect.ProtoMessage)
			if ok {
				msg.Head.Result = uint16(okCode)
				loghlp.Debugf("reply repmsg(%d_%d),isok(%d)[%s]:%+v",
					msg.Head.MsgClass,
					msg.Head.MsgType,
					okCode,
					pbtools.GetFullNameByMessage(pbRep),
					pbRep,
				)
				// pb
				data, err := iframe.ToPbData(pbRep)
				if err == nil {
					msg.Data = data
				}
				if msg.Head.HasSecond > 0 {
					msg.SecondHead.RepID = msg.SecondHead.ReqID
				}
			} else {
				loghlp.Errorf("retData to pb message(%d,%d) fail!!!!!", msg.Head.MsgClass, msg.Head.MsgType)
				msg.Head.Result = protocol.ECodeSysError
				loghlp.Debugf("reply repmsg2(%d_%d),isok(%d)[%s]",
					msg.Head.MsgClass,
					msg.Head.MsgType,
					okCode,
					pbtools.GetFullNameByMessage(pbRep),
				)
				if msg.Head.HasSecond > 0 {
					msg.SecondHead.RepID = msg.SecondHead.ReqID
				}
			}
			session.SendMsg(msg)
		}
	} else {
		return false
	}
	return true
}

func (dsp *FrameMsgDispatcher) genHandlerKey(msgClass int32, msgType int32) int32 {
	return msgClass*1000 + msgType
}

func NewFrameMsgDispatcher(frameObj iframe.ITRFrame) *FrameMsgDispatcher {
	return &FrameMsgDispatcher{
		frameInstance: frameObj,
		handlerList:   make(map[int32]iframe.FrameMsgHandler, 0),
	}
}
