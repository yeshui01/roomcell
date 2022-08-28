package tserver

import "sync"

// netnode mgr
var nodeLock sync.Mutex
var nodeMgr *ServerNodeMgr

func NetNodeMgr() *ServerNodeMgr {
	if nodeMgr == nil {
		nodeLock.Lock()
		if nodeMgr == nil {
			nodeMgr = NewServerNodeMgr()
		}
		nodeLock.Unlock()
	}
	return nodeMgr
}
