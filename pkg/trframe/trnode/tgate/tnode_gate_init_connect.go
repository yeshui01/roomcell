package tgate

import (
	"fmt"
	"time"

	"roomcell/pkg/evhub"
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbframe"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"
)

func (frameNode *FrameNodeGate) InitConnectRoot() bool {
	// 连接root
	frameConfig := frameNode.tframeObj.GetFrameConfig()
	evHub := frameNode.tframeObj.GetEvHub()
	connDo := func() {
		// 发送注册消息
		reqMsg := &pbframe.FrameMsgRegisterServerInfoReq{
			ZoneID:    frameNode.tframeObj.GetFrameConfig().ZoneID,
			NodeType:  trnode.ETRNodeTypeGate,
			NodeIndex: frameNode.nodeIndex,
			NodeDes:   fmt.Sprintf("ETRNodeTypeGate%d", frameNode.nodeIndex),
		}
		cb := func(okCode int32, msgData []byte, env *iframe.TRRemoteMsgEnv) {
			loghlp.Infof("gate register to root callback suss:%d", okCode)
		}
		frameNode.tframeObj.ForwardMessage(
			protocol.EMsgClassFrame,
			protocol.EFrameMsgRegisterServerInfo,
			reqMsg,
			trnode.ETRNodeTypeRoot,
			0,
			cb,
			nil,
		)
	}
	for _, cfg := range frameConfig.RootCfgs {
		if cfg.ListenMode != "unix" {
			userData := &iframe.SessionUserData{
				DataType:  iframe.ESessionDataTypeNetInfo,
				NodeType:  trnode.ETRNodeTypeRoot,
				NodeIndex: frameNode.nodeIndex,
				DesInfo:   fmt.Sprintf("root_%d", frameNode.nodeIndex),
			}
			failCount := 0
			for {
				if !evHub.Connect2(evhub.ListenModeTcp, cfg.ListenAddr, true, userData, connDo) {
					failCount++
					loghlp.Warnf("connect to root fail")
				} else {
					break
				}
				if failCount >= 10 {
					loghlp.Errorf("connect root fail,exit")
					return false
				}
				time.Sleep(time.Second * 1)
			}
		}
	}

	return true
}
