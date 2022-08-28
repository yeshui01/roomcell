package trframe

import (
	"roomcell/pkg/loghlp"
	"roomcell/pkg/pb/pbframe"
	"roomcell/pkg/protocol"
	"roomcell/pkg/trframe/iframe"
	"roomcell/pkg/trframe/trnode"

	"google.golang.org/protobuf/proto"
)

// 协议处理
func handleRegisterNodeInfo(tmsgCtx *iframe.TMsgContext) (isok int32, retData interface{}, rt iframe.IHandleResultType) {
	req := &pbframe.FrameMsgRegisterServerInfoReq{}
	rep := &pbframe.FrameMsgRegisterServerInfoRep{}
	proto.Unmarshal(tmsgCtx.NetMessage.Data, req)
	frameSession := tmsgCtx.Session.(*FrameSession)
	frameSession.nodeType = req.NodeType
	// 关联节点信息
	frameSession.nodeInfo = &trnode.TRNodeInfo{
		ZoneID:    req.ZoneID,
		NodeType:  req.NodeType,
		NodeIndex: req.NodeIndex,
		DesInfo:   req.NodeDes,
	}
	frameCore.frameNodeMgr.AddNode(frameSession)
	loghlp.Infof("handleRegisterNodeInfo,req:%+v", req)
	return protocol.ECodeSuccess, rep, iframe.EHandleContent
}
