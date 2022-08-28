package thalldata

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

func (frameNode *FrameNodeHallData) InitConnectServer() bool {
	// 连接hallserver
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	evHub := frameNode.tframeObj.GetEvHub()

	for idx, cfg := range frameConfig.HallServerCfgs {
		connDo := func() {
			// 发送注册消息
			reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
				ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
				NodeType:  trnode.ETRNodeTypeHallData,
				NodeIndex: frameNode.nodeIndex,
				NodeDes:   fmt.Sprintf("%dETRNodeTypeHallData%d", frameNode.tframeObj.GetFrameConfig().ZoneID, idx),
			}
			cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
				loghlp.Infof("data register to hallserver callback suss:%d", okCode)
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
				DesInfo:        fmt.Sprintf("ETRNodeTypeHallServer%d", frameNode.nodeIndex),
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

	return true
}
