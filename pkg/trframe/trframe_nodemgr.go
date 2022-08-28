/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-06-15 14:14:17
 * @LastEditTime: 2022-06-15 14:14:17
 * @Brief:
 */
package trframe

import (
	"errors"

	"roomcell/pkg/loghlp"
	"roomcell/pkg/trframe/trnode"
)

type FrameNodeMgr struct {
	nodeList []trnode.ITRNodeEntity
}

func NewFrameNodeMgr() *FrameNodeMgr {
	return &FrameNodeMgr{
		nodeList: make([]trnode.ITRNodeEntity, 0),
	}
}

// 增加节点
func (mgr *FrameNodeMgr) AddNode(frameNode trnode.ITRNodeEntity) {
	mgr.nodeList = append(mgr.nodeList, frameNode)
}

// 查找节点
func (mgr *FrameNodeMgr) FindNode(zoneID int32, nodeType int32, nodeIndex int32) (trnode.ITRNodeEntity, error) {
	var err error = nil
	var ret trnode.ITRNodeEntity = nil
	for _, v := range mgr.nodeList {
		if v.Equal(zoneID, nodeType, nodeIndex) {
			ret = v
			err = nil
			break
		}
	}
	if ret == nil {
		err = errors.New("not find node")
	}
	return ret, err
}

// 删除节点
// func (mgr *FrameNodeMgr) RemoveNode(zoneID int32, nodeType int32, nodeIndex int32) {
// 	for idx, v := range mgr.nodeList {
// 		if v.Equal(zoneID, nodeType, nodeIndex) {
// 			// 移动到最后
// 			mgr.nodeList[idx] = mgr.nodeList[len(mgr.nodeList)-1]
// 			mgr.nodeList = mgr.nodeList[0:(len(mgr.nodeList) - 1)]
// 			break
// 		}
// 	}
// }

func (mgr *FrameNodeMgr) RemoveNode2(sessionID int32) {
	for idx, v := range mgr.nodeList {
		if v.GetSessionID() == sessionID {
			// 移动到最后
			mgr.nodeList[idx] = mgr.nodeList[len(mgr.nodeList)-1]
			mgr.nodeList = mgr.nodeList[0:(len(mgr.nodeList) - 1)]
			loghlp.Infof("removeFrameNode2,sessionID:%d,nodeType:%d,nodeIndex:%d,nodeInfo:%s",
				sessionID,
				v.GetNodeInfo().NodeType,
				v.GetNodeInfo().NodeIndex,
				v.GetNodeInfo().DesInfo)
			break
		}
	}
}

// 根据节点类型获取节点列表
func (mgr *FrameNodeMgr) GetNodeListByType(nodeType int32) []trnode.ITRNodeEntity {
	nl := make([]trnode.ITRNodeEntity, 0)
	for _, v := range mgr.nodeList {
		if v.GetNodeInfo().NodeType == nodeType {
			nl = append(nl, v)
		}
	}
	return nl
}
