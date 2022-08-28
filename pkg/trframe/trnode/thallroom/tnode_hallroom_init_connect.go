package thallroom

import (
	"fmt"
	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbframe"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
	"time"
)

func (frameNode *FrameNodeRoom) InitConnectServer() bool {
	// 连接 room_mgr
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	evHub := frameNode.tframeObj.GetEvHub()

	for idx, cfg := range frameConfig.HallRoomMgrCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeHallRoom,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%dETRNodeTypeHallRoom%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("room register to hallroom_mgr callback suss:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeHallRoomMgr,
				int32(idx),
				cb,
				nil,
			)
		}
		if cfg.ListenMode != "unix" {
			userData := &iframe.SessionUserData{
				DataType:       iframe.ESessionDataTypeNetInfo,
				NodeType:       trnode.ETRNodeTypeHallRoomMgr,
				NodeIndex:      int32(idx),
				DesInfo:        fmt.Sprintf("%dETRNodeTypeHallServer%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
				IsServerClient: true,
			}
			failCount := 0
			for {
				if !evHub.Connect2(evhub.ListenModeTcp, cfg.ListenAddr, true, userData, connDo) {
					failCount++
					loghlp.Warnf("connect to hallroom_mgr fail")
				} else {
					break
				}
				if failCount >= 10 {
					loghlp.Errorf("connect hallroom_mgr fail,exit")
					return false
				}
				time.Sleep(time.Second * 1)
			}
			loghlp.Infof("connect hallroom_mgr success:%s", cfg.ListenAddr)
		}
	}

	return true
}
