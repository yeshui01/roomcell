package protocol

import "roomcell/pkg/evhub"

const (
	CellCmdClassWebsocket = evhub.HubCmdUserBase + 1
)

// Websock
const (
	CmdTypeWebsocketConnect = 1
	CmdTypeWebsocketClosed  = 2
	CmdTypeWebsocketMessage = 3
)
