package hallroomhandler

import (
	"roomcell/app/hallroom/hallroommain"
	"roomcell/app/hallroom/hallroommain/gameroom"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbclient"
	"roomcell/pkg/pb/pbserver"
	"roomcell/pkg/protocol"
	"roomcell/pkg/sconst"
	"roomcell/pkg/trframe"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

// 创建房间
func HandleRoomCreate(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomCreateReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgRoomCreateRep{
		ZoneID:    trframe.GetFrameConfig().ZoneID,
		NodeType:  trnode.ETRNodeTypeHallRoom,
		NodeIndex: trframe.GetCurNodeIndex(),
	}
	hallRoomGlobal := roomServe.GetGlobalData()
	roomPlayer := hallRoomGlobal.FindRoomPlayer(req.RoleID)
	if roomPlayer != nil {
		return protocol.ECodeRoomPlayerHasInRoom, rep, iframe.EHandleContent
	}
	roomObj := hallRoomGlobal.CreateRoom(req.RoomID, req.RoomType, req.RoleID)
	if roomObj == nil {
		return protocol.ECodeRoomCreateFail, rep, iframe.EHandleContent
	}

	if roomPlayer == nil {
		roomPlayer = hallroommain.NewRoomPlayer(req.RoleID, req.RoleName)
		hallRoomGlobal.AddRoomPlayer(roomPlayer)
	}
	roomPlayer.GateNode = &trnode.TRNodeInfo{
		ZoneID:    req.GateNode.ZoneID,
		NodeType:  req.GateNode.NodeType,
		NodeIndex: req.GateNode.NodeIndex,
	}
	roomPlayer.HallNode = &trnode.TRNodeInfo{
		ZoneID:    req.HallNode.ZoneID,
		NodeType:  req.HallNode.NodeType,
		NodeIndex: req.HallNode.NodeIndex,
	}
	roomObj.JoinPlayer(roomPlayer)
	roomPlayer.SetRoomID(req.RoomID)
	roomPlayer.SetRoomPtr(roomObj)
	// 房间信息返回
	rep.RoomDetail = roomObj.ToRoomDetail()
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 玩家离开房间
func HandlePlayerRoomLeave(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomLeaveReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgRoomLeaveReq{}
	hallRoomGlobal := roomServe.GetGlobalData()
	hallPlayer := hallRoomGlobal.FindRoomPlayer(req.RoleID)
	if hallPlayer != nil {
		roomID := hallPlayer.GetRoomID()
		roomObj := hallRoomGlobal.FindRoom(roomID)
		if roomObj != nil {
			roomObj.LeavePlayer(req.RoleID)
		} else {
			loghlp.Warnf("not find room(%d)", roomID)
		}
		hallRoomGlobal.DelRoomPlayer(hallPlayer)
	} else {
		loghlp.Warnf("not find room player(%d)", req.RoleID)
	}

	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 加入房间
func HandlePlayerRoomEnter(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgRoomEnterReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgRoomEnterRep{
		ZoneID:    trframe.GetFrameConfig().ZoneID,
		NodeType:  trnode.ETRNodeTypeHallRoom,
		NodeIndex: trframe.GetCurNodeIndex(),
	}
	hallRoomGlobal := roomServe.GetGlobalData()
	roomPlayer := hallRoomGlobal.FindRoomPlayer(req.RoleID)
	if roomPlayer != nil {
		return protocol.ECodeRoomPlayerHasInRoom, rep, iframe.EHandleContent
	}
	roomObj := hallRoomGlobal.FindRoom(req.RoomID)
	if roomObj == nil {
		return protocol.ECodeRoomNotExisted, rep, iframe.EHandleContent
	}
	if roomObj.IsPlayerFull() {
		return protocol.ECodeRoomMaxPlayerNumLimit, rep, iframe.EHandleContent
	}
	// 游戏中途不能进入
	if !roomObj.IsCanJoin() {
		return protocol.ECodeRoomCantjoin, rep, iframe.EHandleContent
	}
	if roomPlayer == nil {
		roomPlayer = hallroommain.NewRoomPlayer(req.RoleID, req.RoleName)
		hallRoomGlobal.AddRoomPlayer(roomPlayer)
	}

	roomPlayer.GateNode = &trnode.TRNodeInfo{
		ZoneID:    req.GateNode.ZoneID,
		NodeType:  req.GateNode.NodeType,
		NodeIndex: req.GateNode.NodeIndex,
	}
	roomPlayer.HallNode = &trnode.TRNodeInfo{
		ZoneID:    req.HallNode.ZoneID,
		NodeType:  req.HallNode.NodeType,
		NodeIndex: req.HallNode.NodeIndex,
	}
	roomObj.JoinPlayer(roomPlayer)
	roomPlayer.SetRoomID(req.RoomID)
	roomPlayer.SetRoomPtr(roomObj)
	// 房间信息返回
	rep.RoomDetail = roomObj.ToRoomDetail()
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 玩家离线
func HandleRoomPlayerOffline(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbserver.ESMsgPlayerDisconnectReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbserver.ESMsgPlayerDisconnectRep{}
	hallRoomGlobal := roomServe.GetGlobalData()
	loghlp.Infof("room player(%d) offline, reason(%d)", req.RoleID, req.Reason)
	roomPlayer := hallRoomGlobal.FindRoomPlayer(req.RoleID)
	if roomPlayer != nil {
		roomObj := hallRoomGlobal.FindRoom(roomPlayer.GetRoomID())
		if roomObj != nil {
			roomObj.OnPlayerOffline(roomPlayer)
		}
		roomPlayer.SetRoomID(0)
		hallRoomGlobal.DelRoomPlayer(roomPlayer)
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}

// 玩家聊天
func HandleRoomPlayerChat(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbclient.ECMsgRoomChatReq{}
	if !trframe.DecodePBMessage(tmsgCtx.NetMessage, req) {
		return protocol.ECodePBDecodeError, nil, iframe.EHandleContent
	}
	trframe.LogMsgInfo(tmsgCtx.NetMessage, req)
	rep := &pbclient.ECMsgRoomChatRsp{}

	hallRoomGlobal := roomServe.GetGlobalData()
	// 屏蔽字处理
	wordsUtil := hallRoomGlobal.GetSensitiveWordsUtil()
	if wordsUtil != nil {
		req.TalkContent = wordsUtil.HandleWord(req.TalkContent, 'x')
	}
	roomPlayer := hallRoomGlobal.FindRoomPlayer(tmsgCtx.NetMessage.SecondHead.ID)
	if roomPlayer != nil {
		if roomPlayer.RoomPtr != nil {
			switch roomPlayer.RoomPtr.GetRoomType() {
			case sconst.EGameRoomTypeDrawGuess:
				{
					roomDrawGuess, ok := roomPlayer.RoomPtr.(*gameroom.RoomDrawGuess)
					if ok {
						if roomDrawGuess.RoomStep == sconst.EDrawGuessStepDraw {
							// 聊天猜词
							return HandlePlayerChatGuessWords(tmsgCtx)
						}
					}
					break
				}
			case sconst.EGameRoomTypeUndercover:
				{
					roomUndercover, ok := roomPlayer.RoomPtr.(*gameroom.RoomUndercover)
					if ok {
						if roomUndercover.RoomStep == sconst.EUndercoverStepTalk {
							// 谁是卧底发言
							return HandlePlayerChatUndertalk(tmsgCtx)
						}
					}
					break
				}
			}
		}

		// 推送聊天
		pushMsg := &pbclient.ECMsgRoomPushChatNotify{
			TalkContent: req.TalkContent,
			Talker: &pbclient.RoomTalker{
				RoleID:   roomPlayer.GetRoleID(),
				Nickname: roomPlayer.Nickname,
				Icon:     roomPlayer.Icon,
			},
		}
		if roomPlayer.RoomPtr != nil {
			roomPlayer.RoomPtr.BroadCastRoomMsg(0, protocol.ECMsgClassRoom, protocol.ECMsgRoomPushChat, pushMsg)
		}
	}
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
