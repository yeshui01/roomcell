package tserver

const (
	TServerNodeTypeGate = 1 // 网关服服务器
	TServerNodeTypeRoom = 1 // 房间服务器
	TServerNodeTypeHall = 2 // 大厅服务器

	TSerrverNodeAccount = 100 // 账号服务器
)

// 节点
type ServerNode struct {
	ZoneID    int32
	NodeType  int32
	NodeIndex int32
	Session   *TServerSession
}

func NewServerNode(zoneID int32, nodeType int32, nodeIndex int32, sess *TServerSession) *ServerNode {
	return &ServerNode{
		ZoneID:    zoneID,
		NodeType:  nodeType,
		NodeIndex: nodeIndex,
		Session:   sess,
	}
}

// 节点管理器
type ServerNodeMgr struct {
	typeNodes map[int32][]*ServerNode
}

func NewServerNodeMgr() *ServerNodeMgr {
	return &ServerNodeMgr{
		typeNodes: make(map[int32][]*ServerNode),
	}
}

func (mgr *ServerNodeMgr) AddNode(nd *ServerNode) {
	slist := mgr.typeNodes[nd.NodeType]
	slist = append(slist, nd)
	mgr.typeNodes[nd.NodeType] = slist
}
func (mgr *ServerNodeMgr) AddServerNode(zoneID int32, nodeType int32, nodeIndex int32, sess *TServerSession) {
	nd := NewServerNode(zoneID, nodeType, nodeIndex, sess)
	slist := mgr.typeNodes[nd.NodeType]
	slist = append(slist, nd)
	mgr.typeNodes[nd.NodeType] = slist
}
func (mgr *ServerNodeMgr) FindNode(zoneID int32, nodeType int32, nodeIndex int32) *ServerNode {
	typeNodeList, ok := mgr.typeNodes[nodeType]
	if !ok {
		return nil
	}
	var nd *ServerNode = nil
	for _, v := range typeNodeList {
		if v.ZoneID == zoneID && v.NodeIndex == v.NodeIndex {
			nd = v
			break
		}
	}
	return nd
}

func (mgr *ServerNodeMgr) RemoveNode(zoneID int32, nodeType int32, nodeIndex int32) {
	typeNodeList, ok := mgr.typeNodes[nodeType]
	if !ok {
		return
	}
	var nd *ServerNode = nil
	for i, v := range typeNodeList {
		if v.ZoneID == zoneID && v.NodeIndex == v.NodeIndex {
			nd = v
			typeNodeList[i] = typeNodeList[len(typeNodeList)-1]
			typeNodeList[len(typeNodeList)-1] = nd
			typeNodeList = typeNodeList[0 : len(typeNodeList)-1] // 移除
			break
		}
	}
}
func (mgr *ServerNodeMgr) RemoveNodeBySession(ssion *TServerSession) {
	typeNodeList, ok := mgr.typeNodes[ssion.nodeType]
	if !ok {
		return
	}
	var nd *ServerNode = nil
	for i, v := range typeNodeList {
		if v.Session.netSession != nil && v.Session.GetSessionID() == ssion.GetSessionID() {
			nd = v
			typeNodeList[i] = typeNodeList[len(typeNodeList)-1]
			typeNodeList[len(typeNodeList)-1] = nd
			typeNodeList = typeNodeList[0 : len(typeNodeList)-1] // 移除
			break
		}
	}
}
