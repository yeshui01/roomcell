package thallgate

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

func (frameNode *FrameNodeHallGate) InitConnectServer() bool {
	// 连接hallserver
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	evHub := frameNode.tframeObj.GetEvHub()

	for idx, cfg := range frameConfig.HallServerCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeHallGate,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%d_ETRNodeTypeHallGate%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("gate register to hallserver callback suss:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeHallServer,
				int32(idx),
				cb,
				nil,
			)
		}
		if cfg.ListenMode != "unix" {
			userData := &iframe.SessionUserData{
				DataType:       iframe.ESessionDataTypeNetInfo,
				NodeType:       trnode.ETRNodeTypeHallServer,
				NodeIndex:      frameNode.nodeIndex,
				DesInfo:        fmt.Sprintf("%dETRNodeTypeHallServer%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
				IsServerClient: true,
			}
			failCount := 0
			for {
				if !evHub.Connect2(evhub.ListenModeTcp, cfg.ListenAddr, true, userData, connDo) {
					failCount++
					loghlp.Warnf("connect to hallserver fail")
				} else {
					break
				}
				if failCount >= 10 {
					loghlp.Errorf("connect hallserver fail,exit")
					return false
				}
				time.Sleep(time.Second * 1)
			}
		}
	}
	// 连接 room_server
	for idx, cfg := range frameConfig.HallRoomCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeHallGate,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%dETRNodeTypeHallGate%d", frameNode.tframeObj.GetFrameConfig().ZoneID, frameNode.nodeIndex),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("gate register to hallroom callback suss:%d", okCode)
			}
			frameNode.tframeObj.ForwardMessage(
				protocol.EMsgClassFrame,
				protocol.EFrameMsgRegisterServerInfo,
				reqMsg,
				trnode.ETRNodeTypeHallRoom,
				int32(idx),
				cb,
				nil,
			)
		}
		if cfg.ListenMode != "unix" {
			userData := &iframe.SessionUserData{
				DataType:       iframe.ESessionDataTypeNetInfo,
				NodeType:       trnode.ETRNodeTypeHallRoom,
				NodeIndex:      int32(idx),
				DesInfo:        fmt.Sprintf("%dETRNodeTypeHallRoom%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
				IsServerClient: true,
			}
			failCount := 0
			for {
				if !evHub.Connect2(evhub.ListenModeTcp, cfg.ListenAddr, true, userData, connDo) {
					failCount++
					loghlp.Warnf("connect to hallroom fail")
				} else {
					break
				}
				if failCount >= 10 {
					loghlp.Errorf("connect hallroom fail,exit")
					return false
				}
				time.Sleep(time.Second * 1)
			}
			loghlp.Infof("connect hallroom%d success:%s", idx, cfg.ListenAddr)
		}
	}
	return true
}
