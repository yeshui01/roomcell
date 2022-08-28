/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-07-15 14:14:17
 * @Brief:远程消息管理
 */
package trframe

import (
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/timeutil"
	"roomcell/pkg/trframe/iframe"
)

// 调用信息
type TRRemoteCallInfo struct {
	env       *iframe.TRRemoteMsgEnv
	cbFun     iframe.MsgCallbackFunc
	msgClass  int32
	msgType   int32
	beginTime int64
	reqID     uint64
}

func MakeInnerMsg(msgClass int32, msgType int32, data []byte) *evhub.NetMessage {
	return evhub.MakeMessage(msgClass, msgType, data)
}
func MakeMsgEnv(srcssID int32, srcMsg *evhub.NetMessage) *iframe.TRRemoteMsgEnv {
	return &iframe.TRRemoteMsgEnv{
		SrcSessionID: srcssID,
		SrcMessage:   srcMsg,
	}
}
func MakeClientMsgEnv(userID int64, srcMsg *evhub.NetMessage, userData interface{}) *iframe.TRRemoteMsgEnv {
	e := &iframe.TRRemoteMsgEnv{
		SrcSessionID: 0,
		SrcMessage:   srcMsg,
		UserData:     userData,
	}
	if srcMsg.SecondHead == nil {
		srcMsg.SecondHead = &evhub.NetMsgSecondHead{
			ID: userID,
		}
	} else {
		srcMsg.SecondHead.ID = userID
	}
	return e
}
func MakeMsgEnv2(iss iframe.ISession, srcMsg *evhub.NetMessage) *iframe.TRRemoteMsgEnv {
	if iss == nil {
		return MakeMsgEnv(0, srcMsg)
	}
	frameSession := iss.(*FrameSession)
	return &iframe.TRRemoteMsgEnv{
		SrcSessionID: frameSession.GetSessionID(),
		SrcMessage:   srcMsg,
	}
}
func setSecondHead(msg *evhub.NetMessage, id int64, reqID uint64, repID uint64) {
	msg.Head.HasSecond = 1
	msg.SecondHead = &evhub.NetMsgSecondHead{
		ID:    id,
		ReqID: reqID,
		RepID: repID,
	}
}

type RemoteMsgCallMgr struct {
	callList    map[uint64]*TRRemoteCallInfo
	reqID       uint64
	frameCore   *TRFrame
	lastUpdTime int64 // 单位毫秒
}

func newRemoteMsgMgr(frameObj *TRFrame) *RemoteMsgCallMgr {
	return &RemoteMsgCallMgr{
		frameCore: frameObj,
		reqID:     0,
		callList:  make(map[uint64]*TRRemoteCallInfo),
	}
}

func (mgr *RemoteMsgCallMgr) update(curTimeMs int64) {
	if curTimeMs-mgr.lastUpdTime < 10000 {
		return
	}
	// 检测超时
	for k, v := range mgr.callList {
		if curTimeMs-v.beginTime >= 10000 {
			loghlp.Errorf("remote call timeout(%d,%d)reqId:%d", v.msgClass, v.msgType, k)
			delete(mgr.callList, k)
		}
	}
}

// 处理回调消息
func (mgr *RemoteMsgCallMgr) checkHandleCallbackMsg(msg *evhub.NetMessage) bool {
	if msg.Head.HasSecond < 1 {
		return false
	}
	if msg.SecondHead.RepID > 0 {
		callInfo, ok := mgr.callList[msg.SecondHead.RepID]
		if ok {
			loghlp.Debugf("handleCallbackMsg(%d_%d),repID:%d", msg.Head.MsgClass, msg.Head.MsgType, msg.SecondHead.RepID)
			callInfo.cbFun(int32(msg.Head.Result), msg.Data, callInfo.env)
			delete(mgr.callList, msg.SecondHead.RepID)
		} else {
			loghlp.Errorf("not find callback info,maybe has timeout,repID:%d", msg.SecondHead.RepID)
		}
		return true
	}
	return false
}

func (mgr *RemoteMsgCallMgr) genReqID() uint64 {
	if mgr.reqID >= 0x00ffffffffffffff {
		mgr.reqID = 0
	}
	mgr.reqID++
	return mgr.reqID
}

// 生成调用信息记录
func (mgr *RemoteMsgCallMgr) makeCallInfo(msgClass int32, msgType int32, cb iframe.MsgCallbackFunc, e *iframe.TRRemoteMsgEnv) *TRRemoteCallInfo {
	callInfo := &TRRemoteCallInfo{
		env:       e,
		reqID:     mgr.genReqID(),
		msgClass:  msgClass,
		msgType:   msgType,
		beginTime: timeutil.NowTimeMs(),
		cbFun:     cb,
	}
	mgr.callList[callInfo.reqID] = callInfo
	return callInfo
}

// 发送注册信息到某个节点
